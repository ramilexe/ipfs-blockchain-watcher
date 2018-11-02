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

package flip_kick

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type FlipKickRepository struct {
	db *postgres.DB
}

func (fkr FlipKickRepository) Create(headerId int64, models []interface{}) error {
	tx, err := fkr.db.Begin()
	if err != nil {
		return err
	}
	for _, model := range models {
		flipKickModel, ok := model.(FlipKickModel)
		if !ok {
			return fmt.Errorf("model of type %T, not %T", model, FlipKickModel{})
		}
		_, err := tx.Exec(
			`INSERT into maker.flip_kick (header_id, bid_id, lot, bid, gal, "end", urn, tab, tx_idx, log_idx, raw_log)
        VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5, $6, $7, $8::NUMERIC, $9, $10, $11)`,
			headerId, flipKickModel.BidId, flipKickModel.Lot, flipKickModel.Bid, flipKickModel.Gal, flipKickModel.End, flipKickModel.Urn, flipKickModel.Tab, flipKickModel.TransactionIndex, flipKickModel.LogIndex, flipKickModel.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, flip_kick_checked)
			VALUES ($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET flip_kick_checked = $2`, headerId, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (fkr FlipKickRepository) MarkHeaderChecked(headerId int64) error {
	_, err := fkr.db.Exec(`INSERT INTO public.checked_headers (header_id, flip_kick_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET flip_kick_checked = $2`, headerId, true)
	return err
}

func (fkr FlipKickRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := fkr.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR flip_kick_checked IS FALSE)
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		fkr.db.Node.ID,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return result, err
	}

	return result, nil
}

func (fkr *FlipKickRepository) SetDB(db *postgres.DB) {
	fkr.db = db
}
