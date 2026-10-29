# Bitstream-level verification as the golden check

On an FPGA, the bitstream is what the hardware actually runs, so it is the most authoritative artifact for checking a flow. Decoding a bitstream through its database and comparing the result with the intended design is the golden check for any change an agent makes to packing, placement, routing or bitstream generation. Knowing which check fits which change keeps agent experiments both fast and safe.

## Why the bits are the ground truth

Every earlier representation, the netlist, the placement, the routed database, is a tool's belief about the hardware. Errors in the bitstream writer or in the feature-to-bit database leave those beliefs intact while the bits are wrong, so simulation of the post-route netlist can pass while the device fails. With a documented database, open flows can decode a bitstream into a netlist of configured primitives and run equivalence checking against the synthesized design. A golden test for an engine change fixes benchmark designs, seeds and architecture, with known-good results that the new build must reproduce or match within tolerance, and a FASM or bitstream diff shows exactly which features changed.

## Matching the check to the change

Different changes need different contracts. A refactor of the bitstream writer that should change nothing is tested by byte-for-byte identity of bitstreams for a fixed set of routed designs. A placer or router change legitimately changes the bits, so bit identity is the wrong test: the contract is functional equivalence of the decoded bitstream, no illegal or shorted routes, and QoR at least as good. A router change that improves the timing report is not accepted on that report alone, because better timing from a wrong route is worthless. A feature diff limited to routing-switch features, with LUT and flip-flop features identical, is direct evidence that a change touched only routing.

## Golden checks in an agent loop

Decoding and formally checking every result is slower than reading a timing report, so an agent uses fast checks to explore and the bitstream-level check to gate promotion, like the accurate engines in a fast and accurate loop. Configurations that change how packing or routing is done must pass it before promotion; changes that leave the bits identical, such as report verbosity, can be shown so directly. A mismatch confined to one tile type while the logic matches points first at the database for that tile. And a reward for an agent that edits the router must be gated on legality and equivalence; otherwise lowering minimum channel width by producing wrong routes is the cheapest way to score.

## Going deeper

The Project IceStorm documentation for icebox_vlog; the Yosys documentation for equivalence checking.
