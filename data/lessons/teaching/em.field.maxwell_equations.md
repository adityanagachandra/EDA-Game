# Maxwell's equations in integral form

Maxwell's equations are the complete classical description of electric and magnetic fields, and every field solver in EDA discretizes some subset of them. Knowing what each law constrains, and which ones survive in the quasi-static limits, tells a tool designer which equations a capacitance extractor, an inductance extractor or a full-wave package solver must actually solve, and what each approximation throws away.

## The four laws and what each constrains

In integral form, for surfaces $S$ bounded by curves $\partial S$ and closed surfaces enclosing a volume:
$$\oint \mathbf{D}\cdot d\mathbf{A} = Q_{enc}, \qquad \oint \mathbf{B}\cdot d\mathbf{A} = 0,$$
$$\oint_{\partial S} \mathbf{E}\cdot d\mathbf{l} = -\frac{d}{dt}\int_S \mathbf{B}\cdot d\mathbf{A}, \qquad \oint_{\partial S} \mathbf{H}\cdot d\mathbf{l} = \int_S \left(\mathbf{J} + \frac{\partial\mathbf{D}}{\partial t}\right)\cdot d\mathbf{A}.$$
Gauss's law ties electric flux to charge; Gauss's law for magnetism says magnetic field lines close; Faraday's law says a changing magnetic flux induces a circulating electric field; the Ampere-Maxwell law says currents and changing electric flux produce a circulating magnetic field. In linear media $\mathbf{D} = \varepsilon\mathbf{E}$ and $\mathbf{B} = \mu\mathbf{H}$. In statics, Faraday's law reduces to $\oint \mathbf{E}\cdot d\mathbf{l} = 0$, which is why voltage is path independent and KVL works.

## Conduction and displacement current

The Ampere-Maxwell law carries two current terms: conduction current density $\mathbf{J} = \sigma\mathbf{E}$ and displacement current density $\partial\mathbf{D}/\partial t$. Displacement current is what flows across a capacitor's gap, where no charge moves; without it the law would give different answers for surfaces cut through the wire and through the gap. Taking the divergence of the Ampere-Maxwell law and using Gauss's law gives the continuity equation $\nabla\cdot\mathbf{J} + \partial\rho/\partial t = 0$, of which KCL is the lumped form.

In sinusoidal steady state the two terms compare as $\sigma$ versus $\omega\varepsilon$. In copper, $\sigma = 5.8 \times 10^7$ S/m while $\omega\varepsilon_0$ at 10 GHz is about 0.56 S/m, so displacement current inside a conductor is negligible. In a good dielectric, $\sigma$ is near zero and displacement current carries everything. Silicon substrate sits in between, which is why its behavior changes with frequency.

## Which laws an extractor solves

When structures are much smaller than a wavelength, the coupling terms between electric and magnetic fields are small and the problem splits. Capacitance extraction is electrostatic: Gauss's law with a curl-free $\mathbf{E} = -\nabla\phi$, giving $\nabla\cdot(\varepsilon\nabla\phi) = 0$ with conductor potentials as boundary conditions. Inductance and resistance extraction is magnetoquasistatic: Ampere's law without displacement current, often through the vector potential and partial inductances, plus current distribution inside conductors at high frequency.

A full-wave solver keeps all four laws coupled, which is needed when a structure is a significant fraction of a wavelength: package planes, connectors, vias at tens of GHz, antennas. On-chip RC extraction, by contrast, almost never needs full-wave physics. Choosing the reduced formulation is the single biggest speed lever in EM tooling, and knowing its limits is what keeps the fast answer correct.

## Going deeper

Griffiths, "Introduction to Electrodynamics" (4th ed.), chapter 7; Haus and Melcher, "Electromagnetic Fields and Energy", chapters 1 and 3 on quasistatics.
