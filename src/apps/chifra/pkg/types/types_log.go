// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package types

// EXISTING_CODE
import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

// EXISTING_CODE

type RawLog struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	// EXISTING_CODE
	// EXISTING_CODE
}

type SimpleLog struct {
	Address          base.Address    `json:"address"`
	ArticulatedLog   *SimpleFunction `json:"articulatedLog,omitempty"`
	BlockHash        base.Hash       `json:"blockHash"`
	BlockNumber      base.Blknum     `json:"blockNumber"`
	CompressedLog    string          `json:"compressedLog,omitempty"`
	Data             string          `json:"data,omitempty"`
	LogIndex         uint64          `json:"logIndex"`
	Timestamp        base.Timestamp  `json:"timestamp,omitempty"`
	Topics           []base.Hash     `json:"topics,omitempty"`
	TransactionHash  base.Hash       `json:"transactionHash"`
	TransactionIndex uint64          `json:"transactionIndex"`
	raw              *RawLog         `json:"-"`
	// EXISTING_CODE
	// EXISTING_CODE
}

func (s *SimpleLog) Raw() *RawLog {
	return s.raw
}

func (s *SimpleLog) SetRaw(raw *RawLog) {
	s.raw = raw
}

func (s *SimpleLog) Model(chain, format string, verbose bool, extraOptions map[string]any) Model {
	var model = map[string]interface{}{}
	var order = []string{}

	// EXISTING_CODE
	model = map[string]interface{}{
		"address":          s.Address,
		"blockHash":        s.BlockHash,
		"blockNumber":      s.BlockNumber,
		"logIndex":         s.LogIndex,
		"timestamp":        s.Timestamp,
		"date":             s.Date(),
		"transactionIndex": s.TransactionIndex,
		"transactionHash":  s.TransactionHash,
	}

	order = []string{
		"blockNumber",
		"transactionIndex",
		"logIndex",
		"blockHash",
		"transactionHash",
		"timestamp",
		"date",
		"address",
		"topic0",
		"topic1",
		"topic2",
		"topic3",
		"data",
	}

	isArticulated := extraOptions["articulate"] == true && s.ArticulatedLog != nil
	var articulatedLog = make(map[string]any)
	if isArticulated {
		articulatedLog["name"] = s.ArticulatedLog.Name
		inputModels := parametersToMap(s.ArticulatedLog.Inputs)
		if inputModels != nil {
			articulatedLog["inputs"] = inputModels
		}
	}

	if format == "json" {
		if len(s.Data) > 2 {
			model["data"] = s.Data
		}
		if isArticulated {
			model["articulatedLog"] = articulatedLog
		}

		model["topics"] = s.Topics

	} else {
		if len(s.Data) > 2 {
			model["data"] = s.Data
		} else {
			model["data"] = ""
		}

		if isArticulated {
			model["compressedLog"] = makeCompressed(articulatedLog)
			order = append(order, "compressedLog")
		}

		model["topic0"] = ""
		if len(s.Topics) > 0 {
			model["topic0"] = s.Topics[0]
		}
		model["topic1"] = ""
		if len(s.Topics) > 1 {
			model["topic1"] = s.Topics[1]
		}
		model["topic2"] = ""
		if len(s.Topics) > 2 {
			model["topic2"] = s.Topics[2]
		}
		model["topic3"] = ""
		if len(s.Topics) > 3 {
			model["topic3"] = s.Topics[3]
		}
	}
	// EXISTING_CODE

	return Model{
		Data:  model,
		Order: order,
	}
}

func (s *SimpleLog) Date() string {
	return utils.FormattedDate(s.Timestamp)
}

// --> cacheable by block as group
type SimpleLogGroup struct {
	BlockNumber      base.Blknum
	TransactionIndex base.Txnum
	Logs             []SimpleLog
}

func (s *SimpleLogGroup) CacheName() string {
	return "Log"
}

func (s *SimpleLogGroup) CacheId() string {
	return fmt.Sprintf("%09d", s.BlockNumber)
}

func (s *SimpleLogGroup) CacheLocation() (directory string, extension string) {
	paddedId := s.CacheId()
	parts := make([]string, 3)
	parts[0] = paddedId[:2]
	parts[1] = paddedId[2:4]
	parts[2] = paddedId[4:6]

	subFolder := strings.ToLower(s.CacheName()) + "s"
	directory = filepath.Join(subFolder, filepath.Join(parts...))
	extension = "bin"

	return
}

func (s *SimpleLogGroup) MarshalCache(writer io.Writer) (err error) {
	return cache.WriteValue(writer, s.Logs)
}

func (s *SimpleLogGroup) UnmarshalCache(version uint64, reader io.Reader) (err error) {
	return cache.ReadValue(reader, &s.Logs, version)
}

