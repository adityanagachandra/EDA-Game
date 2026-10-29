# Linear temporal logic

Linear temporal logic (LTL) states properties of individual executions: something always holds, eventually holds, holds until something else does, or holds infinitely often. A system satisfies an LTL formula when every path from an initial state does. Assertion languages such as SVA and PSL are built on the same operators, so being able to evaluate LTL by hand on a concrete path is the foundation for reading and writing properties.

## Paths and the next, always, eventually operators

LTL formulas are evaluated at a position of an infinite path. Atomic propositions such as $p$ hold where the state's label contains them. $X\, \varphi$ (next) holds at position $i$ if $\varphi$ holds at $i + 1$; $G\, \varphi$ (always) if $\varphi$ holds at every position from $i$ on; $F\, \varphi$ (eventually) if it holds at some position from $i$ on, including $i$ itself. Finite-state counterexamples are lassos: a prefix followed by a loop repeated forever, so evaluating $G$ and $F$ means checking every state the path visits from that point, loop included. On the path s0 $\{p\}$ → s1 $\{q\}$ → s2 $\{p, q\}$ with the loop returning from s2 to s1: $X\, q$ holds at s0, $F\, (p \wedge q)$ holds, $G\, p$ fails (s1 lacks $p$), and $G\, q$ fails (s0 lacks $q$).

## Nesting: infinitely often, eventually always, until

Nesting gives the patterns used most in practice. $G\, F\, p$ means $p$ holds infinitely often; on a lasso, it holds exactly when $p$ holds in some state of the loop. $F\, G\, p$ means $p$ eventually holds forever; it holds exactly when $p$ holds in every state of the loop. $p\ U\ q$ (strong until) holds at $i$ if $q$ holds at some $j \ge i$ and $p$ holds at every position from $i$ up to but not including $j$; $q$ must eventually happen. On s0 $\{p\}$ → s1 $\{p\}$ → s2 $\{q\}$ with the loop from s2 back to s1: $p\ U\ q$ holds at s0 because $p$ holds at s0 and s1 and $q$ at s2; $G\, F\, q$ holds because s2 is in the loop; $F\, G\, p$ fails because s2 in the loop lacks $p$; and $q\ U\ p$ holds at once because $p$ holds at s0.

## Response and implication patterns

Properties of protocols are mostly implications under G. $G(p \rightarrow F\, q)$, response, says every $p$ is eventually followed by $q$ (at the same position or later). $G(p \rightarrow X\, q)$ says $q$ holds right after every $p$. $G \neg(p \wedge q)$ says $p$ and $q$ are never true together, a mutual exclusion. To evaluate them on a lasso, check the implication at every distinct position, remembering that positions in the loop repeat, so a $p$ in the loop needs its $q$ somewhere in the loop as well. On s0 $\{p, q\}$ → s1 $\{q\}$ → s2 $\{\}$ with the loop s1 s2: $G(p \rightarrow X\, q)$ holds, because the only $p$ is at s0 and s1 has $q$; $G(p \rightarrow F\, q)$ holds; but $G \neg(p \wedge q)$ fails at s0, and $F\, G (p \vee q)$ fails because s2 in the loop has neither proposition.

## Going deeper

Pnueli, "The Temporal Logic of Programs" (FOCS 1977); Baier and Katoen, "Principles of Model Checking" (2008), chapter 5.
