# Dennard scaling and its end

In 1974 Dennard and colleagues showed that shrinking every MOSFET dimension and the supply voltage by the same factor keeps the electric field constant, makes circuits faster and keeps power density flat. For three decades this constant-field scaling, together with Moore's observation of doubling transistor counts, drove the industry. Its end, when voltage stopped scaling, is why modern chips are power-limited and why EDA tools optimize power as hard as timing.

## Constant-field scaling rules

Scale every linear dimension (length, width, oxide thickness, junction depth) and the supply voltage by $s < 1$, and increase doping by $1/s$. Then the electric field is unchanged and:
voltage and current scale by $s$;
capacitance scales by $s$ (area $s^2$ over thickness $s$);
gate delay $CV/I$ scales by $s$;
power per circuit $VI$ scales by $s^2$;
area per circuit scales by $s^2$;
power density stays constant;
switching energy $CV^2$ scales by $s^3$.

Example: $s = 0.7$ for one generation. Delay becomes 0.7 of its old value, power per circuit 0.49, energy per switch 0.343, while power density is unchanged.

The constant field kept reliability (oxide stress, hot carriers) under control. Every generation delivered faster, cheaper, lower-energy circuits at no cost in heat per area, which is why it was such a powerful engine.

## Moore's law and density

Moore's observation (1965, revised in 1975) was that the number of components per chip at minimum cost was doubling at a regular interval, which settled near two years. Constant-field scaling supplied the means: shrinking linear dimensions by 0.7 per generation multiplies density by $1/0.49 = 2.04$.

Transistor count grows as
$$N(t) = N_0 \cdot 2^{t/T}.$$
Example: 2 billion transistors with $T = 2$ years gives $2 \times 2^{3} = 16$ billion after 6 years.

Moore's law is an economic and industrial observation, not a physical law, and the doubling period has varied. Density growth continued after Dennard scaling ended, through FinFETs, nanosheets, design-technology co-optimization (smaller cells from fewer tracks, buried or backside power rails) and multi-die packaging. But the benefits per transistor changed: density kept growing while speed and energy per operation improved far more slowly.

## Why Dennard scaling ended

Constant-field scaling requires the threshold voltage to scale with the supply. But off-current depends on $V_t$ through the subthreshold swing, which is fixed at no better than about 60 mV/decade at room temperature. Lowering $V_t$ by 60 to 70 mV costs roughly ten times more leakage. Around the mid-2000s, with $V_{DD}$ near 1 V, further threshold reduction would have made leakage unacceptable, so supply voltage stopped scaling.

With $V_{DD}$ fixed, capacitance per gate scaling by $s$ and frequency rising by $1/s$, the power per gate $CV^2f$ stays constant while density rises by $1/s^2$. Power density then grows by about $2\times$ per generation, so clock frequency plateaued instead.

Example: two generations at $s = 0.7$ with fixed voltage: power density $\times (1/0.49)^2 = 4.2$.

The result was the power wall: multicore designs, aggressive clock and power gating, multiple threshold flavors, dynamic voltage and frequency scaling, and power as a first-class objective in synthesis, placement and timing optimization.

## Going deeper

Dennard and others, "Design of ion-implanted MOSFET's with very small physical dimensions", IEEE Journal of Solid-State Circuits (1974); Moore, "Cramming more components onto integrated circuits", Electronics (1965); Bohr, "A 30 year retrospective on Dennard's MOSFET scaling paper", IEEE SSCS Newsletter (2007).
