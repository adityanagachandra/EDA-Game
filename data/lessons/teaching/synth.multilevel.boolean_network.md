# Boolean networks

Two-level forms are too large for most logic, so synthesis works on multi-level structure. The classic representation is the Boolean network: a directed acyclic graph whose internal nodes each hold a local Boolean function, usually as a sum of products, of their fanins. Every multi-level transformation of the MIS and SIS era, from factoring to don't-care simplification, is an operation on this graph.

## Nodes, edges and local functions

A Boolean network has primary inputs, primary outputs and internal nodes. Each internal node $y$ has a local function $y = f_y$ over its fanins, which may be primary inputs or other nodes, and an edge runs from each fanin to $y$. The graph is acyclic (registers cut cycles in sequential logic). For example, $x = a + b$, $y = c + f$, $z = xy + e$ is a three-node network computing $z = (a + b)(c + f) + e$. The global function of a node is obtained by substituting local functions recursively down to the primary inputs. Keeping local functions small and separate is the point: the network can share $x$ among many fanouts without duplicating it.

## Cost of a network

The usual area cost is the sum over nodes of the literals in each local function, counting a node name used as a fanin as one literal. For $x = a + b$ (2), $y = df + c$ (3), $z = ex + ey$ (4) and $w = fz + y$ (3) the total is 12 literals. Factored-form literals are a better area proxy than SOP literals, and SIS reports both. Delay is estimated by levels: primary inputs at level 0 and each node one more than its highest fanin, so in that example $x$ and $y$ are at level 1, $z$ at 2 and $w$ at 3. Both are computed in one topological pass and updated incrementally when a node changes.

## Transformations on the network

Multi-level optimization is a set of local moves on this graph. **Elimination** (collapsing) substitutes a node into its fanouts and deletes it, trading sharing for flatter logic. **Decomposition**, **extraction** and **substitution** create or reuse nodes for common divisors. **Simplification** minimizes a node's local SOP using don't cares derived from the rest of the network. Eliminating $x = bc + a$ into $z = ex + y$ gives $z = bce + ae + y$; here the total literal count happens to stay at 12, but in general elimination grows logic when the node has many fanouts. Every move must preserve the global functions at the outputs, and scripts alternate moves that grow and shrink the network.

## Going deeper

Brayton, Rudell, Sangiovanni-Vincentelli and Wang, "MIS: A multiple-level logic optimization system", IEEE TCAD (1987); Sentovich et al., "SIS: A system for sequential circuit synthesis", UC Berkeley memorandum UCB/ERL M92/41 (1992).
