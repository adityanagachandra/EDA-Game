# Thermal oxidation and the Deal-Grove model

Heating silicon in oxygen or steam grows silicon dioxide that is consumed from the wafer itself. Thermal oxide was the gate dielectric for decades and still forms interface layers, pad oxides and isolation liners. The Deal-Grove model predicts its thickness from time and temperature with one quadratic equation. It is also the first model a process simulator implements, and it shows how a diffusion-limited process slows down as it proceeds.

## Oxide grows into the silicon

Oxidation converts silicon into SiO2, and the oxide occupies more volume than the silicon it replaces. For every nanometre of oxide grown, about 0.44 nm of silicon is consumed, so 44% of the oxide lies below the original silicon surface and 56% above it. Growing 100 nm of oxide eats 44 nm of silicon. This matters for geometry: features defined in silicon before oxidation shrink, and a step in the silicon surface reappears in the oxide. Dry oxygen grows dense, slow oxide suited to thin, high-quality layers; steam (wet oxidation) grows several times faster and is used for thicker films.

## The Deal-Grove equation

Oxidant must diffuse through the existing oxide and then react at the silicon interface. Balancing those two fluxes gives $x^2 + A x = B (t + \tau)$, where $x$ is oxide thickness, $B$ is the parabolic rate constant, $B/A$ the linear rate constant and $\tau$ accounts for oxide already present. Solve the quadratic: $x = \left(-A + \sqrt{A^2 + 4B(t+\tau)}\right)/2$. With $A = 0.2$ um, $B = 0.1$ um^2/h, $\tau = 0$ and $t = 2$ h: $x = (-0.2 + \sqrt{0.04 + 0.8})/2 = 0.358$ um, which is 358 nm.

## Linear and parabolic regimes

For thin oxides, $x^2 \ll A x$ and growth is linear: $x \approx (B/A)(t + \tau)$, limited by the interface reaction. For thick oxides, $x^2 \gg A x$ and $x \approx \sqrt{B t}$, limited by diffusion through the oxide. To find the time for a target thickness, rearrange: $t = (x^2 + A x)/B - \tau$. With $A = 0.2$ um and $B = 0.1$ um^2/h, 0.5 um needs $(0.25 + 0.1)/0.1 = 3.5$ h. Doubling the target to 1 um needs $(1 + 0.2)/0.1 = 12$ h, more than three times as long. The classic model is known to underpredict the fast initial growth of very thin dry oxides, which is why thin gate interfaces are calibrated separately.

## Going deeper

Deal and Grove, "General Relationship for the Thermal Oxidation of Silicon", Journal of Applied Physics (1965); Plummer, Deal and Griffin, "Silicon VLSI Technology", chapter 6.
