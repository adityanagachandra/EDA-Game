# Inductance and return paths

Inductance is a property of a current loop, not of a wire on its own: it measures the magnetic flux a loop links per ampere. On-chip and package inductance depends on where the return current flows, which is why inductance extraction is harder than capacitance extraction. Closed forms for coaxial lines, wire pairs and wires over planes, plus the link between L and C in a uniform dielectric, give a designer reliable anchors.

## Self-inductance of a coaxial line

Self-inductance is $L = \Lambda/I$, the flux linkage per unit current. For a coaxial line with inner radius $a$ and outer conductor inner radius $b$, Ampere's law gives $B = \mu_0 I/(2\pi r)$ between the conductors. Integrating over $a < r < b$ gives the flux per unit length, so
$$L' = \frac{\mu_0}{2\pi}\ln\frac{b}{a},$$
assuming the current flows on the facing surfaces, as it does at high frequency.

With $\mu_0/2\pi = 2 \times 10^{-7}$ H/m $= 0.2$ nH/mm, a coax with $b/a = 3.5$ has $L' = 0.2\ln 3.5 = 0.251$ nH/mm. The logarithm makes inductance insensitive to geometry: doubling the outer radius adds only $0.2\ln 2 = 0.139$ nH/mm. This weak dependence is a recurring theme, and it is why inductance is hard to reduce by sizing alone.

## Two-wire and wire-over-plane loops

For two parallel round wires of radius $a$ with centres $D$ apart, carrying a current out and back, the exact external loop inductance per unit length is
$$L' = \frac{\mu_0}{\pi}\cosh^{-1}\frac{D}{2a} \approx \frac{\mu_0}{\pi}\ln\frac{D}{a} \quad (D \gg a).$$
A wire of radius $a$ with its centre at height $h$ over a perfectly conducting plane is half of that problem by image theory: $L' = \frac{\mu_0}{2\pi}\cosh^{-1}(h/a)$.

Example: $a = 0.1$ mm and $D = 1$ mm give $L' = 0.4\cosh^{-1} 5 = 0.917$ nH/mm. A wire with $a = 0.1$ mm at $h = 0.4$ mm over a plane has $0.2\cosh^{-1}4 = 0.413$ nH/mm. Moving the return closer shrinks the loop and the inductance. That is the physics of return paths: a ground plane right under a trace gives the smallest loop.

## Inductance from capacitance

For a TEM line whose conductors are surrounded by a single homogeneous dielectric, inductance and capacitance per unit length are linked by
$$L'C' = \mu_0\varepsilon_0\varepsilon_r = \frac{\varepsilon_r}{c^2}.$$
So once a 2D capacitance solver has found $C'$, the external inductance follows without a magnetic solve: $L' = \varepsilon_r/(c^2 C')$.

Example: a line in FR-4 ($\varepsilon_r = 4.3$) with $C' = 130$ pF/m has $L' = 4.3/(9 \times 10^{16} \cdot 1.3 \times 10^{-10}) = 368$ nH/m, and hence $Z_0 = \sqrt{L'/C'} = 53$ ohm. A common trick in 2D extractors is to compute the capacitance a second time with the dielectric replaced by vacuum, $C'_0$, and use $L' = 1/(c^2 C'_0)$, which also works for inhomogeneous dielectrics because inductance does not depend on permittivity.

## Going deeper

Paul, "Inductance: Loop and Partial", Wiley, 2010; Ramo, Whinnery and Van Duzer, "Fields and Waves in Communication Electronics" (3rd ed.), chapters 2 and 5.
