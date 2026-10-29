# Seams and extension points

A seam is the place in an engine where the core stops and customizable behavior begins. It is implemented as an extension point: a narrow, stable interface that the core calls and that plug-ins implement. Seams are what make Tier 2 possible. Where they sit, and how they are built, decides whether customer requests become small plug-ins or repeated core edits.

## Abstract base plus factory

The standard construction is an abstract base class (or an equivalent table of function pointers in C) that declares the few operations the core needs, plus a factory that creates an implementation by name from configuration. The core holds a pointer to the base and never knows which implementation it is calling. A router's cost seam might declare one operation, the cost of a candidate edge given the net and the current congestion. The factory maps the configured name, such as `negotiated` or `analog_guard`, to a class. Adding a new behavior means adding a new class and registering it; nothing in the core changes, and the choice between implementations is a Tier 1 configuration.

## The derate provider: a seam in a timer

Timing derating shows a seam at its best. The core, graph propagation and the delay-model API, needs one thing when it computes an arc delay: the derate for this arc in this context. A `DerateProvider` interface answers that question. It receives the context (arc, cell and Vt class, early or late, rise or fall, logic depth, distance, slew, corner, clock or data, path group) and returns a factor for flat OCV or AOCV, or a mean and sigma adjustment for POCV. Skills implement it: flat OCV, AOCV, POCV and a custom margin layer that stacks on top. Customer differences then become configuration (which provider, which tables, how many sigma) or a small new provider.

## Rules for building a seam

Interfaces at seams are narrow and stable: few methods, simple types, and an extension ABI that is frozen so plug-ins built against it keep working. The hot kernel stays in the core and the plug-in boundary is crossed coarsely, per batch of arcs or nets rather than per element, so the compiler can still inline and vectorize. Data crosses as plain structures (the design state, a batch of contexts), with no hidden or mutable global state and no coupling between skills. Plug-ins must be deterministic and seedable. A seam that violates these rules works in a demo and fails in production, through slowdowns, nondeterministic results or plug-ins that break with every core release.

## Going deeper

Gamma et al., "Design Patterns" (1994), on the strategy and factory patterns; the OpenSTA source on how derates enter delay calculation.
