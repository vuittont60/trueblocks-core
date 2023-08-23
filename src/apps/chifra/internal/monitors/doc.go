// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

// Package monitorsPkg handles the chifra monitors command. It  has two purposes: (1) to --watch a set of addresses. This function is in its early stages and will be better explained elsewhere. Please see an example of what one may do with chifra monitors --watch, and (2) allows one to manage existing monitored addresses. A "monitor" is simply a file on a hard drive that represents the transactional history of a given Ethereum address. Monitors are very small, being only the <block_no><tx_id> pair representing each appearance of an address. Monitor files are only created when a user expresses interest in a particular address. In this way, TrueBlock is able to continue to work on small desktop or even laptop computers. (See chifra list.) You may use the --delete command to delete (or --undelete if already deleted) an address. The monitor is not removed from your computer if you delete it. It is just marked as being deleted making it invisible to the TrueBlocks explorer. Use the --remove command to permanently remove a monitor from your computer. This is an irreversible operation and requires the monitor to have been previously deleted. 
package monitorsPkg