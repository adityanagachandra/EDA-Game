# The device-model seam

The device-model interface is the simulator's most important seam. The core asks one question per device per iteration: given these terminal voltages, temperature and parameters, what are the currents, charges and their derivatives? Everything model-specific lives behind that question. Designing the seam well, and testing what passes through it, is what lets new transistor models arrive as plug-ins instead of core edits.

## What the device interface returns

For each device the core supplies terminal voltages, temperature, instance parameters and the analysis mode. The model returns terminal currents, which form the Newton residual; their derivatives with respect to terminal voltages, which form the Jacobian; terminal charges, which the integrator differentiates so that charge is conserved; charge derivatives, the capacitances used in transient and AC; and, in noise analysis, its noise source densities. A MOSFET returning $g_m$, $g_{ds}$ and $I_{eq}$ at an iterate is exactly this contract. A narrow, stable interface like this lets the core treat a diode and a BSIM transistor identically, as stamps.

## Adding a model as a plug-in

A new model, typically written in Verilog-A, is compiled into code that implements the device interface and loads at run time, leaving the solver untouched: a Tier 2 change. Good plug-ins follow the seam rules: no hidden global state, so results are deterministic and order-independent; a frozen extension ABI, so compiled models survive simulator releases; and coarse calls, evaluating many instances of one model per call. A plug-in whose currents are right but whose derivatives have a sign error still converges to the right answer, because the residual uses the currents, but Newton slows or fails, which is why derivative correctness is a separate test.

## Testing a new device model

Device-level tests check the interface contract: compare the returned Jacobian with finite differences of the returned currents at many bias points; check continuity and source-drain symmetry near $V_{DS} = 0$, where many models kink and cause convergence failures and wrong distortion; and, for a speed-optimized rewrite, compare currents, charges and derivatives with the reference model over a grid of biases and temperatures. Circuit-level golden tests then catch what device tests miss, such as charge errors: if a ring oscillator built from the model runs 5% faster than its golden result, promotion is blocked, and the golden value is not updated to match.

## Going deeper

McAndrew et al., "Best practices for compact modeling in Verilog-A" (IEEE Journal of the Electron Devices Society, 2015); the Xyce Users' Guide, on device models.
