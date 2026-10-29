# Signal probability and activity propagation

Power analysis needs the switching activity of every net, but simulation gives it for only some nets, and early in a project for none. Probabilistic propagation fills the gap: given the probability that each input is 1, compute the probability for every internal node, and from it the activity. It is the engine inside vectorless power analysis and the basis of activity estimates for IR drop and clock gating.

## Signal probability through gates

The signal probability $p$ of a net is the fraction of time, or of cycles, it is 1. With independent inputs, gate functions map probabilities directly: AND multiplies, $p_{ab} = p_a p_b$; OR complements the product of complements, $1 - (1 - p_a)(1 - p_b)$; NOT gives $1 - p$; XOR gives $p_a(1 - p_b) + p_b(1 - p_a)$. For $y = (a \wedge b) \vee c$ with $p_a = p_b = 0.45$ and $p_c = 0.9$, $p_{ab} = 0.2025$ and $p_y = 1 - 0.7975 \cdot 0.1 = 0.9203$. One topological pass computes every node in linear time.

## From probability to activity

The activity $\alpha$ in $P = \alpha C V^2 f$ is the probability of a 0-to-1 transition in a cycle. If a node's values in successive cycles are independent (temporal independence), a rising transition needs a 0 followed by a 1, so $\alpha = (1 - p)p$. For $p_y = 0.7675$, $\alpha = 0.2325 \cdot 0.7675 = 0.178$. The activity peaks at $p = 0.5$, where $\alpha = 0.25$, and vanishes for nets stuck near 0 or 1. Some tools count toggles in both directions instead, $2p(1 - p)$, with a factor of one half in the power formula; the two conventions give the same power.

## Node power and the limits of independence

With $\alpha$ known, a node's power is $\alpha C V^2 f$; femtofarads times volts squared times gigahertz give microwatts. For $p_y = 0.2865$, $\alpha = 0.2044$, $C = 4$ fF, $V = 0.75$ V and $f = 1.5$ GHz, the node burns 0.69 µW. The method's weakness is its assumptions. Reconvergent fanout makes gate inputs correlated, so multiplying probabilities is wrong; glitches add transitions that a per-cycle model cannot see; and real signals are temporally correlated, a counter's low bit toggles every cycle while its high bit rarely does. Tools mitigate with correlation coefficients, partial BDD-based exact computation, or by annotating simulated activity where available.

## Going deeper

Najm, "A survey of power estimation techniques in VLSI circuits", IEEE Transactions on VLSI Systems (1994); Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003).
