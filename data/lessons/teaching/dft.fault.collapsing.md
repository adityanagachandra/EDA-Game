# Fault collapsing

A raw fault list has two faults per line, but many of them are indistinguishable or implied by others. Fault collapsing shrinks the list before ATPG and fault simulation, typically to about half, without losing anything: every removed fault is guaranteed to be detected whenever a kept one is. It is one of the cheapest speedups in the test flow and a first exercise in reasoning about test sets as sets.

## Equivalence

Two faults are **equivalent** if every test for one detects the other, that is, their test sets are identical. Structurally, at a gate: an input stuck at the gate's **controlling** value forces the output exactly as an output fault would. For an AND gate, any input stuck-at-0 is equivalent to the output stuck-at-0; for NAND, input stuck-at-0 equals output stuck-at-1; for OR, input stuck-at-1 equals output stuck-at-1; for NOR, input stuck-at-1 equals output stuck-at-0. A NOT gate pairs input stuck-at-0 with output stuck-at-1 and vice versa. Equivalence is an equivalence relation, so the local rules are merged with union-find into classes, and one representative per class suffices. A 2-input AND has 6 faults and 4 classes.

## Dominance

Fault $f$ **dominates** $g$ if every test for $g$ also detects $f$. Then targeting $g$ suffices and $f$ can be dropped. At an AND gate, the output stuck-at-1 dominates each input stuck-at-1: a test for input a stuck-at-1 sets $a = 0$, the other inputs to 1, and propagates the output, which is exactly a test for the output stuck-at-1, but the output fault has other tests too. The rule generalizes: drop the output fault at the non-controlling response of each multi-input gate. A 2-input AND then keeps 3 faults out of 6. Dominance collapsing must be used with care when the dominated fault turns out untestable, since the dropped fault might still be testable.

## Checkpoints and the collapse ratio

Fanout branches break equivalence: a stem fault affects every branch, a branch fault only one gate, so faults are not merged across fanout. The **checkpoint theorem** says that in a combinational circuit, a test set detecting all single stuck-at faults on primary inputs and fanout branches detects all single stuck-at faults. The **collapse ratio** is the collapsed count divided by the original count. For example, a circuit with 28 faults that collapses to 18 equivalence classes and 13 after dominance has ratio 13/28 = 0.46. Ratios near 0.4 to 0.5 are typical for gate-level netlists.

## Going deeper

Bushnell and Agrawal (2000), chapter 4; McCluskey and Clegg, "Fault equivalence in combinational logic networks," IEEE Transactions on Computers, 1971.
