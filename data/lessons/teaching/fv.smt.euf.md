# Uninterpreted functions and congruence closure

The theory of equality with uninterpreted functions (EUF) lets a solver reason about terms like f(a) without knowing what f computes: the only fact it uses is that equal arguments give equal results. In hardware verification, EUF abstracts datapath units, such as an ALU or a multiplier, so that proofs about control logic do not have to reason about arithmetic. Its decision procedure, congruence closure, is fast and incremental, which makes it a building block of every SMT solver.

## Congruence closure

Given equalities between terms, congruence closure computes all equalities they imply in EUF. It keeps equivalence classes of terms in a union-find structure. First merge the classes of each asserted equality. Then apply congruence: whenever $s$ and $t$ are in one class, merge the classes of $f(s)$ and $f(t)$; repeat until nothing changes. Example with constants $a, b, c, d$ and terms $a, b, c, d, f(a), f(b), f(c)$, assuming $a = b$ and $f(a) = c$. The asserted equalities give classes $\{a, b\}$ and $\{f(a), c\}$. Since $a$ and $b$ are equal, congruence merges $f(a)$ with $f(b)$, giving $\{c, f(a), f(b)\}$. Nothing else merges: $d$ and $f(c)$ stay alone. So the terms fall into 4 classes: $\{a, b\}$, $\{c, f(a), f(b)\}$, $\{d\}$ and $\{f(c)\}$.

## Chains of congruence

Congruence can cascade through nested applications. Assume $a = f(a)$ and consider the terms $a$, $f(a)$, $f(f(a))$ and $b$. The assertion merges $a$ with $f(a)$. Since $a$ and $f(a)$ are now in one class, congruence merges $f(a)$ with $f(f(a))$, so all three are equal; $b$ stays alone. The result is 2 classes. In general, each merge can enable congruences one level up, so the procedure keeps a list of parent terms for every class and re-examines the parents of merged classes; with this bookkeeping it runs in $O(n \log n)$ time for $n$ terms. A term is disequal to another only if a disequality was asserted; when an asserted disequality ends up inside one class, the formula is unsatisfiable, and the chain of merges that caused it is the explanation the solver reports.

## Ackermann's reduction

EUF can also be eliminated: Ackermann's reduction replaces each distinct application $f(t_i)$ by a fresh variable $v_i$ and adds the functional-consistency constraint $t_i = t_j \rightarrow v_i = v_j$ for every pair. For $m$ distinct applications that is $\binom{m}{2} = m(m-1)/2$ constraints: 10 for $m = 5$, 4950 for $m = 100$. The reduction is simple and lets a solver without EUF support decide the formula, but its quadratic growth is why solvers prefer congruence closure inside DPLL(T). In processor verification, Burch and Dill's method compares a pipelined design with its specification by abstracting the ALU and register-file contents with uninterpreted functions, so that the proof depends only on the forwarding and control logic. An uninterpreted function over-approximates the real one, so proofs carry over but counterexamples may be spurious.

## Going deeper

Nelson and Oppen, "Fast Decision Procedures Based on Congruence Closure" (Journal of the ACM, 1980); Burch and Dill, "Automatic Verification of Pipelined Microprocessor Control" (CAV 1994); Kroening and Strichman, "Decision Procedures" (2nd ed., 2016), chapter 4.
