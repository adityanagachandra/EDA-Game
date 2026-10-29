# Current source stamps

Independent current sources are the simplest stamps in MNA: they touch only the right-hand side vector $b$. Yet sign errors in current-source stamps are among the most common bugs in hand-written simulators and in hand analysis, because SPICE's card convention runs opposite to the intuitive one. This concept fixes the convention, then uses it to read off right-hand-side entries for small decks.

## Injected current on the right-hand side

With KCL written as currents leaving through elements on the left and source currents on the right, $G v = b$, the entry $b_i$ is the net current injected into node $i$ by sources. A source that pushes $I$ into node $a$, drawing it from node $b$, adds $+I$ to $b_a$ and $-I$ to $b_b$. A source from ground into node $a$ adds only $+I$ at $a$, since the ground row is deleted. At a node fed by 1 mA from ground and drained by 3.5 mA to ground, with another 1.5 mA pushed in from a neighbouring node, $b = 1 - 3.5 + 1.5 = -1$ mA. The matrix is untouched: an ideal current source has infinite internal resistance, so it adds no conductance.

## The SPICE current-source convention

SPICE's card I1 n+ n- val means positive current flows from n+ through the source to n-. The source therefore takes current out of n+ and delivers it into n-, so it stamps $-I$ at $b_{n+}$ and $+I$ at $b_{n-}$. This looks backwards to anyone who reads n+ as the end the current comes out of, which is exactly the mistake to avoid. For I1 2 0 1m, node 2 gets -1 mA: the source pulls a milliampere out of node 2 to ground. For I2 0 3 2m, node 3 gets $+2$ mA. When reading a deck, rewrite each card as "takes I out of n+, puts it into n-" before stamping.

## Sources that touch ground

Summing $b$ over all non-ground nodes gives a check like the one for $G$. A source between two non-ground nodes adds $+I$ and $-I$ and cancels. Only sources with one terminal at ground survive, so the sum is the net current delivered from ground into the circuit, which by KCL must leave again through resistors to ground. For the cards I1 1 2 2m, I2 0 1 2m and I3 1 2 0.5m, only I2 touches ground, delivering 2 mA into node 1, so the entries of $b$ sum to 2 mA. In a simulator the same idea appears in reverse: the ground row is the one equation dropped, and it is implied by the others.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 2; the ngspice manual, section on independent sources.
