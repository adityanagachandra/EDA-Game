# Sequence pairs

A sequence pair encodes any packing of rectangles, slicing or not, as two permutations of the block names. Introduced by Murata and colleagues in the mid-1990s, it turns the geometric relations between blocks into two constraint graphs whose longest paths give the coordinates. It is the standard non-slicing representation for annealing-based and fixed-outline floorplanners.

## Reading relations from two sequences

A sequence pair is $(\Gamma_+, \Gamma_-)$, two orderings of the $n$ blocks. For two blocks a and b: if a precedes b in both sequences, a is left of b; if a follows b in $\Gamma_+$ but precedes it in $\Gamma_-$, a is below b. Every pair of blocks gets exactly one of these relations (or the mirrored ones), so the pair fully specifies a topology. With $\Gamma_+ = (A\ C\ D\ B\ E)$ and $\Gamma_- = (A\ B\ E\ D\ C)$: A precedes everything in both, so A is left of all; B and E keep order in both, so B is left of E; D precedes B in $\Gamma_+$ but follows it in $\Gamma_-$, so B is below D. There are $(n!)^2$ sequence pairs.

## Packing by longest paths

The left-of relations form a horizontal constraint graph: an edge from a to b with length $w_a$ whenever a is left of b. Each block's x coordinate is the longest path to it from a source at x = 0, and the chip width is the longest path overall. The below relations form a vertical graph with edge lengths $h_a$, giving y coordinates and height. Both graphs are acyclic, and processing blocks in $\Gamma_-$ order is a valid topological order for both, so a simple $O(n^2)$ double loop computes the packing. In the example with A = 6 by 3, B = 8 by 7, C = 7 by 3, D = 2 by 2 and E = 4 by 1, the x coordinates are 0, 6, 6, 6 and 14 and the chip is 18 wide.

## Annealing over sequence pairs

Annealers perturb a sequence pair by swapping two blocks in one sequence, swapping them in both, or rotating a block, then re-pack and evaluate. Any packing can be represented, so the search is complete, but the space has redundancy: many sequence pairs give the same packing. Evaluation cost is the bottleneck. The $O(n^2)$ longest-path packing can be replaced by $O(n \log n)$ methods based on weighted longest common subsequences (the FAST-SP approach), which matters when an annealer evaluates millions of moves. For fixed-outline problems, the cost function adds penalties for width and height beyond the outline, and wirelength is computed from block centres or pin offsets after packing.

## Going deeper

Murata, Fujiyoshi, Nakatake and Kajitani, "VLSI module placement based on rectangle-packing by the sequence-pair" (IEEE TCAD, 1996); Tang, Tian and Wong, "FAST-SP: a fast algorithm for block placement based on sequence pair" (ASP-DAC 2001).
