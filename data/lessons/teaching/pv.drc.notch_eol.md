# Notch, end-of-line and run-length rules

Simple spacing rules treat every pair of facing edges alike. Real decks do not. Edges of the same polygon, wire ends, and long parallel runs of wide wires all have their own rules, because lithography and etch treat them differently. Each variant changes which edge pairs an engine must find and which threshold applies, and they are where many router-produced violations come from.

## Notches

A notch is a gap between two edges of the same polygon that face each other, such as the space between two prongs of a comb. Because the edges belong to one merged polygon, a check that only compares different shapes misses it. For a comb whose prongs span $[0, 130]$, $[160, 290]$ and $[370, 500]$ with a 60 nm notch rule, the first gap (30 nm) violates and the second (80 nm) passes. Engines handle notches by running the spacing search over all edge pairs of the merged layer, including pairs inside one polygon, while skipping adjacent edges that meet at a vertex.

## End-of-line spacing

A wire end narrower than a threshold $W_{eol}$ is a line end. Line ends print shorter than drawn, so they need more spacing in front of them, $S_{eol}$, than the normal spacing $S$. A wire end violates if it is narrow and closer than $S_{eol}$, or closer than $S$ regardless of width. With $S = 40$, $W_{eol} = 60$ and $S_{eol} = 80$ nm, an end 50 nm wide with a 70 nm gap violates, while an end 90 nm wide with the same gap passes. Engines identify line-end edges by their length and their two convex corners, then search a region in front of them.

## Parallel-run-length tables

Spacing tables index the required spacing by the wider wire's width and the length over which the two wires run in parallel. A two-entry example: 40 nm normally, 90 nm when the wider wire is at least 240 nm wide and the parallel run exceeds 500 nm. A pair with width 330, run 850 and spacing 65 needs 90 and fails; a pair with width 390, run 100 and spacing 55 needs only 40 and passes. Computing parallel run length requires the overlap of the two edges' projections, so a table rule needs more than a distance query: it needs the facing edges themselves.
