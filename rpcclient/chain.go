// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpcclient

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrjson/v3"
	"github.com/decred/dcrd/dcrutil/v2"
	"github.com/decred/dcrd/gcs/v2"
	"github.com/decred/dcrd/gcs/v2/blockcf"
	chainjson "github.com/decred/dcrd/rpc/jsonrpc/types/v2"
	"github.com/decred/dcrd/wire"
)

// FutureGetBestBlockHashResult is a future promise to deliver the result of a
// GetBestBlockAsync RPC invocation (or an applicable error).
type FutureGetBestBlockHashResult chan *response

// Receive waits for the response promised by the future and returns the hash of
// the best block in the longest block chain.
func (r FutureGetBestBlockHashResult) Receive() (*chainhash.Hash, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var txHashStr string
	err = json.Unmarshal(res, &txHashStr)
	if err != nil {
		return nil, err
	}
	return chainhash.NewHashFromStr(txHashStr)
}

// GetBestBlockHashAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetBestBlockHash for the blocking version and more details.
func (c *Client) GetBestBlockHashAsync() FutureGetBestBlockHashResult {
	cmd := chainjson.NewGetBestBlockHashCmd()
	return c.sendCmd(cmd)
}

// GetBestBlockHash returns the hash of the best block in the longest block
// chain.
func (c *Client) GetBestBlockHash() (*chainhash.Hash, error) {
	return c.GetBestBlockHashAsync().Receive()
}

// FutureGetBlockResult is a future promise to deliver the result of a
// GetBlockAsync RPC invocation (or an applicable error).
type FutureGetBlockResult chan *response

// Receive waits for the response promised by the future and returns the raw
// block requested from the server given its hash.
func (r FutureGetBlockResult) Receive() (*wire.MsgBlock, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var blockHex string
	err = json.Unmarshal(res, &blockHex)
	if err != nil {
		return nil, err
	}

	// Decode the serialized block hex to raw bytes.
	serializedBlock, err := hex.DecodeString(blockHex)
	if err != nil {
		return nil, err
	}

	// Deserialize the block and return it.
	var msgBlock wire.MsgBlock
	err = msgBlock.Deserialize(bytes.NewReader(serializedBlock))
	if err != nil {
		return nil, err
	}
	return &msgBlock, nil
}

// GetBlockAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetBlock for the blocking version and more details.
func (c *Client) GetBlockAsync(blockHash *chainhash.Hash) FutureGetBlockResult {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}

	cmd := chainjson.NewGetBlockCmd(hash, dcrjson.Bool(false), nil)
	return c.sendCmd(cmd)
}

// GetBlock returns a raw block from the server given its hash.
//
// See GetBlockVerbose to retrieve a data structure with information about the
// block instead.
func (c *Client) GetBlock(blockHash *chainhash.Hash) (*wire.MsgBlock, error) {
	return c.GetBlockAsync(blockHash).Receive()
}

// FutureGetBlockVerboseResult is a future promise to deliver the result of a
// GetBlockVerboseAsync RPC invocation (or an applicable error).
type FutureGetBlockVerboseResult chan *response

// Receive waits for the response promised by the future and returns the data
// structure from the server with information about the requested block.
func (r FutureGetBlockVerboseResult) Receive() (*chainjson.GetBlockVerboseResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal the raw result into a BlockResult.
	var blockResult chainjson.GetBlockVerboseResult
	err = json.Unmarshal(res, &blockResult)
	if err != nil {
		return nil, err
	}
	return &blockResult, nil
}

// GetBlockVerboseAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetBlockVerbose for the blocking version and more details.
func (c *Client) GetBlockVerboseAsync(blockHash *chainhash.Hash, verboseTx bool) FutureGetBlockVerboseResult {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}

	cmd := chainjson.NewGetBlockCmd(hash, dcrjson.Bool(true), &verboseTx)
	return c.sendCmd(cmd)
}

