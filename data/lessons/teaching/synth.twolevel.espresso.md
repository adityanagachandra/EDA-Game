# The Espresso heuristic minimizer

Espresso, developed at IBM and Berkeley in the early 1980s, is the two-level minimizer that made PLA synthesis practical and still runs inside multi-level flows for node simplification. Instead of generating all primes, it improves a cover in place with expand, irredundant and reduce, using cube operations and tautology checks that scale to many inputs.

## The overall loop

Espresso first computes the OFF-set by complementing ON plus don't cares. It then runs EXPAND (make cubes prime, absorbing others) and IRREDUNDANT (drop unnecessary cubes), and extracts the essential primes, which it sets aside since they appear in every cover. The main loop repeats REDUCE, EXPAND, IRREDUNDANT while the cost (cubes, then literals) improves. When the loop stalls, LAST_GASP tries a different reduction to escape. The essentials are then added back. Every step preserves the function, so the result is always a valid cover; only its size is heuristic. The final cover is prime and irredundant, a local minimum.

## Why REDUCE helps

A prime and irredundant cover is stuck: expanding cannot grow its cubes further, and none can be removed. REDUCE makes progress possible by shrinking each cube, in order, to the smallest cube that still covers the minterms only it covers, given the rest of the cover. The reduced cover is no longer prime, but it represents the same function. The next EXPAND can now grow each cube in a different direction, possibly absorbing other cubes, and IRREDUNDANT may then remove more cubes. The sequence reduce, expand, irredundant moves the cover between local minima, and the order in which cubes are reduced matters as much as the expansion order.

## Scope and limits

Espresso handles multiple outputs by adding an output part to each cube, so that one product term can be shared by several outputs, exactly the row sharing a PLA rewards. It runs on cube lists with unate recursion, so functions with 30 or more inputs and thousands of cubes are routine. Its answer is not guaranteed minimum; Espresso-exact and later exact methods (based on implicit prime generation and covering) exist for small or structured functions. In multi-level synthesis Espresso, or an equivalent minimizer, simplifies the SOP of each network node using that node's don't cares, which is where it lives on in SIS-style flows.

## Going deeper

Brayton, Hachtel, McMullen and Sangiovanni-Vincentelli, "Logic Minimization Algorithms for VLSI Synthesis" (1984); Rudell, "Logic Synthesis for VLSI Design", PhD thesis, UC Berkeley (1989).
