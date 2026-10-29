# Scoreboards and reference models

A scoreboard is the testbench component that decides pass or fail. It holds the expected results produced by a reference model and matches them against what monitors observe at the DUT outputs. Its matching policy has to fit the DUT: in-order, out-of-order, or keyed by an identifier. Its sizing and end-of-test checks catch lost and duplicated transactions, which are among the most common bugs in interconnect, queues and caches.

## In-order scoreboards

For a DUT that preserves order, such as a FIFO or a pipeline, the scoreboard keeps a queue of expected transactions. Each input transaction goes through the reference model and its prediction is pushed; each output transaction pops the oldest expected one and compares. Example: expected 12, 40, 33, 71, 58 and actual 12, 40, 35, 71, 58: one comparison fails, at the third position. A dropped transaction is worse: every later comparison shifts by one and fails, so a single bug produces a cascade of mismatches. Good scoreboards report the first mismatch with both transactions in full and stop or resynchronize. At the end of the test the queue must be empty; leftover entries mean outputs never appeared.

## Out-of-order matching

Many DUTs reorder legitimately: a memory controller completes requests by bank, a network router by port, a bus with transaction IDs by response latency. An in-order comparison would flag every reordering as an error. Instead the scoreboard matches each output against any outstanding expected item, usually by a key such as a transaction ID or address, and removes the match. Without a unique key it treats expected and actual items as multisets. Example: expected 12, 40, 33, 71 and actual 71, 12, 35, 40: matching removes 71, 12 and 40, and 35 is unmatched, so one actual output finds no match and 33 remains expected. The cost is a search per output, which a hash map keyed by ID makes constant time.

## Sizing queues and choosing reference models

The expected queue holds every transaction that has entered the DUT but not yet left. Its peak size is the maximum number of transactions in flight. If transactions enter at cycles 1, 2, 4, 5 and leave at cycles 6, 7, 8, 10, then at cycle 5 four are outstanding (none has left), which is the maximum. A peak far above the DUT's buffering capacity points to a lost or stalled transaction long before the end of the test. The reference model that fills the queue is usually written at transaction level: a SystemVerilog class, a C or C++ model called through DPI-C, or a Python model in cocotb. It must be independent of the RTL, otherwise a shared misunderstanding of the specification hides the bug in both.

## Going deeper

Spear and Tumbush, "SystemVerilog for Verification" (3rd ed., 2012), chapter on building a testbench; the UVM 1.2 class reference (uvm_in_order_comparator and analysis FIFOs).
