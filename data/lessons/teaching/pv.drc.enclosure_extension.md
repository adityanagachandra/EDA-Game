# Enclosure and extension rules

Enclosure and extension rules relate two layers. They keep a via cut inside the metal it lands on despite misalignment, and keep a poly gate long enough past diffusion that the transistor does not leak at its ends. They are computed with Boolean operations and sizing, and modern processes make them direction-dependent, which turns a simple margin test into a small logical condition.

## Enclosure

The enclosure of a via by metal is the smallest distance from a via edge to the corresponding metal edge. For a via $[115, 185] \times [190, 260]$ inside metal $[100, 207] \times [179, 275]$ the margins are left 15, right 22, bottom 11 and top 15 nm, so the enclosure is 11 nm. A negative margin means the via sticks out. Engines check enclosure by sizing: grow the via by the required margin $e$ and test whether the grown shape lies inside the metal, that is, whether (VIA SIZED BY $e$) NOT METAL is empty.

## Two-sided enclosure rules

At advanced nodes enclosure is often asymmetric: the via needs a large margin $b$ on both sides along one axis and only a small margin $a$ along the other, in either orientation. A via with enclosures (left, right, bottom, top) = (35, 33, 33, 11) and the rule $b = 30$, $a = 4$ passes: the horizontal pair gives min 33 at least 30, and the vertical pair min 11 at least 4. A via with (18, 32, 33, 27) fails both orientations. Such rules let line-end metal stop close to a via in the direction of the wire while protecting the sides, and the check becomes min-over-pairs plus an OR of two conditions.

## Extension and endcaps

Extension measures how far one layer continues beyond another. Poly must extend past diffusion at both ends of a gate by the endcap rule $E$, because line-end shortening and rounding would otherwise leave the channel partly ungated. With $E = 110$ nm, a gate whose poly runs from 296 to 857 over diffusion from 370 to 730 has extensions 74 and 127 nm, so its bottom end needs 36 nm more. Repairing short endcaps is a common automated fix, and the total added length tells how much area the fix costs.
