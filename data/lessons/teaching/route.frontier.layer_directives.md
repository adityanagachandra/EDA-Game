# Layer directives as Tier 1 changes

Many routing requests are about layers: route this block only on M2 to M6, keep clocks above M5, reserve capacity on M2 for pin access, give these nets double width. These are layer directives, and a well-designed router makes them Tier 1 configuration that the core already enforces. An agent should never write code for them, but it must still verify their consequences.

## Directives as configuration

A layer directive restricts or reshapes the resources the router may use, through settings the core already understands. Minimum and maximum routing layers bound the whole design. Per-net-class layer ranges restrict, for example, clock nets to M5 and above. Layer adjustments reduce a layer's capacity by a percentage. Non-default rules assign wider wires and larger spacing to named nets. Routing blockages forbid layers over a region.

Each is a typed setting, validated when it is read and enforced by the core as legality. That makes them Tier 1: no recompile, no new code, and the same behaviour every time the setting is used. OpenROAD, for instance, exposes routing layer ranges and per-layer adjustments as commands in its flow scripts.

## Tier 1 versus Tier 2

A request is Tier 1 when an existing setting expresses it. Forbid routing below M7 over an analog region: a set of routing blockages, Tier 1. Make vias on timing-critical nets more expensive only when slack is below a threshold: a conditional cost that no setting expresses, so it needs a cost plug-in, Tier 2. Allow wires off the track grid: a change to legality, Tier 3.

Directives must be validated, not silently ignored. A net restricted to M2 through M4 whose pin sits only on M6 cannot be routed; the router must report the conflict rather than drop the pin or quietly use M6. Common, well-defined requests belong in Tier 1 precisely because a typed setting is safer and faster than new code, and because an agent can check it mechanically.

## Directives still need verification

Configuration changes behaviour. Lowering the maximum routing layer from M8 to M6 removes all capacity on M7 and M8, so overflow and detours rise and long critical nets lose their low-resistance layers. Restricting all nets to M2 and M3 removes most of the stack's capacity and will overflow heavily. So an agent applying a directive still runs the golden designs: DRC-clean, overflow within limits, timing within tolerance.

Recurring Tier 2 requests are a signal to grow Tier 1. If teams keep asking for different via costs per layer, the engine should expose a per-layer via-cost table. Legality on the chosen layers, direction, spacing and via rules, remains the core's job whatever the directive says.

## Going deeper

The OpenROAD grt documentation for routing layers and layer adjustments.