// GetBlockVerbose returns a data structure from the server with information
// about a block given its hash.
//
// See GetBlock to retrieve a raw block instead.
func (c *Client) GetBlockVerbose(blockHash *chainhash.Hash, verboseTx bool) (*chainjson.GetBlockVerboseResult, error) {
	return c.GetBlockVerboseAsync(blockHash, verboseTx).Receive()
}

// FutureGetBlockCountResult is a future promise to deliver the result of a
// GetBlockCountAsync RPC invocation (or an applicable error).
type FutureGetBlockCountResult chan *response

// Receive waits for the response promised by the future and returns the number
// of blocks in the longest block chain.
func (r FutureGetBlockCountResult) Receive() (int64, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return 0, err
	}

	// Unmarshal the result as an int64.
	var count int64
	err = json.Unmarshal(res, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetBlockCountAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetBlockCount for the blocking version and more details.
func (c *Client) GetBlockCountAsync() FutureGetBlockCountResult {
	cmd := chainjson.NewGetBlockCountCmd()
	return c.sendCmd(cmd)
}

// GetBlockCount returns the number of blocks in the longest block chain.
func (c *Client) GetBlockCount() (int64, error) {
	return c.GetBlockCountAsync().Receive()
}

// FutureGetDifficultyResult is a future promise to deliver the result of a
// GetDifficultyAsync RPC invocation (or an applicable error).
type FutureGetDifficultyResult chan *response

// Receive waits for the response promised by the future and returns the
// proof-of-work difficulty as a multiple of the minimum difficulty.
func (r FutureGetDifficultyResult) Receive() (float64, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return 0, err
	}

	// Unmarshal the result as a float64.
	var difficulty float64
	err = json.Unmarshal(res, &difficulty)
	if err != nil {
		return 0, err
	}
	return difficulty, nil
}

// GetDifficultyAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetDifficulty for the blocking version and more details.
func (c *Client) GetDifficultyAsync() FutureGetDifficultyResult {
	cmd := chainjson.NewGetDifficultyCmd()
	return c.sendCmd(cmd)
}

// GetDifficulty returns the proof-of-work difficulty as a multiple of the
// minimum difficulty.
func (c *Client) GetDifficulty() (float64, error) {
	return c.GetDifficultyAsync().Receive()
}

// FutureGetBlockChainInfoResult is a future promise to deliver the result of a
// GetBlockChainInfoAsync RPC invocation (or an applicable error).
type FutureGetBlockChainInfoResult chan *response

// Receive waits for the response promised by the future and returns the info
// provided by the server.
func (r FutureGetBlockChainInfoResult) Receive() (*chainjson.GetBlockChainInfoResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a getblockchaininfo result object.
	var blockchainInfoRes chainjson.GetBlockChainInfoResult
	err = json.Unmarshal(res, &blockchainInfoRes)
	if err != nil {
		return nil, err
	}

	return &blockchainInfoRes, nil
}

// GetBlockChainInfoAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetBlockChainInfo for the blocking version and more details.
func (c *Client) GetBlockChainInfoAsync() FutureGetBlockChainInfoResult {
	cmd := chainjson.NewGetBlockChainInfoCmd()
	return c.sendCmd(cmd)
}

// GetBlockChainInfo returns information about the current state of the block
// chain.
func (c *Client) GetBlockChainInfo() (*chainjson.GetBlockChainInfoResult, error) {
	return c.GetBlockChainInfoAsync().Receive()
}

// FutureGetBlockHashResult is a future promise to deliver the result of a
// GetBlockHashAsync RPC invocation (or an applicable error).
type FutureGetBlockHashResult chan *response

// Receive waits for the response promised by the future and returns the hash of
// the block in the best block chain at the given height.
func (r FutureGetBlockHashResult) Receive() (*chainhash.Hash, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal the result as a string-encoded sha.
	var txHashStr string
	err = json.Unmarshal(res, &txHashStr)
	if err != nil {
		return nil, err
	}
	return chainhash.NewHashFromStr(txHashStr)
}

// GetBlockHashAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetBlockHash for the blocking version and more details.
func (c *Client) GetBlockHashAsync(blockHeight int64) FutureGetBlockHashResult {
	cmd := chainjson.NewGetBlockHashCmd(blockHeight)
	return c.sendCmd(cmd)
}

