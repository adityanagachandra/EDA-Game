# The D-algorithm

The D-algorithm was the first complete algorithm for combinational test generation: for any detectable stuck-at fault it finds a test, and for an undetectable one it eventually proves none exists. It works directly on the netlist with five-valued values, assigning internal lines as well as inputs, and organizes its search around two frontiers. Those frontiers are the bookkeeping every later structural ATPG inherits.

## The D-frontier

The algorithm starts by placing D or $\overline{D}$ at the fault site, using a **primitive D-cube of failure** for the gate driving it (for an AND output stuck-at-0, inputs 1, 1 and output D). The error must then be driven forward. The **D-frontier** is the set of gates whose output is still X while at least one input carries D or $\overline{D}$. Those are the only places where the error can advance. **D-drive** picks a frontier gate and applies a **propagation D-cube**: side inputs set to non-controlling values, so the output becomes D or $\overline{D}$. An empty D-frontier with no error at an output means every path is blocked, so the algorithm must backtrack.

## The J-frontier

Assigning side inputs and internal values creates obligations. A line assigned a value that its gate's inputs do not yet imply must be **justified**. The **J-frontier** is the set of gates whose output is assigned but not implied by their current inputs. For an AND with output assigned 0 and both inputs X, justification can choose either input to be 0; for output 1 it must set every input to 1. After each decision, **implication** propagates every forced value forward and backward; a line implied to both 0 and 1 is a conflict.

## Search and completeness

The algorithm alternates D-drive and justification, implying after every assignment, until an output carries D or $\overline{D}$ and the J-frontier is empty: the primary input values, with X as don't-care, are then a test. On a conflict it backtracks to the last decision with an untried alternative. Because it branches on internal lines, its decision space is much larger than the $2^n$ input space, and many decisions turn out unjustifiable. Exhausting all alternatives proves the fault untestable. Its multiple-path handling is complete: the D-frontier formulation naturally sensitizes several paths when a fault needs it.

## Going deeper

Roth, "Diagnosis of automata failures," IBM Journal, 1966; Bushnell and Agrawal (2000), chapter 7.
