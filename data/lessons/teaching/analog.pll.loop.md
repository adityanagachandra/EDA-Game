# Phase-locked loops

A phase-locked loop locks a VCO's output to a multiple of a reference clock. A phase detector compares the reference with the divided VCO output, a loop filter smooths the error, and the VCO corrects its frequency until the phases align. PLLs generate every clock on a chip and every local oscillator in a radio, and their behaviour is a feedback-control problem in the Laplace domain.

## Integer-N synthesis

At lock the divided output matches the reference in frequency and phase, so $f_{out} = N f_{ref}$. A 26 MHz reference with $N = 92$ gives 2392 MHz. The output can only move in steps of $f_{ref}$, so channel spacing equals the reference frequency: a 1 MHz channel grid needs a 1 MHz reference, which forces a low loop bandwidth (typically below a tenth of $f_{ref}$ for stability of a sampled loop) and therefore slow settling and poor suppression of VCO noise. This coupling between resolution and bandwidth is the limitation that fractional-N synthesis removes. The divider is a digital counter, so even an integer-N PLL is a mixed-signal circuit whose verification needs both analog and event-driven models.

## Fractional-N division

A fractional-N PLL changes the division ratio between $N$ and $N + 1$ so that its average is fractional. Dividing by $N + 1$ for $K$ out of every $M$ reference cycles gives an average ratio $N + K/M$ and $f_{out} = f_{ref}(N + K/M)$. With a 40 MHz reference, $N = 60$, $K = 3$ and $M = 8$, $f_{out} = 40 \cdot 60.375 = 2415$ MHz. Fine resolution no longer needs a low reference, so the loop can be fast. The price is that the instantaneous phase error is periodic, producing fractional spurs. Delta-sigma modulators randomize the division sequence and push this quantization noise to high frequencies where the loop filters it, the same noise-shaping idea used in delta-sigma ADCs.

## Lock and the loop as a control system

In lock the control voltage settles wherever the VCO needs it: $V_c = (f_{out} - f_0)/K_{VCO}$. A VCO with free-running frequency 2.2 GHz and $K_{VCO} = 400$ MHz/V, locked at 2.4 GHz, sits at 0.5 V. If that voltage lies outside the charge pump's linear range the loop cannot lock, which is why tuning range and $K_{VCO}$ are checked against the supply. For dynamics, the loop is linearized in phase: the phase detector gives a current or voltage proportional to phase error, the loop filter has transfer function $Z(s)$, the VCO is $K_{VCO}/s$ and the divider is $1/N$. The open-loop gain is their product, and stability, bandwidth and noise transfer all follow from it by Laplace-domain analysis.

## Going deeper

Razavi, "Design of CMOS Phase-Locked Loops" (2020); Gardner, "Phaselock Techniques" (3rd ed., 2005).
