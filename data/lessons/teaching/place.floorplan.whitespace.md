# Floorplan area, whitespace and aspect ratio

A floorplan arranges blocks, such as macros, hard IP or soft partitions, inside an outline. Its first measurable qualities are the area it uses, the whitespace it wastes and its shape. Every floorplanning algorithm, from annealing over Polish expressions to fixed-outline sequence-pair packers, scores candidates with these numbers before any wirelength term is added.

## Outline area and whitespace

The floorplan's area is the area of the bounding box of all blocks, or the given outline in a fixed-outline flow. Whitespace, also called dead space, is the outline area not covered by blocks:
$$\text{whitespace} = \frac{A_{outline} - \sum_b A_b}{A_{outline}}.$$
Three blocks of areas 200, 130 and 442 um$^2$ in a 38 by 24 um outline (912 um$^2$) leave 140 um$^2$, which is 15.4% whitespace. Some whitespace is wanted: it becomes routing channels, room for buffers and space for the standard cells that surround hard macros. But whitespace that comes from a poor arrangement is simply wasted die area, and the floorplanner tries to minimise it subject to the other goals.

## Packing and its bounding box

When blocks are packed by a representation, the outline is their bounding box, so whitespace depends on how the shapes fit. Two blocks of 49 by 37 and 24 by 28 placed side by side, bottoms aligned, have a bounding box 73 wide and 37 tall, 2701 um$^2$, of which the blocks cover 2485 um$^2$: 8.0% whitespace, all of it above the shorter block. Stacking the same blocks instead gives a 49 by 65 box, 3185 um$^2$ and 22.0% whitespace. Mismatched heights in a row or mismatched widths in a column are the main source of waste, which is why floorplanners consider block rotation and, for soft blocks, variable aspect ratios.

## Fixed outlines and aspect ratio

Classical floorplanning minimises area freely, but in a hierarchical design the block's outline is dictated by the parent floorplan, so modern floorplanners solve a fixed-outline problem: fit all blocks inside a given W by H with a maximum whitespace. Given total block area $A$, a fixed height $H$ and allowed whitespace fraction $\gamma$, the outline width must be at least $A/((1-\gamma)H)$. For 9000 um$^2$ of blocks, height 91 um and 5% whitespace, the width is $9000/(0.95 \cdot 91) = 104.11$ um. The aspect ratio $H/W$ is then a constraint rather than an outcome, and annealers add a penalty for exceeding the outline. Fixed-outline packing is markedly harder than free packing, because many candidate packings that are compact are still the wrong shape.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 3; Adya and Markov, "Fixed-outline floorplanning: enabling hierarchical design" (IEEE TVLSI, 2003).
