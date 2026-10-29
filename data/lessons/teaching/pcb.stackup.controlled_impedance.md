# Stackups and controlled impedance

A fast edge on a board trace travels as a wave on a transmission line made of the trace and its reference plane. Its characteristic impedance $Z_0 = \sqrt{L/C}$ per unit length is set by the trace width, the copper thickness, the dielectric height to the plane and the dielectric constant. Controlled impedance means choosing those numbers so every critical net has, say, 50 ohms single-ended, and having the fabricator build to that target. Board tools embed closed-form approximations and 2D field solvers for this; the closed forms show the trends.

## Microstrip impedance

A **microstrip** is a trace on an outer layer over a plane, with dielectric below and air above. IPC-2141 gives the approximation $Z_0 = \frac{87}{\sqrt{\varepsilon_r + 1.41}} \ln\frac{5.98 h}{0.8 w + t}$ for $h$ the height to the plane, $w$ the width and $t$ the copper thickness, all in one length unit, reasonable for $w/h$ between about 0.1 and 2. Example: $h = 5$ mil, $w = 8$ mil, $t = 1.4$ mil, $\varepsilon_r = 4.2$ gives $87/\sqrt{5.61} = 36.7$ and $\ln(29.9 / 7.8) = 1.344$, so $Z_0 \approx 49.4$ ohms. The formula is a fit, accurate to several percent; signoff uses a 2D field solver that also models solder mask, which lowers a microstrip's impedance by a few ohms.

## Stripline impedance

A **stripline** is a trace buried between two planes, fully surrounded by dielectric. For a trace centered between planes $b$ apart, IPC-2141 gives $Z_0 = \frac{60}{\sqrt{\varepsilon_r}} \ln\frac{1.9 b}{0.8 w + t}$. Example: $b = 16$ mil, $w = 5$ mil, $t = 0.7$ mil, $\varepsilon_r = 4.0$ gives $30 \cdot \ln(30.4 / 4.7) = 30 \cdot 1.867 \approx 56$ ohms. Two planes add capacitance, so a stripline needs a narrower trace than a microstrip for the same impedance, and its field is entirely in the dielectric, which slows it and removes the odd and even velocity difference that causes far-end crosstalk. Real stacks often place the trace off-center (asymmetric stripline) or route two signal layers between a pair of planes, which needs the field solver.

## Sensitivities and fabrication tolerance

The formulas encode the trends a designer must know: impedance **rises with height** (log of $h$ or $b$), **falls with width** and thickness (log of $0.8w + t$), and **falls with the dielectric constant** (as $1/\sqrt{\varepsilon_r}$ roughly). Because the dependence is logarithmic, a 1 mil error on a 4 mil dielectric moves impedance far more than on a 10 mil one. Fabricators quote controlled impedance with a tolerance, commonly 10%, and adjust trace width and prepreg choice within limits to hit it, which is why the drawing names target impedances and reference layers rather than freezing widths. A tool should store the target and recompute the width when the stackup changes, not store a width alone.

## Going deeper

IPC-2141A, "Design Guide for High-Speed Controlled Impedance Circuit Boards". Bogatin, "Signal and Power Integrity Simplified" (3rd ed., 2018), chapter 7.
