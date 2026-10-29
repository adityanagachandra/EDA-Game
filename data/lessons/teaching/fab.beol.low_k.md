# Low-k dielectrics and air gaps

Wire capacitance is proportional to the permittivity of the insulator around it. Replacing silicon dioxide, with relative dielectric constant about 3.9 to 4.2, by carbon-doped or porous films with $k$ near 2.5 to 3 cuts capacitance, delay and dynamic power. The cost is mechanical weakness and integration difficulty, and real stacks mix low-k films with higher-k caps and etch stops.

## Capacitance scales with k

At fixed geometry, every capacitance between conductors scales with the relative permittivity $k$ of the dielectric between them. A wire with 0.20 fF/um in oxide of $k = 4.0$ would have $0.20 \times 2.7/4.0 = 0.135$ fF/um in a $k = 2.7$ film. Since delay of a wire-dominated path is roughly proportional to $RC$ and dynamic power to $C V^2 f$, a lower $k$ helps both. Low-k films achieve it by adding carbon and hydrogen (SiCOH) to break up the polar Si-O network, and ultra-low-k films add porosity, since air has $k = 1$.

## Effective k of a real stack

Between two metal levels there is not one film but several: an etch-stop or copper cap of higher $k$, often 4.5 to 7, and the low-k bulk. For vertical coupling these act as capacitors in series, so $k_{eff} = (t_1 + t_2)/(t_1/k_1 + t_2/k_2)$. A 10 nm cap at $k = 5.5$ under 90 nm of low-k at $k = 2.6$ gives $100/(1.82 + 34.6) = 2.75$, higher than the bulk film. A thin high-k layer costs more than its thickness suggests. Lateral coupling between neighbours on the same layer sees mostly the low-k fill, plus any liner or damage layer on the trench sidewalls.

## Air gaps and the limits of low-k

Air gaps remove dielectric between closely spaced lines on selected layers, lowering the effective lateral $k$ toward 1. If a fraction $f$ of a wire's capacitance is lateral and the gap lowers lateral $k$ from $k$ to $k_g$, total capacitance falls by $f(1 - k_g/k)$. With $f = 0.6$, $k = 2.7$ and $k_g = 1.4$, the saving is $0.6 \times 0.48 = 28.9\%$. Porous and gapped films are mechanically weak, crack under packaging stress, absorb moisture and are damaged by plasma etch, which raises their $k$ again. Extraction must know where gaps exist, usually through layer-specific rules, and reliability rules limit where they can be used.

## Going deeper

Grill, Gates, Ryan, Nguyen and Priyadarshini, "Progress in the development and understanding of advanced low k and ultralow k dielectrics for very large-scale integrated interconnects", Applied Physics Reviews (2014).
