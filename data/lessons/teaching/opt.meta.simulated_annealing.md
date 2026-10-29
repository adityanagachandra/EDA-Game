# Simulated annealing

Simulated annealing is local search that sometimes accepts a worse move, with a probability that falls as a temperature parameter cools. Early on it explores freely; late, it behaves like descent. It powered TimberWolf placement, classic floorplanners on Polish expressions and sequence pairs, and VPR's FPGA placer, which still uses it. Its behaviour is governed by one formula, the Metropolis acceptance probability.

## The Metropolis acceptance rule

Propose a random move with cost change $\Delta$. If $\Delta \le 0$ accept it. If $\Delta > 0$ accept it with probability
$$P = e^{-\Delta / T}.$$
A placement swap that raises HPWL by 34 um at $T = 20$ um is accepted with probability $e^{-1.7} = 0.183$. The same swap at $T = 4$ is accepted with probability $e^{-8.5} = 0.0002$, effectively never.

Temperature has the units of cost. Small uphill moves are cheap at moderate temperature and large ones only at high temperature, so the search crosses low barriers easily and high ones rarely. At fixed T the accepted moves form a Markov chain whose stationary distribution is proportional to $e^{-\text{cost}/T}$: low-cost states are exponentially more likely, and as T falls the distribution concentrates on the global minima.

## Choosing a temperature

Inverting the rule gives the temperature at which a typical uphill move of size $\Delta$ is accepted with probability p:
$$T = -\frac{\Delta}{\ln p}.$$
For $\Delta = 7$ area units and $p = 0.4$: $T = -7/\ln 0.4 = 7.64$. For half of such moves to be accepted, $T = 7/\ln 2 = 10.1$.

This calibration matters because cost scales vary wildly between designs. A placer whose HPWL changes by 1 um per swap and one whose changes are 1000 um need temperatures three orders of magnitude apart. Engines therefore sample random moves before annealing starts, measure the typical cost change, and set the initial temperature from a target acceptance ratio. The same scale-awareness applies when the cost mixes terms such as wirelength and timing: their weights decide which moves are cheap and so which barriers the annealer can cross.

## Expected acceptances in a batch

Each proposal is accepted independently with its own probability, so the expected number of accepted moves in a batch is the sum of the probabilities, 1 for each non-increasing move and $e^{-\Delta/T}$ for each increasing one. At $T = 9$ with changes 11, 27, -2, 3 and 7: $0.295 + 0.050 + 1 + 0.717 + 0.459 = 2.52$ accepted moves out of 5.

The acceptance ratio is the annealer's main health signal. Near 1, the temperature is too high and the search is a random walk. Near 0, it is too cold and the search is frozen in a local minimum. Adaptive schedules, such as VPR's, watch the ratio and cool slowly while it is moderate, where most improvement happens. VPR also shrinks the move range limit as the acceptance ratio falls, so proposals stay local enough to be accepted at a useful rate.

## Going deeper

Kirkpatrick, Gelatt and Vecchi, "Optimization by simulated annealing" (Science, 1983); Sechen and Sangiovanni-Vincentelli, "The TimberWolf placement and routing package" (IEEE JSSC, 1985); Betz and Rose, "VPR: a new packing, placement and routing tool for FPGA research" (FPL, 1997).
