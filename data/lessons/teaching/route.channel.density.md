# Channel density

Channel routing is the classical detailed-routing problem: a rectangular region with pins on its top and bottom edges, nets that enter at those pins, horizontal trunks on tracks and vertical branches down to the pins. Its central lower bound is density, the most nets that must cross one column. The same interval reasoning reappears in track assignment and panel routing inside modern routers.

## Local density

Columns are numbered left to right; each column may hold a top pin and a bottom pin, each labelled with a net or empty. A net's trunk must span from its leftmost to its rightmost pin column. The local density at column $c$ is the number of nets whose span includes $c$; those nets all need a horizontal track at that column, and two of them cannot share one.

Example: with top pins 0, 4, 0, 1, 0, 0, 2, 1 and bottom pins 1, 2, 3, 2, 3, 3, 0, 4, the spans are net 1 [1, 8], net 2 [2, 7], net 3 [3, 6] and net 4 [2, 8], so column 7 has local density 3 (nets 1, 2 and 4).

## Channel density as a lower bound

The channel density $d$ is the maximum local density over all columns. Any routing that gives each net a single horizontal trunk needs at least $d$ tracks, because at the densest column $d$ trunks run in parallel. Computing it is a plane sweep: sort span endpoints, add one at each left end, subtract one after each right end, and track the maximum.

Example: spans [6, 7], [2, 6], [1, 5], [1, 4] and [7, 8] reach density 3 at columns 2 to 4. Density counts horizontal demand only. Vertical constraints between nets can force more tracks than density, which is why the left-edge algorithm achieves density only when there are no vertical constraints.

## Bottleneck columns

Several columns may attain the density. Each is a bottleneck: no rearrangement of trunks reduces the tracks needed there, so a channel router cannot use fewer than $d$ tracks, and every column at density pins down a set of nets that must be stacked. Example: spans [3, 7], [4, 5], [6, 8], [1, 8] and [1, 2] reach density 3 at four columns.

Channel heights in building-block layouts were set by these numbers, and reducing density meant changing the placement of pins. In modern flows the same quantity appears as gcell-edge demand in global routing and as panel density in track assignment. A router designer should read density as the horizontal capacity a region must supply, before any vertical interaction is considered.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.), on channel routing; Sherwani, "Algorithms for VLSI Physical Design Automation", chapter on channel routing.
