# Clock meshes

A clock mesh is a grid of horizontal and vertical wires shorted together and driven at many points, usually by the leaves of a tree. Every sink connects to its nearest mesh wire with a short stub. Because the mesh has loops, each sink sees many parallel paths, which averages out driver mismatch and variation. Meshes give the lowest skew of any common topology, which is why high-performance processors use them, but they pay in wire, power and analysis effort.

## Mesh geometry and wirelength

A mesh over a $D \times D$ region with $m$ horizontal and $m$ vertical full-span wires has wirelength $2mD$. For a 3 mm square with 24 wires each way that is 144 mm, before counting the stubs to the sinks and the tree that drives the mesh. The mesh pitch $D/m$ sets how long the stubs are, and therefore how much residual skew the stubs add; a finer mesh lowers skew and raises wire. Unlike a tree, a mesh's wirelength does not depend on where the sinks are, which is both its robustness and its waste.

## The power cost

All of the mesh toggles every cycle. A 52 by 52 mesh over a 4 mm square at 0.28 pF/mm has $2 \cdot 52 \cdot 4 \cdot 0.28 = 116.5$ pF of wire; with 9000 pins of 1.2 fF it switches 127.3 pF, which at 3 GHz and 0.75 V is $127.3 \cdot 0.5625 \cdot 3 = 214.8$ mW. Mesh drivers also fight briefly whenever their input edges are skewed, adding short-circuit current. These costs are why meshes are reserved for the regions where skew matters most.

## Comparing a mesh with an H-tree

On the same region, a four-level H-tree uses $1.5D(2^4 - 1) = 22.5D$ of wire, so a mesh with 16 wires each way ($32D$) uses $9.5D$ more; at $D = 4$ mm and 0.24 pF/mm that is about 9.1 pF, which costs $9.1 V^2 f$ in extra power. The mesh buys lower skew and tolerance of variation, since a weak driver or a slow branch is compensated by its neighbours. The analysis cost is the other price: a mesh is not an RC tree, so Elmore's closed form does not apply, and its delays must come from solving the mesh as a linear circuit or from transient simulation.

## Going deeper

Restle et al., "A clock distribution network for microprocessors", IEEE Journal of Solid-State Circuits (2001); Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the clock chapter.
