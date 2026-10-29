# Contact and via resistance

Current enters a transistor through a contact to silicide on its source or drain and climbs to the wiring through a chain of vias. At advanced nodes these vertical connections contribute as much resistance as the transistor channel itself for short paths. Their values, and the effect of using several in parallel, feed cell characterization, extraction and power-grid analysis.

## Specific contact resistivity

The resistance of a metal-semiconductor interface scales inversely with its area: $R_c = \rho_c/A$, where $\rho_c$ is the specific contact resistivity in ohm cm^2. Good silicide contacts reach about $10^{-9}$ ohm cm^2. A contact interface of 20 nm by 20 nm has area 400 nm^2 $= 4 \times 10^{-12}$ cm^2, so $R_c = 10^{-9}/(4 \times 10^{-12}) = 250$ ohm. Halving the area doubles it. This is why contact area shrinkage hurts so much, why trench contacts run the full width of the active region, and why wrap-around contacts on epitaxial source-drain regions are attractive.

## Via stacks in series and parallel

A route from metal 1 to metal 6 crosses five via levels in series; each level may use one via or several in parallel. Total resistance is the sum over levels of $R_{via}/k$, where $k$ is the number of parallel vias at that level, ignoring the short metal landings. Levels of 40 ohm single, 30 ohm doubled and 20 ohm doubled give $40 + 15 + 10 = 65$ ohm. Lower vias are smaller and more resistive than upper ones. Power grids use arrays of vias for this reason, and DFM flows double single vias where space allows, which lowers both resistance and the probability of an open.

## When contacts dominate

For a short local connection, the path from one transistor to another may be two contacts, a few vias and a short wire. With 150 ohm contacts, 60 ohm of vias and 100 ohm of wire, contacts carry $300/460 = 65\%$ of the resistance. Wire length optimization cannot touch that share; only process changes (better silicide, larger contact area, new MOL metals) or architecture changes (fewer contacts in series, backside contacts) can. Timing extraction must therefore model MOL resistance in detail, and DTCO studies treat contact resistivity as a primary knob alongside pitch.

## Going deeper

Schroder, "Semiconductor Material and Device Characterization", chapter 3 on contact resistance; the IRDS interconnect chapter on via and contact resistance trends.
