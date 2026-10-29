# Pin and I/O assignment

Pins on a block's boundary, whether chip I/O or the ports of a hierarchical partition, must be assigned to legal slots. Their locations pull the connected logic toward them, so a bad assignment lengthens every net that leaves the block. Pin assignment is a matching problem, solved exactly by the Hungarian algorithm or approximately by fast heuristics inside floorplanning loops.

## Pins, slots and costs

A block boundary offers legal slots: positions on allowed layers at the right pitch, outside blockages. Each pin connects to logic inside, or to the outside world, at some location. The cost of putting pin $i$ in slot $j$ is an estimate of the wirelength this creates, most simply the Manhattan distance from the slot to the pin's connected cell or to the centroid of its connected cells. For a 100 by 100 um block with pin p1 connecting to a cell at (57, 29) and a slot at (37, 0), the cost is $20 + 29 = 49$. The assignment's total cost is the sum over pins. Additional constraints include pin groups that must stay together, spacing, and pins restricted to one side.

## Optimal assignment as matching

Assigning $n$ pins to $n$ slots with one pin per slot and minimum total cost is the linear assignment problem, a minimum-cost perfect matching in a bipartite graph. The Hungarian algorithm solves it in $O(n^3)$; for a handful of pins, checking all $n!$ permutations works too. With more slots than pins, dummy pins with zero cost make the problem square. For three pins there are six assignments, and evaluating each and keeping the cheapest gives the optimum. The improvement over a naive assignment, such as pins in netlist order around the boundary, can be large because a single badly placed pin drags its logic across the block.

## Pin assignment in the flow

Pin assignment and placement depend on each other: pins should go near their logic, but the logic is placed near its pins. Flows break the cycle by iterating: place with initial pins, reassign pins to the placed logic, then place again. For chip-level I/O, constraints from the package and the board often fix the order, leaving only fine adjustment. For hierarchical blocks, pins on both sides of a boundary must agree, so assignment is done at the parent level. OpenROAD's ioPlacer assigns pins with a matching formulation over slots on the die boundary and supports constraints such as allowed edges and pin groups. Because it is fast, an agent can rerun it after any floorplan change.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 3.7 on pin assignment; the OpenROAD documentation for ioPlacer.
