# Netlists as hypergraphs

A netlist is a set of cells joined by nets, and a net can join any number of cells. The natural mathematical object is therefore a hypergraph, not a graph: vertices are cells, hyperedges are nets. Partitioners, placers and routers all start from this representation, so a tool designer must know what it counts, what it hides, and how its basic statistics (pins, degrees, neighbourhoods) drive runtime and memory.

## Cells, nets and pins

A hypergraph $H = (V, E)$ has a vertex set $V$ and a family $E$ of hyperedges, each a subset of $V$ with at least two members. In a netlist, $V$ is the set of cells (or of pins, in a pin-level model) and each net $e \in E$ is the set of cells it touches. The incidence between a cell and a net is a **pin**. The total pin count $P = \sum_{e} |e|$ is the true size of a netlist: almost every algorithm in physical design runs in time proportional to $P$, not to the number of cells or nets.

For example, nets $n_1 = \{c_1, c_2, c_3\}$, $n_2 = \{c_2, c_4\}$ and $n_3 = \{c_1, c_3, c_4, c_5\}$ have sizes 3, 2 and 4, so $P = 9$ pins over 5 cells and 3 nets. Real designs average roughly three to four pins per net, but a few clock and reset nets have thousands.

## Degrees and neighbourhoods

The **degree** of a cell is the number of nets incident to it, and the size of a net is its own degree as a hyperedge. Two cells are **neighbours** if some net contains both. Neighbourhoods are what local algorithms touch: a Fiduccia-Mattheyses move updates the gains of the neighbours of the moved cell, and a placer's gradient for a cell depends on the nets it is on.

Take $n_1 = \{c_1, c_2, c_3\}$, $n_3 = \{c_1, c_3, c_4, c_5\}$. Cell $c_1$ has degree 2, and its neighbourhood is the union of its nets minus itself: $\{c_2, c_3, c_4, c_5\}$, four cells. Cell $c_3$ shares both nets with $c_1$ but is counted once. Summing all cell degrees gives $P$ again, which is the double-counting identity $\sum_v \deg(v) = \sum_e |e|$.

## What the hypergraph abstracts away

The hypergraph keeps connectivity and nothing else. It forgets which pin of a cell a net uses, which pin drives the net, cell sizes, and geometry. Partitioning needs only connectivity and cell area, so the plain hypergraph (with vertex weights) is its exact input. Timing needs direction, so it uses a directed pin graph instead. Placement needs pin offsets inside cells.

The abstraction has a real cost: converting a hyperedge to ordinary graph edges, as many algorithms want, distorts the objective, and a large net of $k$ pins becomes $k(k-1)/2$ edges if expanded naively. A tool designer chooses the representation to match what the algorithm must compute, and keeps hyperedges native where the objective (cut nets, net wirelength) is defined per net.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 1 and chapter 2; Karypis and Kumar's hMETIS manual for the hypergraph input format.
