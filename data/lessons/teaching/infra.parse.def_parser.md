# Parsing LEF and DEF

LEF and DEF are the physical interchange formats of digital implementation. LEF describes the technology and the abstract of each cell; DEF describes one design's floorplan, instances, pins, nets and routing. A design database such as OpenDB follows their structure closely, so reading them is how a physical tool comes to life. A reader must map each section onto database objects and resolve every cross-reference between the two files.

## Technology LEF and cell LEF

LEF comes in two parts. Technology LEF defines units (`DATABASE MICRONS`), the manufacturing grid, routing and cut layers (with DIRECTION, PITCH, WIDTH and SPACING), via definitions and placement SITEs. Cell LEF defines one MACRO per library cell: its CLASS (CORE, PAD, BLOCK...), its SIZE w BY h in microns, the SITE it fits, and its PINs with DIRECTION, USE and PORT shapes on layers. OBS lists internal metal that routes must avoid. LEF is deliberately an abstract: it gives a router the pin access shapes and blockages but not the transistors inside, and it holds no timing. Pin capacitance and delay come from Liberty, which is why a timer needs both files.

## DEF sections and what they build

A DEF begins with VERSION, DESIGN and `UNITS DISTANCE MICRONS n`, then DIEAREA. ROW statements create the placement rows from the LEF site, TRACKS the legal routing tracks per layer and GCELLGRID the global-routing tiles. COMPONENTS creates instances: `- u7 NAND2_X1 + PLACED ( 3800 5600 ) FS ;` gives the master, a status (PLACED, FIXED, UNPLACED) and an orientation (N, S, FN, FS and so on). PINS creates the block's ports with their shapes. NETS gives connectivity as lists of (instance pin) pairs, plus routing once routed; SPECIALNETS carries power straps and rings; BLOCKAGES carries placement and routing blockages. Each section maps onto a database object kind, which is why a DEF reader is a sequence of section handlers.

## Reading order and cross-references

DEF refers to LEF by name: components name masters, ROWs name sites, routing names layers and vias. So LEF must be read first, and every name must resolve. When a NETS entry says `( u7 A3 )` but the master has only A1, A2 and ZN, the reader should report the error with a line number and keep going, never invent a pin. Units must also agree: the DEF resolution may not be finer than the LEF `DATABASE MICRONS`, so a 2000 LEF with a 1000 DEF is fine. Scale matters too. A large DEF can be tens of gigabytes, mostly NETS. The Si2 reference parsers, which OpenDB builds on, are callback-driven: the parser hands each component or net to the database as it is read, instead of building an intermediate tree.

## Going deeper

The LEF/DEF Language Reference (distributed with the Si2 LEF/DEF parsers); the OpenDB documentation in OpenROAD; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 1.
