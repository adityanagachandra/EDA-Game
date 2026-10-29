# Isolation and level-shifting requirements

Every signal that crosses a power-domain boundary must be checked against the power states. If its driver can be off while its receiver is on, it needs isolation; if the two domains run at different voltages, it needs level shifting; sometimes both, in one enable level-shifter cell. Tools derive these requirements automatically from the UPF, and verification checks that no crossing was missed. The rules are short but each has a reason rooted in transistor behaviour.

## When a crossing needs protection

An off domain's outputs float or drift. If the receiving logic is powered, a floating input can sit near the switching threshold, turning on both transistors of the input stage and causing crowbar current, and it can be read as a random value. So a crossing needs isolation whenever some legal state has the source off and the sink on. It does not need it if the source is always on, or if the power states guarantee that the sink is off whenever the source is. Separately, any crossing between different supply voltages needs level shifting, whatever the on and off states.

## Choosing the clamp value

An isolation cell holds its output at a constant while the source is off. The constant must be the receiving pin's inactive value, so that the live logic sees a quiet input: 1 for an active-low reset, chip select or write enable; 0 for an active-high request, interrupt, valid flag or clock-gate enable. Clamping to the wrong value can hold a block in reset, start a spurious transfer, or open a clock gate. Some cells latch the last value instead, which suits data but can freeze an asserted control signal. The isolation cell is powered from an always-on supply, usually in the receiving domain.

## Level-shifter direction

Driving a gate supplied at 0.9 V with a 0.7 V high level leaves its PMOS partly on, which causes static current and a degraded output. A low-to-high shifter, typically a cross-coupled structure, uses both supplies to restore full swing. Going from high to low voltage, the input overdrives the receiver, so a buffer on the destination supply restores the lower swing; flows still insert an explicit high-to-low cell so the crossing is visible to checks. Where a crossing needs both isolation and shifting, an enable level shifter combines them. Placement rules then decide in which domain each cell sits and which supplies it connects to.

## Going deeper

Keating et al., "Low Power Methodology Manual" (2007); IEEE Std 1801 (UPF).
