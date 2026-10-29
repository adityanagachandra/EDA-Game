# Carry-save reduction trees

Adding many operands, as in a multiplier's partial products or a sum of products, is slow if each addition propagates carries. Carry-save adders postpone carries: a 3:2 compressor turns three rows of bits into two in one full-adder delay, regardless of width. A tree of them reduces any number of operands to two, and a single fast adder finishes. Synthesis recognizes multi-operand sums and builds such trees.

## Levels of a 3:2 reduction tree

A word-level carry-save adder takes three $n$-bit rows and produces a sum row and a carry row, using $n$ full adders in parallel, with a delay of one full adder. Each level reduces the row count from $r$ to about $\lceil 2r/3 \rceil$. Working backwards from two rows, the largest row count reducible in $j$ levels follows the Dadda sequence: 2, 3, 4, 6, 9, 13, 19, 28, 42, each term the floor of 1.5 times the previous. So 17 operands need 6 levels (17 is at most 19), and 9 operands need 4. Depth grows logarithmically, with base 1.5, in the number of operands.

## How many carry-save adders

Each 3:2 carry-save adder consumes three rows and produces two, a net reduction of exactly one row. Reducing $k$ rows to two therefore needs exactly $k - 2$ word-level carry-save adders, however they are arranged: a linear chain and a balanced tree use the same count. With $n$-bit words that is about $(k - 2)n$ full adders, ignoring the growth of word width. The arrangement changes delay, not area: the chain has $k - 2$ levels, the tree logarithmically many. This separation, fixed count but variable depth, is why reduction trees are a pure delay optimization.

## Total delay and synthesis

After the tree, a carry-propagate adder adds the last two rows. With a ripple-carry final adder of $n$ bits, the total is about (tree levels) + $n$ full-adder delays; a parallel-prefix final adder makes the second term logarithmic. A 7-operand, 24-bit sum needs 4 levels plus 24 ripple stages, 28 full-adder delays, against far more for a chain of six carry-propagate adders. Wallace and Dadda trees differ in how greedily they compress: Wallace reduces as much as possible at each level, while Dadda reduces only as much as the next level requires, which typically needs fewer adder cells. In synthesis, Yosys's alumacc gathers sums of products into a multiply-accumulate cell before lowering, and commercial tools re-balance carry-save trees by arrival times.

## Going deeper

Wallace, "A suggestion for a fast multiplier", IEEE Transactions on Electronic Computers (1964); Dadda, "Some schemes for parallel multipliers", Alta Frequenza (1965); Ercegovac and Lang, "Digital Arithmetic" (2004).
