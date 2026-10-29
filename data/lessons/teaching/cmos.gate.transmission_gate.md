# Transmission gates and pass-transistor logic

A transmission gate is an nMOS and a pMOS in parallel, driven by complementary controls, that passes a signal in either direction without restoring it. It is the workhorse inside library multiplexers, XORs, latches and flip-flops. A tool designer meets it in characterization (internal nodes with no gain), in switch-level and static checks (degraded levels), and in timing (RC chains through switches).

## Strong and weak levels

A single pass transistor is a poor switch for one of the two levels. An nMOS with its gate at $V_{DD}$ passes a 0 perfectly, but when passing a 1 its source rises until $V_{GS}$ falls to $V_{tn}$ and then it turns off: the output settles at $V_{DD} - V_{tn}$. With $V_{DD} = 0.9$ V and $V_{tn} = 0.3$ V that is 0.6 V. Body effect makes it worse, because the rising source raises the threshold. A pMOS with its gate at 0 V passes a strong 1 but a weak 0, stopping at $|V_{tp}|$. Putting both in parallel fixes this: when one device starts to cut off, the other is still fully on, so a transmission gate passes both rails fully. A chain of nMOS pass devices with gates tied to $V_{DD}$ loses only one threshold, not one per stage; drops accumulate only when one pass device's output drives the next device's gate.

## Pass networks in library cells

Transmission gates build compact functions. A 2:1 multiplexer is two transmission gates sharing an output, one controlled by S and one by its complement: four transistors plus an inverter for the select. An XOR can be made from two transmission gates and inverters. Latches use one transmission gate to admit data and another to close a feedback loop. Because a transmission gate has no gain, its output is driven through whatever drives its input. Library cells therefore surround pass networks with restoring inverters at their boundaries, so that each input pin presents a fixed gate capacitance and each output has a defined drive. That buffering is what lets characterization reduce the cell to ordinary pin-to-pin arcs.

## Costs that tools must see

Three effects make pass logic tricky outside a cell. A chain of $n$ transmission gates is an RC ladder whose Elmore delay grows as $n^2$, so long chains need buffers. A weak level, such as $V_{DD} - V_t$ from an nMOS, feeding a static inverter leaves the inverter's pMOS slightly on, wasting static current and eroding the next stage's noise margin. And a transmission gate exposed at a cell pin couples the external driver to internal nodes, so the pin capacitance depends on the state of the cell, breaking the fixed-load abstraction a timer assumes. Static checkers flag weak levels and unbuffered pass gates, and a timer handles internal pass networks only through the characterized arcs of the cell that contains them.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), sections 1.4.6 and 9.2.5; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), section 6.3.
