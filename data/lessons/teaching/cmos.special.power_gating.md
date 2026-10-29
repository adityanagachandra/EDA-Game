# Power switches and power gating

Power gating cuts a block's leakage by disconnecting it from the supply when it is idle. A network of switch cells, header pMOS devices between the true supply and a virtual supply rail (or footer nMOS devices on ground), is placed across the block. Sizing that network is a resistance and charge calculation that implementation tools perform when they insert and verify switches.

## Switch on-resistance and IR drop

When the block is on, its current flows through the switches, which act as resistors. With $N$ identical switches in parallel, each of on-resistance $R_{on}$, and a block current $I$ shared evenly, the drop across the network is
$$\Delta V = \frac{I R_{on}}{N}.$$
Example: 150 mA through 2,000 switches of 1,200 ohms each gives $150 \times 1200 / 2000 = 90$ mV, since milliamperes times ohms is millivolts. That drop comes straight off the block's effective supply, slowing every cell in it, so it must be included in timing, usually by analyzing the block at its virtual-rail voltage or with IR-drop-aware timing. Current is not shared evenly in practice; switches near current hot spots carry more, which is why switch placement is verified with a power-grid analysis.

## Sizing the switch network

Turning the formula around, the number of switches needed to keep the drop within a budget $\Delta V_{max}$ is
$$N \ge \frac{I R_{on}}{\Delta V_{max}},$$
rounded up. For 150 mA, 1,200 ohms and a 30 mV budget, $N \ge 6000$. More switches mean less drop but more area and more leakage through the switches themselves when off, since an off switch still leaks; the switch device is usually a high-Vt transistor to keep that leakage small. Switches are placed in columns or a checkerboard across the domain so that every cell is near one, and the power grid connects the true supply above to the virtual supply below. This sizing is done per domain, against the peak current the block can draw.

## Rush current and wake-up time

When a gated block turns on, its virtual rail and all its internal capacitance charge from near 0 V. Turning all switches on at once would draw a huge rush current, causing supply droop that can upset neighbouring active blocks. Switch networks are therefore enabled in a daisy chain: a weak first stage charges the rail slowly, and the strong switches turn on later. If the charging current is limited to $I_{max}$, the wake-up time is roughly
$$t = \frac{C V_{DD}}{I_{max}}.$$
With 8 nF of domain capacitance, 0.8 V and a 40 mA limit, $t = 160$ ns. Wake-up latency is the price of limiting rush current, and it sets how often power gating is worth using: a block must stay off long enough for the leakage saved to exceed the energy spent recharging it.

## Going deeper

Keating, Flynn, Aitken, Gibbons and Shi, "Low Power Methodology Manual" (2007), chapters on power gating; IEEE 1801 (Unified Power Format).
