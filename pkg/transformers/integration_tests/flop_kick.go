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

package integration_tests

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("FlopKick Transformer", func() {
	It("fetches and transforms a FlopKick event from Kovan chain", func() {
		blockNumber := int64(8672119)
		config := flop_kick.Config
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockchain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockchain.Node())
		test_config.CleanTestDB(db)

		err = persistHeader(db, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		initializer := flop_kick.FlopKickTransformerInitializer{Config: config}
		transformer := initializer.NewFlopKickTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []flop_kick.Model
		err = db.Select(&dbResult, `SELECT bid, bid_id, "end", gal, lot FROM maker.flop_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("0"))
		Expect(dbResult[0].BidId).To(Equal("1"))
		Expect(dbResult[0].End.Equal(time.Unix(1536726768, 0))).To(BeTrue())
		Expect(dbResult[0].Gal).To(Equal("0x9B870D55BaAEa9119dBFa71A92c5E26E79C4726d"))
		// this very large number appears to be derived from the data including: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		Expect(dbResult[0].Lot).To(Equal("115792089237316195423570985008687907853269984665640564039457584007913129639935"))
	})

	It("fetches and transforms another FlopKick event from Kovan chain", func() {
		blockNumber := int64(8955611)
		config := flop_kick.Config
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockchain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockchain.Node())
		test_config.CleanTestDB(db)

		err = persistHeader(db, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		initializer := flop_kick.FlopKickTransformerInitializer{Config: config}
		transformer := initializer.NewFlopKickTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []flop_kick.Model
		err = db.Select(&dbResult, `SELECT bid, bid_id, "end", gal, lot FROM maker.flop_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("10000000000000000000000"))
		Expect(dbResult[0].BidId).To(Equal("2"))
		Expect(dbResult[0].End.Equal(time.Unix(1538810564, 0))).To(BeTrue())
		Expect(dbResult[0].Gal).To(Equal("0x3728e9777B2a0a611ee0F89e00E01044ce4736d1"))
		Expect(dbResult[0].Lot).To(Equal("115792089237316195423570985008687907853269984665640564039457584007913129639935"))
	})
})