# Parasitic delay

A gate's delay has a part that depends on its load and a part that does not. The load-independent part, the parasitic delay, comes from the gate charging its own diffusion capacitance. It is the reason upsizing a gate cannot make it arbitrarily fast and the reason a chain of buffers has an optimal length. Every delay model, from logical effort to a Liberty table's zero-load intercept, contains it.

## Self-loading from output diffusion

Every transistor whose drain touches the output node adds diffusion capacitance there. That capacitance must be charged and discharged through the gate's own resistance, just like the external load. In the RC model, a gate with resistance $R$, output diffusion $C_{par}$ and load $C_L$ has delay $R(C_{par} + C_L)$; the term $R C_{par}$ is present even with no load. If the gate is scaled up by a factor $k$, its resistance becomes $R/k$ and its diffusion $kC_{par}$, so $R C_{par}$ is unchanged: parasitic delay is independent of gate size. In logical-effort units it is written $p$, with the inverter's value $p_{inv}$ close to 1 when diffusion and gate capacitance per unit width are similar.

## Computing p from widths

Normalize to the inverter: with nMOS width 1 and pMOS width $\gamma$, the inverter's output diffusion is proportional to $1 + \gamma$. A gate's parasitic delay is its total output diffusion width divided by $1 + \gamma$, times $p_{inv}$. A NAND$n$ output touches $n$ parallel pMOS of width $\gamma$ and one nMOS of width $n$: total $n\gamma + n = n(1 + \gamma)$, so $p = n\, p_{inv}$. A NOR$n$ output touches $n$ nMOS of width 1 and one pMOS of width $n\gamma$, again $n(1 + \gamma)$, so $p = n\, p_{inv}$. A gate whose output touches nMOS diffusion of total width 3 and pMOS of total width 4, with $\gamma = 2$, has $p = 7/3 \approx 2.33$. These estimates ignore internal stack nodes, which add more.

## Why parasitic delay limits sizing

Stage delay in logical-effort units is $d = gh + p$, with $h$ the electrical effort. When a gate drives a light load, $gh$ is small and $p$ dominates: a NAND3 with $g = 5/3$, $p = 3$ and $h = 1$ spends $3/4.67 = 64\%$ of its delay on its own diffusion. When it drives a heavy load, the effort term dominates and $p$ is minor. Two consequences follow. Upsizing a gate reduces $gh$ (its input grows, so its own $h$ shrinks) but never $p$, so it has diminishing returns. And adding a stage to a path costs at least that stage's $p$, which is what makes the optimal number of stages finite. In Liberty, the parasitic delay appears as the near-zero-load delay at fast input slew.

## Going deeper

Sutherland, Sproull and Harris, "Logical Effort: Designing Fast CMOS Circuits" (1999), chapter 1; Weste and Harris, "CMOS VLSI Design" (4th ed.), section 4.5.
