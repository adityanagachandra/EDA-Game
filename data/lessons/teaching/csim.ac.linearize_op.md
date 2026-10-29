# Small-signal models from the operating point

AC analysis begins where the DC operating point ends. Every nonlinear device is replaced by its small-signal model, the derivatives of its currents and charges at the bias point, which are exactly the conductances Newton already computed in its last iteration. For a MOSFET those are the transconductance and output conductance. Getting them right by hand is how you check an AC result against intuition.

## Transconductance at the bias point

For a square-law device in saturation, $I_D = \tfrac{k}{2} V_{ov}^2$ with $V_{ov} = V_{GS} - V_T$, so $g_m = k V_{ov}$, which can be written without $k$ as
$$g_m = \frac{2 I_D}{V_{ov}}.$$
With $I_D = 0.5$ mA and $V_{ov} = 0.2$ V, $g_m = 5$ mS. The form matters: at fixed current, a smaller overdrive gives more transconductance, which is why low-power analog designs bias near threshold. The simulator does not use this formula; it takes $\partial I_D/\partial V_{GS}$ from the compact model at the converged operating point, which agrees with the square law only for long-channel devices in strong inversion.

## Output resistance

Channel-length modulation makes the drain current rise slowly with drain voltage, $I_D \propto (1 + \lambda V_{DS})$, so the output conductance is $g_{ds} \approx \lambda I_D$ and the small-signal output resistance is
$$r_o \approx \frac{1}{\lambda I_D}.$$
With $\lambda = 0.1$ 1/V and $I_D = 0.5$ mA, $r_o = 20$ k$\Omega$ (one over millisiemens is kilohms). In AC analysis $r_o$ stamps as a conductance between drain and source, $g_m$ as a VCCS controlled by the gate-source voltage, and the device capacitances, derived from the charge model, as $j\omega C$ entries. All are frozen at the operating point: AC analysis is linear by construction.

## Intrinsic gain

The product $g_m r_o$ is the largest voltage gain a single transistor can give, its intrinsic gain. With the square-law forms, $g_m r_o = \tfrac{2 I_D}{V_{ov}} \cdot \tfrac{1}{\lambda I_D} = \tfrac{2}{\lambda V_{ov}}$: the current cancels. For $\lambda = 0.1$ 1/V and $V_{ov} = 0.2$ V, $g_m r_o = 100$, or 40 dB. Shorter channels have larger $\lambda$ and lower intrinsic gain. A designer reading AC output should expect a common-source stage's DC gain to stay below $g_m r_o$; a simulated gain far above it usually means a mistaken bias point, such as a device in the wrong region, rather than a good design.

## Going deeper

Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., McGraw-Hill, 2017), chapters 2 and 3; Najm, "Circuit Simulation" (Wiley, 2010), chapter on AC analysis.
