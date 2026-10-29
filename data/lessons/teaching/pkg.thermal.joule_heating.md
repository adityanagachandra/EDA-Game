# Joule heating and self-heating

Joule heating is the power $I^2 R$ dissipated when current flows through resistance. In packages it heats power bumps, balls, TSVs and planes; on chip it heats wires; and in transistors the channel's own power heats the device, called self-heating. Each local temperature rise feeds electromigration and resistance models, so even small heat sources matter for reliability signoff.

## Power in a conductor

A conductor of resistance $R$ carrying current $I$ dissipates $P = I^2 R$. A power bump of 6 milliohm carrying 0.15 A dissipates $0.0225 \times 0.006 = 0.135$ mW. Across thousands of bumps the total is modest, but a wire of 50 ohm carrying 2 mA dissipates 0.2 mW in a tiny volume. The quadratic dependence on current is the key: doubling current quadruples the heat, which is why current crowding at the entrance of a bump or via, where local current density is several times the average, produces local hot spots that limit electromigration lifetime.

## Temperature rise of a wire

A heat source of power $P$ connected to the substrate through thermal resistance $R_{th}$ rises $\Delta T = P R_{th}$. On-chip wires sit in low-conductivity dielectric, so their thermal resistance to the substrate can be large, thousands of K/W per wire segment. A segment dissipating 0.2 mW with 5000 K/W rises 1 K; one at 2 mW rises 10 K. Upper metal layers are farthest from the silicon, have the most dielectric beneath them, and carry the largest currents, so they heat most. Low-k dielectrics conduct heat worse than silicon dioxide, which makes wire heating worse as interconnect scales, and electromigration rules include a wire temperature that accounts for it.

## Transistor self-heating

A transistor dissipates $P = I_D V_{DS}$ in its channel. In planar bulk devices the heat flows easily into the substrate, but FinFETs and nanosheets are narrow structures surrounded by oxide, so their thermal resistance is higher and the channel runs hotter than the substrate. With $R_{th}$ in K per milliwatt, a device dissipating 200 uW with 30 K/mW rises 6 K. Self-heating lowers mobility and so drive current, and it raises the temperature seen by nearby wires and by aging models. Compact models such as BSIM include a self-heating thermal node, but they need the right $R_{th}$, which depends on layout and neighbours, a known modelling gap.

## Going deeper

Banerjee and Mehrotra, "Global (interconnect) warming", IEEE Circuits and Devices Magazine (2001); the BSIM-CMG manual on self-heating.
