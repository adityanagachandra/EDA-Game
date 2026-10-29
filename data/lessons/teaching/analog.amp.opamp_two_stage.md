# The two-stage op-amp

The Miller-compensated two-stage op-amp is the canonical benchmark circuit of analog design automation: a differential input stage for gain and common-mode rejection, a common-source second stage for more gain and output swing, and a compensation capacitor for stability. Almost every sizing paper and tool demonstrates on it, so its gain, bandwidth and pole structure are the first things a tool designer should be able to compute.

## Gain of two cascaded stages

The first stage is a differential pair with mirror load, gain $A_1 = g_{m1}(r_{o2} \parallel r_{o4})$. The second is a common-source stage with a current-source load, gain $A_2 = g_{m6}(r_{o6} \parallel r_{o7})$. The DC open-loop gain is $A_0 = A_1 A_2$, and in dB the stage gains add. With $A_1 = 40$ and $A_2 = 50$, $A_0 = 2000$, which is 66 dB. Each stage contributes roughly half an intrinsic gain, so in a process with intrinsic gain around 50, two stages land between 60 and 70 dB. Higher gain needs cascoded stages, longer channels or a third stage, each with a cost in swing, speed or stability.

## Gain-bandwidth set by gm1 and Cc

With Miller compensation the compensation capacitor $C_c$ between the second stage's input and output dominates the frequency response. The first stage's output current $g_{m1}v_{id}$ flows into $C_c$, which behaves as an integrator near crossover, so the unity-gain frequency is $\mathrm{GBW} = g_{m1}/(2\pi C_c)$. With $g_{m1} = 300\ \mu$S and $C_c = 1$ pF, $\mathrm{GBW} = 47.7$ MHz (microsiemens over picofarads is $10^6$ per second). The second stage does not appear: as long as the non-dominant poles are far above GBW, the product of gain and bandwidth is fixed by these two quantities. That is why sizing procedures start from GBW and $C_c$ to set $g_{m1}$, then pick $g_m/I_D$ to set the input pair's current.

## The dominant pole

A single-pole response with DC gain $A_0$ and gain-bandwidth GBW has its dominant pole at $f_{p1} = \mathrm{GBW}/A_0$. For 80 dB ($A_0 = 10^4$) and 50 MHz, $f_{p1} = 5$ kHz. The pole sits at the first-stage output, where the Miller-multiplied capacitance $C_c(1 + A_2)$ meets the first stage's output resistance. That the dominant pole is so low, kilohertz for a megahertz amplifier, is what makes the op-amp useful: it trades open-loop bandwidth for a gain that feedback converts into accuracy. The remaining poles, from the output node and the mirror, must sit well above GBW for the loop to be stable, which is the subject of compensation and phase margin.

## Going deeper

Allen and Holberg, "CMOS Analog Circuit Design" (3rd ed., 2011), chapter 6; Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., 2017), chapter 10.