func (s *SimpleLog) MarshalCache(writer io.Writer) (err error) {
	// Address
	if err = cache.WriteValue(writer, s.Address); err != nil {
		return err
	}

	// ArticulatedLog
	optArticulatedLog := &cache.Optional[SimpleFunction]{
		Value: s.ArticulatedLog,
	}
	if err = cache.WriteValue(writer, optArticulatedLog); err != nil {
		return err
	}

	// BlockHash
	if err = cache.WriteValue(writer, &s.BlockHash); err != nil {
		return err
	}

	// BlockNumber
	if err = cache.WriteValue(writer, s.BlockNumber); err != nil {
		return err
	}

	// CompressedLog
	if err = cache.WriteValue(writer, s.CompressedLog); err != nil {
		return err
	}

	// Data
	if err = cache.WriteValue(writer, s.Data); err != nil {
		return err
	}

	// LogIndex
	if err = cache.WriteValue(writer, s.LogIndex); err != nil {
		return err
	}

	// Timestamp
	if err = cache.WriteValue(writer, s.Timestamp); err != nil {
		return err
	}

	// Topics
	if err = cache.WriteValue(writer, s.Topics); err != nil {
		return err
	}

	// TransactionHash
	if err = cache.WriteValue(writer, &s.TransactionHash); err != nil {
		return err
	}

	// TransactionIndex
	if err = cache.WriteValue(writer, s.TransactionIndex); err != nil {
		return err
	}

	return nil
}

func (s *SimpleLog) UnmarshalCache(version uint64, reader io.Reader) (err error) {
	// Address
	if err = cache.ReadValue(reader, &s.Address, version); err != nil {
		return err
	}

	// ArticulatedLog
	optArticulatedLog := &cache.Optional[SimpleFunction]{
		Value: s.ArticulatedLog,
	}
	if err = cache.ReadValue(reader, optArticulatedLog, version); err != nil {
		return err
	}
	s.ArticulatedLog = optArticulatedLog.Get()

	// BlockHash
	if err = cache.ReadValue(reader, &s.BlockHash, version); err != nil {
		return err
	}

	// BlockNumber
	if err = cache.ReadValue(reader, &s.BlockNumber, version); err != nil {
		return err
	}

	// CompressedLog
	if err = cache.ReadValue(reader, &s.CompressedLog, version); err != nil {
		return err
	}

	// Data
	if err = cache.ReadValue(reader, &s.Data, version); err != nil {
		return err
	}

	// LogIndex
	if err = cache.ReadValue(reader, &s.LogIndex, version); err != nil {
		return err
	}

	// Timestamp
	if err = cache.ReadValue(reader, &s.Timestamp, version); err != nil {
		return err
	}

	// Topics
	s.Topics = make([]base.Hash, 0)
	if err = cache.ReadValue(reader, &s.Topics, version); err != nil {
		return err
	}

	// TransactionHash
	if err = cache.ReadValue(reader, &s.TransactionHash, version); err != nil {
		return err
	}

	// TransactionIndex
	if err = cache.ReadValue(reader, &s.TransactionIndex, version); err != nil {
		return err
	}

	s.FinishUnmarshal()

	return nil
}

func (s *SimpleLog) FinishUnmarshal() {
	// EXISTING_CODE
	// EXISTING_CODE
}

// EXISTING_CODE
//

func (s *SimpleLog) getHaystack() string {
	haystack := make([]byte, 66*len(s.Topics)+len(s.Data))
	haystack = append(haystack, s.Address.Hex()[2:]...)
	for _, topic := range s.Topics {
		haystack = append(haystack, topic.Hex()[2:]...)
	}
	haystack = append(haystack, s.Data[2:]...)
	return string(haystack)
}

func (s *SimpleLog) ContainsAny(addrArray []base.Address) bool {
	haystack := s.getHaystack()
	for _, addr := range addrArray {
		if strings.Contains(string(haystack), addr.Hex()[2:]) {
			return true
		}
	}
	return false
}

func (s *SimpleLog) ContainsAddress(addr base.Address) bool {
	haystack := s.getHaystack()
	return strings.Contains(string(haystack), addr.Hex()[2:])
}

func (r *RawLog) RawToSimple(vals map[string]any) (SimpleLog, error) {
	hash, ok := vals["hash"].(base.Hash)
	if !ok {
		logger.Fatal("Hash not found in raw log values")
	}

	log := SimpleLog{
		Address:          base.HexToAddress(r.Address),
		BlockNumber:      utils.MustParseUint(r.BlockNumber),
		BlockHash:        base.HexToHash(r.BlockHash),
		TransactionIndex: utils.MustParseUint(r.TransactionIndex),
		TransactionHash:  hash,
		LogIndex:         utils.MustParseUint(r.LogIndex),
		Data:             r.Data,
		raw:              r,
	}
	for _, topic := range r.Topics {
		log.Topics = append(log.Topics, base.HexToHash(topic))
	}

	if ts, ok := vals["timestamp"].(base.Timestamp); ok && ts != utils.NOPOSI {
		log.Timestamp = ts
	}

	return log, nil
}

// EXISTING_CODE
