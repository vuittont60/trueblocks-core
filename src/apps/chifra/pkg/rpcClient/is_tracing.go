package rpcClient

// TODO: BOGUS This needs to be implemented in a cross-chain, cross-client manner

// IsNodeTracing returns true if the node is an archive node. Note currently always returns true.
func (options *Options) IsNodeTracing(testMode bool) bool {
	// TODO: We can test this with a unit test
	if testMode && options.Chain == "non-tracing" {
		return false
	}
	return true
}
