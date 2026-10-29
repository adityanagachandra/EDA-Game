# Silicon on insulator

Silicon on insulator builds transistors in a thin silicon film on a buried oxide instead of in a bulk wafer. The insulator removes most junction capacitance and the latch-up path, and in fully depleted SOI it gives a thin body with good electrostatics and a usable back gate. SOI also brings its own modeling problems, floating bodies and self-heating, that compact models and characterization flows must handle.

## Partially and fully depleted SOI

In partially depleted SOI (PD-SOI), the silicon film is thicker than the maximum depletion depth under the gate. A neutral body region remains, and unless it is tied to a contact it floats. Charge from impact ionization and junction leakage collects there and shifts the body potential, and through the body effect the threshold. Delays then depend on switching history (the history effect), and output characteristics can show a kink.

In fully depleted SOI (FD-SOI), the film is thin enough, a few nanometres to around ten, that the gate depletes it completely. There is no neutral body to float. The thin film limits how far the drain field can reach, giving good short-channel control without heavy channel doping, which also reduces random dopant fluctuation.

In both, the buried oxide removes most bottom-wall junction capacitance and breaks the well-substrate thyristor, so latch-up largely disappears.

## Back-gate bias

In FD-SOI with a thin buried oxide, the substrate under the oxide acts as a second gate. Biasing it shifts the threshold with a coupling coefficient $r$, often quoted in mV per volt of back bias:
$$V_t = V_{t0} - r\,V_{BB} \quad \text{(NMOS, forward back bias positive)}.$$

Example: $V_{t0} = 380$ mV and $r = 80$ mV/V. A forward back bias of 1.0 V gives $V_t = 300$ mV; a reverse bias of 1.0 V gives 460 mV.

This is a far stronger knob than bulk body bias, because it couples through a thin oxide rather than through depletion charge, and it works over a wider voltage range since there are no forward-biased junctions to limit it. Designs use it adaptively: forward bias to speed up slow dies or high-performance modes, reverse bias to cut leakage in standby.

## What SOI means for models and flows

Compact models for PD-SOI (such as BSIM-SOI) carry an internal body node with slow charge and discharge time constants, so transient simulations must run long enough for the body state to settle, and timing can depend on activity history.

The buried oxide conducts heat roughly a hundred times worse than silicon, so SOI devices self-heat: the channel runs hotter than the ambient under sustained current, lowering mobility. SOI models include a thermal node with a thermal resistance and capacitance.

FD-SOI back biasing changes threshold, speed and leakage after fabrication, so libraries are characterized at several back-bias points, and implementation tools must sign off timing and power at each operating point the chip will use. Models with independent front and back gates (BSIM-IMG) exist for exactly this. Finally, thin films bring film-thickness variation and quantum confinement into threshold, which statistical models must capture.

## Going deeper

Colinge, "Silicon-on-Insulator Technology: Materials to VLSI" (3rd ed.); the BSIM-SOI and BSIM-IMG technical manuals.
