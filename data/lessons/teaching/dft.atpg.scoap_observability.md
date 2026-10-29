# SCOAP observability

Controllability says how hard it is to set a line; observability says how hard it is to see it, that is, to propagate a change on the line to a primary output. SCOAP computes combinational observability $CO$ in one backward pass, using the controllabilities already computed. The pair together gives a cheap estimate of how hard each stuck-at fault is to detect.

## Backward propagation

Primary outputs have $CO = 0$. To observe an input of an AND or NAND gate, the output must be observed and every other input must be at the non-controlling value 1:
$$CO(i) = CO(out) + \sum_{j \neq i} CC1(j) + 1.$$
For OR and NOR the other inputs must be 0, so the sum uses $CC0$. An inverter or buffer adds 1. For a 2-input XOR, the other input may take either value, so it adds $\min(CC0, CC1)$ of the other input plus 1. Gates are processed in reverse topological order, so each output's observability is final before its inputs are computed.

## Fanout stems

A stem with several branches is observable through whichever branch is easiest: $CO(stem) = \min_b CO(b)$. For example, if c feeds $g_1 = $ NAND(b, c) with $CO(g_1) = 4$ and $CC1(b) = 1$, the branch into $g_1$ has observability $4 + 1 + 1 = 6$; if c also feeds an inverter whose output has $CO = 1$, that branch has $1 + 1 = 2$, so $CO(c) = 2$. Taking the minimum is optimistic when the cheap branch is blocked by the very assignments needed elsewhere, again a consequence of ignoring correlation.

## Detection difficulty

A stuck-at fault needs activation and observation, so SCOAP estimates the difficulty of detecting $x$ stuck-at-0 as $CC1(x) + CO(x)$, and of $x$ stuck-at-1 as $CC0(x) + CO(x)$. Large values flag faults that are likely to be aborted by ATPG or missed by random patterns. Test-point insertion uses exactly this: an observation point at a line with large $CO$ resets its observability to near zero, and a control point fixes a large $CC$. Sequential SCOAP adds sequential depth counts for designs without full scan.

## Going deeper

Goldstein (1979); Bushnell and Agrawal (2000), chapter 6.
