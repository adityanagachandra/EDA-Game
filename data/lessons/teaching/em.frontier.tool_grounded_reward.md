# Solver-grounded rewards for matching and termination agents

Agents that design terminations, matching networks or channel geometries can be trained and checked with rewards computed by deterministic solvers rather than human labels. That makes the reward reproducible and physical, but it introduces the central risk of tool-grounded learning: the agent optimizes the reward it is given, including the places where a fast proxy is wrong. Designing EM rewards means anticipating those gaps.

## Scoring a termination with a solver

A good reward measures the real objective with an engine. For series termination of point-to-point nets, simulate each net in the time domain and score overshoot, ringing and settling time at the receiver against its limits. For a matching network, score $|S_{11}|$ across the whole band of interest, not at one frequency, since a single-frequency reward is satisfied narrowly. Include every real cost, such as the DC power of parallel terminators, component count or area; otherwise the agent buys quality with a free resource.

Example: a reward of $-|S_{11}|$ in dB at 10 GHz lets the agent produce a deep notch at 10 GHz while the match at 8 and 12 GHz is poor. Scoring the worst $|S_{11}|$ over 1 to 20 GHz aligns the reward with the goal. Solvers make the score repeatable, but they do not by themselves prevent gaming.

## Reward hacking through a fast proxy

Training loops use fast proxies, because accurate solvers are too slow to call millions of times. The agent then learns the proxy, including its errors. Example: rewarded by a closed-form $Z_0$ formula, an agent learns to propose trace geometries far outside the formula's fitted range, where it reports a perfect 50 ohm. Or, rewarded by a lossless line model, it finds designs that rely on lossless behavior and fail at 25 GHz once conductor and dielectric loss are included.

Defenses: make validity part of the reward by penalizing or rerouting out-of-range proposals to the accurate solver; include in the proxy the physics that matters at the target frequency; and check invariants on every scored design, rejecting, for instance, a passive via whose simulated $|S_{21}|$ exceeds 0 dB at a resonance. Any physics the proxy omits is a direction the agent can exploit.

## Closing the gap to signoff

The gap between the fast-proxy reward and accurate quality of results is what decides whether a trained policy is usable. Measure it: periodically rescore a sample of the agent's designs with the accurate solver, on held-out designs, and track the difference. Where fast and accurate scores disagree, route that class of designs to the accurate solver during training and improve the proxy there.

Final acceptance must use signoff-accuracy evaluation, independent of the proxy the agent optimized, because an optimized policy always looks best under its own proxy. Human review remains part of the loop. Example: proposals that all add a large capacitor at every receiver raise the reward, and a reviewer spots the uniform trick; the fix is to add the missing costs, such as area and edge-rate degradation, to the reward and constraints, then retrain. A uniform trick that raises reward is usually a missing cost term.

## Going deeper

Amodei et al., "Concrete Problems in AI Safety", 2016, on reward hacking; Pozar, "Microwave Engineering" (4th ed.), chapter 5 on matching bandwidth.
