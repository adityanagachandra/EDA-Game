# Short-channel effects

When a transistor's channel becomes comparable to the depletion depths of its source and drain, the gate loses exclusive control of the channel. Threshold falls with length, drain voltage lowers the source barrier, current saturates through velocity rather than pinch-off, and in the extreme the device punches through. These short-channel effects drove every device architecture change of the last two decades and fill much of every modern compact model.

## How the gate loses control

Three related effects appear as $L$ shrinks.

Threshold roll-off: the source and drain junction depletion regions already hold part of the charge under the gate, so the gate needs less voltage to invert. $V_t$ falls as $L$ shrinks, even at low drain bias.

Drain-induced barrier lowering (DIBL): the drain's field reaches the source end of the channel and lowers the barrier there. $V_t$ then falls as $V_{DS}$ rises, raising off-current.

Punch-through: at extreme drain bias the source and drain depletion regions merge below the surface, where the gate is weakest, and current flows uncontrolled.

All three are governed by the ratio of channel length to an electrostatic natural length $\lambda$, set by oxide thickness and the depth of silicon the gate must control. Short-channel effects scale roughly as $e^{-L/(2\lambda)}$, so the cure is to shrink $\lambda$, not merely to lengthen the gate.

## Countermeasures and companion effects

Planar processes fought short-channel effects with doping. Halo (pocket) implants raise the doping near the source and drain, which raises threshold most in short devices and flattens roll-off; overcompensation produces the reverse short-channel effect, where $V_t$ first rises as $L$ shrinks. Shallow source-drain extensions reduce charge sharing. Thinner and high-k dielectrics increase gate control.

Two further effects accompany short channels without being caused by charge sharing. Velocity saturation makes current grow nearly linearly with overdrive. Vertical-field mobility degradation lowers effective mobility at high gate bias, as carriers are pressed against the rough interface; simple models write $\mu_{eff} = \mu_0/(1 + \theta(V_{GS} - V_t))$.

When doping and oxide scaling ran out, the remaining lever was geometry: a thin body with gates on several sides, which is what FinFETs, nanosheets and FD-SOI provide.

## Consequences for EDA

Short-channel effects make the electrical behavior of a transistor depend on its exact length, its drain voltage during switching, and its neighbors.

Threshold depends steeply on $L$ near minimum length, so gate-length variation becomes threshold and delay variation. This is why statistical and LVF data are needed, and why designers use longer-than-minimum channels for matched analog devices.

Drive current depends on $V_{DS}$ through DIBL and output conductance, and $V_{DS}$ is the output voltage during a transition. Current-source cell models tabulate current against input and output voltage for this reason.

Off-current is exponential in threshold, so a compact model that underestimates roll-off or DIBL at minimum length can be wrong by an order of magnitude in leakage. Extraction must therefore fit short devices at high drain bias explicitly, not just long-channel curves, and model qualification checks leakage across length.

## Going deeper

Taur and Ning, "Fundamentals of Modern VLSI Devices" (2nd ed.), section 3.2; Yan, Ourmazd and Lee, "Scaling the Si MOSFET: from bulk to SOI to bulk", IEEE Transactions on Electron Devices (1992).
