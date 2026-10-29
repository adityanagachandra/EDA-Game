# Constrained-random stimulus

Constrained-random verification replaces hand-written stimulus with random values restricted to legal inputs. The testbench declares random variables and constraints; a solver picks values that satisfy them. Random stimulus reaches combinations nobody wrote down, while constraints keep it inside the protocol. For a tool designer the important part is the semantics: what distribution the solver must produce, which is uniform over the set of legal solutions, and what that implies for each individual field.

## Random variables and constraints define a solution space

In SystemVerilog a class declares `rand` variables and `constraint` blocks. Calling `randomize()` asks the solver for an assignment to all random variables that makes every active constraint true. The set of all such assignments is the **solution space**. Example: `rand bit [2:0] x; rand bit [2:0] y; constraint c { x + y < 5; }` restricts the 64 possible $(x, y)$ pairs to those whose sum is below 5. Constraints are declarative and solved together, not executed in order, so `x < y` and `y < 3` in the same block jointly restrict both variables. If no assignment satisfies them, `randomize()` returns 0 and the variables keep their old values, which a testbench must check.

## Counting solutions

For small fields, count the solution space directly. With $x, y \in \{0, \ldots, 7\}$ and $x + y < 5$: for $x = 0$ there are 5 legal $y$ values (0 to 4), for $x = 1$ there are 4, then 3, 2 and 1, and none for $x \ge 5$. The total is $5 + 4 + 3 + 2 + 1 = 15$ legal pairs out of 64. Adding a second constraint $x \ne y$ removes the pairs $(0,0), (1,1), (2,2)$, which satisfy the sum, leaving 12. Counting is how you predict what random stimulus will do, and it is the quantity a BDD-based solver computes internally to sample uniformly. Overflow in narrow arithmetic is a classic source of unexpected solutions, so operand widths in constraints deserve care.

## Uniform over solutions, not uniform per variable

Without `dist` or `solve...before`, the standard requires that legal solutions be equally likely: each of the 15 pairs above has probability $1/15$. That does not make each variable uniform. The probability that $x = 0$ is the number of solutions with $x = 0$ divided by the total, $5/15 \approx 0.333$, while $x = 4$ has probability $1/15 \approx 0.067$ and $x = 5$ never occurs. Variables with many legal partners are favored. This is often what the engineer did not intend: a constraint like "if the packet is short, the payload must be zero" makes short packets rare, because they have few legal payloads. The fix is to shape the distribution explicitly with `dist` weights or an ordering directive.

## Going deeper

IEEE Std 1800-2017, clause 18 (constrained random value generation); Spear and Tumbush, "SystemVerilog for Verification" (3rd ed., 2012), chapter on randomization.
