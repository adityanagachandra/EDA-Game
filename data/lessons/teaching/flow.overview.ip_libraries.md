# IP cores and standard-cell libraries as products

Almost nothing on a modern chip is designed from scratch. Standard cells, memories, I/O cells, PHYs and processor cores are bought or reused as intellectual-property (IP) products. A design team assembles them, and EDA tools consume the views each product ships with. A tool designer must know what those views are, because every engine reads them.

## Soft, firm and hard IP

Soft IP is delivered as synthesizable RTL: a processor core, a bus fabric, a controller. It is portable across processes but must be synthesized, placed and timed by the buyer. Hard IP is delivered as finished layout for one specific process: a SerDes PHY, a PLL, an SRAM macro, an analog block. It has proven timing and power but cannot be moved to another process without redesign. Firm IP sits between, for example a gate netlist with placement guidance. The more physical the delivery, the less freedom the integrator has and the less risk they carry.

## The standard-cell library

A standard-cell library is a set of small logic cells (inverters, NANDs, multiplexers, flip-flops, often several drive strengths and threshold-voltage flavours each) drawn to a common height so they abut in rows with shared power rails. It is a product with many views: Liberty files give timing, power and function per corner; LEF gives the cell outline, pin shapes and blockages for placement and routing; GDSII gives the full layout; SPICE netlists give the transistor circuit for LVS and characterization; Verilog models give simulation behaviour. A missing or inconsistent view breaks a specific stage: a Liberty file without a corner means that corner cannot be timed.

## Characterization and reuse economics

Library and IP vendors characterize each cell by running thousands of circuit simulations across input slew, output load, voltage and temperature, and tabulating the results into Liberty. That work is paid back by reuse: one library serves every design on the process. The same logic drives IP licensing: a team pays a licence fee and often a royalty per unit because designing and verifying a PCIe PHY or a CPU core itself would cost more and take longer. The integration risk moves to the interfaces, so IP ships with integration guides, constraints and verification models.

## Going deeper

Chinnery and Keutzer, "Closing the Gap Between ASIC and Custom" (2002), on how library quality limits ASIC performance.
