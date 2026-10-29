# Standards bodies and open formats

EDA works only because tools from different vendors can read each other's files and agree on what a language means. That agreement comes from standards: hardware description languages, constraint and parasitic formats, test access protocols, board data formats, mask layout formats. A handful of bodies own them, and a tool designer implements against their documents.

## IEEE and Accellera

The IEEE standardizes the hardware languages and many test standards: Verilog as IEEE 1364, SystemVerilog as IEEE 1800, VHDL as IEEE 1076, SystemC as IEEE 1666, the UVM verification library as IEEE 1800.2, the Unified Power Format for power intent as IEEE 1801, boundary scan (JTAG) as IEEE 1149.1 and embedded core test as IEEE 1500. Accellera Systems Initiative is an industry consortium that develops many of these before handing them to the IEEE for formal standardization: SystemVerilog, UVM, SystemC and the Portable Stimulus Standard began there. Because the languages are standardized, a simulator, a synthesizer and a formal tool from three different vendors can all read the same SystemVerilog source.

## Si2, SEMI and IPC

Si2 (Silicon Integration Initiative) is an industry consortium that maintains interoperability formats and databases: it distributes the LEF/DEF reference and stewards OpenAccess, a design database API used by custom and analog tools. SEMI is the global association of the semiconductor manufacturing supply chain; its standards cover equipment, materials and data, including the OASIS layout format (SEMI P39). IPC sets standards for printed boards and electronics assembly: design guidelines (IPC-2221), assembly acceptability (IPC-A-610) and the IPC-2581 board data exchange format. Their standards matter to tool designers whenever a tool reads or writes data that crosses a company boundary, from mask shop to board assembler.

## De facto standards

Several of the most important formats were created by one vendor and opened later. Liberty and SDC came from Synopsys and are published for anyone to implement. LEF and DEF came from Cadence. GDSII came from Calma. Gerber, the classic board fabrication format, is maintained by Ucamco. De facto formats evolve faster but can carry vendor-specific extensions, which is why a parser that handles the published reference may still fail on a real file. A tool designer should write parsers against the published specification and test them against files produced by the major tools.

## Going deeper

The IEEE 1800 SystemVerilog standard; the Si2 LEF/DEF reference; SEMI P39 for OASIS.
