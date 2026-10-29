# The four layers in code

The four-layer engine model (immutable core, strategy or policy, parameters, observation hooks) is only useful if it is visible in the source tree. This concept turns it into code: which classes and files of an OpenROAD-like engine belong to each layer, how the seams between them are built, and where an incoming change request lands. An agent editing an engine needs that map to make the smallest safe change.

## Mapping code to the four layers

In a placer, the core is the code that every result must pass through: the legalizer that guarantees cells sit on legal sites without overlap, the row and site model, and the density and wirelength kernels. Policies are the classes that decide how: a cost function, a move set, an ordering for rip-up in a router, each implementing a common interface and selected by name. Parameters are values read at run time from configuration, such as target density, net weights and iteration limits; changing them needs no recompile. Observation hooks are callbacks or query points where an agent reads live state (overflow per iteration, WNS after an incremental update) and can inject a decision, such as stopping early. A quick test for a piece of code: if a wrong change breaks an invariant, it is core; if it only changes quality, it is policy or parameter.

## Narrow interfaces at the seams

A seam is built as an abstract interface plus a factory that creates implementations by name, so the core depends only on the interface and never on any particular plug-in. Keep the interface narrow: a cost plug-in should receive arrays of the values it needs, such as cell widths, positions and net pins, not a pointer to the whole database, which would let it depend on and mutate anything. Plain data at the boundary is easy to version, test and run in parallel. Keep the hot numerical kernel inside the core and call plug-ins coarsely, once per batch of cells rather than once per cell, so the compiler can still inline and vectorize the inner loop. A seam that passes the world, or that is called a million times per pass, has the shape of an interface but not its benefits.

## Placing a change request

Given a request, ask what must change. "Raise the congestion weight from 0.3 to 0.5 on this design" changes a value: a parameter, no recompile. "Keep macros 20 microns from the die edge" adds behavior; if the constraint system has an extension point, it is a policy plug-in written against that interface, and the core is untouched. "Stop global placement when overflow falls below 10%" needs live state and a control: an observation hook. "Change what counts as a legal site" alters the core invariant, so it is a core edit that needs human approval and full regression. If requests of some common kind keep landing in the core, the seam is in the wrong place, and the fix is to move that variability behind a new extension point once, under review.

## Going deeper

The agent-first engine notes in this curriculum (agent_first_eda); the OpenROAD source tree as an example of tools as modules over one database; Ousterhout, "A Philosophy of Software Design", on deep modules.
