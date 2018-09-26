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

package drip_drip

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_drip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockDripDripConverter struct {
	Err       error
	PassedLog types.Log
}

func (converter *MockDripDripConverter) ToModel(ethLog types.Log) (drip_drip.DripDripModel, error) {
	converter.PassedLog = ethLog
	return test_data.DripDripModel, converter.Err
}

func (converter *MockDripDripConverter) SetConverterError(e error) {
	converter.Err = e
}