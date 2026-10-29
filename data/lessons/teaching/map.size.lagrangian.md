# Lagrangian relaxation for gate sizing

Sizing every gate of a netlist to meet timing at minimum area is a large constrained problem. Lagrangian relaxation moves the timing constraints into the objective with one multiplier per timing arc, which decouples the problem into one tiny subproblem per gate. The multipliers then act as a criticality measure, and their update is the outer loop. It is the backbone of several academic and industrial sizers.

## Multipliers and flow conservation

Minimize area subject to arrival constraints $a_u + d_{uv} \le a_v$ on every arc and $a_v \le T$ at outputs. Attach multiplier $\lambda_{uv} \ge 0$ to each arc. Chen, Chu and Wong showed that at the optimum the multipliers satisfy **flow conservation**: at every gate, the sum on incoming arcs equals the sum on outgoing arcs, like a flow from outputs back to inputs. Using this, the arrival-time variables drop out of the Lagrangian. Example: a gate with outgoing multipliers 3.4, 1.3 and 1.0 and incoming multipliers 1.0, 0.8 and one unknown needs the unknown to be 5.7 - 1.8 = 3.9.

## The local subproblem

For fixed multipliers the Lagrangian separates by gate. Gate i's size x appears in its own area, in its own delay (weighted by $\mu_i$, the sum of multipliers on its arcs) and in its driver's delay through its input capacitance (weighted by $\nu_i$). With delay $r C_L / x$ and input capacitance $c x$, the subproblem is $\min_x a x + \mu r C_L / x + \nu c x$, whose solution is
$$x^* = \sqrt{\frac{\mu r C_L}{a + \nu c}}.$$
With a = 3, $\mu$ = 0.8, r = 1.7 kOhm, $C_L$ = 47 fF, $\nu$ = 1.8 and c = 3 fF: $x^* = \sqrt{63.92/8.4} = 2.76$. Sweeping this local rule over all gates until sizes settle solves the inner problem.

## Updating the multipliers

The outer loop adjusts multipliers by subgradient steps: $\lambda \leftarrow \max(0, \lambda + \theta (a - q))$, where $a - q$ is the arc's timing violation (arrival minus required). A multiplier of 2.5 with step 0.01 per ps on an arc violating by 33 ps becomes 2.83; an arc with slack loses weight and may fall to zero. After each update the multipliers are projected back onto the flow-conservation space. Step sizes shrink over iterations to guarantee convergence of the dual. For discrete libraries the local rule picks the best available size instead of $x^*$, which keeps the method practical.

## Going deeper

Chen, Chu and Wong, "Fast and exact simultaneous gate and wire sizing by Lagrangian relaxation", IEEE TCAD 1999; Ozdal, Burns and Hu, "Algorithms for gate sizing and device parameter selection for high-performance designs", IEEE TCAD 2012.
