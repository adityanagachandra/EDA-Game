# Hybrid bonding

Hybrid bonding joins two dies or wafers face to face with no solder: polished dielectric surfaces bond directly, and recessed copper pads expand during an anneal to form metal joints. Because there is no solder to wet or bridge, the pitch can drop below 10 um and toward 1 um in wafer-to-wafer bonding. The result is interconnect density orders of magnitude above microbumps, limited mainly by alignment and surface flatness.

## Connection density

Density at pitch $p$ is $(1000/p)^2$ per mm^2. At 9 um pitch that is about 12346 per mm^2; at 3 um, 111111; at 1 um, a million. A 50 mm^2 overlap at 3 um pitch offers about 5.6 million connections. Such density changes the design question from whether to cross a die boundary to how to use the crossings: memory stacked on logic with direct per-bank connections, or logic split between dies at the block or even cell level. The bonding pads also have far lower capacitance than microbumps, so each crossing costs little energy. The fabrication cost is a near-perfectly flat surface from chemical-mechanical polishing and particle-free bonding.

## Gain over microbumps

Because density scales as $1/p^2$, the gain from microbumps at pitch $p_{mb}$ to hybrid bonding at $p_{hb}$ is $(p_{mb}/p_{hb})^2$. From 40 um to 9 um the gain is $(40/9)^2 = 19.8$; from 36 to 3 um, 144. A designer should translate this into bandwidth only with care: bandwidth equals connections times data rate, and dense hybrid-bond interfaces often run each wire slower, because simple low-swing drivers save energy. Still, area per connection falls so far that the die-to-die interface can shrink from a large edge PHY to a small fraction of the die, which is what enables fine-grained 3D partitioning.

## Overlay and pad overlap

At fine pitch, alignment error is a large fraction of the pad. Two square pads of side $a$ offset by $e_x$ and $e_y$ overlap over $(a - e_x)(a - e_y)$, a fraction $(a-e_x)(a-e_y)/a^2$ of the pad. With 2 um pads and a 0.3 um by 0.2 um offset, overlap is $1.7 \times 1.8/4 = 76.5\%$. Contact resistance rises as overlap falls, and if the pad offset brings copper close to a neighbour's pad, the dielectric gap shrinks and leakage or shorts follow. Pitch therefore scales with alignment capability. Die-to-wafer bonding has looser alignment than wafer-to-wafer bonding, which is why the finest pitches are reached wafer to wafer, and why design rules couple pad size, pitch and the bonder's overlay specification.

## Going deeper

Lau, "Semiconductor Advanced Packaging" (2021), chapter on hybrid bonding; Enquist and others, "Direct bond interconnect", IEEE 3DIC conference (2009).
