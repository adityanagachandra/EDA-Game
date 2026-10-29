# Design databases

A design database is the in-memory model that every engine of a physical design tool reads and edits: placement, routing, timing, extraction and verification all work on the same objects. Older flows passed files between point tools, re-parsing the design at every step. A tool designer cares because the database fixes what each engine can see, how fast it can query, how changes propagate to other engines, and how large a design fits in memory.

## Why tools share one database

When each tool reads files from the previous one, every step pays to tokenize LEF, DEF, Verilog and Liberty, resolve names and build its own model. On a design with 10 million instances, parsing alone can take minutes per step, and each tool's private model can drift from the others. A shared database is built once and then edited in place: the placer moves instances, the resizer swaps masters, the router adds wires, and every later engine sees those edits immediately. The database becomes the single source of truth. That only works if engines edit the shared objects rather than keeping private copies, and if derived results (arrival times, congestion maps) are kept in the engines and refreshed when the design changes. File formats remain, but as interchange at the edges of the flow: reading the inputs, writing the outputs, and handing data to tools outside the database.

## What the database stores and links

A physical design database keeps three kinds of data apart. Technology data comes from technology LEF: layers, directions, pitches, widths, spacing rules, vias and sites. Library data describes master cells: size, pins and their shapes, obstructions. Design data describes this chip: instances, each pointing to its master, with location and orientation; nets; instance terminals joining an instance pin to a net; block ports; rows; and routed wires. The links run both ways. From an instance the tool reaches its terminals and from each terminal its net; from a net it reaches all its terminals and their instances. So the question "which nets touch U17?" is a walk over the four terminals of U17, not a scan of a million nets. Masters are stored once and shared: 40,000 NAND2 instances hold only their placement and a reference, while pin geometry lives once in the master.

## Derived results, change and persistence

Not everything an engine computes belongs in the database. Placement and wires are design data. Arrival times, slacks, congestion per gcell and wirelength gradients are derived: they can be rebuilt from the design, so they live inside the engine that owns them. That split forces a mechanism: when the placer moves 200 cells, the timer's cached arrivals near those cells become stale, so the database must notify interested engines of edits, and they mark the affected results dirty. Persistence is the other half. A binary checkpoint of the object tables reloads far faster than re-parsing DEF and keeps objects and properties that DEF cannot express. Two further design choices recur: references between objects are stable handles, so that a deleted object is detected rather than dereferenced, and names are interned once in a string table so that hot objects stay small.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 1; the OpenROAD documentation for OpenDB.
