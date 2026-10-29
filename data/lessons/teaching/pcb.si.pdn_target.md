# Power delivery: target impedance and decoupling

A power delivery network (PDN) carries current from the regulator through planes, capacitors, package and die to every transistor. When the load current changes, the rail voltage moves by the current step times the PDN impedance. The target-impedance method turns a ripple budget into a maximum impedance over frequency; decoupling capacitors, planes and package capacitance are chosen to stay below it. Board PI tools compute the PDN impedance profile from these parts.

## Target impedance

If a rail of voltage $V$ may move by a fraction $r$ while the load current steps by $\Delta I$, the PDN impedance must satisfy $|Z(f)| \le Z_t = V r / \Delta I$ across the frequencies the current step contains. Example: a 1.0 V core with 3% ripple and a 10 A step needs $Z_t = 0.03/10 = 3$ milliohms. Lower voltages, tighter tolerances and larger currents all shrink the target, which is why modern processor rails need milliohm-level PDNs. The method is a simplification: it assumes a flat target and a worst-case current spectrum, and real designs refine it with the actual load current profile. It still sets the first-order budget that every decoupling choice is measured against.

## Capacitors and their resonance

A real capacitor is a series RLC: capacitance $C$, a small ESR and an inductance $L$ from its body and, often dominant, its mounting loop (pads, traces and vias to the planes). Its impedance falls as $1/(2\pi f C)$ to a minimum at the series resonance $f = 1/(2\pi\sqrt{LC})$ and rises as $2\pi f L$ above it. Example: 100 nF with 1 nH resonates at $1/(2\pi\sqrt{10^{-16}}) \approx 15.9$ MHz. Above resonance the capacitor is an inductor, so a larger value does not help at high frequency; only less mounting inductance does. Designers combine values to cover a band, and must watch for **anti-resonance** peaks where one capacitor's inductive region meets another's capacitive region in parallel.

## Counting capacitors for an inductive limit

Above resonance, $n$ identical capacitors in parallel present about $L/n$, so meeting a target $Z_t$ at frequency $f$ needs $2\pi f L / n \le Z_t$, or $n \ge 2\pi f L / Z_t$. Example: 0.8 nH per mounted capacitor, a 10 milliohm target at 50 MHz: one capacitor gives $2\pi \cdot 50 \times 10^6 \cdot 0.8 \times 10^{-9} \approx 251$ milliohms, so $n = 26$. That count is why boards carry dozens of capacitors per rail and why mounting inductance (via placement, via-in-pad, short pad connections) is engineered. Above roughly 100 MHz board capacitors become ineffective regardless of count, because the spreading inductance of planes and the package inductance dominate; on-package and on-die capacitance take over there.

## Going deeper

Swaminathan and Engin, "Power Integrity Modeling and Design for Semiconductors and Systems" (2007). Novak and Miller, "Frequency-Domain Characterization of Power Distribution Networks" (2007).
