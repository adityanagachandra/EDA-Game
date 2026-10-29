# Linearizing a MOSFET at a Newton iterate

A MOSFET is a three-port nonlinearity: its drain current depends on both gate-source and drain-source voltages. Linearizing it for Newton produces a transconductance, an output conductance and an equivalent current, the same companion idea as for the diode with one more partial derivative. Device models in SPICE return exactly these derivatives, and their accuracy decides how fast Newton converges.

## Transconductance from the square law

In saturation the square law with channel-length modulation is $I_D = \tfrac{k}{2}(V_{GS} - V_T)^2(1 + \lambda V_{DS})$. The transconductance is the partial derivative with respect to the gate-source voltage at the iterate:
$$g_m = \frac{\partial I_D}{\partial V_{GS}} = k(V_{GS} - V_T)(1 + \lambda V_{DS}).$$
For $k = 1$ mA/V$^2$, $V_T = 0.5$ V, $V_{GS} = 0.9$ V, $V_{DS} = 1$ V and $\lambda = 0.1$ 1/V, the overdrive is 0.4 V and $g_m = 1 \times 0.4 \times 1.1 = 0.44$ mS. It stamps like a VCCS controlled by $v_{GS}$ from drain to source.

## Output conductance

The partial derivative with respect to the drain-source voltage is
$$g_{ds} = \frac{\partial I_D}{\partial V_{DS}} = \tfrac{k}{2}(V_{GS} - V_T)^2 \lambda.$$
In the example, $g_{ds} = 0.5 \times 0.16 \times 0.1 = 0.008$ mS, or 8 uS, an output resistance of 125 k$\Omega$. It stamps like a resistor between drain and source. Real compact models compute $g_{ds}$ from far more physics, but they must keep it continuous across the boundary between triode and saturation: a jump in $g_{ds}$ there makes Newton hop back and forth between regions, a classic convergence failure in older models.

## The MOSFET equivalent current

The linearized drain current is $I_D \approx g_m V_{GS} + g_{ds} V_{DS} + I_{eq}$, where
$$I_{eq} = I_D(V_k) - g_m V_{GS} - g_{ds} V_{DS}$$
is the intercept of the tangent plane. In the example $I_D = 0.5 \times 0.16 \times 1.1 = 0.088$ mA, so $I_{eq} = 0.088 - 0.44 \times 0.9 - 0.008 \times 1 = -0.316$ mA. The full stamp is the VCCS $g_m$, the conductance $g_{ds}$ and a current source $I_{eq}$ from drain to source. Simulators also stamp gate and bulk derivatives and charge derivatives for transient. A model with a wrong derivative still converges to the right answer, because the residual uses the true current, but it converges slowly.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 4; Tsividis and McAndrew, "Operation and Modeling of the MOS Transistor" (3rd ed., Oxford, 2011).
