# Differential pairs

A differential pair carries one signal as the voltage difference between two traces. The receiver rejects noise that appears equally on both, the return current largely flows in the partner trace, and the fields of the two traces partly cancel, so pairs radiate less. Board interfaces specify a differential impedance, typically 85, 90 or 100 ohms, and the layout tool must translate it into trace width and gap. That translation goes through the odd and even mode impedances, and a TDR measurement goes back through them.

## From odd mode to differential impedance

A pair's differential impedance is twice its odd-mode impedance, $Z_{diff} = 2 Z_{odd}$, so a 100 ohm interface needs each trace to present 50 ohms in the odd mode. That is not the same as 50 ohms single-ended. When traces are tightly coupled, $Z_{odd}$ is well below the uncoupled $Z_0$, so each trace must be made narrower than a 50 ohm single-ended trace to bring $Z_{odd}$ back up to 50. A field solver sweeps width and gap to hit the target; the board tool stores the pair rule (width, gap, target) in a net class. Example: a solver reports $Z_{odd} = 46.1$ ohms for a candidate geometry, giving $Z_{diff} = 92.2$ ohms, close to a 90 ohm target and short of a 100 ohm one.

## Common-mode impedance

The common-mode impedance, $Z_{cm} = Z_{even}/2$, is what common-mode noise and the common-mode part of a signal see. It matters because receivers and connectors often specify a common-mode range or termination, and because unterminated common mode rings and radiates. For a loosely coupled pair, $Z_{even} \approx Z_{odd} \approx Z_0$, so $Z_{cm} \approx Z_0/2$ and $Z_{diff} \approx 2 Z_0$. Example: a pair with $Z_{odd} = 41.8$ and $Z_{even} = 44.8$ ohms is loosely coupled; $Z_{cm} = 22.4$ ohms. Tighter coupling raises $Z_{cm}$ and lowers $Z_{odd}$, separating the two numbers. A complete pair specification therefore names both, or names $Z_{diff}$ together with a geometry from which the solver derives $Z_{cm}$.

## Coupling from measured impedances

A TDR in differential mode measures $Z_{diff}$, and in common mode measures $Z_{cm}$. From them, $Z_{odd} = Z_{diff}/2$ and $Z_{even} = 2 Z_{cm}$, and the coupling coefficient $k = (Z_{even} - Z_{odd})/(Z_{even} + Z_{odd})$ summarizes how tightly the traces interact. Example: $Z_{diff} = 100$ and $Z_{cm} = 30$ ohms give $Z_{odd} = 50$, $Z_{even} = 60$ and $k = 10/110 \approx 0.091$. Small $k$ (wide gap) makes the impedance insensitive to gap variation and lets each trace be designed near 50 ohms single-ended, at the cost of area; large $k$ saves area and improves noise immunity between the pair's own traces, but etch variation in the gap moves the impedance more. The choice is a parameter a constraint system exposes, not a law.

## Going deeper

Bogatin, "Signal and Power Integrity Simplified" (3rd ed., 2018), chapter 11, on differential pairs.
