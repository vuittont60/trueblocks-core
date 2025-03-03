// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package blocksPkg

import (
	"context"
	"errors"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/ethereum/go-ethereum"
)

func (opts *BlocksOptions) HandleTraces() error {
	chain := opts.Globals.Chain

	ctx, cancel := context.WithCancel(context.Background())
	fetchData := func(modelChan chan types.Modeler[types.RawTrace], errorChan chan error) {
		for _, br := range opts.BlockIds {
			blockNums, err := br.ResolveBlocks(chain)
			if err != nil {
				errorChan <- err
				if errors.Is(err, ethereum.NotFound) {
					continue
				}
				cancel()
				return
			}

			for _, bn := range blockNums {
				var traces []types.SimpleTrace
				traces, err = opts.Conn.GetTracesByBlockNumber(bn)
				if err != nil {
					errorChan <- err
					if errors.Is(err, ethereum.NotFound) {
						continue
					}
					cancel()
					return
				}

				for _, trace := range traces {
					// Note: This is needed because of a GoLang bug when taking the pointer of a loop variable
					trace := trace
					modelChan <- &trace
				}
			}
		}
	}

	extra := map[string]interface{}{
		"uncles":     opts.Uncles,
		"logs":       opts.Logs,
		"traces":     opts.Traces,
		"addresses":  opts.Uniq,
		"articulate": opts.Articulate,
	}
	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOptsWithExtra(extra))
}
