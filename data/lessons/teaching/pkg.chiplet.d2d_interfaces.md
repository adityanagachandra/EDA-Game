# Die-to-die interface choices

A die-to-die interface is the circuit and protocol that moves data between chiplets in one package. The design space runs from wide, slow, parallel interfaces over very short dense wires to narrow, fast serial links over longer package traces. The right choice depends on reach, the package technology, bandwidth per millimetre of die edge, energy per bit and latency, and it constrains the floorplan and the package.

## Parallel and serial interfaces

Parallel interfaces send many single-ended signals at moderate rates, forwarded with a clock, over wires a few millimetres long on an interposer, bridge or fine-pitch substrate. Because the wires are short and dense, they need no equalization or clock recovery, so they achieve low energy and latency. Serial SerDes links send fewer differential lanes at tens of gigabits per second each, with equalization and clock-data recovery, over longer reaches on ordinary organic substrates. They need fewer bumps but spend more energy per bit and add latency for serialization and recovery. The choice follows the package: dense advanced packaging favours wide parallel interfaces; a cheap laminate with coarse bumps favours serial links.

## Shoreline and reach

Die-to-die interfaces live on die edges, and edge length, the shoreline, is scarce: a die has four edges, and memory, I/O and other chiplets compete for them. The relevant figure of merit is bandwidth per millimetre of edge, which multiplies bump density, signal share and data rate. Reach matters because loss and crosstalk grow with wire length: a short unterminated parallel link stops working cleanly beyond a few millimetres, when transmission-line effects and RC attenuation take over and termination becomes necessary, which costs power. A floorplanner must therefore place communicating chiplets with facing edges close together, and a package designer must keep those routes short and on the finest layers.

## Protocol, latency and test

Above the electrical layer, an interface carries a protocol. Some interfaces map existing protocols such as PCIe or CXL onto the die-to-die link so that chiplets from different vendors interoperate; others carry a raw streaming protocol for one vendor's chiplets. Latency comes from serialization, clock-domain crossing, retry and protocol layers, and it matters for cache-coherent traffic far more than bandwidth does. Interfaces also need lane repair, because a single bad microbump among thousands should not scrap the package: spare lanes and remapping logic replace failed lanes found at assembly test. A tool designer must model these redundancy schemes, because they change the bump map and the netlist that verification sees.

## Going deeper

Lau, "Chiplet Design and Heterogeneous Integration Packaging" (2023); the Open Compute Project Bunch of Wires specification; the UCIe specification overview.
