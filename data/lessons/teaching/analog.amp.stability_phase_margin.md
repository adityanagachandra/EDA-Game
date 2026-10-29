# Stability and phase margin

An amplifier in negative feedback is stable only if its loop gain loses enough phase before reaching unity magnitude. Phase margin measures that distance and predicts overshoot and ringing in the closed-loop step response. It is a hard spec in every op-amp sizing problem and the quantity most easily destroyed by layout parasitics, so a tool designer must be able to compute it from pole and zero positions.

## Phase margin from the loop gain

Feedback is unstable when the loop gain $T(j\omega)$ reaches -1: magnitude one with phase -180 degrees. Phase margin is $\mathrm{PM} = 180 + \angle T(j\omega_u)$, where $\omega_u$ is the unity-gain frequency. For left-half-plane poles the phase is a sum of arctangents: $\angle T = -\sum_i \arctan(\omega_u/p_i)$. Take poles at 0.01 MHz and 100 MHz with crossover at 50 MHz: the first contributes $\arctan(5000) \approx 89.99$ degrees, the second $\arctan(0.5) = 26.57$ degrees, so $\mathrm{PM} = 180 - 89.99 - 26.57 = 63.4$ degrees. The argument of each factor of the transfer function adds, so phase computation is just complex-argument arithmetic at one frequency.

## Zeros and the right-half plane

A left-half-plane zero at $z$ adds $+\arctan(\omega_u/z)$ of phase. A right-half-plane zero adds magnitude like a normal zero but subtracts $\arctan(\omega_u/z)$, like a pole. In a Miller op-amp the uncancelled zero at $g_{m6}/C_c$ is often only a few times GBW. With poles at 0.01 and 100 MHz, crossover at 50 MHz and a right-half-plane zero at 250 MHz, the zero costs $\arctan(0.2) = 11.3$ degrees and the margin falls from 63.4 to 52.1 degrees. This is why the nulling resistor matters, and why a sizing tool that ignores the zero overestimates phase margin. The rule generalizes: any element whose phase lags at crossover, pole or right-half-plane zero or delay, subtracts from the margin.

## Designing for a target margin

With a dominant pole far below crossover, it contributes almost exactly -90 degrees, so $\mathrm{PM} \approx 90 - \arctan(\mathrm{GBW}/f_{p2})$. Inverting gives the required non-dominant pole: $f_{p2} = \mathrm{GBW}/\tan(90 - \mathrm{PM})$. For 60 degrees, $f_{p2} = \mathrm{GBW}/\tan 30 = 1.73\,\mathrm{GBW}$; for 45 degrees, $f_{p2} = \mathrm{GBW}$. Many sizing procedures use $f_{p2} \ge 2.2\,\mathrm{GBW}$, which gives about 65.6 degrees. Since $f_{p2} \approx g_{m6}/(2\pi C_L)$ and $\mathrm{GBW} = g_{m1}/(2\pi C_c)$, the condition becomes $g_{m6}/C_L \ge 2.2\,g_{m1}/C_c$: a relation among sizes that an equation-based sizer can impose directly. A margin around 60 degrees gives a well-damped step; much less rings, much more is slow and wastes power.

## Going deeper

Gray, Hurst, Lewis and Meyer, "Analysis and Design of Analog Integrated Circuits" (5th ed., 2009), chapter 9; Franklin, Powell and Emami-Naeini, "Feedback Control of Dynamic Systems", chapter 6.
