# Time-dependent dielectric breakdown

Time-dependent dielectric breakdown is the wear-out of a gate dielectric under sustained field: defects accumulate until a conducting path forms and the gate shorts. Breakdown times are statistical and follow a Weibull distribution, and a chip fails at its weakest gate, so a chip's lifetime depends on its total gate area. Reliability limits on supply voltage, and overdrive limits in design rules, come from these statistics.

## The Weibull failure distribution

Breakdown times of nominally identical capacitors follow a Weibull distribution:
$$F(t) = 1 - \exp\left[-\left(\frac{t}{\eta}\right)^{\beta}\right],$$
where $F$ is the cumulative fraction failed, $\eta$ is the characteristic life (the time at which 63.2% have failed) and $\beta$ is the shape parameter. Thin dielectrics have small $\beta$, often near 1 to 2, meaning a wide spread.

Example: $\eta = 100$ years and $\beta = 1.5$. At 10 years, $(0.1)^{1.5} = 0.0316$ and $F = 1 - e^{-0.0316} = 3.1\%$.

Plotting $\ln(-\ln(1 - F))$ against $\ln t$ gives a straight line of slope $\beta$, the Weibull plot reliability engineers use to fit data. The Weibull form arises naturally from weakest-link statistics: a dielectric fails when its first percolation path completes, which is the minimum of many independent local failure times.

## Time to a target failure fraction

Products are specified at low failure fractions, such as 0.01% or 0.1%, not at 63%. Inverting the Weibull distribution,
$$t_F = \eta\left[-\ln(1 - F)\right]^{1/\beta}.$$
For small $F$, $-\ln(1 - F) \approx F$.

Example: $\eta = 100$ years, $\beta = 1.5$, $F = 0.1\%$. $t_F = 100 \times (0.0010005)^{1/1.5} = 100 \times 0.01001 = 1.0$ year.

The low-percentile lifetime is far below $\eta$, and the smaller $\beta$ is, the further below. That is the danger of thin dielectrics: their wide distributions push early failures much earlier for a given characteristic life. Reliability qualification therefore extracts both $\eta$ and $\beta$ carefully, since an error in $\beta$ moves the low-percentile prediction by orders of magnitude.

## Scaling from test structures to a chip

Test capacitors are small; a chip has a large total gate area. If breakdown is a weakest-link process over area, the Weibull characteristic life scales as
$$\eta_{chip} = \eta_{test}\left(\frac{A_{test}}{A_{chip}}\right)^{1/\beta},$$
with $\beta$ unchanged.

Example: test structures of $10^{-4}$ cm$^2$ with $\eta = 5000$ years at use conditions and $\beta = 1.5$. A chip with 0.1 cm$^2$ of gate area has $\eta_{chip} = 5000 \times (10^{-3})^{0.667} = 5000 \times 0.01 = 50$ years.

Area scaling is why chip reliability is worse than device reliability, and why it worsens as chips grow. With small $\beta$ the penalty is severe. For EDA, this links reliability to design: total biased gate area, the fraction of time gates are under stress, and local overvoltage (such as from ringing on I/O) all enter. Reliability signoff tools check voltage limits per device type and stress time, derived from these models.

## Going deeper

Degraeve and others, "New insights in the relation between electron trap generation and the statistical properties of oxide breakdown", IEEE Transactions on Electron Devices (1998); JEDEC JEP122, "Failure Mechanisms and Models for Semiconductor Devices".
