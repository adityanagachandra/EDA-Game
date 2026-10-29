# Analytical FPGA placement

Analytical placers replace random moves with mathematics: write wirelength as a smooth function of block coordinates, minimize it by solving linear systems or by gradient descent, then spread and legalize. They scale better than annealing to very large FPGAs, and nextpnr's HeAP and the analytical placers in VTR follow this pattern, with legalization made harder by heterogeneous block columns.

## Quadratic placement of one block

With quadratic wirelength, the cost of a movable block at $x$ connected to fixed pins at $x_i$ with weights $w_i$ is $\sum_i w_i (x - x_i)^2$. Setting the derivative $2\sum_i w_i(x - x_i)$ to zero gives the weighted mean $x^* = \sum_i w_i x_i / \sum_i w_i$. For pins at $x = 12$ with weight 2 and $x = 19$ with weight 1, $x^* = (24 + 19)/3 = 14.33$. The $x$ and $y$ problems separate, so each is solved on its own. Quadratic cost overweights long nets, so real placers reweight each connection by the inverse of its current length (the bound-to-bound model) and iterate, which approximates linear, half-perimeter-like wirelength.

## Many movable blocks

With several movable blocks, the zero-gradient conditions form a linear system $Ax = b$: $A$ is the weighted graph Laplacian of the movable connections plus a diagonal term for each block's fixed-pin weights, and $b$ collects the fixed-pin contributions. For two blocks, where block 1 connects to fixed pins at 12 and 6, block 2 to fixed pins at 10 and 18, and the blocks to each other with weight 4, the system is $6x_1 - 4x_2 = 18$, $-4x_1 + 6x_2 = 28$, so $x_1 = 11$ and $x_2 = 12$. The matrix is symmetric positive definite whenever every connected group touches a fixed pin, so conjugate gradients solve it quickly even for millions of blocks. Without fixed pins, the trivial solution would collapse every block onto one point.

## Spreading and legalizing on an FPGA

The unconstrained solution clumps blocks together, so analytical placers alternate solving with spreading: they compute where blocks should go to respect density and add pseudo-connections pulling them there, then re-solve. FPGAs add a twist: block types live in fixed columns. LUT clusters may sit almost anywhere, but DSP blocks and block RAMs may sit only in their columns, so each type is spread and legalized separately. A DSP block that the solver puts at $x = 28.5$, with DSP columns at 4, 10, 16, 22, 28 and so on, moves 0.5 columns to $x = 28$; one at $x = 25$ moves 3 columns either way. Large legalization displacements undo the solver's work, so heterogeneous placers often finish with a short annealing or detailed-placement pass.

## Going deeper

Gort and Anderson, "Analytical Placement for Heterogeneous FPGAs", FPL (2012), the HeAP paper; Spindler, Schlichtmann and Johannes, "Kraftwerk2", IEEE TCAD (2008), for the bound-to-bound net model.
