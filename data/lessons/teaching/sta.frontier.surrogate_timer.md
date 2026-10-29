# Surrogate fast timers in an agent loop

An agent exploring placements, sizings or ECOs may evaluate thousands of candidates, too many for a signoff timer. A fast timer, a reduced-accuracy engine or a learned surrogate, can rank candidates in milliseconds. The engineering question is when its answer can be trusted and when the loop must recurse to the accurate timer. A validated error bound turns that into arithmetic.

## Uncertain endpoints

Suppose validation shows the surrogate's slack error never exceeds $\varepsilon$ on designs like this one. The true slack of an endpoint with predicted slack $p$ then lies in $[p - \varepsilon, p + \varepsilon]$. If $p > \varepsilon$, the endpoint surely passes; if $p < -\varepsilon$, it surely fails; if $|p| < \varepsilon$, its sign is uncertain, and the accurate timer must decide. With predictions 7.5, -2.5, 2.5, 46.5, 33.5, 50.5, -31.5, -26.5 and -23.5 ps and $\varepsilon = 15$ ps, three endpoints are uncertain. Re-timing only those, rather than the whole design, is what makes a fast-then-accurate loop efficient.

## Bounding the true WNS

The same bound gives guarantees on summaries. The true worst slack is at least the smallest prediction minus $\varepsilon$: if the worst prediction is -37.5 ps and $\varepsilon = 4$ ps, the true WNS is no lower than -41.5 ps. An agent can use such bounds to prune candidates, rejecting any whose best possible WNS is worse than the current solution's, without accurate timing. The bound holds only on the distribution where the error was validated; on a design unlike the validation set, a learned surrogate's error can be much larger, and the loop must detect drift, for example by re-timing a random sample of endpoints accurately and checking the error.

## Choosing the error bound

Often validation gives an error distribution rather than a hard maximum: say a sigma of 6 ps. The loop then chooses a multiple $k$, so $\varepsilon = k\sigma$. With $k = 3$, $\varepsilon = 18$ ps, and every endpoint with $|p| < 18$ ps is re-timed. A larger $k$ is safer and costs more accurate runs; a smaller $k$ is cheaper and occasionally trusts a wrong sign. The final decision is never left to the surrogate: acceptance of a change rests on the accurate engine. Reward functions built on surrogate slack invite gaming, since an optimizer will find inputs where the surrogate is wrong in its favor, which is another reason to re-time accepted candidates accurately.

## Going deeper

Kahng, "Machine learning applications in physical design: recent results and directions" (ISPD, 2018).
