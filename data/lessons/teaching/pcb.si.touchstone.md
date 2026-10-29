# Touchstone files and channel S-parameters

Passive board structures at high frequency, such as traces, vias, connectors and packages, are described by scattering parameters measured or simulated over frequency. The Touchstone file format (.sNp) is the common container: one frequency per row, $N^2$ complex S-parameters for an $N$-port network. Channel simulators cascade these blocks between transmitter and receiver models. Reading the format and converting magnitudes to decibels are everyday operations for a board SI tool.

## The Touchstone format

A Touchstone file's extension names the port count: .s2p for 2 ports, .s4p for 4 (a differential pair has 4 single-ended ports). An option line such as "# GHz S MA R 50" gives the frequency unit, the parameter type (S), the number format (MA magnitude-angle, DB decibel-angle, or RI real-imaginary) and the reference impedance. Each frequency point then lists the frequency and the $N^2$ complex entries, so $2N^2 + 1$ real numbers. A 4-port file carries 16 S-parameters, 33 numbers, per frequency. The order is row-major except for one well-known quirk: 2-port files list $S_{11}, S_{21}, S_{12}, S_{22}$. Version 2.0 of the format adds keywords for mixed reference impedances and other metadata, which parsers must handle.

## Insertion loss in decibels

$S_{21}$ is the ratio of the voltage wave leaving port 2 to the wave entering port 1, with all ports terminated in the reference impedance. **Insertion loss** is $-20\log_{10}|S_{21}|$ in dB. Example: $|S_{21}| = 0.5$ is 6.02 dB of loss; 20 dB is $|S_{21}| = 0.1$. Board traces lose more at higher frequency, from conductor loss (skin effect, roughness) and dielectric loss (the laminate's loss tangent), so channel budgets are stated at the Nyquist frequency of the link, half the symbol rate. For matched sections the dB losses add, which is why loss is budgeted in dB per inch for each layer and material, and why a longer or lossier channel needs more equalization.

## Return loss and mismatch

$S_{11}$ is the wave reflected back into port 1. For a one-port load $Z_L$ measured with reference $Z_0$, $S_{11}$ equals the reflection coefficient $(Z_L - Z_0)/(Z_L + Z_0)$. **Return loss** is $-20\log_{10}|S_{11}|$; larger means a better match. Example: a 75 ohm load on a 50 ohm reference gives $|S_{11}| = 0.2$, a return loss of 13.98 dB; a 55 ohm load gives $|S_{11}| = 0.0476$, about 26.4 dB. Return loss of a channel reveals its discontinuities: vias, connector transitions and stubs show as peaks in $|S_{11}|$ and notches in $|S_{21}|$. Data quality matters too: a simulator needs passive and causal data, so tools check and enforce passivity and causality before using a Touchstone file in time-domain simulation.

## Going deeper

The Touchstone file format specification, published through the IBIS Open Forum. Pozar, "Microwave Engineering" (4th ed., 2011), chapter 4, on scattering parameters.
