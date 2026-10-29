# Main effects and interactions

Once the runs of a two-level design are simulated, the analysis is arithmetic on contrasts. A main effect says how much the response changes when one factor goes from low to high; an interaction says how much that change depends on another factor. These numbers are what decide which process parameters a sensitivity model must keep and whether a simple additive model is enough.

## Main effects

Code each factor's levels as -1 and $+1$. In a $2^2$ design with responses $y_{--}, y_{+-}, y_{-+}, y_{++}$ (first sign for A, second for B), the main effect of A is the average response at A high minus the average at A low:
$$E_A = \frac{y_{+-} + y_{++}}{2} - \frac{y_{--} + y_{-+}}{2}.$$
If a delay is 50, 58, 54 and 66 ps, $E_A = 62 - 52 = 10$ ps and $E_B = 60 - 54 = 6$ ps. Every run is used in every estimate, which is why factorial designs are efficient: each effect is averaged over all levels of the other factors rather than estimated one factor at a time.

## Interactions

The AB interaction compares the runs where the product of the coded levels is $+1$ with those where it is -1:
$$E_{AB} = \frac{y_{--} + y_{++}}{2} - \frac{y_{+-} + y_{-+}}{2}.$$
In the example, $E_{AB} = 58 - 56 = 2$ ps. Equivalently, the effect of A is $66 - 54 = 12$ ps when B is high and $58 - 50 = 8$ ps when B is low, and the interaction effect is half that difference. A large interaction means a model that adds separate parameter effects will be wrong, and in a variation model it means a linear sensitivity at nominal under-reports risk at the combined extremes.

## The fitted regression model

For a two-level design in coded units, the least-squares model $y = b_0 + b_A x_A + b_B x_B + b_{AB} x_A x_B$ has a closed form: $b_0$ is the grand mean and each other coefficient is half its effect, because the coded variable spans 2 units from -1 to $+1$. In the example, $b_0 = 57$, $b_A = 5$, $b_B = 3$ and $b_{AB} = 1$, and the prediction at $x_A = 0.5$, $x_B = -0.5$ is $57 + 2.5 - 1.5 - 0.25 = 57.75$ ps. The columns of a full factorial are orthogonal, so each coefficient is estimated independently of the others and adding or dropping a term does not change the rest. The model is exact at the four runs and assumes linearity between them.

## Going deeper

Montgomery, "Design and Analysis of Experiments" (9th ed., 2017), chapter 6; Box, Hunter and Hunter, "Statistics for Experimenters" (2nd ed., 2005), chapter 5.
