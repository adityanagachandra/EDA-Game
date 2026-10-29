# Cubes in positional notation

Two-level minimizers do not manipulate algebraic strings. They store each product term as a cube in positional notation, one symbol per variable, and implement intersection, containment, distance and consensus as bitwise operations. This is the data structure under Quine-McCluskey and Espresso, and the cube calculus is the vocabulary of every two-level algorithm that follows.

## Cubes as strings of 0, 1 and dash

Over variables $a, b, c, d$ the cube $a\bar{c}$ is written `1-0-`: position $i$ holds 1 for an uncomplemented literal, 0 for a complemented one, and a dash when the variable is absent. A cube with $k$ dashes covers $2^k$ minterms, so `1-0-` covers 4. Intersection works position by position: two specified symbols must agree, and the more specific symbol wins otherwise. `1-0-` meets `--00` in `1-00`, which covers 2 minterms; `1---` meets `0---` in the empty cube. In Espresso each position is two bits (01 for 0, 10 for 1, 11 for dash), so intersection is a bitwise AND and an empty field 00 signals an empty result.

## Distance between cubes

The **distance** between two cubes is the number of positions where one has 0 and the other has 1. Dashes never conflict. Distance 0 means the cubes intersect. For `010-` and `0-11`, position 3 has 0 against 1 and nothing else conflicts, so the distance is 1. Distance is the test behind every merge: in Quine-McCluskey two cubes with identical dash positions and distance 1 merge into one cube with a new dash. Distance 2 or more means the cubes are far apart: no single cube generated from them can cover both without other help. In Espresso the distance is computed as the number of empty fields of the bitwise AND.

## Consensus of distance-one cubes

When two cubes are at distance exactly 1, their **consensus** is formed by dropping the conflicting variable and intersecting the remaining positions. For `-00-1` and `10000` over $a$ to $e$ the conflict is in $e$; intersecting the rest gives `1000-`, a cube with 4 literals. Algebraically this is the consensus theorem $xy + \bar{x}z = xy + \bar{x}z + yz$: the consensus is implied by the two cubes together though contained in neither. Iterated consensus, with absorption of contained cubes, generates exactly the prime implicants of a cover. That is an alternative to Quine-McCluskey that starts from any cover rather than from minterms.

## Going deeper

Brayton, Hachtel, McMullen and Sangiovanni-Vincentelli, "Logic Minimization Algorithms for VLSI Synthesis" (1984); Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms" (1996), the two-level chapters.
