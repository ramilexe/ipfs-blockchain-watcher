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

package tend

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type TendRepository struct {
	db *postgres.DB
}

func (repository TendRepository) Create(headerId int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}

	for _, model := range models {
		tend, ok := model.(TendModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, TendModel{})
		}

		_, err = tx.Exec(
			`INSERT into maker.tend (header_id, bid_id, lot, bid, guy, tic, log_idx, tx_idx, raw_log)
        	VALUES($1, $2, $3::NUMERIC, $4::NUMERIC, $5, $6::NUMERIC, $7, $8, $9)`,
			headerId, tend.BidId, tend.Lot, tend.Bid, tend.Guy, tend.Tic, tend.LogIndex, tend.TransactionIndex, tend.Raw,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, tend_checked)
			VALUES ($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET tend_checked = $2`, headerId, true)

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository TendRepository) MarkHeaderChecked(headerId int64) error {
	_, err := repository.db.Exec(`INSERT INTO public.checked_headers (header_id, tend_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET tend_checked = $2`, headerId, true)
	return err
}

func (repository TendRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
				LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR tend_checked IS FALSE)
               	AND headers.block_number >= $1
               	AND headers.block_number <= $2
               	AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		repository.db.Node.ID,
	)

	return result, err
}

func (repository *TendRepository) SetDB(db *postgres.DB) {
	repository.db = db
}