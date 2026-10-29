# Technology mapping as library covering

Logic synthesis ends with a technology-independent network, usually an and-inverter graph. Technology mapping turns it into a netlist of real library cells. For a tool designer the important point is that mapping is a covering problem with a cost: choose a set of library matches that together implement every node, wire each match's inputs to other matches or primary inputs, and minimize area, delay or power. Everything in this domain refines that one formulation.

## Subject graphs and pattern graphs

The mapper rewrites both sides into one primitive base, classically two-input NAND and inverter. The network becomes the **subject graph**. Each library cell becomes one or more **pattern graphs** over the same base: NAND2 is NAND(x, y), AND2 is INV(NAND(x, y)), AOI21 is INV(NAND(NAND(a, b), INV(c))). With a common base, the question "can cell C implement this piece of logic?" becomes "does C's pattern appear, rooted at this node, in the subject graph?". The base must be functionally complete, which is why INV and NAND2 are always kept: any subject graph can then be covered at least one node at a time.

## Covers and their costs

A **match** is a pattern placed at a subject node. A **cover** is a set of matches such that every primary output is the root of a match, and every input of a chosen match is either a primary input or the root of another chosen match. Internal nodes of a match need no gate of their own: they are absorbed into the cell. The cost of a cover is the sum of the chosen cells' areas, or the latest arrival time at the outputs. Example: NAND(INV(NAND(a, b)), c) can be covered by NAND2 + INV + NAND2 or by one NAND3. If NAND3 has area 4 and the other cells areas 3, 2, 3, the single cell wins, 4 against 8.

## What the formulation leaves out

Covering is exact only relative to its inputs. It cannot find a cover that does not exist in the chosen subject graph, so the decomposition fixes the search space. Structural matching cannot see a cell whose function is present but whose shape is not, which is the motivation for Boolean matching and for structural choices in ABC. Some cells, such as XOR2, have no tree-shaped pattern at all. And a pure area cover ignores load and wiring, which is why delay-oriented mapping and post-mapping sizing follow. A mapper's answer is therefore always "the best cover of this subject graph with these patterns and this cost model", which is the claim its tests must check.

## Going deeper

Giovanni De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 10; Kurt Keutzer, "DAGON: Technology Binding and Local Optimization by DAG Matching", DAC 1987; Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms", chapter on technology mapping.
