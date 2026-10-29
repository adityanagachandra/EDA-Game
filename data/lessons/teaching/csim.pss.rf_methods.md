# RF simulation beyond periodic steady state

Periodic steady state is the base on which RF simulation builds. Mixers, oscillators and switched circuits need small-signal and noise analyses that account for a large periodic drive, and multitone or modulated signals need methods that avoid simulating millions of carrier cycles. Knowing which analysis answers which question, and when shooting or harmonic balance is the better base, is the core of RF simulation practice.

## Periodic small-signal analyses

A mixer's small RF input sees a circuit whose operating point is modulated by the large local oscillator. Linearizing about the periodic steady state gives a linear time-varying system, which translates a small input at $f$ to outputs at $f + k f_{LO}$ for every harmonic $k$. Periodic AC (PAC) computes those conversion gains; periodic transfer-function analysis (PXF) uses an adjoint solve to get the transfer from many inputs and sidebands to one output; periodic noise (PNoise) handles cyclostationary noise, including noise folded from the image and other sidebands onto the output. All three linearize around the PSS solution, so a poorly converged PSS gives wrong gains and noise.

## Multitone and envelope methods

Two tones at 1.000 and 1.001 GHz beat with a 1 us period while the carrier period is 1 ns, so a transient must integrate a thousand fast cycles to see one beat. Multitone harmonic balance represents both tones and their mixing products $m f_1 + n f_2$ directly, truncated by total order; the number of frequencies, and the system size, grows roughly quadratically with the order. A modulated carrier whose envelope varies over microseconds suits envelope following instead: a periodic problem at the carrier is solved at each step of a slow time axis. Both methods exploit the separation of time scales that plain transient cannot.

## Choosing shooting or harmonic balance

Shooting integrates in time, so it handles sharp switching edges with adaptive steps and few unknowns, which suits switched-capacitor circuits, switching converters and hard-driven mixers. Harmonic balance works in frequency, so it stamps frequency-domain models such as S-parameters directly, represents distributed elements exactly at each harmonic, and handles mildly nonlinear circuits with few harmonics, which suits low-noise amplifiers and power amplifiers with package models. Strongly nonlinear waveforms need many harmonics in HB and become expensive. Autonomous oscillators need an extra unknown in either method: the period, fixed by a phase condition that removes the arbitrary time shift.

## Going deeper

Kundert, White and Sangiovanni-Vincentelli, "Steady-State Methods for Simulating Analog and Microwave Circuits" (Kluwer, 1990); Telichevesky, Kundert and White, "Efficient steady-state analysis based on matrix-free Krylov-subspace methods" (DAC, 1995).
