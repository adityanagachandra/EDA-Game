# Carry chains

A carry chain is dedicated hard logic that passes the carry of an adder from one bit to the next through fast local wires, bypassing the LUTs and general routing. It turns the slowest part of soft-logic arithmetic into a few tens of picoseconds per bit, and it constrains placement, since chained blocks must stack in one column. Mappers, packers and placers all treat it specially.

## Why soft carries are slow

In a ripple adder each carry is $c_{i+1} = g_i + p_i c_i$. Without hard support, every carry lives in its own LUT, so each bit costs one LUT delay plus one general-routing hop. With 200 ps per LUT and 250 ps per hop, a 32-bit ripple takes $32 \times 450 = 14{,}400$ ps. A dedicated chain computes propagate and generate in the LUT and then ripples the carry through a hard multiplexer per bit on a dedicated wire. If entering and leaving the chain costs 160 ps and each bit 25 ps, a 27-bit adder takes $160 + 27 \times 25 = 835$ ps. The chain is linear in $n$, but its per-bit constant is about twenty times smaller than the soft one.

## Chains occupy columns

Each logic block's carry logic handles a fixed number of bits, for example 4 or 8, and the carry-out drives the carry-in of the block directly above. An $n$-bit adder with $b$ bits per block therefore needs $\lceil n/b \rceil$ vertically adjacent blocks: a 64-bit adder with 4 bits per block occupies 16 blocks in one column. The synthesis tool emits carry primitives, the packer keeps each slice of the chain together, and the placer must treat the whole chain as a rigid macro that moves as one. A chain longer than a column, or one that collides with a hard-block column, cannot be placed as drawn and must be split, which adds a slow soft connection at the break. The speed comes at the price of placement freedom.

## When soft prefix adders win

A parallel-prefix adder (Kogge-Stone, Brent-Kung) has depth logarithmic in $n$, which wins in ASICs. On an FPGA every prefix level costs a LUT and a routing hop, a few hundred picoseconds, while the chain costs tens of picoseconds per bit. Model the chain as $t_0 + t_c n$ and the soft prefix adder as $(\lceil \log_2 n \rceil + 1) t_\ell$. With $t_0 = 180$, $t_c = 36$ and $t_\ell = 430$ ps, the chain is faster up to $n = 90$, and the prefix adder first wins at $n = 91$. In practice the crossover is wide, so synthesis tools map ordinary additions onto the chain and use carry-select or prefix structures only for very wide operands. Algorithms that are optimal for ASICs can be wrong for FPGAs because the cost model differs.

## Going deeper

Kuon, Tessier and Rose, "FPGA Architecture: Survey and Challenges" (2008), section on arithmetic; Hauck, Hosler and Fry, "High-Performance Carry Chains for FPGAs", IEEE Transactions on VLSI Systems (2000).
