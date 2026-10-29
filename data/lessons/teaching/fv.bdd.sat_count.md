# Counting satisfying assignments on a BDD

Once a function is a BDD, counting its satisfying assignments takes one pass over the graph, even when the count is astronomically large. The same pass gives signal probabilities, uniform sampling of solutions and the weights used in probabilistic analysis. The only subtlety is that a reduced BDD skips levels, and every skipped variable doubles the number of assignments a path stands for.

## Counting paths to 1

Number the levels of the order $0, \dots, n - 1$ and give both terminals level $n$. Process nodes bottom-up. For the terminals, the count is 1 for terminal 1 and 0 for terminal 0. For a node $u$ at level $\ell(u)$ with children $L$ and $H$, $\text{count}(u) = \text{count}(L) \cdot 2^{\ell(L) - \ell(u) - 1} + \text{count}(H) \cdot 2^{\ell(H) - \ell(u) - 1}$, which counts assignments to the variables from $u$'s level downward. The total is $\text{count}(\text{root}) \cdot 2^{\ell(\text{root})}$. Example over $x_1 < x_2 < x_3$: the root tests $x_1$ with low child the terminal 1 and high child a node testing $x_3$ (low 0, high 1). The $x_3$ node, at level 2, has count $0 + 1 = 1$. The root, at level 0, has count $1 \cdot 2^{3-0-1} + 1 \cdot 2^{2-0-1} = 4 + 2 = 6$. The function is $\neg x_1 \vee x_3$, true on 6 of 8 inputs.

## Skipped levels

Skipped levels are where hand counts go wrong. An edge from level $\ell(u)$ to a child at level $\ell(c)$ skips $\ell(c) - \ell(u) - 1$ variables, each free, so it multiplies the child's count by $2^{\ell(c) - \ell(u) - 1}$. Take order $x_1 < x_2 < x_3 < x_4$, root testing $x_1$ with low child a node $a$ testing $x_4$ (low 0, high 1) and high child a node $b$ testing $x_2$ (low 0, high 1). Node $a$ is at level 3: count 1. Node $b$ is at level 1, and its high edge to the terminal skips $x_3$ and $x_4$: count $1 \cdot 2^{4-1-1} = 4$. The root's low edge to $a$ skips $x_2$ and $x_3$: $1 \cdot 2^{3-0-1} = 4$; its high edge to $b$ skips nothing: $4 \cdot 1 = 4$. Total 8, and indeed $\neg x_1 x_4 + x_1 x_2$ is true on 4 + 4 inputs. Work is linear in the number of nodes.

## From counts to probabilities

Dividing the count by $2^n$ gives the probability that a uniformly random input satisfies $f$: for $\neg x_1 \vee x_3$, $6/8 = 0.75$. If each input $x_i$ is 1 with its own probability $p_i$, the same bottom-up pass computes $P(u) = (1 - p_i) P(L) + p_i P(H)$ for a node testing $x_i$, with $P = 1$ and $P = 0$ at the terminals and skipped variables contributing a factor of 1. That is the signal probability that power estimators and random-pattern testability tools need. Counting also enables uniform sampling: walk down from the root and take each child with probability proportional to its weighted count, which picks every satisfying assignment with equal probability. Constrained-random stimulus generators in verification use BDD-based sampling of this kind for small constraint sets, and switch to SAT- or SMT-based samplers when BDDs grow too large.

## Going deeper

Knuth, "The Art of Computer Programming", volume 4A, section 7.1.4; Bryant, "Graph-Based Algorithms for Boolean Function Manipulation" (1986).
