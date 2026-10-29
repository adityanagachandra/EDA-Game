# Simulating a March test on a faulty memory

Knowing which faults an algorithm detects is a theorem about March tests; checking it is a simulation. Running a March test by hand on a few cells with one injected fault shows exactly which operation catches it, why address order matters, and why some short algorithms miss coupling faults. MBIST diagnosis uses the same reasoning to map the failing operation back to a fault type.

## Executing the algorithm

Number the operations in execution order: element by element, and within an element address by address, applying all the element's operations to one address before moving on. For an element with $k$ operations on $n$ cells, the operation at address position $j$ (counting from 0 in that element's order) and operation $i$ has global index (operations in earlier elements) $+ jk + i$. Simulate the memory: a write sets the cell (unless a fault interferes), a read compares the stored value with the expected one. The first read that mismatches is where the test fails. With MATS++ on 5 cells and cell 4 stuck-at-1, the first element uses operations 1 to 5, and in $\Uparrow(r0, w1)$ address 4 is reached at operation 14, an r0 that reads 1.

## Cell faults

A **stuck-at** cell fails the first read that expects the opposite value. A **down transition fault** passes until the cell holds 1 and a w0 is attempted: the write leaves it at 1 and the next r0 of that cell fails. In March Y on 6 cells with cell 3 having a down transition fault, the w0 in $\Downarrow(r1,w0,r0)$ fails silently and the r0 immediately after it, operation 33, catches it. That r0 right after the write is why March Y adds a read to each element.

## Coupling faults and order

A coupling fault fires when the aggressor's cell makes the triggering transition, and is seen when the victim is read afterward. Whether that happens depends on address order. With March C- on 4 cells and a falling transition in cell 1 forcing cell 3 to 0, the $\Uparrow(r1,w0)$ element writes 0 into cell 1 (falling) at operation 16, then reaches cell 3 at operation 19 with r1, which reads 0. If the victim had a lower address, the descending elements would catch it instead. Short tests such as MATS+ lack the needed combinations for some coupling faults and report no failure, which is the answer 0.

## Going deeper

van de Goor, "Testing Semiconductor Memories" (1991); Hamdioui, "Testing Static Random Access Memories" (2004).
