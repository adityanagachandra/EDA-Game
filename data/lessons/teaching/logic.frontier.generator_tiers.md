# Change tiers for parameterized RTL generators

Much RTL is not written by hand but generated: parameterized modules, Chisel or other hardware-construction-language generators, memory compilers, arithmetic generators. A generator is an engine in the four-layer sense, with a core (the elaboration framework and its invariants), policies (architectures), parameters (widths, depths) and hooks. Placing each change request at the right tier is what keeps agent edits safe.

## Tier 1: parameters

Most requests to a generator are parameter changes: data width, FIFO depth, number of pipeline stages, number of read ports, column-mux factor of a memory, an option to register the outputs. These need no code change: the generator is re-run with new values, and the existing golden tests (parameterized to match) re-verify the result.

A well-designed generator exposes these explicitly and validates them (depth a power of two, width at least 1) so that a configuration change cannot produce an illegal design. An agent that edits generator source to change a width is working at the wrong tier.

## Tier 2: a new policy at a seam

Adding a new adder architecture to an arithmetic generator, a new FSM encoding to a controller generator, or a new replacement policy to a cache generator is a plug-in change. The seam is an interface: for an adder generator, "given $n$ and the $(g, p)$ signals, return a prefix network", with the generate-propagate logic and the sum XORs staying in the core. The new policy (say a Han-Carlson network) implements only the prefix network.

The proof is differential: the new architecture must be combinationally equivalent to the existing ones at every width tested, and its QoR report shows where it sits on the depth-area trade-off.

## Tier 3: the core

Core changes alter the generator's invariants or framework: how clocks and resets are emitted, how modules are named and connected, the intermediate representation, the semantics of a parameter. They can affect every design the generator has ever produced, so they need a human-approved diff and full regression across all configurations.

Some requests look small but are core: "make all generated flip-flops use asynchronous reset" changes reset semantics everywhere and interacts with timing and CDC checks. The tier is set by what the change can affect, not by how many lines it touches.

## Going deeper

Bachrach et al., "Chisel: Constructing Hardware in a Scala Embedded Language" (DAC, 2012); Shacham et al., "Rethinking Digital Design: Why Design Must Change" (IEEE Micro, 2010).
