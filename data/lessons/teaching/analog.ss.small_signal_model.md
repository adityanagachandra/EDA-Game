# The MOSFET small-signal model

Every analog hand calculation and every AC simulation starts by linearizing the transistors around their DC operating point. The small-signal model replaces a saturated MOSFET by a transconductance $g_m$, an output resistance $r_o$ and capacitances. A simulator computes these from the compact model at the operating point; a designer estimates them from the square law. Both views matter to a tool designer: the first is what AC analysis does, the second is what sizing equations assume.

## Linearizing around the operating point

In saturation the drain current is a nonlinear function $I_D(V_{GS}, V_{DS})$. For small signals around a bias point, a first-order Taylor expansion gives $i_d = g_m v_{gs} + g_{ds} v_{ds}$, with $g_m = \partial I_D/\partial V_{GS}$ and $g_{ds} = 1/r_o = \partial I_D/\partial V_{DS}$. The circuit becomes a voltage-controlled current source $g_m v_{gs}$ in parallel with $r_o$, plus gate capacitances $C_{gs}$, $C_{gd}$ and junction capacitances. This is exactly what an AC analysis does: it solves the DC operating point, evaluates the Jacobian of the device model there and solves a linear complex system at each frequency. The model is valid only while signals are small enough that the curvature of $I_D$ does not matter, typically a few millivolts at the gate.

## Square-law estimates of gm and ro

In the long-channel square law, $I_D = \frac{1}{2}\mu_n C_{ox}(W/L)V_{OV}^2$ with overdrive $V_{OV} = V_{GS} - V_T$. Differentiating gives three equivalent forms: $g_m = \mu_n C_{ox}(W/L)V_{OV} = 2I_D/V_{OV} = \sqrt{2\mu_n C_{ox}(W/L)I_D}$. Channel-length modulation adds a factor $(1+\lambda V_{DS})$, so $r_o \approx 1/(\lambda I_D)$. For $I_D = 100\ \mu$A at $V_{OV} = 0.2$ V, $g_m = 2 \cdot 100/0.2 = 1000\ \mu$S, or 1 mS. With $\lambda = 0.1\ \mathrm{V}^{-1}$, $r_o = 1/(0.1 \cdot 10^{-4}) = 100$ kilohms. Note the scaling: at fixed current, $g_m$ grows only as the square root of $W/L$, while at fixed overdrive it grows linearly with current.

## Where the square law fails

Modern devices are not long-channel. Velocity saturation makes $I_D$ closer to linear in $V_{OV}$ at high overdrive, so $g_m$ saturates. In weak inversion the current is exponential in $V_{GS}$ and $g_m = I_D/(nU_T)$, larger than $2I_D/V_{OV}$ would predict. Output resistance depends on drain-induced barrier lowering and is a strong function of channel length, not a single $\lambda$. A compact model such as BSIM captures all of this, and the simulator's small-signal parameters come from it. For a tool designer the consequence is direct: an optimizer built on square-law equations finds the optimum of the wrong function. That is the motivation for the gm/Id methodology and lookup tables, which take $g_m$, $r_o$ and capacitances from simulated device data instead of from a formula.

## Going deeper

Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., 2017), chapters 2 and 3; Tsividis and McAndrew, "Operation and Modeling of the MOS Transistor" (3rd ed., 2011).
