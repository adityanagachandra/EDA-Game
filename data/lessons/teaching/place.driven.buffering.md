# Buffering and sizing during placement

Placement and logic optimisation are not separate in modern flows. Long wires that placement cannot shorten are split by buffers, and gates driving heavy loads are upsized, during and right after global placement, so that timing estimates stay realistic. This placement-time optimisation, done in OpenROAD by the Resizer's repair_design and repair_timing, feeds back into where cells must go.

## Delay of a long unbuffered wire

For a driver of resistance $R_d$ driving a wire of length $L$ into a load $C_L$, the Elmore delay is $R_d(cL + C_L) + rL(cL/2 + C_L)$, whose wire term grows with $L^2$. A 2400 um wire at 0.5 ohm/um and 0.1 fF/um driven through 100 ohm into 6 fF has a driver term of $100 \cdot 246 = 24600$ ohm-fF and a wire term of $1200 \cdot 126 = 151200$ ohm-fF, 175.8 ps in total, dominated by the wire. No placement move can fix this if the two endpoints must be far apart, because they connect to logic or pins at opposite ends of the block.

## A repeater at the midpoint

Inserting a buffer of input capacitance $C_b$, output resistance $R_b$ and intrinsic delay $t_b$ at the midpoint splits the net into two stages of length $L/2$. Stage 1: $R_d(cL/2 + C_b) + r(L/2)(cL/4 + C_b)$. Stage 2: $t_b + R_b(cL/2 + C_L) + r(L/2)(cL/4 + C_L)$. Each half's wire term is a quarter of the original, so the two together are about half, at the cost of the buffer's delay and its own drive. For a 1800 um wire with $R_d = 300$ ohm and a buffer with $R_b = 150$ ohm and $t_b = 14$ ps, the delay drops from 276 ps to 183.5 ps. For short wires the buffer's intrinsic delay outweighs the saving, so buffering pays only beyond a critical length.

## Buffering and sizing inside the placement flow

Optimal buffer positions and counts on trees are found by van Ginneken's dynamic program, and gate sizing chooses drive strengths, but both need placement to know wire lengths. Flows therefore interleave: global placement, then repair of long wires, slew and capacitance violations by buffering and sizing, then incremental placement or legalization of the new cells, then timing repair. Some placers use virtual buffering, estimating buffered delays during timing-driven placement without inserting cells. The added cells need whitespace, which is one reason placement targets less than full density. A buffer inserted far from where legalization can put it loses its benefit, so buffer placement must respect density and blockages.

## Going deeper

van Ginneken, "Buffer placement in distributed RC-tree networks for minimal Elmore delay" (ISCAS 1990); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8; the OpenROAD Resizer documentation.