// GetBlockHash returns the hash of the block in the best block chain at the
// given height.
func (c *Client) GetBlockHash(blockHeight int64) (*chainhash.Hash, error) {
	return c.GetBlockHashAsync(blockHeight).Receive()
}

// FutureGetBlockHeaderResult is a future promise to deliver the result of a
// GetBlockHeaderAsync RPC invocation (or an applicable error).
type FutureGetBlockHeaderResult chan *response

// Receive waits for the response promised by the future and returns the
// blockheader requested from the server given its hash.
func (r FutureGetBlockHeaderResult) Receive() (*wire.BlockHeader, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal the result as a string.
	var bhHex string
	err = json.Unmarshal(res, &bhHex)
	if err != nil {
		return nil, err
	}

	serializedBH, err := hex.DecodeString(bhHex)
	if err != nil {
		return nil, err
	}

	// Deserialize the blockheader and return it.
	var bh wire.BlockHeader
	err = bh.Deserialize(bytes.NewReader(serializedBH))
	if err != nil {
		return nil, err
	}

	return &bh, nil
}

// GetBlockHeaderAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetBlockHeader for the blocking version and more details.
func (c *Client) GetBlockHeaderAsync(hash *chainhash.Hash) FutureGetBlockHeaderResult {
	cmd := chainjson.NewGetBlockHeaderCmd(hash.String(), dcrjson.Bool(false))
	return c.sendCmd(cmd)
}

// GetBlockHeader returns the hash of the block in the best block chain at the
// given height.
func (c *Client) GetBlockHeader(hash *chainhash.Hash) (*wire.BlockHeader, error) {
	return c.GetBlockHeaderAsync(hash).Receive()
}

// FutureGetBlockHeaderVerboseResult is a future promise to deliver the result of a
// GetBlockHeaderAsync RPC invocation (or an applicable error).
type FutureGetBlockHeaderVerboseResult chan *response

// Receive waits for the response promised by the future and returns a data
// structure of the block header requested from the server given its hash.
func (r FutureGetBlockHeaderVerboseResult) Receive() (*chainjson.GetBlockHeaderVerboseResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	var bh chainjson.GetBlockHeaderVerboseResult
	err = json.Unmarshal(res, &bh)
	if err != nil {
		return nil, err
	}
	return &bh, nil
}

// GetBlockHeaderVerboseAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetBlockHeaderVerbose for the blocking version and more details.
func (c *Client) GetBlockHeaderVerboseAsync(hash *chainhash.Hash) FutureGetBlockHeaderVerboseResult {
	cmd := chainjson.NewGetBlockHeaderCmd(hash.String(), dcrjson.Bool(true))
	return c.sendCmd(cmd)
}

// GetBlockHeaderVerbose returns a data structure of the block header from the
// server given its hash.
//
// See GetBlockHeader to retrieve a raw block header instead.
func (c *Client) GetBlockHeaderVerbose(hash *chainhash.Hash) (*chainjson.GetBlockHeaderVerboseResult, error) {
	return c.GetBlockHeaderVerboseAsync(hash).Receive()
}

// FutureGetBlockSubsidyResult is a future promise to deliver the result of a
// GetBlockSubsidyAsync RPC invocation (or an applicable error).
type FutureGetBlockSubsidyResult chan *response

// Receive waits for the response promised by the future and returns a data
// structure of the block subsidy requested from the server given its height
// and number of voters.
func (r FutureGetBlockSubsidyResult) Receive() (*chainjson.GetBlockSubsidyResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	var bs chainjson.GetBlockSubsidyResult
	err = json.Unmarshal(res, &bs)
	if err != nil {
		return nil, err
	}
	return &bs, nil
}

// GetBlockSubsidyAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetBlockSubsidy for the blocking version and more details.
func (c *Client) GetBlockSubsidyAsync(height int64, voters uint16) FutureGetBlockSubsidyResult {
	cmd := chainjson.NewGetBlockSubsidyCmd(height, voters)
	return c.sendCmd(cmd)
}

