# Sizing header and footer switches

Power switches sit in series with the supply of a gated domain, so their on-resistance adds IR drop whenever the domain runs, and their off-state leakage remains when it sleeps. The implementation tool inserts many small switch cells in parallel, distributed as a grid or a ring through the domain. Choosing how many is a direct trade between run-time voltage drop and sleep-mode leakage and area.

## Drop across the switch network

$N$ identical switches with on-resistance $R_{on}$ in parallel present $R_{on}/N$. A domain drawing 130 mA through 5000 switches of $1450\ \Omega$ sees $130 \cdot 0.29 = 37.7$ mV of drop across them. This drop adds to the grid's own IR drop, so the switch network's budget is carved out of the domain's total. Distributing switches through the domain, rather than at its edge, also shortens the virtual-rail paths and keeps the drop uniform.

## Switch count for a budget

Requiring $I R_{on}/N \le \Delta V$ gives $N \ge I R_{on}/\Delta V$. For a peak of 130 mA, $R_{on} = 1950\ \Omega$ and a 30 mV switch budget, $N \ge 8450$. The peak current, not the average, sets the count, because the voltage must hold through the worst cycle. Switch cells come in sizes; a tool may mix strong and weak switches, which also helps control the rush current at wake-up.

## The leakage floor

Off switches still leak. With 1620 switches leaking 2 nA each from 0.95 V, the domain still draws $1620 \cdot 2 \cdot 0.95 = 3.08$ µW in sleep. More switches mean less run-time drop but more sleep leakage and area, and they lengthen the break-even sleep time. High-threshold devices are used for switches to cut this leakage, at the cost of higher on-resistance per unit width, which the sizing must then absorb.

## Going deeper

Keating et al., "Low Power Methodology Manual" (2007); Shi and Howard, "Challenges in sleep transistor design and implementation in low-power designs", Design Automation Conference (2006).
