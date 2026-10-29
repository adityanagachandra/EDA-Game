# Model cards and device instances

A netlist separates what a device is from where it is. An instance card, such as M1 d g s b nch W=1u L=0.1u, places one transistor and sets its geometry. A .model card, such as .model nch nmos (level=... ), carries the process parameters that many instances share. Foundry libraries, subcircuits and corner sections build on this split. Knowing it is essential for anyone who writes a netlist parser, a model loader or a Monte Carlo engine.

## Instances and model cards

An instance card lists terminals in a fixed order: a diode D1 a k dmod has anode then cathode; a MOSFET M1 d g s b nch has drain, gate, source and bulk, the bulk terminal setting the body effect. Instance parameters (W, L, multiplicity) describe that one device. The .model card names a model type (d, nmos, pmos) and gives parameters shared by every instance that references it. Two hundred diodes naming dmod share one saturation current: change is on the card and all two hundred change. Internally the simulator keeps one model structure and many instance structures pointing to it, so device evaluation code reads model parameters once and instance geometry per device.

## Model levels and parameter sets

The level parameter (or the model type name) selects which equations are evaluated. Level 1 is the Shichman-Hodges square law: $I_D = \tfrac{k'}{2}\tfrac{W}{L}(V_{GS} - V_T)^2(1 + \lambda V_{DS})$ in saturation, with a handful of parameters such as vto, kp and lambda. Modern processes use BSIM-class compact models with hundreds of parameters, because short-channel effects (velocity saturation, drain-induced barrier lowering, mobility degradation) bend the curves far from the square law. Foundries also bin models: several cards cover ranges of length and width (lmin, lmax, wmin, wmax), and the simulator picks the bin containing the instance geometry. A tool must apply binning consistently, or two identical devices could get different parameters.

## Libraries, corners and subcircuits

Foundry model files group parameter sets into named library sections, one per process corner, selected with a line such as .lib models.lib tt or .lib models.lib ss. Switching the section switches every model at once. Subcircuits add hierarchy: .subckt inv in out vdd vss defines a body, and each X instance expands it with its own internal nodes, renamed by instance path so that X1's internal n1 and X2's n1 stay distinct. Two inverter instances of two transistors each flatten to four MOSFETs. For statistical simulation, randomness placed on a shared model card moves all devices together (global variation); local mismatch needs an independent draw per instance.

## Going deeper

The ngspice manual, chapters on device models and subcircuits; the BSIM4 and BSIM-CMG technical manuals from the Berkeley BSIM group.
