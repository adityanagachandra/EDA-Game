# Partitioning a subject DAG into trees

A real subject graph is a DAG: a gate output may feed several gates. The classic mapper (DAGON) handles this by cutting the DAG into trees at every multiple-fanout point and covering each tree optimally. The partition is fast and simple, but it fixes which signals must exist as gate outputs, so it bounds the quality of the whole mapping.

## Where the DAG is cut

A node becomes a **tree root** if it is a primary output or if it drives two or more gate inputs. Each root, together with the single-fanout nodes beneath it down to primary inputs or other roots, forms one tree. Inside a tree every node has exactly one consumer, so the tree's internal signals are private to it. In the network g1 = NAND(a, b), g2 = INV(g1), g3 = NAND(g1, c), g4 = NAND(g2, g3) with output g4, the gate g1 has fanout two, so it is a root; the trees are {g1} and {g2, g3, g4}.

## Counting trees and tree sizes

The number of trees equals the number of roots: primary outputs plus multi-fanout gates, each counted once even if it is both. To size a tree, walk down from its root and stop at primary inputs and at other roots, counting the gates you visit. A topological pass computes all fanouts in time linear in the number of edges, and a second pass assigns every gate to the tree of the root it reaches first going upward. The roots become the leaves of the trees that consume them, so after mapping each root is a real gate output with its own wire.

## What partitioning costs

Forcing every multi-fanout node to be a gate output forbids any match that would swallow it. If a shared NAND2 feeds three consumers that could each have used one NAND3, the tree mapper must build the shared NAND2 plus three smaller covers. Duplicating the shared logic into each consumer is sometimes cheaper in delay and occasionally in area, but it lies outside the partition's search space. The loss grows with the density of reconvergent fanout. That is the gap that DAG covering and cut-based mapping close: they treat the partition as unnecessary and let matches cross fanout points.

## Going deeper

Keutzer, "DAGON: Technology Binding and Local Optimization by DAG Matching", DAC 1987; De Micheli, "Synthesis and Optimization of Digital Circuits", section 10.3.
