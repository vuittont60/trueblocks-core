## chifra state

Use this tool to retrieve the balance of an address (or list of addresses) at the given block (or blocks). Specify multiple addresses and/or multiple blocks if you wish, but you must specify at least one address. If no block is specified, the latest block is reported.

You may also query to see if an address is a smart contract as well as retrieve a contract's byte code.

### usage

`Usage:`    chifra state [-p|-c|-n|-v|-h] &lt;address&gt; [address...] [block...]  
`Purpose:`  Retrieve account balance(s) for one or more addresses at given block(s).

`Where:`

{{<td>}}
|          | Option                          | Description                                                                                                                |
| -------- | ------------------------------- | -------------------------------------------------------------------------------------------------------------------------- |
|          | addrs                           | one or more addresses (0x...) from which to retrieve<br/>balances (required)                                               |
|          | blocks                          | an optional list of one or more blocks at which to<br/>report balances, defaults to 'latest'                               |
| &#8208;p | &#8208;&#8208;parts &lt;val&gt; | control which state to export, one or more of [none,<br/>some, all, balance, nonce, code, storage, deployed,<br/>accttype] |
| &#8208;c | &#8208;&#8208;changes           | only report a balance when it changes from one block<br/>to the next                                                       |
| &#8208;n | &#8208;&#8208;no_zero           | suppress the display of zero balance accounts                                                                              |
| &#8208;x | &#8208;&#8208;fmt &lt;val&gt;   | export format, one of [none, json, txt, csv, api]                                                                          |
| &#8208;v | &#8208;&#8208;verbose           | set verbose level (optional level defaults to 1)                                                                           |
| &#8208;h | &#8208;&#8208;help              | display this help screen                                                                                                   |
{{</td>}}

`Notes:`

- An `address` must start with '0x' and be forty-two characters long.
- `blocks` may be a space-separated list of values, a start-end range, a `special`, or any combination.
- If the queried node does not store historical state, the results are undefined.
- `special` blocks are detailed under `chifra when --list`.
- `balance` is the default mode. To select a single mode use `none` first, followed by that mode.
- You may specify multiple `modes` on a single line.

**Source code**: [`tools/getState`](https://github.com/TrueBlocks/trueblocks-core/tree/master/src/tools/getState)
