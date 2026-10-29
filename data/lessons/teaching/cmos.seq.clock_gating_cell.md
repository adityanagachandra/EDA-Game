# Integrated clock-gating cells

The clock has an activity factor of 1 and a large capacitance, so stopping it where registers would hold their values anyway is the single most effective dynamic-power optimization in digital design. Synthesis tools insert clock gates automatically, using a dedicated library cell, the integrated clock gate (ICG). Its circuit explains both why it is glitch-free and why timers apply a special check to its enable.

## Why a plain AND fails

Gating a positive-edge clock with an AND gate, GCLK = CLK AND EN, works only if EN never changes while CLK is high. If EN falls in the middle of a high phase, GCLK is cut short, producing a runt pulse that may or may not clock the flops. If EN rises while CLK is high, GCLK rises in mid-cycle, creating a spurious edge. Because EN comes from logic launched by the same clock, it naturally changes shortly after the rising edge, exactly when CLK is high. An unprotected AND gate would therefore glitch on almost every enable change. The fix must guarantee that EN is stable throughout every high phase of CLK.

## The latch-based ICG

An ICG for positive-edge flops contains a latch transparent while CLK is low, followed by an AND with CLK. While CLK is low the latch passes EN, and the AND output is low anyway. When CLK rises the latch closes and holds EN constant for the whole high phase, so the AND passes either a complete clock pulse or nothing. A test-enable input is usually ORed with EN before the latch so that scan shifting can force the clock on. Liberty marks such a cell with the clock_gating_integrated_cell attribute (for example latch_posedge, or latch_posedge_precontrol with a test pin), which tells synthesis and clock-tree tools what the cell does. The enable path to the latch is timed with a clock-gating setup and hold check, just like a flop's data.

## Where and when gating pays

A clock gate replaces the hold multiplexers in front of a bank of flops: instead of each flop recirculating its value when EN is low, the whole bank's clock stops. That saves the mux area and the clock power of the flops and of the clock-tree branch below the gate. The saving is proportional to the fraction of cycles in which EN is off; a gate whose enable is almost always on saves nothing and costs an ICG. Placing gates higher in the clock tree stops more buffers but requires an enable valid for more flops. Clock-tree synthesis treats ICGs as part of the clock path, balancing their delay against ungated branches, and moving an ICG changes its enable timing check, so gating is decided jointly with the tree.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), section 10.3; the Liberty Reference Manual, clock_gating_integrated_cell; Keating and others, "Low Power Methodology Manual" (2007), chapter on clock gating.
