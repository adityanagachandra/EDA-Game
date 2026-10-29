# Lookup tables

A lookup table is the logic cell of an FPGA: a small SRAM that stores a truth table, read out by a multiplexer whose select lines are the logic inputs. Every function of $k$ inputs costs one $k$-LUT, which is why FPGA synthesis counts LUTs and levels rather than gates. A tool designer needs the LUT's bit cost, its function count and how modern fracturable LUTs split, because packing and bitstream generation are built on them.

## Configuration bits and the multiplexer tree

A $k$-input LUT holds one configuration bit for each input combination, $2^k$ bits in all. The inputs drive a binary tree of 2:1 multiplexers whose leaves are those bits, so the tree has $2^k - 1$ multiplexers and depth $k$. Reading the LUT is addressing a memory: with inputs $(a, b, c)$ and $a$ the most significant select, the output is the bit at address $4a + 2b + c$. If a 3-LUT stores 0 1 1 0 1 0 0 1 at addresses 0 to 7 (three-input XOR), then $(1, 0, 1)$ selects address 5 and outputs 0. A device with 1000 6-LUTs carries $1000 \times 64 = 64{,}000$ LUT bits before any routing configuration is counted. The SRAM is written once at configuration and then held, so the LUT is combinational logic in operation.

## How many functions a LUT implements

Because every assignment of the $2^k$ bits is a distinct truth table, one $k$-LUT implements all $2^{2^k}$ Boolean functions of $k$ inputs: 16 for $k = 2$, 256 for $k = 3$, 65,536 for $k = 4$. This is what makes FPGA mapping unlike standard-cell mapping: there is no library of shapes to match, only a limit on the number of inputs. A function of $j < k$ inputs still fits in one $k$-LUT. The router ties the $k - j$ unused inputs to constants, so only the $2^j$ addresses that agree with those constants are ever read. Tools usually replicate the $2^j$-bit truth table across the remaining addresses, so that whatever constant the router chooses, the output is right. The cost of a LUT is therefore set by its input count, not by the complexity of its function.

## Fracturable LUTs

Commercial 6-input LUTs are fracturable: the 64-bit table is two 32-bit halves, so the element can produce one function of up to 6 inputs or two functions that together use at most 5 distinct inputs, each on its own output. The second function is nearly free area, so the packer tries to pair small functions. Suppose $F_1$ uses $\{a, b, c\}$ and $F_2$ uses $\{b, c, d, e\}$. Their union has 5 signals, so they share one physical LUT. If $F_2$ instead used $\{d, e, f\}$, the union would have 6 signals and the pair would need two LUTs. A greedy packer that walks the functions in order and pairs each with the first compatible partner gives a quick upper bound; the true minimum is a matching problem on the compatibility graph. This is why mapped LUT counts and packed LUT counts differ.

## Going deeper

Betz, Rose and Marquardt, "Architecture and CAD for Deep-Submicron FPGAs" (1999), chapter 2; Kuon, Tessier and Rose, "FPGA Architecture: Survey and Challenges", Foundations and Trends in EDA (2008).
