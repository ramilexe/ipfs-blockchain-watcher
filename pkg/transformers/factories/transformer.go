// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package factories

import (
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Transformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Repository Repository
	Fetcher    shared.SettableLogFetcher
}

func (transformer Transformer) NewTransformer(db *postgres.DB, bc core.BlockChain) shared.Transformer {
	transformer.Repository.SetDB(db)
	transformer.Fetcher.SetBC(bc)
	return transformer
}

func (transformer Transformer) Execute() error {
	transformerName := transformer.Config.TransformerName
	config := transformer.Config
	topics := [][]common.Hash{{common.HexToHash(config.Topic)}}
	missingHeaders, err := transformer.Repository.MissingHeaders(config.StartingBlockNumber, config.EndingBlockNumber)
	if err != nil {
		log.Printf("Error fetching missing headers in %v transformer: %v \n", transformerName, err)
		return err
	}
	log.Printf("Fetching %v event logs for %d headers \n", transformerName, len(missingHeaders))
	for _, header := range missingHeaders {
		logs, err := transformer.Fetcher.FetchLogs(config.ContractAddresses, topics, header.BlockNumber)
		if err != nil {
			log.Printf("Error fetching matching logs in %v transformer: %v", transformerName, err)
			return err
		}

		if len(logs) < 1 {
			err = transformer.Repository.MarkHeaderChecked(header.Id)
			if err != nil {
				log.Printf("Error marking header as checked in %v: %v", transformerName, err)
				return err
			}

			continue
		}

		entities, err := transformer.Converter.ToEntities(config.ContractAbi, logs)
		if err != nil {
			log.Printf("Error converting logs to entities in %v: %v", transformerName, err)
			return err
		}

		models, err := transformer.Converter.ToModels(entities)
		if err != nil {
			log.Printf("Error converting entities to models in %v: %v", transformerName, err)
			return err
		}

		err = transformer.Repository.Create(header.Id, models)
		if err != nil {
			log.Printf("Error persisting %v record: %v", transformerName, err)
			return err
		}

	}

	return nil
}