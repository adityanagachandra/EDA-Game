# False and multicycle paths

By default a timer assumes every path launched on one edge must be captured on the next. Real designs violate that assumption on purpose: some paths never carry a meaningful transition, some have several cycles to settle, and some clocks have no phase relationship at all. Timing exceptions tell the timer so. Each one removes or relaxes checks, so each is a claim about the design that must be justified.

## False paths and clock groups

set_false_path removes the timing checks on matching paths, for example from a configuration register that is written once at boot and held constant while the datapath runs. A wrong false path hides a real failure, so every one needs a reviewable reason. Between whole clocks, set_clock_groups is the clearer tool. -asynchronous declares that clocks from unrelated sources have no phase relationship, so crossings are not timed; synchronizers handle metastability. -logically_exclusive declares clocks that can never be active together on the same logic, such as two inputs of a clock multiplexer. One clock-group statement covers every crossing in both directions, which pairwise false paths easily miss, and timers also use it when analyzing crosstalk between the clocks.

## Multicycle paths and absolute delays

set_multicycle_path $N$ -setup moves the setup check to the $N$-th capture edge, for logic that is sampled only every $N$ cycles, typically under an enable. It must normally be paired with set_multicycle_path $N-1$ -hold, because moving setup also moves the default hold check forward. set_max_delay and set_min_delay replace the edge-derived requirement with an absolute bound, which is useful on asynchronous crossings that still need a bounded delay. Exceptions select paths with -from, -to and -through; -through filters on the pins a path visits, which keeps an exception narrow. When several exceptions match a path, precedence is by type: false paths override max and min delays, which override multicycle paths.

## Exceptions as claims

An exception is a statement about behavior that the timer cannot verify. If a multicycle path's logic does not actually wait $N$ cycles, the report is optimistic and the failure moves to silicon. Exceptions are also mode-specific: a path false in mission mode may be real during scan shift, so each mode's SDC carries its own justified set. Good practice keeps exceptions few, narrow and documented, checks that each one matches the intended paths, and reviews any that were added to make a failing path disappear. Formal tools can prove some exceptions, for instance that a register is constant or that a path is sampled only under an enable, which turns a claim into a checked fact.

## Going deeper

Gangadharan and Churiwala, "Constraining Designs for Synthesis and Timing Analysis" (2013), chapters on exceptions and clock groups; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 8.
