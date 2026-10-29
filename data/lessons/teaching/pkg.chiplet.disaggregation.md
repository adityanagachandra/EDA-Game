# Chiplet disaggregation and yield

Disaggregation splits what would be one large die into several smaller chiplets, joined in a package. The main driver is yield: the probability that a die is free of killer defects falls exponentially with its area, so small dies waste less silicon per good system. The gain is offset by die-to-die interface area, assembly yield and package cost, and the trade-off is computed before any architecture is committed.

## Poisson yield of a large die

With killer defects scattered at random with density $D_0$ per unit area, the number on a die of area $A$ is Poisson with mean $A D_0$, and the die is good with probability $Y = e^{-A D_0}$. Keep units consistent: with $D_0$ per cm^2, convert the area from mm^2 by dividing by 100. A 600 mm^2 die at 0.1 defects per cm^2 has $Y = e^{-0.6} = 54.88\%$. The same process yields $e^{-0.15} = 86.07\%$ for a 150 mm^2 die. Real fabs fit clustered models such as negative binomial, which are less pessimistic for large dies, but the exponential is the standard first estimate and the one used here.

## Chiplet yield with interface overhead

Splitting the die into $k$ chiplets does not divide the area exactly: each chiplet needs die-to-die PHYs and duplicated infrastructure, an overhead fraction $o$ of its logic area. Each chiplet has area $a = A(1+o)/k$ and yield $e^{-a D_0}$. Splitting the 600 mm^2 design into four chiplets with 10% overhead gives $a = 165$ mm^2 and a yield of $e^{-0.165} = 84.79\%$. Chiplets are tested before assembly, so bad chiplets are discarded individually rather than taking three good ones with them. That independence is the point; the overhead is the price, and it grows with the number of cuts, so splitting into many tiny chiplets eventually loses.

## Silicon per good system

Compare the manufactured silicon needed for one good system. Monolithic: $A/Y_{mono}$. Chiplets with known-good-die testing: $k a / Y_{chip} = A(1+o)/Y_{chip}$. For the example, monolithic needs $600/0.5488 = 1093$ mm^2 and chiplets $660/0.8479 = 778$ mm^2, a ratio of 1.40 in favour of chiplets. This ignores assembly yield, package cost and test cost, all of which favour monolithic, and the performance and energy cost of crossing die boundaries. It also ignores the reuse benefit: a chiplet design can serve several products, and an I/O chiplet can stay on an older, cheaper process. Real decisions add those terms to the same comparison.

## Going deeper

Stow, Xie, Siddhu and Loh, "Cost-effective design of scalable high-performance systems using active and passive interposers", ICCAD (2017); Feng and Ma, "Chiplet actuary", DAC (2022).
