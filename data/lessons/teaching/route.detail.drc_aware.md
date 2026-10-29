# DRC-aware detailed routing

A detailed router's output must satisfy every rule in the deck: spacing, width, end-of-line, minimum area, via enclosure and more. DRC-aware routing builds those rules into the search, as costs on candidate moves and as blocked resources, so that the router produces clean wires rather than finding violations afterwards. This is the core of what makes a detailed router hard to write.

## Pricing violations

During search a detailed router lets a route pass through states that would create a violation, but charges for it: a route's cost is its length plus a cost per via plus a large penalty per design-rule violation it would create. Example: 18 steps, 5 vias at 5 each and one violation at 100 cost $18 + 25 + 100 = 143$.

Why not forbid violations outright? Because early in routing a short route with a violation may be the only way through, and the violation can be fixed later by moving a different net. Soft penalties keep the search complete while steering it toward legality, and the penalty grows during repair so that the final result is clean.

## Choosing among candidates

Given candidates for one connection, the router keeps the cheapest under its costs. Example: route A has 50 steps, 3 vias and no violation; B has 31 steps, 4 vias and 2 violations; C has 32 steps, 2 vias and none. With via cost 6 and violation cost 50, the costs are 68, 155 and 44, so C wins: clean and short.

The relative weights encode policy. A violation penalty much larger than any plausible detour makes the router prefer clean routes almost always. A penalty too small lets it leave violations everywhere, and one too large freezes nets that should move. Production routers tune these weights per stage and grow them across repair iterations.

## Rules as blocked tracks

Some rules are cheaper to model as blocked resources than to check per move. A wide wire needs more spacing to its neighbours than a normal wire; a spacing table gives the required value by width. With track pitch $p$, normal width $w$ and a wide wire of width $W$, the edge gap to a wire on the $k$-th track over is $kp - W/2 - w/2$. Every track where that gap is below the wide-wire spacing is blocked.

Example: $p = 48$ nm, $w = 24$ nm, $W = 96$ nm and a required spacing of 94 nm. The gaps to tracks 1 to 4 are -12, 36, 84 and 132 nm (a negative gap means the wires would overlap), so tracks 1 to 3 on each side are blocked and track 4 is the first usable one. Precomputing such blockages turns rule checking into graph editing, which keeps the search kernel fast.

## Going deeper

The TritonRoute paper (IEEE TCAD, 2021); the OpenROAD drt documentation; the LEF reference for SPACINGTABLE and related rules.
