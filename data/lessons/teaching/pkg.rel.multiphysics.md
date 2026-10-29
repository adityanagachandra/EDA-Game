# Electro-thermal multiphysics coupling

Electrical, thermal and mechanical behaviour are coupled. Current heats conductors, heat raises resistance, higher resistance changes the heat, and temperature drives both stress and electromigration. A multiphysics simulation solves these fields together, either in one monolithic system or by iterating between single-physics solvers. The simplest case, a resistor heating itself, shows how the coupling converges, and how it can run away.

## Iterating the coupling

A conductor with resistance $R(T) = R_0(1 + \beta(T - T_0))$ carries current $I$ and has thermal resistance $\theta$ to an ambient $T_a$, with $T_0 = T_a$. Its temperature satisfies $T = T_a + \theta I^2 R(T)$. Fixed-point iteration alternates the two solvers: compute power from the current temperature, then temperature from the power. With $R_0 = 0.05$ ohm, $I = 3$ A, $\theta = 40$ K/W, $\beta = 0.0039$ per K and $T_a = 50$ °C, iteration one gives $50 + 40 \times 9 \times 0.05 = 68$ °C, and iteration two uses $R = 0.05(1 + 0.0039 \times 18) = 0.0535$ ohm to give 69.26 °C. This loose coupling is how many tools combine an electrical and a thermal solver.

## The converged temperature

Because the equation is linear in $T$, it has a closed form. With $T_0 = T_a$, the rise is $\Delta T = \theta I^2 R_0/(1 - g)$, where $g = \theta I^2 R_0 \beta$ is the loop gain. In the example $g = 40 \times 9 \times 0.05 \times 0.0039 = 0.0702$ and $\Delta T = 18/0.9298 = 19.36$ K, so the converged temperature is 69.36 °C. The iteration converges because $g < 1$, and its error shrinks by a factor $g$ each step, so after two steps it is within 0.1 K. That rate is the general lesson: loose coupling converges fast when the coupling is weak and slowly when it is strong, when a monolithic Newton solve is better.

## Thermal runaway

When the loop gain reaches 1, the denominator vanishes and there is no steady state: each degree of heating adds more than a degree, and the temperature runs away until something melts or a limit intervenes. The critical current is $I_c = 1/\sqrt{\theta R_0 \beta}$; for the example $1/\sqrt{40 \times 0.05 \times 0.0039} = 11.32$ A. Real systems include other feedbacks: leakage power grows exponentially with temperature, so a chip can run away even when wire heating is benign. Multiphysics signoff checks that the coupled solution exists and converges with margin, because a solver that stops after a fixed number of iterations can report a plausible temperature for a system that has no steady state.

## Going deeper

Pedram and Nazarian, "Thermal modeling, analysis, and management in VLSI circuits: principles and methods", Proceedings of the IEEE (2006); Zienkiewicz and others, "The Finite Element Method", on coupled problems.