// GetBlockSubsidy returns a data structure of the block subsidy
// from the server given its height and number of voters.
func (c *Client) GetBlockSubsidy(height int64, voters uint16) (*chainjson.GetBlockSubsidyResult, error) {
	return c.GetBlockSubsidyAsync(height, voters).Receive()
}

// FutureGetCoinSupplyResult is a future promise to deliver the result of a
// GetCoinSupplyAsync RPC invocation (or an applicable error).
type FutureGetCoinSupplyResult chan *response

// Receive waits for the response promised by the future and returns the
// current coin supply
func (r FutureGetCoinSupplyResult) Receive() (dcrutil.Amount, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return 0, err
	}

	// Unmarshal the result
	var cs int64
	err = json.Unmarshal(res, &cs)
	if err != nil {
		return 0, err
	}
	return dcrutil.Amount(cs), nil
}

// GetCoinSupplyAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetCoinSupply for the blocking version and more details.
func (c *Client) GetCoinSupplyAsync() FutureGetCoinSupplyResult {
	cmd := chainjson.NewGetCoinSupplyCmd()
	return c.sendCmd(cmd)
}

// GetCoinSupply returns the current coin supply
func (c *Client) GetCoinSupply() (dcrutil.Amount, error) {
	return c.GetCoinSupplyAsync().Receive()
}

// FutureGetRawMempoolResult is a future promise to deliver the result of a
// GetRawMempoolAsync RPC invocation (or an applicable error).
type FutureGetRawMempoolResult chan *response

// Receive waits for the response promised by the future and returns the hashes
// of all transactions in the memory pool.
func (r FutureGetRawMempoolResult) Receive() ([]*chainhash.Hash, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal the result as an array of strings.
	var txHashStrs []string
	err = json.Unmarshal(res, &txHashStrs)
	if err != nil {
		return nil, err
	}

	// Create a slice of ShaHash arrays from the string slice.
	txHashes := make([]*chainhash.Hash, 0, len(txHashStrs))
	for _, hashStr := range txHashStrs {
		txHash, err := chainhash.NewHashFromStr(hashStr)
		if err != nil {
			return nil, err
		}
		txHashes = append(txHashes, txHash)
	}

	return txHashes, nil
}

// GetRawMempoolAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetRawMempool for the blocking version and more details.
func (c *Client) GetRawMempoolAsync(txType chainjson.GetRawMempoolTxTypeCmd) FutureGetRawMempoolResult {
	cmd := chainjson.NewGetRawMempoolCmd(dcrjson.Bool(false),
		dcrjson.String(string(txType)))
	return c.sendCmd(cmd)
}

// GetRawMempool returns the hashes of all transactions in the memory pool for
// the given txType.
//
// See GetRawMempoolVerbose to retrieve data structures with information about
// the transactions instead.
func (c *Client) GetRawMempool(txType chainjson.GetRawMempoolTxTypeCmd) ([]*chainhash.Hash, error) {
	return c.GetRawMempoolAsync(txType).Receive()
}

// FutureGetRawMempoolVerboseResult is a future promise to deliver the result of
// a GetRawMempoolVerboseAsync RPC invocation (or an applicable error).
type FutureGetRawMempoolVerboseResult chan *response

// Receive waits for the response promised by the future and returns a map of
// transaction hashes to an associated data structure with information about the
// transaction for all transactions in the memory pool.
func (r FutureGetRawMempoolVerboseResult) Receive() (map[string]chainjson.GetRawMempoolVerboseResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal the result as a map of strings (tx shas) to their detailed
	// results.
	var mempoolItems map[string]chainjson.GetRawMempoolVerboseResult
	err = json.Unmarshal(res, &mempoolItems)
	if err != nil {
		return nil, err
	}
	return mempoolItems, nil
}

// GetRawMempoolVerboseAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
//
// See GetRawMempoolVerbose for the blocking version and more details.
func (c *Client) GetRawMempoolVerboseAsync(txType chainjson.GetRawMempoolTxTypeCmd) FutureGetRawMempoolVerboseResult {
	cmd := chainjson.NewGetRawMempoolCmd(dcrjson.Bool(true),
		dcrjson.String(string(txType)))
	return c.sendCmd(cmd)
}

// GetRawMempoolVerbose returns a map of transaction hashes to an associated
// data structure with information about the transaction for all transactions in
// the memory pool.
//
// See GetRawMempool to retrieve only the transaction hashes instead.
func (c *Client) GetRawMempoolVerbose(txType chainjson.GetRawMempoolTxTypeCmd) (map[string]chainjson.GetRawMempoolVerboseResult, error) {
	return c.GetRawMempoolVerboseAsync(txType).Receive()
}

// FutureVerifyChainResult is a future promise to deliver the result of a
// VerifyChainAsync, VerifyChainLevelAsyncRPC, or VerifyChainBlocksAsync
// invocation (or an applicable error).
type FutureVerifyChainResult chan *response

// Receive waits for the response promised by the future and returns whether
// or not the chain verified based on the check level and number of blocks
// to verify specified in the original call.
func (r FutureVerifyChainResult) Receive() (bool, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return false, err
	}

	// Unmarshal the result as a boolean.
	var verified bool
	err = json.Unmarshal(res, &verified)
	if err != nil {
		return false, err
	}
	return verified, nil
}

// FutureGetChainTipsResult is a future promise to deliver the result of a
// GetChainTipsAsync RPC invocation (or an applicable error).
type FutureGetChainTipsResult chan *response

// Receive waits for the response promised by the future and returns slice of
// all known tips in the block tree.
func (r FutureGetChainTipsResult) Receive() ([]chainjson.GetChainTipsResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal the result.
	var chainTips []chainjson.GetChainTipsResult
	err = json.Unmarshal(res, &chainTips)
	if err != nil {
		return nil, err
	}
	return chainTips, nil
}

// GetChainTipsAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetChainTips for the blocking version and more details.
func (c *Client) GetChainTipsAsync() FutureGetChainTipsResult {
	cmd := chainjson.NewGetChainTipsCmd()
	return c.sendCmd(cmd)
}

// GetChainTips returns all known tips in the block tree.
func (c *Client) GetChainTips() ([]chainjson.GetChainTipsResult, error) {
	return c.GetChainTipsAsync().Receive()
}

// VerifyChainAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See VerifyChain for the blocking version and more details.
func (c *Client) VerifyChainAsync() FutureVerifyChainResult {
	cmd := chainjson.NewVerifyChainCmd(nil, nil)
	return c.sendCmd(cmd)
}

// VerifyChain requests the server to verify the block chain database using
// the default check level and number of blocks to verify.
//
// See VerifyChainLevel and VerifyChainBlocks to override the defaults.
func (c *Client) VerifyChain() (bool, error) {
	return c.VerifyChainAsync().Receive()
}

// VerifyChainLevelAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See VerifyChainLevel for the blocking version and more details.
func (c *Client) VerifyChainLevelAsync(checkLevel int64) FutureVerifyChainResult {
	cmd := chainjson.NewVerifyChainCmd(&checkLevel, nil)
	return c.sendCmd(cmd)
}

// VerifyChainLevel requests the server to verify the block chain database using
// the passed check level and default number of blocks to verify.
//
// The check level controls how thorough the verification is with higher numbers
// increasing the amount of checks done as consequently how long the
// verification takes.
//
// See VerifyChain to use the default check level and VerifyChainBlocks to
// override the number of blocks to verify.
func (c *Client) VerifyChainLevel(checkLevel int64) (bool, error) {
	return c.VerifyChainLevelAsync(checkLevel).Receive()
}

// VerifyChainBlocksAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See VerifyChainBlocks for the blocking version and more details.
func (c *Client) VerifyChainBlocksAsync(checkLevel, numBlocks int64) FutureVerifyChainResult {
	cmd := chainjson.NewVerifyChainCmd(&checkLevel, &numBlocks)
	return c.sendCmd(cmd)
}

