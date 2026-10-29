# Wire and via resistance

A routed net is a chain of wire segments on several layers joined by vias. Its resistance is the series sum of every segment and every via level, with parallel via cuts dividing a via level's resistance. At advanced nodes the vias, not the wires, often dominate short nets. The extractor computes this sum from fractured geometry and the layer stack, and the result feeds every delay and IR-drop calculation downstream, so a misplaced via array or a wrong per-cut value propagates straight into timing.

## Resistance from dimensions

The resistance of one wire segment is $R = \rho L/(W T)$, with the effective resistivity $\rho$ of its layer. Keep all lengths in micrometres and $\rho$ in ohm-um, and the units cancel to ohms. A copper segment with $\rho = 0.025$ ohm-um, $W = 0.04$ um, $T = 0.05$ um and $L = 120$ um has $R = 0.025 \times 120 / (0.04 \times 0.05) = 1500$ ohm. That is large: minimum-width lower metal at advanced nodes reaches tens of ohms per micrometre, which is why long nets are promoted to thicker upper layers. The extractor reads $W$ from the drawn geometry, corrected for the process bias, and $T$ from the layer stack in the technology file; both shift with the extraction corner.

## Vias in series and in parallel

Each via level a net crosses adds the resistance of its cut, or of its cut array. Identical cuts in parallel divide: four cuts of 50 ohm give 12.5 ohm. A route that runs 190 squares on M2 at 1.5 ohm/sq (285 ohm), crosses one via level through four 50 ohm cuts (12.5 ohm), and runs 60 squares on M3 at 0.5 ohm/sq (30 ohm) totals $285 + 12.5 + 30 = 327.5$ ohm. A net that climbs from M1 to M5 crosses four via levels on the way up and four on the way down, so its via resistance is eight per-level values. Redundant vias, inserted by the router or a post-route step, both lower resistance and protect against a single failed cut, which is why via doubling is a yield and reliability step as much as a timing one.

## When vias dominate

On a short connection the vias can outweigh the wire. A net that climbs two single-cut via levels of 55 ohm each, runs 170 squares at 1 ohm/sq, and descends two levels has 170 ohm of wire and 220 ohm of vias: 56% of its resistance is in the vias. Two consequences follow for tool design. First, a router that minimizes wirelength alone can make a net slower by adding layer changes, so via count must be in the cost function. Second, the extractor must model each via level separately, with the cut count and the landing geometry, rather than lumping vias into a per-net constant. Via resistance also has its own corner variation, independent of metal width, which signoff corners must cover.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), the chapters on routing and timing; the IRDS interconnect roadmap tables for resistance trends.
