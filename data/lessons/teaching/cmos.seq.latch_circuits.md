# Transistor-level latches and flip-flops

Every register in a digital design is a library flip-flop, and every flip-flop is built from latches made of inverters and transmission gates. Its setup time, hold time and clock-to-q are not abstract parameters: they come from specific internal nodes settling. A tool designer who knows the circuit can read a library's sequential arcs, understand why hold can be negative, and see what characterization is searching for.

## The static latch

A latch has two states. When transparent, its output follows its input; when opaque, it holds the last value. A static CMOS latch implements this with an input transmission gate and a feedback loop of two inverters closed by a second transmission gate. A positive-level latch is transparent while the clock is high: the input transmission gate is on and the feedback is open. When the clock falls, the input gate turns off and the feedback gate closes, so the cross-coupled inverters regenerate and hold the value indefinitely. A dynamic latch omits the feedback and stores the value as charge on a node, which is smaller but loses the value to leakage if the clock stops.

## Master-slave flip-flops

An edge-triggered flip-flop is two latches of opposite polarity in series. For a positive-edge flop, the master is transparent while the clock is low and the slave while the clock is high. Before the rising edge the master follows D and the slave holds the old Q. At the edge the master closes, capturing D, and the slave opens, passing the captured value to Q. The flop is never transparent from D to Q, which is what makes it edge-triggered. Library flops buffer the clock internally to make the complementary clocks locally, so the CK pin presents one small fixed capacitance, and they buffer Q so its drive does not depend on the internal nodes. Scan flops add a 2:1 mux on the input selecting functional data or scan data.

## Where setup, hold and clock-to-q come from

Setup time is how long before the edge D must be stable so that the new value has propagated around the master's loop before the input switch closes. Hold time is how long after the edge D must stay stable, because the internal clock reaches the master's input switch later than the CK pin and turns it off gradually. If the internal data path into the master is slower than the internal clock path, the hold requirement can be negative. Clock-to-q is the delay from the edge, through the slave's transparency and the output buffer, to Q. These three numbers trade against each other: data arriving closer to the edge pushes clock-to-q out, which is exactly what characterization measures when it searches for setup and hold.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), chapter 10; Stojanovic and Oklobdzija, "Comparative analysis of master-slave latches and flip-flops for high-performance and low-power systems" (IEEE JSSC, 1999).