// VerifyChainBlocks requests the server to verify the block chain database
// using the passed check level and number of blocks to verify.
//
// The check level controls how thorough the verification is with higher numbers
// increasing the amount of checks done as consequently how long the
// verification takes.
//
// The number of blocks refers to the number of blocks from the end of the
// current longest chain.
//
// See VerifyChain and VerifyChainLevel to use defaults.
func (c *Client) VerifyChainBlocks(checkLevel, numBlocks int64) (bool, error) {
	return c.VerifyChainBlocksAsync(checkLevel, numBlocks).Receive()
}

// FutureGetTxOutResult is a future promise to deliver the result of a
// GetTxOutAsync RPC invocation (or an applicable error).
type FutureGetTxOutResult chan *response

// Receive waits for the response promised by the future and returns a
// transaction given its hash.
func (r FutureGetTxOutResult) Receive() (*chainjson.GetTxOutResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// take care of the special case where the output has been spent already
	// it should return the string "null"
	if string(res) == "null" {
		return nil, nil
	}

	// Unmarshal result as a gettxout result object.
	var txOutInfo *chainjson.GetTxOutResult
	err = json.Unmarshal(res, &txOutInfo)
	if err != nil {
		return nil, err
	}

	return txOutInfo, nil
}

// GetTxOutAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetTxOut for the blocking version and more details.
func (c *Client) GetTxOutAsync(txHash *chainhash.Hash, index uint32, mempool bool) FutureGetTxOutResult {
	hash := ""
	if txHash != nil {
		hash = txHash.String()
	}

	cmd := chainjson.NewGetTxOutCmd(hash, index, &mempool)
	return c.sendCmd(cmd)
}

// GetTxOut returns the transaction output info if it's unspent and
// nil, otherwise.
func (c *Client) GetTxOut(txHash *chainhash.Hash, index uint32, mempool bool) (*chainjson.GetTxOutResult, error) {
	return c.GetTxOutAsync(txHash, index, mempool).Receive()
}

// FutureRescanResult is a future promise to deliver the result of a
// RescanAsynnc RPC invocation (or an applicable error).
type FutureRescanResult chan *response

// Receive waits for the response promised by the future and returns the
// discovered rescan data.
func (r FutureRescanResult) Receive() (*chainjson.RescanResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	var rescanResult *chainjson.RescanResult
	err = json.Unmarshal(res, &rescanResult)
	if err != nil {
		return nil, err
	}

	return rescanResult, nil
}

// RescanAsync returns an instance of a type that can be used to get the result
// of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See Rescan for the blocking version and more details.
func (c *Client) RescanAsync(blockHashes []chainhash.Hash) FutureRescanResult {
	hashes := make([]string, len(blockHashes))
	for i := range blockHashes {
		hashes[i] = blockHashes[i].String()
	}

	cmd := chainjson.NewRescanCmd(hashes)
	return c.sendCmd(cmd)
}

// Rescan rescans the blocks identified by blockHashes, in order, using the
// client's loaded transaction filter.  The blocks do not need to be on the main
// chain, but they do need to be adjacent to each other.
func (c *Client) Rescan(blockHashes []chainhash.Hash) (*chainjson.RescanResult, error) {
	return c.RescanAsync(blockHashes).Receive()
}

// FutureGetCFilterResult is a future promise to deliver the result of a
// GetCFilterAsync RPC invocation (or an applicable error).
type FutureGetCFilterResult chan *response

// Receive waits for the response promised by the future and returns the
// discovered rescan data.
func (r FutureGetCFilterResult) Receive() (*gcs.FilterV1, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	var filterHex string
	err = json.Unmarshal(res, &filterHex)
	if err != nil {
		return nil, err
	}
	filterBytes, err := hex.DecodeString(filterHex)
	if err != nil {
		return nil, err
	}

	return gcs.FromBytesV1(blockcf.P, filterBytes)
}

// GetCFilterAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See GetCFilter for the blocking version and more details.
func (c *Client) GetCFilterAsync(blockHash *chainhash.Hash, filterType wire.FilterType) FutureGetCFilterResult {
	var ft string
	switch filterType {
	case wire.GCSFilterRegular:
		ft = "regular"
	case wire.GCSFilterExtended:
		ft = "extended"
	default:
		return futureError(errors.New("unknown filter type"))
	}

	cmd := chainjson.NewGetCFilterCmd(blockHash.String(), ft)
	return c.sendCmd(cmd)
}

// GetCFilter returns the committed filter of type filterType for a block.
func (c *Client) GetCFilter(blockHash *chainhash.Hash, filterType wire.FilterType) (*gcs.FilterV1, error) {
	return c.GetCFilterAsync(blockHash, filterType).Receive()
}

// FutureGetCFilterHeaderResult is a future promise to deliver the result of a
// GetCFilterHeaderAsync RPC invocation (or an applicable error).
type FutureGetCFilterHeaderResult chan *response

// Receive waits for the response promised by the future and returns the
// discovered rescan data.
func (r FutureGetCFilterHeaderResult) Receive() (*chainhash.Hash, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	var filterHeaderHex string
	err = json.Unmarshal(res, &filterHeaderHex)
	if err != nil {
		return nil, err
	}

	return chainhash.NewHashFromStr(filterHeaderHex)
}

// GetCFilterHeaderAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See GetCFilterHeader for the blocking version and more details.
func (c *Client) GetCFilterHeaderAsync(blockHash *chainhash.Hash, filterType wire.FilterType) FutureGetCFilterHeaderResult {
	var ft string
	switch filterType {
	case wire.GCSFilterRegular:
		ft = "regular"
	case wire.GCSFilterExtended:
		ft = "extended"
	default:
		return futureError(errors.New("unknown filter type"))
	}

	cmd := chainjson.NewGetCFilterHeaderCmd(blockHash.String(), ft)
	return c.sendCmd(cmd)
}

// GetCFilterHeader returns the committed filter header hash of type filterType
// for a block.
func (c *Client) GetCFilterHeader(blockHash *chainhash.Hash, filterType wire.FilterType) (*chainhash.Hash, error) {
	return c.GetCFilterHeaderAsync(blockHash, filterType).Receive()
}

// FutureEstimateSmartFeeResult is a future promise to deliver the result of a
// EstimateSmartFee RPC invocation (or an applicable error).
type FutureEstimateSmartFeeResult chan *response

// Receive waits for the response promised by the future and returns a fee
// estimation for the given target confirmation window and mode.
func (r FutureEstimateSmartFeeResult) Receive() (float64, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return 0, err
	}

	// Unmarshal the result as a float64.
	var dcrPerKB float64
	err = json.Unmarshal(res, &dcrPerKB)
	if err != nil {
		return 0, err
	}
	return dcrPerKB, nil
}

// EstimateSmartFeeAsync returns an instance of a type that can be used to get
// the result of the RPC at some future time by invoking the Receive function on
// the returned instance.
//
// See EstimateSmartFee for the blocking version and more details.
func (c *Client) EstimateSmartFeeAsync(confirmations int64, mode chainjson.EstimateSmartFeeMode) FutureEstimateSmartFeeResult {
	cmd := chainjson.NewEstimateSmartFeeCmd(confirmations, &mode)
	return c.sendCmd(cmd)
}

// EstimateSmartFee returns an estimation of a transaction fee rate (in dcr/KB)
// that new transactions should pay if they desire to be mined in up to
// 'confirmations' blocks.
//
// The mode parameter (roughly) selects the different thresholds for accepting
// an estimation as reasonable, allowing users to select different trade-offs
// between probability of the transaction being mined in the given target
// confirmation range and minimization of fees paid.
//
// As of 2019-01, only the default conservative mode is supported by dcrd.
func (c *Client) EstimateSmartFee(confirmations int64, mode chainjson.EstimateSmartFeeMode) (float64, error) {
	return c.EstimateSmartFeeAsync(confirmations, mode).Receive()
}
