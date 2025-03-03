package whenPkg

import (
	"context"
	"errors"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/tslib"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/ethereum/go-ethereum"
)

func (opts *WhenOptions) HandleList() error {
	chain := opts.Globals.Chain

	ctx, cancel := context.WithCancel(context.Background())
	fetchData := func(modelChan chan types.Modeler[types.RawNamedBlock], errorChan chan error) {
		results, err := tslib.GetSpecials(chain)
		if err != nil {
			errorChan <- err
			if errors.Is(err, ethereum.NotFound) {
				return
			}
			cancel()
			return
		}

		for _, result := range results {
			// Note: This is needed because of a GoLang bug when taking the pointer of a loop variable
			result := result
			modelChan <- &result
		}
	}

	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOpts())
}
