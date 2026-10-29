# Static single assignment

Static single assignment (SSA) is the form in which every variable is assigned exactly once. Renaming each assignment to a fresh version makes use-def chains explicit, and phi functions merge versions where control flow joins. Nearly every modern compiler IR, including LLVM and MLIR, is SSA, and hardware netlists are naturally close to it, so SSA is the bridge between procedural HDL code and circuits.

## Renaming straight-line code

In straight-line code, give each assignment to a variable the next version number and rewrite each read to name the most recent version above it. With x's incoming value as $x_0$, the code x = x + w; x = z + z; z = y + 4; becomes $x_1 = x_0 + w_0$, $x_2 = z_0 + z_0$, $z_1 = y_0 + 4$. A later read of z refers to $z_1$, and a read of x to $x_2$.

This is exactly what blocking assignments mean in a combinational always block: each write creates a new value, and later statements read it. Synthesis tools perform this renaming when they turn a process into logic: every version becomes a net, and the final versions drive the block's outputs. Nonblocking assignments, in contrast, read the pre-edge values, all version 0.

## Phi functions at joins

When control joins, a variable may arrive with different versions from different predecessors. After if (sel) z = y + y; else w = z + 9;, z arrives as $z_1$ from the then arm and $z_0$ from the else arm. SSA inserts $z_2 = \phi(z_1, z_0)$ at the join, which selects the version from whichever edge was taken. In hardware this phi is exactly a multiplexer controlled by sel.

Minimal SSA places a phi at the join for every variable assigned in at least one arm, because two different definitions meet there; a variable assigned in neither arm needs none. In an if/else where the arms assign {z, y} and {w}, minimal SSA places three phis: for w, y and z.

## Pruned SSA

Minimal SSA may place phis whose results are never used. Pruned SSA places a phi only for a variable that is live at the join, meaning some path from the join reads the variable before overwriting it. If the code after the join is x = y + z; z = z + x; and nothing later reads w, then w's phi is dead: pruned SSA omits it, while y and z, if assigned in an arm, keep theirs because they are read.

Computing pruned SSA requires liveness, a backward data-flow analysis, which costs compile time but saves IR size. Many compilers build minimal SSA and let dead-code elimination remove useless phis afterwards. In hardware terms, a dead phi is a multiplexer driving nothing, which opt_clean in Yosys would sweep.

## Going deeper

Cytron, Ferrante, Rosen, Wegman and Zadeck, "Efficiently Computing Static Single Assignment Form and the Control Dependence Graph" (TOPLAS 1991); Braun et al., "Simple and Efficient Construction of Static Single Assignment Form" (CC 2013).
