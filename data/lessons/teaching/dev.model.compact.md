# Compact transistor models

A compact model is the closed-form description of a transistor that a circuit simulator evaluates at every Newton iteration for every device. It must be accurate across bias, geometry and temperature, fast, smooth, and charge conserving. Its equations and parameters are the contract between a foundry and every design team, and every timing library is characterized through it. Knowing what a compact model must capture, and why, is basic for anyone building simulators or characterization tools.

## The contract with the simulator

At each Newton iteration, SPICE asks each device model: given terminal voltages (and temperature), what are the terminal currents and charges, and their derivatives? The derivatives are the conductances $\partial I/\partial V$ and capacitances $\partial Q/\partial V$ that fill the Jacobian.

Three properties follow. Smoothness: currents and charges must be continuous with continuous derivatives across all regions (subthreshold, linear, saturation, accumulation), because a kink in $g_m$ or $g_{ds}$ makes Newton oscillate between regions; analog distortion analysis needs higher derivatives smooth too. Charge conservation: terminal charges must be defined first and capacitances derived from them; the older Meyer capacitance model was not charge conserving and could create or destroy charge in switched circuits. Speed: evaluation cost multiplies by devices, iterations and time steps.

Unlike TCAD, a compact model does not solve PDEs. It uses physics-derived closed forms with fitted parameters, which is why it can be fast enough.

## Formulations and qualification

Compact MOSFET models fall into three families by their core variable. Threshold-voltage-based models (BSIM3, BSIM4) build regional expressions around $V_t$ and blend them with smoothing functions. Surface-potential-based models (PSP, BSIM-CMG) solve for the surface potential from the gate bias and electrostatics, then derive charge and current, giving one continuous description from weak to strong inversion. Inversion-charge-based models (EKV-style) use the inversion charge as the core variable.

A model must scale with geometry: one parameter set must predict all legal widths and lengths, through built-in dependences or binning. It must include temperature, since corners span cold to hot.

Qualification checks go beyond fitting $I_D$. The Gummel symmetry test drives the device symmetrically around $V_{DS} = 0$ and checks that source-drain swap leaves the physics unchanged and higher derivatives are smooth, which passive mixers and switches require. Other checks cover $g_m/I_D$ behavior, output conductance, C-V curves and noise.

## Engineering a compact model

Model evaluation can be a large fraction of SPICE run time on post-layout netlists with millions of devices, so expressions are chosen for speed as well as accuracy.

Why not a table of measured I-V points? Tables are exact at their points, but derivatives between points depend on the interpolant, can be discontinuous, and say nothing outside the table. Physics-based expressions extrapolate sensibly, stay smooth and conserve charge.

A model is qualified over a range of bias, geometry and temperature. Outside it, the equations still evaluate but nobody has verified them, so simulators can warn on out-of-range operation.

Models also carry noise (channel thermal noise and flicker noise, both bias and geometry dependent), parasitic resistances and capacitances, gate and junction leakage, layout-dependent stress effects and, through extensions, aging. Each is a separately parameterized piece. For a tool designer, the important point is that the model is a module behind a narrow interface, which is exactly why the frontier treats it as a plug-in seam.

## Going deeper

Tsividis and McAndrew, "Operation and Modeling of the MOS Transistor" (3rd ed.); McAndrew, "Validation of MOSFET model source-drain symmetry", IEEE Transactions on Electron Devices (2006); Gildenblat (ed.), "Compact Modeling: Principles, Techniques and Applications" (2010).
