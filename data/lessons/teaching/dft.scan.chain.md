# Scan chains

Sequential ATPG is hard because flip-flop states can only be reached through many clock cycles. Scan design removes that difficulty: in test mode every flip-flop becomes a stage of a shift register, so the tester can load any state and read any captured state directly. The circuit then looks combinational to ATPG, which is why full scan is the foundation of every modern digital test flow, and why scan shift dominates test time.

## The mux-D scan cell

A **mux-D scan flip-flop** puts a 2-to-1 multiplexer in front of an ordinary D flip-flop. The scan-enable signal SE selects between the functional data input D (SE = 0) and the scan input SI (SE = 1). Stitching each flop's output to the next flop's SI forms a **scan chain** from a scan-in pin to a scan-out pin. With SE = 1, each clock shifts the chain by one position; with SE = 0, the flops capture their functional next state. The cost is one mux per flop, a little delay on the D path, and the routing of the chain; alternatives such as clocked-scan and level-sensitive scan design (LSSD) trade those costs differently.

## A scan test cycle

Each pattern runs in three phases: **shift** the stimulus into the chain with SE = 1, which takes as many cycles as the longest chain; apply primary inputs and **capture** the circuit's response with one functional clock at SE = 0; then **shift out** the response while shifting in the next pattern. Chains shift in parallel, so a design with 8333 scan flops in 4 balanced chains needs $\lceil 8333/4 \rceil = 2084$ shift cycles per load. The combinational logic between flops is tested exactly like a combinational circuit whose pseudo-inputs are flop outputs and pseudo-outputs are flop inputs.

## Test time

Because unload of one response overlaps the load of the next, $p$ patterns on chains of length $L$ take
$$p(L + 1) + L$$
cycles: $L$ shift cycles plus one capture per pattern, plus a final unload. With 1563 patterns and $L = 703$: $1563 \times 704 + 703 = 1{,}101{,}055$ cycles. At a 10 MHz shift clock, 935 patterns on $L = 1426$ take about 134 ms. Shift frequency is limited by power (every flop toggles during shift) and by timing on the scan path, so the main levers on test time are shorter chains, which need more scan pins or compression, and fewer patterns.

## Going deeper

Bushnell and Agrawal (2000), chapter 14; Eichelberger and Williams, "A logic design structure for LSI testability," Design Automation Conference, 1977.
