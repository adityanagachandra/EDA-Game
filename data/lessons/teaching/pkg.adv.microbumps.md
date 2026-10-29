# Microbumps

Microbumps are small solder-capped copper pillars, at pitches of roughly 25 to 55 um, that join dies to interposers, bridges or other dies. They are the workhorse of 2.5D and of TSV-based 3D stacks. Their density sets die-to-die bandwidth, their current limit sets how many must carry power, and their resistance enters power-integrity analysis.

## Bump density from pitch

On a square grid at pitch $p$, the bump density is $(1000/p)^2$ per square millimetre with $p$ in micrometres. At 40 um this is 625 per mm^2; at 55 um about 331; at 25 um, 1600. Compare 44 per mm^2 for a 150 um flip-chip bump. Density grows with the inverse square of the pitch, so every pitch step matters twice. Microbumps cannot shrink indefinitely: the solder volume falls with the cube of the size, so small bumps become intermetallic-dominated, wetting and bridging control get harder, and alignment tolerance becomes a large fraction of the pad. Below about 10 um pitch, bumpless hybrid bonding takes over.

## Power bumps from current limits

Each bump carries a maximum current, set by electromigration in the solder and the intermetallic layers at operating temperature. If a die draws $I$ amps through bumps each limited to $I_{max}$, it needs at least $\lceil I / I_{max} \rceil$ power bumps, and as many ground bumps. A die drawing 60 A with 0.1 A per bump needs 600 power and 600 ground bumps, before any margin for uneven current sharing. Current does not share evenly: bumps above high-activity regions carry more, so analysis must solve the power network to find the worst bump rather than divide the total. The power bumps then compete with signal bumps for the same area, which is why bump maps are planned early.

## Bump resistance

A microbump joint is roughly a cylinder of solder and intermetallic with height $h$ and diameter $d$, so $R = \rho h/(\pi d^2/4)$. Tin-silver solder has resistivity around 12 to 15 micro-ohm cm, close to ten times copper. A joint 15 um tall and 20 um wide at 13 micro-ohm cm, which is 0.13 ohm um, has $R = 0.13 \times 15/314.2 = 6.2$ milliohm. Copper pillars carry most of the height with low resistance, so the solder cap dominates the series resistance. Thousands of such joints in parallel give a small total, but the per-bump value matters for the current sharing that sets the electromigration limit.

## Going deeper

Lau, "Semiconductor Advanced Packaging" (2021), chapter on microbumps; Tu, "Electronic Thin-Film Reliability" (2010), on electromigration in solder.
