# Elmore delay of an RC tree

A routed net is not a wire with one delay. It is a tree: one driver, a trunk, branches and several sinks, each piece carrying resistance along its length and capacitance to ground. Timing analysis needs a delay number for every driver-to-sink arc of that tree, millions of times per run, so the formula must be cheap. The Elmore delay is that formula. It is the first estimate every timer computes and the baseline every better model is measured against.

## The first moment of the impulse response

Model the net as resistors in a tree, rooted at the driver, with a capacitor from every node to ground. Drive the root with an impulse and watch the voltage at a sink node $i$. Its response $h_i(t)$ rises and decays. W. C. Elmore proposed in 1948 that the delay to node $i$ be the mean of that response, its first moment $T_i = \int_0^\infty t\, h_i(t)\, dt$. For an RC tree the response is unimodal and the first moment upper-bounds the true 50% delay of a step input. That is why the estimate is safe to use in early optimization: it never claims a path is faster than it is.

## Shared resistance times downstream capacitance

The first moment has a closed form that needs no integration. For every capacitor $C_k$ in the tree, find the resistance that the path from the root to $k$ shares with the path from the root to the sink $i$. Call it $R_{ik}$. Then

$$T_i = \sum_k R_{ik} C_k.$$

Take a driver resistance of $100\ \Omega$ into node $a$ with $C_a = 10$ fF, then $200\ \Omega$ from $a$ to node $b$ with $C_b = 20$ fF. For sink $b$: $C_a$ shares only the $100\ \Omega$, and $C_b$ shares all $300\ \Omega$. So $T_b = 100 \cdot 10 + 300 \cdot 20 = 7000\ \Omega\cdot\text{fF} = 7$ ps.

## Computing it in one tree traversal

Summing over every capacitor for every sink is quadratic. A timer instead walks the tree twice. A post-order pass computes the downstream capacitance $C_{down}(v)$ of each node, the total capacitance at and below it. A pre-order pass then sets $T(v) = T(\text{parent}) + R_{edge}\, C_{down}(v)$. Each edge resistance charges everything beneath it, which is the same sum regrouped. The work is linear in the size of the net. Units matter: one ohm times one femtofarad is one femtosecond, so $1000\ \Omega\cdot\text{fF}$ is one picosecond.
