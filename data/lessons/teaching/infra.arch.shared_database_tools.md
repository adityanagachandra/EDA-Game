# One database, many engines

An EDA flow chains engines: floorplanner, placer, clock-tree builder, router, extractor, timer. They can be separate programs exchanging files, or modules of one program sharing one in-memory design database, as in OpenROAD. The choice shapes everything else: speed of optimization loops, how engines learn of each other's edits, and what can go wrong when they share state. Tool designers must know both architectures and their costs.

## Point tools that exchange files

In the classic architecture each step is a separate program. The placer writes DEF; the router reads it and writes routed DEF; the extractor writes SPEF; the timer reads Verilog, DEF, SPEF, Liberty and SDC. The formats are the contract, so any tool can be swapped for another vendor's. The cost is serialization at every boundary: a design with ten million instances may take minutes to write and parse, so a loop that wants timing feedback after every few thousand moves is impractical. A second cost is loss: anything the format cannot express, such as a tool-specific attribute, disappears in the round trip unless both sides agree on an extension.

## Engines on one in-memory database

OpenROAD links its engines into one application around OpenDB. Floorplanning, placement, resizing, CTS, routing, extraction and timing read and edit the same objects, so an engine's results are visible to the next immediately, with no file in between. This is what makes tight loops possible: the resizer swaps a master, the timer updates incrementally, and the resizer asks again, thousands of times. Each engine still keeps its own derived structures, such as the timer's graph with arrival times or the router's grid, built from the database. So sharing data creates a new duty: when one engine edits the database, every engine caching derived results must hear about it, typically through callbacks, and update incrementally.

## Consistency when engines share state

Shared mutable state brings classic hazards. If two engines edit concurrently without coordination, one can read a half-applied edit or corrupt a shared list, so tools allow one writer at a time, partition the design into regions, or use transactions. If the router deletes a wire while the timer holds a raw pointer to it, the timer crashes later; stable ids or handles, plus deletion callbacks, avoid that. Every path into the database must notify, including user Tcl edits; otherwise a manual change leaves the timer's cache stale and its next report wrong. Files do not disappear: shared-database tools still write DEF, Verilog and SPEF at stage boundaries as checkpoints and for hand-off to sign-off tools.

## Going deeper

Ajayi et al., "OpenROAD: Toward a Self-Driving, Open-Source Digital Layout Implementation Tool Chain" (GOMACTech 2019); the OpenROAD developer guide; the OpenDB documentation.
