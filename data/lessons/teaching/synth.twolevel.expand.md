# Expanding cubes against the off-set

EXPAND is the step that makes Espresso's cubes prime. It grows each cube by raising literals to dashes as long as the cube stays inside the ON-set plus don't cares, that is, as long as it does not meet the OFF-set. Bigger cubes have fewer literals and swallow other cubes of the cover, which then disappear.

## The legality test for raising a literal

Raising a literal of cube $c$ doubles it: `0011` with $a$ raised becomes `-011`, adding minterm `1011`. The move is legal only if the new cube does not intersect the OFF-set $R$, because every minterm a cube covers must be ON or don't care. Espresso keeps $R$ explicitly as a cover (computed once by complementing ON plus DC) and tests intersection cube by cube. Testing literals one at a time from the original cube tells which raises are individually possible; for `0011` with OFF minterms 0, 1, 4, 5, 6, 9, 10, 11, 14, raising $a$ adds 1011 (minterm 11, OFF) and fails, while raising $b$ adds 0111 (minterm 7, ON) and succeeds.

## Greedy expansion to a prime

Raising literals one after another, keeping each legal raise, ends at a cube where no further literal can be raised: a **prime implicant**. The result depends on the order. With OFF-set minterms 0, 2, 6, 9, 11, 12, 15 and the cube `1010`: raising $a$ gives `-010`, which adds minterm 2 (OFF), so it fails; raising $b$ gives `1-10`, adding 14 (ON), so keep it; raising $c$ from `1-10` gives `1--0`, adding 8 and 12, and 12 is OFF, so it fails; raising $d$ gives `1-1-`, adding 11 and 15, both OFF. The final cube is `1-10`, with 3 literals. Another order could reach a different prime.

## Steering expansion to shrink the cover

Espresso does not expand in alphabetical order. It orders literals so that the expanded cube covers as many other cubes of the current cover as possible, using a blocking matrix (which OFF cubes each raise would hit) and a covering matrix (which ON cubes the cube could absorb). Any cube fully contained in an expanded cube is deleted, so EXPAND both reduces literals and reduces the cube count. The quality of the heuristic is in that ordering: a cube expanded toward a region already covered gains nothing, while one expanded toward isolated cubes can absorb several. Cubes are processed from smallest to largest, since small cubes are the hardest to cover.

## Going deeper

Brayton et al., "Logic Minimization Algorithms for VLSI Synthesis" (1984), on EXPAND; Rudell and Sangiovanni-Vincentelli, "Multiple-valued minimization for PLA optimization", IEEE TCAD (1987).
