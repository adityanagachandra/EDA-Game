# Deductive fault simulation

Deductive fault simulation computes, in one pass of the fault-free circuit, the complete set of faults that each pattern detects. Instead of simulating faulty circuits, it propagates at every line a fault list: the set of faults that would complement that line. It is a beautiful application of set algebra and a direct model for reasoning about which faults a vector reaches.

## Fault lists

For the current vector, every line $x$ has a fault-free value and a **fault list** $L_x$: the faults that, if present, would make $x$ take the opposite value. A primary input's list contains only its own stuck-at fault opposite its value. At each gate, the output's list is computed from the input lists and the input values, and then the output's own local fault (stuck at the complement of its fault-free value) is added. At a primary output, $L_z$ is the set of faults this vector detects. One pass gives all detected faults, rather than one simulation per fault.

## Propagation rules

Let $c$ be the gate's controlling value and $S$ the inputs at $c$. If **no input** is at $c$ (for AND, all inputs are 1), complementing any input complements the output, so $L_{out} = \bigcup_i L_i$. If some inputs are at $c$, the output flips only if a fault complements **every** input in $S$ and **none** of the others:
$$L_{out} = \Big(\bigcap_{i \in S} L_i\Big) \setminus \Big(\bigcup_{i \notin S} L_i\Big).$$
Take a 3-input AND with a = 0, $L_a = \{5, 6, 8\}$, b = 1, $L_b = \{1, 2, 5\}$, c = 1, $L_c = \{1, 3, 5\}$. Only a is at 0, so $L_{out} = \{5,6,8\} \setminus \{1,2,3,5\} = \{6, 8\}$, before adding the local fault. Fault 5 flips a to 1 but also flips b and c to 0, so the output stays 0. For XOR, a fault flips the output if it flips an odd number of inputs.

## Strengths and limits

Deductive simulation is fast for combinational two-valued circuits because set operations on sorted lists or bit vectors are cheap. Its weaknesses are memory (lists can be large near the inputs of big cones), and awkward handling of unknown values and timing: with X values, list membership becomes uncertain, and with delays the lists must be recomputed per event. Those limits are why concurrent fault simulation, which evaluates faulty gate copies explicitly, took over for sequential and timing-aware use, while PPSFP took over for full-scan combinational use.

## Going deeper

Armstrong, "A deductive method for simulating faults in logic circuits," IEEE Transactions on Computers, 1972; Bushnell and Agrawal (2000), chapter 5.
