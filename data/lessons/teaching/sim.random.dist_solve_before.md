# Distribution weights and solve-before

Uniform sampling over the solution space is a precise rule, and often the wrong distribution for verification. SystemVerilog gives two tools to reshape it: `dist` assigns weights to values or ranges, and `solve...before` changes the order in which variables are chosen. Both change probabilities without changing which assignments are legal. Computing their effect is the way to check that a constraint block really produces the traffic mix a test plan intends.

## Weights with := and :/

A `dist` constraint attaches weights to values or ranges. With `:=`, each value of a range gets the full weight. With `:/`, the weight is divided equally among the values of the range. Example: `x dist { 0 := 2, [1:4] := 1 }` gives weights 2, 1, 1, 1, 1, total 6, so $P(x = 0) = 2/6 \approx 0.333$ and $P(x = 3) = 1/6$. Change it to `x dist { 0 := 2, [1:4] :/ 1 }` and the range shares weight 1, so each of 1 to 4 gets $1/4$, the total is 3, $P(x = 0) = 2/3 \approx 0.667$ and $P(x = 3) = (1/4)/3 \approx 0.083$. The two operators differ by exactly the size of the range, which is why mixing them up on a large range silently skews a test.

## solve...before changes the distribution, not the space

Take `rand bit s; rand bit [3:0] d; constraint c { s -> d == 0; }`. Legal solutions: $s = 0$ with any of 16 values of $d$, and $s = 1$ with $d = 0$ only: 17 in all. Uniform sampling gives $P(s = 1) = 1/17 \approx 0.0588$, so the interesting mode is rare. Adding `solve s before d;` makes the solver first choose $s$ uniformly over the values that have at least one solution, then choose $d$ uniformly among the values legal for that $s$. Now $P(s = 1) = 1/2$, and $P(d = 0) = \tfrac{1}{2} \cdot \tfrac{1}{16} + \tfrac{1}{2} \cdot 1 = 0.53125$. The set of legal solutions is the same 17; only their probabilities changed. With a weaker implication $s \to d < c$, the uniform probability is $P(s = 1) = c / (16 + c)$.

## Combining ordering and weights

Ordering and weights compose by total probability. Suppose `op dist { READ := 3, WRITE := 1 }; (op == READ) -> len inside {[1:4]}; (op == WRITE) -> len inside {[1:8]}; solve op before len;`. The solver picks READ with probability $3/4$ and WRITE with $1/4$, then $len$ uniformly in its range, so $P(len = 1) = \tfrac{3}{4} \cdot \tfrac{1}{4} + \tfrac{1}{4} \cdot \tfrac{1}{8} = 0.21875$. Without the ordering, the probabilities would instead depend on how many legal lengths each operation has. A few pitfalls: `solve...before` cannot name `randc` variables, and a `dist` over a huge range with `:=` makes that range dominate. When in doubt, compute the distribution on paper and confirm it with a histogram from a few thousand calls to `randomize()`.

## Going deeper

IEEE Std 1800-2017, clauses 18.5.4 (distribution) and 18.5.10 (variable ordering); Spear and Tumbush, "SystemVerilog for Verification" (3rd ed., 2012), chapter 6.
