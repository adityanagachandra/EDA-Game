# Voltage-controlled oscillators

A voltage-controlled oscillator produces a periodic signal whose frequency is set by a control voltage. It is the frequency-generating core of every PLL, clock generator and RF synthesizer. Two families dominate: ring oscillators built from inverters, and LC oscillators built around an inductor-capacitor tank. Their frequency equations and tuning gain are the first numbers a PLL designer and a PLL model need.

## Ring oscillators

An odd number $N$ of inverting stages in a loop has no stable DC state, so it oscillates. An edge must pass through every stage twice per period, once rising and once falling, so the period is $T = 2Nt_d$ and the frequency is $f = 1/(2Nt_d)$, where $t_d$ is the stage delay. Five stages of 20 ps give $T = 200$ ps and $f = 5$ GHz. Tuning changes $t_d$, by varying the supply, a load capacitance or a current limit. Ring oscillators are compact, need no inductor and give multiple phases, but their phase noise is relatively high because every stage adds noise to the timing of every edge.

## LC oscillators

An LC tank resonates at $f = 1/(2\pi\sqrt{LC})$. A cross-coupled transistor pair supplies negative resistance that cancels the tank's losses and sustains oscillation. With $L = 2$ nH and $C = 1$ pF, $\sqrt{LC} = \sqrt{2 \times 10^{-21}} = 4.47 \times 10^{-11}$ s and $f = 3.56$ GHz. Tuning uses a varactor whose capacitance depends on the control voltage, often with switched capacitor banks for coarse bands. Because the tank filters noise and the quality factor of the inductor sets the loss, LC oscillators achieve much lower phase noise than rings, which is why RF synthesizers use them. The inductor is large and needs electromagnetic modelling, a separate EDA problem in itself.

## Tuning gain

Near its operating point a VCO is modelled as linear: $f = f_0 + K_{VCO}V_c$, with tuning gain $K_{VCO}$ in Hz per volt. From two measured points, 2.40 GHz at 0.4 V and 2.46 GHz at 0.6 V, $K_{VCO} = 60\ \mathrm{MHz}/0.2\ \mathrm{V} = 300$ MHz/V. In loop analysis it is used in rad/s per volt, $2\pi K_{VCO}$, and since phase is the integral of frequency, the VCO acts as an integrator from control voltage to phase: $K_{VCO}/s$. A large $K_{VCO}$ covers a wide range but converts noise on the control line into large phase noise, which is why designs combine coarse band switching with a small fine-tuning gain. Real tuning curves are nonlinear, so $K_{VCO}$ varies across the range and the loop dynamics vary with it.

## Going deeper

Razavi, "Design of CMOS Phase-Locked Loops" (2020), chapters on oscillators; Hajimiri and Lee, "A general theory of phase noise in electrical oscillators" (IEEE JSSC, 1998).
