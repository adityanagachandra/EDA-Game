# The OpenROAD tool suite

OpenROAD is one application that links a set of physical-design tools onto one OpenDB database, driven from a Tcl shell or Python. Each tool has a short name and one job: initialize the floorplan, place pins, build the power grid, place cells, repair timing, build the clock tree, route, extract and time. Knowing which tool owns which job, and what it leaves in the database for the next, is the map a tool builder needs before modifying any of them.

## Floorplan and placement tools

ifp initializes the floorplan: it sets die and core areas, either explicitly or from a target utilization, and builds rows of sites from the library's site definition. ppl places top-level IO pins on the die boundary, usually to shorten wires to the cells they connect to. tapcell inserts well taps and endcaps that the process requires, and pdn builds the power distribution network of rings, straps and rails. gpl is the global placer, from the RePlAce line of electrostatics-based placers: it treats cells as charges, solves for a density potential, and moves cells with a Nesterov-type method toward low wirelength and uniform density. Its result still overlaps. dpl, the detailed placer, snaps every cell to a legal site in a row without overlap and makes local improvements. Each tool reads and writes the same database, so the next tool simply sees the updated objects.

## Optimization, clock and routing tools

rsz, the resizer, runs repair design (fixing maximum slew, capacitance and fanout by buffering and sizing) and repair timing (sizing and buffering to fix setup and hold). cts, from TritonCTS, builds a buffered clock tree toward balanced latency at the flip-flop clock pins. grt, the global router from the FastRoute line, divides the die into gcells, assigns each net a path through them and through layers, and reports congestion as demand versus capacity per gcell edge. It writes route guides. drt, from TritonRoute, turns guides into exact wires and vias that obey the design rules, iterating with rip-up and reroute and reporting any violations it cannot fix. Order matters: grt's guides confine drt's search, and the clock tree must exist before signal routing claims the tracks near the flip-flops.

## Extraction, timing and the shared database

OpenRCX extracts resistance and capacitance from the routed geometry using calibrated rules and writes them as SPEF or loads them for timing. OpenSTA is the static timer: it reads Liberty and SDC and times the design using estimated or extracted parasitics. Because all tools live in one process on one database, a tool can call another without writing a file. The resizer asks OpenSTA for slack after every batch of edits; a timing-driven gpl run asks for net criticality between iterations and reweights nets; global routing can feed estimated parasitics back to the timer. In a flow of separate point tools, each of those exchanges would be a DEF and SPEF write and read, minutes on a large design. This is why OpenROAD's architecture is a shared in-memory database with tools as modules, and why an incremental timer is part of it.

## Going deeper

The OpenROAD documentation for each tool (ifp, ppl, pdn, gpl, rsz, dpl, cts, grt, drt, OpenRCX, OpenSTA); Ajayi et al. on the OpenROAD project.
