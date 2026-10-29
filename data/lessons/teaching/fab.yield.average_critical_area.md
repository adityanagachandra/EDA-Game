# Average critical area over defect sizes

Critical area depends on defect size, and defects come in a distribution of sizes. Averaging critical area over that distribution gives one number per layer and mechanism, the average critical area, which multiplies the total defect density to give the expected fault count. For the common $1/x^3$ size distribution, a wire pair has a clean closed form that makes the dependence on spacing explicit.

## The defect size distribution

Measured defect size distributions fall off steeply for sizes above the smallest that inspection or physics produces. A standard model is $f(x) = 2x_0^2/x^3$ for $x \ge x_0$, normalized because $\int_{x_0}^\infty 2x_0^2 x^{-3}\, dx = 1$. The fraction of defects larger than $x$ is $P(X > x) = (x_0/x)^2$. With $x_0 = 20$ nm, 25% of defects exceed 40 nm and only 4% exceed 100 nm. Small defects are common but usually harmless; large ones are rare but can short wires that are far apart. The weighting of the two is what the average captures.

## Averaging the short critical area

For two wires with spacing $s \ge x_0$ and parallel length $L$, the average short critical area is $\bar A = \int_s^\infty L(x - s) \frac{2x_0^2}{x^3} dx$. The antiderivative of $(x - s)/x^3$ is $-1/x + s/(2x^2)$, which runs from $-1/(2s)$ at $x = s$ to 0 at infinity, so the integral is $1/(2s)$ and $\bar A = L x_0^2/s$. For $L = 1000$ um, $x_0 = 20$ nm and $s = 40$ nm: $\bar A = 1000 \times 400/40$ nm um $= 10{,}000$ nm um $= 10$ um^2. Average critical area is inversely proportional to spacing: doubling the space halves it.

## From average critical area to yield

The expected number of faults on a layer is $\lambda = D_0 \bar A$, where $D_0$ is the total density of the defect type, and with Poisson statistics $Y = e^{-\lambda}$. Suppose a die has 100 m of minimum-spaced wire at 40 nm spacing and $x_0 = 20$ nm. Then $\bar A = 10^8$ um $\times 400/40$ nm $= 10^9$ nm um $= 10^6$ um^2 $= 0.01$ cm^2, and with $D_0 = 0.5$ per cm^2, $\lambda = 0.005$ and $Y = 99.5\%$ for this layer and mechanism. Multiplying such yields over layers and mechanisms gives the random-defect yield. This calculation is why wire spreading targets minimum-spaced runs: they dominate $\bar A$.

## Going deeper

Stapper, "Modeling of integrated circuit defect sensitivities" (1983); Ferris-Prabhu, "Role of defect size distribution in yield modeling", IEEE Transactions on Electron Devices (1985).
