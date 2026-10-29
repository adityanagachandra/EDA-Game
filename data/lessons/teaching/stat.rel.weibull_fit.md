# Fitting Weibull life data

Life tests produce a list of failure times, often with some units still running when the test stops. Turning that list into Weibull parameters is the step that connects a reliability test to a lifetime prediction. The classical tool is the Weibull plot, a linearization that makes the fit a straight-line regression and makes departures from the model visible.

## The Weibull plot

Take logs twice. From $1 - F(t) = e^{-(t/\eta)^\beta}$,
$$\ln(-\ln(1 - F)) = \beta\ln t - \beta\ln\eta.$$
On axes $x = \ln t$ and $y = \ln(-\ln(1 - F))$, Weibull data fall on a line of slope $\beta$ that crosses $y = 0$ at $t = \eta$. For $F = 0.1$, $-\ln 0.9 = 0.1054$ and $y = -2.25$; for $F = 0.632$, $y \approx 0$. A bend in the plotted points signals a mixture of mechanisms (two slopes), a threshold time before which no failures occur, or a wrong distribution, which is exactly the diagnostic value of the plot.

## The shape from two points

The slope between two plotted points is the shape estimate:
$$\beta = \frac{y_2 - y_1}{\ln t_2 - \ln t_1}.$$
If 5 percent have failed at 200 hours and 30 percent at 500 hours, $y_1 = \ln(0.0513) = -2.970$, $y_2 = \ln(0.3567) = -1.031$, and $\beta = 1.939/0.916 = 2.12$, a wear-out mechanism. A full fit regresses all points, usually with $\ln t$ treated as the response because the time measurements carry the noise, and maximum likelihood is preferred when there are censored units, since regression cannot use the information that a unit survived to the end of the test.

## Plotting positions

The observed failure times do not come with $F$ values; they must be assigned. Sorting $n$ failure times and plotting the $i$-th at $F = i/n$ puts the last failure at $F = 1$, which has $y = +\infty$. Benard's median-rank approximation,
$$F_i \approx \frac{i - 0.3}{n + 0.4},$$
estimates the median of the $i$-th order statistic of a uniform sample and is the standard choice. For $n = 10$, the first failure plots at $0.7/10.4 = 0.067$ and the last at $9.7/10.4 = 0.933$. The choice matters most for small samples and for the lowest ranks, which are exactly the points that drive the early-life prediction.

## Going deeper

Abernethy, "The New Weibull Handbook" (5th ed., 2006); Meeker and Escobar, "Statistical Methods for Reliability Data" (1998), chapters 6 and 8.
