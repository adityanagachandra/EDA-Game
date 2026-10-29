# Current mirrors and cascodes

Current mirrors copy and scale bias currents across a chip and serve as active loads. Their accuracy and output resistance set bias precision, gain and power-supply rejection. Cascoding raises output resistance by a factor of one intrinsic gain. These are among the most frequently recognized structures in analog layout tools, which must match their devices.

## Copying current through equal gate voltage

A diode-connected transistor M1 carries a reference current and develops the $V_{GS}$ needed to sink it. A second transistor M2 with the same $V_{GS}$ and the same $V_T$ and $L$ carries a current scaled by the width ratio: $I_{OUT} = I_{REF}(W_2/L_2)/(W_1/L_1)$. A 2:5 width ratio copies 40 microamperes into 100. Ratios are built from identical unit fingers (two units and five units) rather than different widths, so that edge effects scale with the ratio and the devices can be interdigitated. Changing $L$ between devices is avoided, because short-channel effects make $L$ ratios inaccurate. The mirror's accuracy is then limited by $V_{DS}$ mismatch and by device mismatch.

## Channel-length modulation and copy error

Equal $V_{GS}$ does not give equal currents when $V_{DS}$ differs, because $I_D \propto 1 + \lambda V_{DS}$. The diode device sits at $V_{DS1} = V_{GS}$; the output device sits wherever the load puts it. So $I_{OUT} = I_{REF}(1+\lambda V_{DS2})/(1+\lambda V_{DS1})$. For $\lambda = 0.2\ \mathrm{V}^{-1}$, $V_{GS} = 0.6$ V and $V_{DS2} = 1.2$ V the error is $1.24/1.12 = 1.107$, an 11% copy error. The same effect gives the mirror a finite output resistance $r_o$, which limits its use as a load. The cures are longer channels (smaller $\lambda$) and cascodes that hold $V_{DS2}$ close to $V_{DS1}$.

## The cascode multiplies output resistance

Stacking a common-gate device M2 on top of the current-source device M1 shields M1's drain from the output. A change $\Delta V$ at the output barely moves M1's drain, because M2 needs only a tiny gate-source change to absorb it. The small-signal result is $R_{out} \approx g_{m2} r_{o2} r_{o1}$: the cascode multiplies $r_{o1}$ by its own intrinsic gain. With $g_{m2} = 1$ mS and both $r_o = 50$ kilohms, $R_{out} \approx 50 \cdot 50$ kilohms $= 2.5$ megohms, 50 times the simple mirror. The cost is headroom: the output must stay above two saturation voltages. As supplies fall, that headroom becomes scarce, which is why low-voltage designs and migration flows must track stacked $V_{DS,sat}$ budgets explicitly.

## Going deeper

Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., 2017), chapter 5.
