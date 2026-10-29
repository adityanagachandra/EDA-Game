# Power gating

Power gating cuts a block off from its supply when it is idle, removing leakage almost entirely. Header switches (PMOS, between VDD and a virtual VDD) or footer switches (NMOS, between virtual VSS and VSS) connect the block when it is on. The block loses its state when off, and waking it costs energy and time. Deciding when gating pays is a break-even computation, and implementing it touches the grid, isolation, retention and the control sequence.

## The wake-up energy

During sleep the virtual rail discharges toward zero. On wake-up the switches recharge its capacitance $C$ to $V$, drawing charge $Q = CV$ from the supply at voltage $V$, so the energy drawn is $QV = CV^2$: half is stored on the capacitor and half is dissipated in the switch resistance. For 5.5 nF charged to 0.6 V that is $5.5 \cdot 0.36 = 1.98$ nJ (nanofarads times volts squared give nanojoules). Restoring state and re-initialising the block add more.

## Break-even sleep time

Sleep saves leakage power minus the switches' own off-state leakage, $P_{leak} - P_{sw}$, for as long as it lasts. It pays once that saving exceeds the wake-up energy: $t_{be} = E_{wake}/(P_{leak} - P_{sw})$. For 480 µW of block leakage, 24 µW of switch leakage and 3.84 nJ per wake-up, $t_{be} = 3.84/456$ nJ/µW, which is 8.42 µs. Power-management firmware or hardware uses this number: idle periods shorter than break-even should not gate. Header PMOS switches are common because they keep the ground continuous across domains.

## Net energy over an interval

For a sleep of length $t$, the net saving is $(P_{leak} - P_{sw})t - E_{wake}$. With 378 µW saved and a 2.34 nJ wake-up, a 390 µs sleep saves $147.4 - 2.3 = 145.1$ nJ, while a 5 µs sleep would lose energy. Gating also has costs that this sum hides: the wake-up latency delays the next task, the rush current during wake-up disturbs neighbouring domains, and the switches' on-resistance adds IR drop whenever the block runs. These are the next concepts.

## Going deeper

Keating et al., "Low Power Methodology Manual" (2007); Shi and Howard, "Challenges in sleep transistor design and implementation in low-power designs", Design Automation Conference (2006).
