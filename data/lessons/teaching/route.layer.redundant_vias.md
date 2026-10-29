# Redundant via insertion

A via is a small hole filled with metal, and a fraction of them fail open in manufacturing. A chip has billions of vias, so even a tiny failure probability adds up. Redundant via insertion places a second cut next to the first wherever space allows, so the connection fails only if both cuts fail. It is a standard post-routing step for yield.

## Expected via failures

If each of $N$ vias fails independently with probability $p$, the expected number of failures is $Np$ by linearity of expectation. Example: 15 million single-cut vias with $p = 5 \times 10^{-7}$ give 7.5 expected failures per die. Any one failure can kill the die if it opens a net, so this number is a direct yield risk.

Failure probabilities for single vias depend on the process, the via size and the surrounding layout; the independence assumption is a first-order model. Even so, it shows the scale: with $N$ in the billions, $p$ has to be extraordinarily small, or the vias must be made redundant.

## Doubling vias

A doubled via opens only if both of its cuts fail, with probability about $p^2$, which is negligible. If a fraction $r$ of $N$ via locations is doubled, the expected number of opens is $N[(1 - r)p + r p^2] \approx N(1 - r)p$. Example: 40 million vias at $p = 10^{-7}$ give 4 expected opens; doubling 52% of them leaves about 1.92.

So the benefit is set almost entirely by the doubling rate, and the remaining single vias dominate. That is why redundant via insertion tries to maximise the doubled fraction. Each candidate second cut needs room on both metal layers, so the problem is a maximum independent set or matching over conflicting candidates, solved heuristically or with integer programming on small windows.

## Via-limited yield

The probability that no via opens is $Y = (1 - p)^{N_s}(1 - p^2)^{N_d}$ for $N_s$ single and $N_d$ doubled vias, which is close to $e^{-N_s p}$. Example: 21 million vias at $p = 2 \times 10^{-7}$: with no doubling $Y = e^{-4.2} \approx 1.5\%$; doubling 63% leaves 7.77 million single vias and $Y \approx e^{-1.55} \approx 21\%$.

The router helps before insertion even runs. Via-aware detailed routing leaves space next to vias for a second cut, prefers via positions where a double fits, and avoids long via stacks. The cost is a little extra metal and some routing flexibility, which is why the doubling rate is a quality metric reported next to wirelength and DRC count.

## Going deeper

K.-Y. Lee and T.-C. Wang, "Post-Routing Redundant Via Insertion for Yield/Reliability Improvement" (2006); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.), on detailed routing.
