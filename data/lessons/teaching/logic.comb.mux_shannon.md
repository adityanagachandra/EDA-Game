# Multiplexer trees and lookup tables

A multiplexer is Shannon expansion built in hardware: $f = \bar{x} f_{\bar{x}} + x f_x$ is a 2-to-1 mux with select $x$ choosing between two cofactors. Repeating the expansion gives a mux tree that implements any function, and an FPGA lookup table is exactly that tree with its leaves stored in configuration memory. Technology mappers for FPGAs reason in these terms.

## Any function as a mux tree

Expand $f(x_1, \dots, x_n)$ on $x_1$, then each cofactor on $x_2$, and so on. Each expansion is one 2-to-1 mux. With constant leaves (the truth table values) a full binary tree over $n$ variables needs

$$1 + 2 + \dots + 2^{n-1} = 2^n - 1$$

two-input muxes. Feeding the last variable into the data inputs saves a level: each leaf is one of $0, 1, x_n, \bar{x}_n$, and the tree over $n - 1$ select variables needs $2^{n-1} - 1$ muxes. For $n = 4$ that is 15 muxes with constant leaves, or 7 with the last variable at the leaves.

## The LUT and its configuration word

A $k$-input lookup table stores $2^k$ configuration bits and uses its inputs as the select lines of a $2^k$-to-1 mux. Programming a LUT means writing the function's truth table into those bits.

Fix the convention: row index $i = 4a + 2b + c$ for a 3-input LUT, and bit $i$ of the configuration word is $f$ on that row. For $f = a \oplus (b c)$: rows 0 to 7 give $f = 0,0,0,1,1,1,1,0$. Bits 3, 4, 5 and 6 are set, so the word is $2^3 + 2^4 + 2^5 + 2^6 = 120$. FPGA vendors expose this word as the LUT's INIT parameter, and a bitstream is, to first order, a long list of such words plus routing switch settings.

## Mapping wide functions onto small LUTs

A $k$-LUT implements any function of up to $k$ inputs in one level, but a wider function needs several. For an AND (or any associative function) of $n$ inputs, each LUT absorbs $k$ signals and emits one, reducing the signal count by $k - 1$. The tree needs

$$\left\lceil \frac{n - 1}{k - 1} \right\rceil$$

LUTs and $\lceil \log_k n \rceil$ levels. A 16-input AND in 6-LUTs takes $\lceil 15/5 \rceil = 3$ LUTs and 2 levels. An arbitrary function is far worse: Shannon expansion of a $(k+1)$-input function takes two $k$-LUTs plus a mux. LUT mapping algorithms such as FlowMap (depth-optimal) and the cut-based mappers in ABC enumerate $k$-feasible cuts of the network to choose which subcircuits each LUT absorbs.

## Going deeper

Cong and Ding, "FlowMap: An Optimal Technology Mapping Algorithm for Delay Optimization in Lookup-Table Based FPGA Designs" (IEEE TCAD, 1994); the ABC documentation for the if mapper.
