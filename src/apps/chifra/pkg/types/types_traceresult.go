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
	"io"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EXISTING_CODE

type RawTraceResult struct {
	Address string `json:"address"`
	Code    string `json:"code"`
	GasUsed string `json:"gasUsed"`
	Output  string `json:"output"`
	// EXISTING_CODE
	// EXISTING_CODE
}

type SimpleTraceResult struct {
	Address base.Address    `json:"address,omitempty"`
	Code    string          `json:"code,omitempty"`
	GasUsed base.Gas        `json:"gasUsed,omitempty"`
	Output  string          `json:"output,omitempty"`
	raw     *RawTraceResult `json:"-"`
	// EXISTING_CODE
	// EXISTING_CODE
}

func (s *SimpleTraceResult) Raw() *RawTraceResult {
	return s.raw
}

func (s *SimpleTraceResult) SetRaw(raw *RawTraceResult) {
	s.raw = raw
}

func (s *SimpleTraceResult) Model(chain, format string, verbose bool, extraOptions map[string]any) Model {
	var model = map[string]interface{}{}
	var order = []string{}

	// EXISTING_CODE
	if format == "json" {
		if s.GasUsed > 0 {
			model["gasUsed"] = s.GasUsed
			order = append(order, "gasUsed")
		}
		if len(s.Output) > 2 { // "0x" is empty
			model["output"] = s.Output
			order = append(order, "output")
		}
		if !s.Address.IsZero() {
			model["address"] = s.Address
			order = append(order, "address")
		}
		if extraOptions["traces"] != true && len(s.Code) > 2 { // "0x" is empty
			model["code"] = utils.FormattedCode(verbose, s.Code)
			order = append(order, "code")
		}
	} else {
		model = map[string]interface{}{
			"gasUsed": s.GasUsed,
			"output":  s.Output,
		}

		order = []string{
			"gasUsed",
			"output",
		}
		if !s.Address.IsZero() {
			model["address"] = hexutil.Encode(s.Address.Bytes())
			order = append(order, "address")
		}
		// if len(s.Output) > 0 && s.Output != "0x" {
		// 	model["output"] = s.Output
		// 	order = append(order, "output")
		// }
	}
	// EXISTING_CODE

	return Model{
		Data:  model,
		Order: order,
	}
}

// --> marshal_only
func (s *SimpleTraceResult) MarshalCache(writer io.Writer) (err error) {
	// Address
	if err = cache.WriteValue(writer, s.Address); err != nil {
		return err
	}

	// Code
	if err = cache.WriteValue(writer, s.Code); err != nil {
		return err
	}

	// GasUsed
	if err = cache.WriteValue(writer, s.GasUsed); err != nil {
		return err
	}

	// Output
	if err = cache.WriteValue(writer, s.Output); err != nil {
		return err
	}

	return nil
}

func (s *SimpleTraceResult) UnmarshalCache(version uint64, reader io.Reader) (err error) {
	// Address
	if err = cache.ReadValue(reader, &s.Address, version); err != nil {
		return err
	}

	// Code
	if err = cache.ReadValue(reader, &s.Code, version); err != nil {
		return err
	}

	// GasUsed
	if err = cache.ReadValue(reader, &s.GasUsed, version); err != nil {
		return err
	}

	// Output
	if err = cache.ReadValue(reader, &s.Output, version); err != nil {
		return err
	}

	s.FinishUnmarshal()

	return nil
}

func (s *SimpleTraceResult) FinishUnmarshal() {
	// EXISTING_CODE
	// EXISTING_CODE
}

// EXISTING_CODE
// EXISTING_CODE
