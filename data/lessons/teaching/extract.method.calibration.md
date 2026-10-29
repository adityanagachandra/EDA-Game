# Calibrating extraction to a golden solver

A fast extractor is trusted only because it has been compared with a golden field solver on many structures and adjusted until the two agree within a stated tolerance. That calibration step is a small fitting problem: measure the error, fit correction coefficients, and check what error remains. It is also the step that foundries qualify, and the one that most often needs repeating when a process, a layout style or a metal stack changes.

## Measuring the error

Run the fast extractor and the golden solver on the same set of test structures or nets, and compute each relative error $|x_i - g_i| / g_i$, with $x_i$ the fast value and $g_i$ the golden one. The figures that matter are the mean error, which reveals systematic bias, and the worst-case error, which bounds what a timer can trust. For four nets with fast and golden values of 10.2 vs 9.88, 5.6 vs 5.11, 7.4 vs 6.82 and 9.9 vs 9.21 fF, the worst error is 9.59% on the second net, and every net reads high: a sign of bias, not scatter. Foundries publish accuracy targets for qualified extractors, typically a few percent on total and coupling capacitance, measured exactly this way.

## Least-squares scale factor

The simplest correction multiplies fast values by one factor $k$ per pattern class. Choosing $k$ to minimize $\sum (k x_i - g_i)^2$ gives, by setting the derivative to zero, $k = \sum x_i g_i / \sum x_i^2$. It is least squares through the origin, a one-parameter linear regression. With five nets totalling $\sum x_i g_i = 294.03$ and $\sum x_i^2 = 300.99$, $k = 0.9769$, a 2.3% downward correction. Real calibration fits several coefficients per pattern (for example separate area, fringe and coupling terms) by linear least squares in the same way, using the normal equations or a QR factorization. A fit that minimizes relative rather than absolute error weights small nets more heavily, which suits timing better when small nets matter.

## Residual error after calibration

After applying $k$, recompute the errors. A single factor removes systematic bias but cannot remove scatter that varies from net to net. For four nets with fast values 11.1, 7.3, 7.2 and 5.4 fF against golden 10.77, 7.15, 7.34 and 5.4 fF, the fitted $k = 0.9854$ leaves a worst error of 3.34%, on the third net, which was low before calibration and is lower now. Residual scatter of this kind points to a pattern the library models poorly, and the fix is a new or finer pattern, not a different factor. Calibration must also be validated on nets not used in the fit; a factor that only matches its own training nets has been overfitted and will mislead the timer on the rest of the chip.

## Going deeper

Kao, Lo, Basel and Singh, "Parasitic extraction: current state of the art and future trends", Proceedings of the IEEE (2001); Golub and Van Loan, "Matrix Computations" (4th ed., 2013), on linear least squares.
