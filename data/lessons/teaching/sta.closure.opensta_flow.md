# The OpenSTA flow: inputs and reports

OpenSTA is the open-source static timer used inside OpenROAD and in many research flows. Its inputs are exactly the abstractions this domain studies: Liberty for cells, a Verilog netlist for connectivity, SDC for constraints and SPEF for parasitics. Reading its flow end to end shows what each file contributes and what each abstracts away, and it is the reference implementation to read when building a timer.

## The four inputs

read_liberty loads cell libraries: timing arcs with their delay and slew tables, pin capacitances, setup and hold check tables, and timing senses. read_verilog loads a gate-level netlist, and link_design binds its instances to library cells and elaborates the top module into the design the timer analyzes. read_sdc applies the constraints: clocks, I/O delays, uncertainty and exceptions. read_spef loads the extracted RC network of every routed net. Without parasitics, the timer can only estimate wires; with them, net delays and driver loads reflect the actual layout. Each input preserves something and drops something: Liberty keeps tables, not transistor waveforms; SPEF keeps resistors and capacitors, not geometry.

## Running the analysis and reading results

After loading, clocks are ideal unless set_propagated_clock is applied, which a routed design with a clock tree needs. report_checks prints timing paths with their arrival, required time and slack: -path_delay max selects setup checks, min selects hold, and options choose how many paths per group or endpoint and which fields to show. report_wns, report_tns and report_worst_slack print summary numbers for scripts. Net delays come from a delay calculator that reduces each net's parasitics and models the driver from its Liberty tables, with effective capacitance to account for resistive shielding. The reported numbers are therefore the composition of all four inputs, and an error in any one appears as a wrong slack.

## A timer embedded in a flow

In OpenROAD, OpenSTA runs in the same process as the placer, the resizer and the clock-tree and routing tools, which query slack in memory. After a resize or buffer insertion, the timer invalidates only the affected part of the graph and recomputes on the next query, so optimizers can call it after every move. That design choice, an incremental timer behind an API rather than a program that writes files, is what makes timing-driven optimization possible in an open flow. It is also the natural place for the seams of an agent-first timer: derating, path reporting and fast modes can be changed without touching the propagation core.

## Going deeper

The OpenSTA documentation and source code; Ajayi et al., "OpenROAD: Toward a self-driving, open-source digital layout implementation tool chain" (GOMACTech, 2019).
