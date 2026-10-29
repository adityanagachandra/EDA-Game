# The commercial and open-source tool landscape

Engineers name tools by product, but a tool designer should name them by the problem class they solve: simulation, synthesis, place and route, static timing, extraction, physical verification, test, custom layout. Each major product is one vendor's answer to one class. Mapping the products onto the flow tells you which engines exist, which ones compete, and where open-source tools now offer a full alternative.

## Front-end tools

Logic simulators execute RTL and gate netlists: Synopsys VCS, Cadence Xcelium and Siemens Questa. Formal verification and equivalence checking have their own products, such as Cadence JasperGold and Conformal and Synopsys Formality. Logic synthesis turns RTL into gates: Synopsys Design Compiler and Cadence Genus. DFT tools insert scan and generate test patterns: Siemens Tessent, and test products from Synopsys and Cadence. FPGA flows are vendor-owned: AMD Vivado and Altera Quartus. Most of these products compete in pairs or triples across the three large vendors, so a team usually picks one vendor's simulator, synthesizer and equivalence checker to keep its methodology in one ecosystem. Each class corresponds to a different engine: an event-driven or compiled simulator, a SAT- or BDD-based prover, a logic optimizer and mapper.

## Back-end and signoff tools

Place-and-route systems take a netlist to routed layout: Cadence Innovus and Synopsys IC Compiler II, while Synopsys Fusion Compiler combines synthesis and place and route in one database. Static timing signoff is dominated by Synopsys PrimeTime, with Cadence Tempus as its competitor. Parasitic extraction includes Synopsys StarRC and Cadence Quantus. Physical verification (DRC and LVS) is led by Siemens Calibre, with Synopsys IC Validator and Cadence Pegasus. Power-integrity analysis includes Ansys RedHawk and Cadence Voltus. Custom and analog design runs in Cadence Virtuoso with circuit simulators such as Cadence Spectre and Synopsys HSPICE.

## The open-source landscape

An open flow now exists for every major class. Yosys and ABC handle synthesis and logic optimization; OpenROAD covers floorplanning, placement, CTS and routing; OpenSTA does static timing; KLayout and Magic handle layout viewing, editing and DRC; Netgen does LVS; ngspice and Xyce simulate circuits; Verilator and Icarus Verilog simulate RTL; SymbiYosys runs formal checks; nextpnr and VPR place and route FPGAs. Open PDKs let these tools run real tape-outs at mature nodes. For a tool designer, open tools are the place to read real engine code and test changes end to end.

## Going deeper

The OpenROAD project documentation and the Yosys manual; vendor product documentation for the commercial tools.
