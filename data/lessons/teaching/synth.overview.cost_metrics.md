# Literal count, gate count and depth as cost

Before a library is chosen, synthesis needs numbers to optimize. The classic ones are the literal count of a sum of products or factored form, the number of two-input gates or AIG nodes, and the logic depth. They are cheap to compute and update incrementally, and they correlate with mapped area and delay well enough to steer restructuring. Knowing exactly how each is counted is the first computational skill of the domain.

## Counting literals

A **literal** is a variable or its complement. A cube (product term) such as $a\bar{b}c$ has three literals, and the literal count of a sum of products is the total over its cubes: $f = a\bar{b}c + \bar{a}d$ has $3 + 2 = 5$ literals. Complemented and uncomplemented occurrences count the same, because in CMOS each literal is roughly one transistor pair in the gate it feeds. Literal count of a factored form counts leaf occurrences: $(a + b)(c + d)$ has 4 literals while its SOP $ac + ad + bc + bd$ has 8. This is why multi-level synthesis prefers factored-form literals as its area measure: they track the transistor count of a static CMOS implementation better than SOP literals do.

## Two-input gate count

Build an SOP with two-input AND and OR gates and free inverters. A cube with $l$ literals needs $l - 1$ ANDs, and $k$ cubes need $k - 1$ ORs. The total is $\sum_i (l_i - 1) + (k - 1) = L - 1$, where $L$ is the literal count. For $f = a\bar{b}c + \bar{a}d + bd$: $L = 7$, so 6 two-input gates (2 + 1 + 1 ANDs and 2 ORs). The same identity holds for any factored form: a form with $L$ leaves has $L - 1$ binary operators. It is the reason an AIG built from a factored form has (leaves minus one) AND nodes, and the reason literal count and two-input gate count move together.

## Depth as the delay proxy

Depth is the number of gate levels on the longest input-to-output path. Built from balanced two-input trees, a cube of $l$ literals has depth $\lceil \log_2 l \rceil$ and an OR of $k$ cubes adds $\lceil \log_2 k \rceil$, so a two-level SOP has depth $\lceil \log_2 l_{max} \rceil + \lceil \log_2 k \rceil$. For three cubes of widths 4, 3, 3 that is $2 + 2 = 4$. Depth tracks delay before cells and wires exist, but it ignores fanout load and wire length. Area and depth conflict: factoring reduces literals but can lengthen paths, and balancing reduces depth but can duplicate logic. Synthesis scripts interleave both kinds of move.

## Going deeper

Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms" (1996); De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), the chapter on multiple-level optimization.
