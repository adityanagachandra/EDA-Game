# Multi-FPGA partitioning and pin multiplexing

A large SoC does not fit in one FPGA, so emulators and prototypes split it across many. Every net cut by the partition becomes an inter-FPGA signal, and FPGAs have far fewer pins than a large design has cut nets. Time-division multiplexing shares pins among signals at the cost of speed. The resulting system clock is set directly by the cut size, which makes this a graph-partitioning problem with a visible, measurable objective.

## Cut nets become inter-FPGA signals

Model the design as a hypergraph: blocks are vertices and each net is a hyperedge connecting the blocks it touches. Assign each block to an FPGA. A net is **cut** if its blocks sit on more than one FPGA, and each cut net needs a signal between FPGAs. Example: blocks B1, B2, B3 on FPGA A and B4, B5 on FPGA B, with nets {B1, B2}, {B2, B4}, {B3, B5, B1}, {B4, B5} and {B1, B4}. The cut nets are {B2, B4}, {B3, B5, B1} and {B1, B4}: 3 inter-FPGA signals. Moving B4 to FPGA A would cut {B4, B5} instead of {B2, B4} and {B1, B4}, reducing the cut to 2. Partitioners minimize the cut subject to each FPGA's logic capacity, using multilevel methods such as those in hMETIS.

## Time-division multiplexing

Real partitions of large designs cut tens of thousands of signals, while two FPGAs may share only a few hundred pins. **Time-division multiplexing** (TDM) sends several signals over one pin in successive time slots, serializing them at the sender and deserializing at the receiver. The required ratio is $r = \lceil S / P \rceil$ for $S$ cut signals and $P$ available pins. Example: 4,100 signals over 480 pins gives $\lceil 8.54 \rceil = 9$ signals per pin. Ratios differ by link, since each FPGA pair has its own pin budget, and the worst link sets the pace. High-speed serial transceivers raise the per-pin data rate, which allows higher ratios at the same system speed.

## The system clock and why the cut matters

All multiplexed signals must cross before the next system clock edge, so one system cycle needs at least $r$ link cycles, plus overhead for alignment and synchronization. With link clock $f_{link}$, ratio $r$ and $o$ overhead cycles, the system clock is at most $f_{link}/(r + o)$. Example: 500 MHz links, $r = 16$ and 4 overhead cycles give $500/20 = 25$ MHz. Halving the cut halves $r$: with $r = 8$ the clock rises to $500/12 \approx 41.67$ MHz. Paths that cross several FPGAs within one cycle pay this delay more than once. This is a direct line from partition quality to prototype speed, and it is why partitioners for emulation optimize cut size, balance and the number of FPGA hops on timing-critical paths together.

## Going deeper

Karypis and Kumar, "Multilevel k-way hypergraph partitioning", VLSI Design (2000); Hauck and DeHon (eds.), "Reconfigurable Computing" (2008), the chapter on multi-FPGA systems; Babb, Tessier and Agarwal, "Virtual wires: overcoming pin limitations in FPGA-based logic emulators", FCCM (1993).
