# Control-data-flow graphs

Real programs have branches and loops, which a data-flow graph cannot express. The control-data-flow graph (CDFG) is the two-level representation HLS tools use instead: a control-flow graph of basic blocks, each holding an acyclic DFG. It decides where a scheduler can move operations, where the controller needs states and transitions, and which operations can share hardware because they never run together.

## Basic blocks inside a control-flow graph

A basic block is a maximal straight-line run of operations entered only at the top and left only at the bottom. The CDFG has one node per basic block, with control edges for every possible transfer: fall-through, the two arms of a branch, a loop's back edge and its exit. Inside each block sits that block's DFG, which is acyclic. A loop such as 'for i in 0..N-1: s += a[i]' becomes an entry block, a header that tests $i < N$, a body that loads $a[i]$ and adds, and an exit; the cycle lives in the control level as the back edge from body to header. A function with blocks entry, header, body and exit therefore carries four DFGs. The controller FSM that HLS later emits mirrors this structure: states for the steps of each block, transitions along the control edges.

## Branches, if-conversion and exclusive arms

A branch can stay as control, with the controller choosing which arm's states to enter, or be converted into data. If-conversion computes both arms and selects the result with a multiplexer driven by the condition, turning a control dependence into a data dependence and merging the arms into one larger basic block that the scheduler packs freely. The cost is work: both arms execute every time, so two multiplications in different arms both consume multiplier slots even though one result is discarded. When the branch stays as control, the arms are mutually exclusive, and a CDFG-aware binder exploits that: an operation in the then-arm and one in the else-arm can share a single unit even in the same control step, since they never run in the same execution. Loops need a known trip count, or a bound, before the tool can report a fixed latency.

## What the representation lets the scheduler do

A scheduler that works one basic block at a time treats block boundaries as barriers, and loses the chance to fill idle steps with operations from neighbouring blocks. Global techniques move operations across those boundaries: code hoisting, trace and path-based scheduling, and speculation, which executes an operation before the branch that guards it. Speculation is safe only for operations without side effects whose result can be dropped if the branch goes the other way; a store to memory must stay guarded. Loop-carried values, such as an accumulator $s$ updated every iteration, enter the header through a phi node fed from the end of the body. They are not edges inside the body's DFG, which stays acyclic, but they are exactly the recurrences that bound a pipelined loop's initiation interval later. A while-loop with a data-dependent exit forces the controller to test every iteration, so its latency is only known as a range.

## Going deeper

Gajski, Dutt, Wu and Lin, "High-Level Synthesis" (1992), chapter on design representation; Coussy and Morawiec (eds.), "High-Level Synthesis: From Algorithm to Digital Circuit" (2008), chapters on front-end compilation.
