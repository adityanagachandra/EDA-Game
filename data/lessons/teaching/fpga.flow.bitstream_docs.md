# Documenting bitstreams: IceStorm, Trellis, X-Ray

Vendors do not publish the bit-level format of their FPGAs. Projects such as IceStorm (Lattice iCE40), Trellis (Lattice ECP5) and X-Ray (AMD Xilinx 7-series) recovered it, producing databases that map every configurable feature to its bits. Those databases are what make open place and route and open bitstream generation possible, and they are a lesson in how to recover a hidden mapping.

## What a database records

A bitstream database lists, for every tile type, each configurable feature (a LUT's truth-table bits, a flip-flop's mode, a routing multiplexer's input, an I/O standard) and the configuration bits that enable it, as frame and bit offsets relative to the tile's base address. Tiles of the same type share one entry, and each tile instance records its base, which keeps the database compact. With it, a tool can write a bitstream for any design by setting the bits of the features the design uses, and can read a bitstream back by asking which features its bits enable. IceStorm covers iCE40, Trellis covers ECP5, and X-Ray covers 7-series parts such as Artix-7.

## Fuzzing

The mapping is found by fuzzing. Generate many small designs that differ in a controlled way, such as whether one routing switch is used, build each with the vendor tools, and compare the bitstreams. Bits that change exactly when the feature changes encode it. Vendor tools make side changes, so a single pair of designs is not enough: a toggle that flips eight bits, three of which also flip in unrelated runs, leaves five candidates, and more varied designs are needed to isolate the real ones. Fuzzers therefore build hundreds or thousands of designs per feature class and solve for the bits that are consistent across all of them. The process is statistical and incremental, and coverage grows feature class by feature class.

## Limits and uses

A database is only as good as its coverage and correctness. A wrong entry produces a bitstream that loads and misbehaves, so databases are validated against vendor-built bitstreams; rarely used tile features, such as unusual I/O modes, are where gaps tend to remain. Beyond place and route, a documented bitstream allows readback into a netlist (icebox_vlog does this for iCE40), which lets the configured hardware be compared with the intended design. For a new family with no database, the first part of an open flow that can be built is synthesis to the documented primitives; packing, placement, routing and bitstream generation all need the device database.

## Going deeper

The Project IceStorm, Project Trellis and Project X-Ray documentation; the F4PGA documentation for the VPR-based 7-series flow.
