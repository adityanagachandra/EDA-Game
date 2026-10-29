# Rush current and wake-up

When a power-gated domain wakes, its virtual rail must be recharged from near zero. Turning on every switch at once would draw a large current spike from the shared grid, drooping the always-on neighbours and possibly corrupting their state. Wake-up is therefore staged: switches turn on in a controlled sequence. Sizing that sequence is a charge computation under a current limit, and it sets the domain's wake-up latency.

## Average rush current

Recharging capacitance $C$ by $\Delta V$ in time $\Delta t$ needs an average current $I = C\,\Delta V/\Delta t$. For 5.5 nF charged to 0.75 V in 140 ns, $I = 4.125$ nC/140 ns $= 29.5$ mA. Units: nanocoulombs per nanosecond are amperes. This current flows through the always-on grid, so it appears as dynamic drop on neighbours, and through the package, where its rate of change creates $L\,di/dt$ noise.

## The shortest safe ramp

If the rush current must stay below $I_{max}$ on average, the ramp must last at least $t = CV/I_{max}$. For 6 nF at 0.6 V and a 500 mA limit, $t \ge 3.6/500$ nC/mA, which is 7.2 ns. A real ramp is not linear: with all switches on, the current is largest at first, when the rail is far from VDD, so peak current exceeds the average and the limit must account for it. Wake-up latency, which software sees, is at least this ramp plus the time to restore state.

## Staged switch turn-on

Switches are daisy-chained: each group's enable passes through a delay to the next, so the rail charges over roughly $k\,t_d$ for $k$ groups spaced $t_d$ apart. Keeping $CV/(k\,t_d) \le I_{max}$ needs $k \ge CV/(I_{max}t_d)$; for 2.5 nF at 0.8 V, a 170 mA limit and 5 ns spacing, $k \ge 2.35$, so 3 groups. A common refinement first turns on a few weak (high-resistance) switches to charge the rail slowly, then the strong ones once it is nearly at VDD, which limits the peak while keeping the final on-resistance low. The power controller sequences this with isolation and retention.

## Going deeper

Keating et al., "Low Power Methodology Manual" (2007); Shi and Howard, "Challenges in sleep transistor design and implementation in low-power designs", Design Automation Conference (2006).
