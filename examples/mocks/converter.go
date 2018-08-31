// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mocks

import (
	et1 "github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	et2 "github.com/vulcanize/vulcanizedb/examples/generic/event_triggered"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockERC20Converter struct {
	WatchedEvents      []*core.WatchedEvent
	TransfersToConvert []et1.TransferEntity
	ApprovalsToConvert []et1.ApprovalEntity
	BurnsToConvert     []et2.BurnEntity
	MintsToConvert     []et2.MintEntity
	block              int64
}

func (mlkc *MockERC20Converter) ToTransferModel(entity *et1.TransferEntity) *et1.TransferModel {
	mlkc.TransfersToConvert = append(mlkc.TransfersToConvert, *entity)
	return &et1.TransferModel{}
}

func (mlkc *MockERC20Converter) ToTransferEntity(watchedEvent core.WatchedEvent) (*et1.TransferEntity, error) {
	mlkc.WatchedEvents = append(mlkc.WatchedEvents, &watchedEvent)
	e := &et1.TransferEntity{Block: watchedEvent.BlockNumber}
	mlkc.block++
	return e, nil
}

func (mlkc *MockERC20Converter) ToApprovalModel(entity *et1.ApprovalEntity) *et1.ApprovalModel {
	mlkc.ApprovalsToConvert = append(mlkc.ApprovalsToConvert, *entity)
	return &et1.ApprovalModel{}
}

func (mlkc *MockERC20Converter) ToApprovalEntity(watchedEvent core.WatchedEvent) (*et1.ApprovalEntity, error) {
	mlkc.WatchedEvents = append(mlkc.WatchedEvents, &watchedEvent)
	e := &et1.ApprovalEntity{Block: watchedEvent.BlockNumber}
	mlkc.block++
	return e, nil
}

func (mlkc *MockERC20Converter) ToBurnEntity(watchedEvent core.WatchedEvent) (*et2.BurnEntity, error) {
	mlkc.WatchedEvents = append(mlkc.WatchedEvents, &watchedEvent)
	e := &et2.BurnEntity{Block: watchedEvent.BlockNumber}
	mlkc.block++
	return e, nil
}

func (mlkc *MockERC20Converter) ToBurnModel(entity *et2.BurnEntity) *et2.BurnModel {
	mlkc.BurnsToConvert = append(mlkc.BurnsToConvert, *entity)
	return &et2.BurnModel{}
}

func (mlkc *MockERC20Converter) ToMintEntity(watchedEvent core.WatchedEvent) (*et2.MintEntity, error) {
	mlkc.WatchedEvents = append(mlkc.WatchedEvents, &watchedEvent)
	e := &et2.MintEntity{Block: watchedEvent.BlockNumber}
	mlkc.block++
	return e, nil
}

func (mlkc *MockERC20Converter) ToMintModel(entity *et2.MintEntity) *et2.MintModel {
	mlkc.MintsToConvert = append(mlkc.MintsToConvert, *entity)
	return &et2.MintModel{}
}