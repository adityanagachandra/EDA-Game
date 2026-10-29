# Scan insertion and chain balancing

Scan insertion is the synthesis-stage transformation that replaces flip-flops with scan cells, fixes scan design-rule violations, and stitches the cells into chains. Its main quantitative decision is how many chains to build and how to balance them, because test time is set by the longest chain. A placer or DFT tool designer must also respect clock domains, which constrain which flops may share a chain.

## Design rules and stitching

Before stitching, the tool checks **scan DRC**. Every scan flop's clock must be controllable from a test clock, asynchronous sets and resets must be held inactive during shift, and there must be no combinational loops or uncontrolled latches in the shift path. Violating flops are fixed (for example by adding test-mode gating) or excluded as non-scan. Chains are then stitched, typically grouped by clock domain and edge. Where a chain must cross between domains or edges, a **lockup latch** is inserted to avoid hold violations during shift. Equivalence checking with SE tied to 0 confirms the functional behaviour is unchanged.

## Balancing

Chains shift in parallel, so the longest one sets the shift length. Ideally $F$ flops in $C$ chains give length $\lceil F/C \rceil$: 20,826 flops in 27 chains give at most 772. Clock-domain grouping often produces unbalanced chains. With lengths 424, 778, 398, 447 and 450 (2497 flops in 5 chains), the longest is 778; rebalanced, it would be $\lceil 2497/5 \rceil = 500$. For 1406 patterns the test drops from $1406 \times 779 + 778 = 1{,}096{,}052$ cycles to $1406 \times 501 + 500 = 704{,}906$, saving 391,146 cycles. Balancing across domains needs lockup latches, which is the price of the saving.

## Choosing the chain count

Chain count is limited by scan pins on the package and channels on the tester. Given a cycle budget, the smallest number of chains $C$ that meets it is found from $p(\lceil F/C \rceil + 1) + \lceil F/C \rceil \le$ budget. Since test time scales roughly as $pF/C$, doubling the chains halves it until pins run out; beyond that, scan compression decouples the internal chain count from the external pin count.

## Going deeper

Bushnell and Agrawal (2000), chapter 14; Wang, Wu and Wen (2006), chapter 2; the Fault and OpenROAD-flow DFT documentation for open-source scan insertion.
