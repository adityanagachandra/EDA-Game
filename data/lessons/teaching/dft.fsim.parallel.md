# Parallel-fault and parallel-pattern simulation

A machine word holds 32 or 64 bits, and a bitwise AND evaluates all of them at once. Parallel fault simulation packs many two-valued simulations into one word so that one pass of the netlist does the work of many. The two classic packings, many faults per word or many patterns per word, dominate fault simulation of large scan designs and show how data layout, not algorithmic cleverness, buys an order of magnitude.

## Parallel-fault simulation

In **parallel-fault** simulation each bit of a word is a different machine for the same pattern. One bit holds the fault-free circuit and the other $W - 1$ bits hold faulty circuits. Each fault is injected by forcing its bit at the fault site with masks: OR with a mask to force 1, AND with an inverted mask to force 0. One pass evaluates $W - 1$ faults, so each pattern needs $\lceil F/(W-1) \rceil$ passes. With 16-bit words, 3036 faults and 86 patterns, that is $\lceil 3036/15 \rceil = 203$ passes per pattern and 17,458 in all without fault dropping. Comparing each faulty bit with the fault-free bit at the outputs (an XOR) reveals detection.

## Parallel-pattern single-fault simulation

**PPSFP** packs $W$ different patterns into one word and simulates one fault at a time. The fault-free circuit is simulated once per word of patterns, then each fault once per word: $\lceil P/W \rceil (F + 1)$ passes in all without dropping. PPSFP is preferred in practice because it combines naturally with fault dropping (a fault is dropped as soon as any bit of its output word differs) and with event-driven evaluation that only simulates the fault's fan-out cone, where faulty values can differ. Full scan makes the circuit combinational, which is exactly the setting where PPSFP shines.

## Bitwise evaluation

The engine of both methods is that two-valued gates map to machine instructions: AND, OR, XOR and NOT on whole words. If input A carries 10100111 and input B carries 10011010 for eight patterns, an OR gate outputs 10111111 in one instruction; a NAND outputs the complement of A AND B. Three-valued logic needs two bits per signal and a few more instructions per gate, which is why X-handling costs speed. Gates are evaluated in levelized order, as in compiled simulation, so the inner loop is a tight sequence of word operations.

## Going deeper

Bushnell and Agrawal (2000), chapter 5; Waicukauski et al., "Fault simulation for structured VLSI," VLSI Systems Design, 1985.
