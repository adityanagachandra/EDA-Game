# Board bring-up and hardware debugging

Bring-up is the first power-on and functional test of a new board, where design errors that every check missed finally appear. It follows a dependency order: rails, then clocks and resets, then debug access, then firmware and interfaces, so that each failure is isolated to one layer. Debugging relies on measurement discipline as much as on design knowledge. For a tool designer, bring-up is the ground truth that shows which checks and models were missing upstream.

## A safe power-up order

Before power, measure the **resistance from each rail to ground**: a solder bridge or rotated part often shows as a short, cheap to find unpowered. Power first from a **current-limited** bench supply, so a fault draws limited current instead of burning traces; watch each rail's voltage and current against expectations and check sequencing. Then confirm clocks and resets, then **JTAG**: reading every device's IDCODE proves the chain is wired, devices are powered and the right parts are fitted. Only then load firmware and test interfaces, starting with memory. The order follows dependency: an interface failure on a board with a bad rail wastes days. Each stage's expected values (rail currents, clock frequencies, IDCODEs) should be written down before power-on.

## Measurement and localization

Debugging needs trustworthy measurements. A scope probe with a long ground lead forms an inductive loop that rings with the probe capacitance, so a short ground spring is used before blaming the board. To find a short on a rail, inject a limited current and look for the hot spot with a thermal camera, or follow the voltage-drop gradient along the plane with a milliohm meter. A **TDR** locates impedance discontinuities by time: a dip is a capacitive feature (via, pad, wide section), a rise an inductive one (neck-down, gap in the reference plane). Open-drain buses such as I2C with slow edges point to pull-ups too large for the bus capacitance, since rise time is set by $RC$.

## Closing the loop

A failure that appears only at temperature, only in the enclosure or only on some boards is a **margin** problem, in timing, signal integrity or power, not a connectivity error, which would fail everywhere. Those lead back to the models: re-run SI or timing at corners, or compare memory-training results across temperature. Every fix found at bring-up, such as a bodge wire, a changed resistor or a cut trace, must go back into the schematic and layout as an ECO for the next revision, with its reason recorded, or the bug returns. Test points and debug headers planned during design are what make bring-up possible at all on BGA-heavy boards. Bring-up findings are the best feedback for improving rule decks, review checklists and agent skills.

## Going deeper

Bogatin, "Signal and Power Integrity Simplified" (3rd ed., 2018), chapter on measurement. Oscilloscope vendor application notes on probing high-speed signals.
