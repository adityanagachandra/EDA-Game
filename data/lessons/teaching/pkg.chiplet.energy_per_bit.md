# Energy per bit and link power

Energy per bit is the energy an interface spends to move one bit from one die to another, usually quoted in picojoules per bit. Multiplied by bandwidth it gives link power, which can rival the power of the logic it connects. It is the main reason heterogeneous integration favours short, dense, simple die-to-die links over off-package SerDes, and it is computed from wire and pad capacitance.

## Link power

Power is energy per bit times bit rate: $P = E_b \times BW$. One picojoule per bit at one terabit per second is one watt, since $10^{-12} \times 10^{12} = 1$. A 0.5 pJ/bit die-to-die link carrying 8 Tb/s draws 4 W. Count both directions if both are active. These watts add to the package's thermal budget and to its power-delivery current, and they sit at the die edge, so they also create local hot spots under the PHYs. Architects therefore budget link power with logic power, and data-movement-heavy designs such as accelerators can spend a large share of their power on moving data rather than computing.

## On-package against off-package

Off-package SerDes over a board or cable spends energy on equalization, clock recovery and driving lossy channels, typically several picojoules per bit, while short on-package parallel links can stay well below one picojoule per bit. Moving 4 Tb/s at 6 pJ/bit costs 24 W, and at 0.4 pJ/bit 1.6 W, saving 22.4 W. This saving, more than raw bandwidth, justifies putting memory and compute in one package. The comparison must be like for like: some quoted figures count only the PHY and others include the adapter and protocol logic, and quoted values come from different process nodes and reaches, so a designer computes them for the actual channel.

## Energy from capacitance

For a full-swing, unterminated link, the dominant energy is charging the wire, pads and receiver input. Each 0-to-1 transition draws $C V^2$ from the supply. With random data, a given bit is a 0-to-1 transition with probability one quarter, so the average energy per bit is $C V^2/4$. A channel of 150 fF at 0.75 V spends $150 \times 0.5625/4 = 21.1$ fJ per bit, or 0.021 pJ/bit for the wire alone; drivers, receivers, clocking and the adapter add more. The formula shows the levers: lower capacitance, from shorter wires and smaller pads such as hybrid bonds, and lower swing, which helps quadratically but needs more sensitive receivers.

## Going deeper

Horowitz, "Computing's energy problem", ISSCC (2014); Lau, "Chiplet Design and Heterogeneous Integration Packaging" (2023).
