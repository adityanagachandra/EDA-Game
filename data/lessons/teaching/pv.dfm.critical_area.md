# Critical area of a layout

Random particle defects cause shorts and opens. Whether a defect of a given size kills the chip depends on where it lands: the critical area is the area of positions where its centre causes a failure. Averaged over the defect size distribution and multiplied by defect density, it predicts yield. Critical-area analysis turns layout geometry into a yield number and drives DFM fixes such as wire spreading.

## Critical area at one defect size

For two long parallel wires with gap $s$ and run length $L$, a circular conducting defect of diameter $x$ shorts them when it touches both. For $x \le s$ that is impossible; for $x > s$ its centre must fall in a band of width $x - s$ between the wires, so $A_c(x) = L (x - s)$, ignoring end effects. With $L = 150$ um, $s = 135$ nm and $x = 220$ nm, the band is 85 nm wide and $A_c = 12.75$ square um. Opens are the dual: a missing-material defect of diameter $x$ cuts a wire of width $w$ when $x > w$, with $A_c = L(x - w)$.

## Averaging over defect sizes

Defect sizes follow a distribution, commonly $f(x) = 2 x_0^2 / x^3$ for $x \ge x_0$. The average short critical area of one gap is $\int_s^\infty L (x - s) f(x)\, dx = L x_0^2 / s$ for $s \ge x_0$. A bus of 4 wires has 3 gaps; with $L = 470$ um, $x_0 = 25$ nm and $s = 75$ nm the average is $3 \cdot 470 \cdot 625 / 75 / 1000 = 11.75$ square um, counting each gap separately. The $1 / s$ dependence explains why widening gaps pays off most where they are smallest.

## From critical area to yield

With defect density $D$ for a mechanism and average critical area $A_c$, the Poisson model gives yield $Y = e^{-D A_c}$. Units must agree: 13.3 square mm is 0.133 square cm, so with $D = 0.39$ defects per square cm, $D A_c = 0.05187$ and $Y = 94.95$ percent. Yields of independent mechanisms and layers multiply, so a chip's random-defect yield is the product over layers of shorts and opens. Clustered defects are modelled with negative binomial yield instead. Critical area is computed by Boolean and sizing operations, so it runs in the same engines as DRC.

## Going deeper

Stapper, "Modeling of integrated circuit defect sensitivities" (IBM Journal of Research and Development, 1983).
