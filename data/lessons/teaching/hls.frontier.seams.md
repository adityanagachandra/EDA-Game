# Seams and change tiers in an HLS engine

An HLS engine built for agents separates what must never change casually, the rules that make every schedule and binding correct, from what customers and agents tune: priorities, cost functions, directive values and observation points. Applying the four-layer engine model and the change tiers to HLS tells a tool designer where each kind of request lands and what proof it needs before it ships.

## The four layers of an HLS engine

The immutable core holds correctness: dependence analysis, including loop-carried and memory dependences, the precedence and modulo reservation checks, and the rule that every binding is a proper coloring of its conflict graph. Changing it needs human approval and full regression. The strategy layer holds policies that choose among legal results: list-scheduling priorities, force-directed or SDC objectives, binding cost functions that weigh multiplexers against units. The parameter layer holds values: target clock, II targets, unroll and partition factors, resource limits, all changeable without recompiling. The observation layer exposes live state: per-loop achieved II, the limiting resource or recurrence, the chosen schedule, register and multiplexer counts. A request to log each scheduling decision is a hook; a request to favour store-feeding operations is policy; raising an unroll cap is a parameter.

## Which tier a request lands in

Requests are handled at the smallest safe tier. Pipelining every loop at II = 2 with unrolling capped at 4 is Tier 1: directive and parameter values the engine already reads. A binding cost that penalizes sharing a multiplier between operands from different memory banks is Tier 2: new behavior, but it only re-ranks legal bindings, so it belongs at the binder's cost-function seam. Letting the scheduler ignore a loop-carried memory dependence to reach a lower II is Tier 3, if it is allowed at all, because it changes what function the hardware computes. Customer-specific operator latencies belong behind an operator-characterization interface that the scheduler, chaining and binder query per operation, width and target, so that a new device or library becomes a Tier 1 table or a Tier 2 provider, never a core edit.

## Designing the seams and proving changes

A good seam is narrow and stable: an abstract interface with a factory, called coarsely, per batch of operations or per loop rather than once per operation inside the binder's innermost loop, so the hot kernel stays inlinable. Plug-ins receive plain data, keep no hidden global state and are deterministic and seedable; a plug-in that reads state left over from a previous run makes differential tests irreproducible. Every change is proven against a golden kernel suite: all kernels still pass C/RTL co-simulation, and a QoR differential report shows that latency, II, registers and area moved only where the change intended. A new priority that shortens one kernel and silently lengthens three others fails that review, even though every output is still correct.

## Going deeper

For the engine model, the agent_first_eda lessons on the four-layer model, change tiers and seams; for HLS internals, Coussy and Morawiec (eds.), "High-Level Synthesis: From Algorithm to Digital Circuit" (2008).
