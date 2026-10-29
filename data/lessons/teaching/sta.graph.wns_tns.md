# Worst and total negative slack

A chip has millions of endpoints, each with a slack. Two summary numbers reduce them to something a team can track from run to run: worst negative slack and total negative slack. They answer different questions, and timing-driven tools optimize both. Misreading which one moved, or forgetting how positive slacks enter, leads to wrong conclusions about whether a change helped.

## Worst negative slack

WNS is the most negative endpoint slack in the design or path group: the single worst violation. If the endpoint slacks are -12, -3, 4 and 9 ps, WNS is -12 ps. It sets the minimum clock period the design reaches today: the period would have to grow by 12 ps for every endpoint to pass. Many reports clip it at zero, so a passing design shows a WNS of 0 rather than its smallest positive slack. WNS is fragile as an objective: an optimizer that fixes only the worst path can improve WNS by a few ps while hundreds of nearly-as-bad paths stay untouched, and the next run simply reports a different endpoint as the worst.

## Total negative slack

TNS sums the negative endpoint slacks; positive slacks contribute zero. For -12, -3, 4 and 9 ps, TNS is -15 ps. It measures how widespread the violations are. Two designs can share a WNS of -8 ps with TNS values of -9 ps and -4000 ps: the first has a local problem, the second a systematic one such as a bad floorplan or an over-tight constraint. TNS is a smoother objective than WNS because every failing endpoint contributes, so optimizers that improve many paths a little see progress. The number of failing endpoints is usually reported alongside it.

## How a fix changes TNS

When an ECO improves some endpoints, recompute TNS from the new slacks rather than adding the improvement directly. An endpoint that crosses zero contributes nothing, so it saves only what it was failing by. Suppose endpoints at -53, -7 and -27 ps each gain 15 ps. They become -38, $+8$ and -12 ps. The total improvement in TNS is $15 + 7 + 15 = 37$ ps, not $3 \times 15 = 45$ ps, because the endpoint at -7 ps could only give back 7. The same clipping explains why fixing a few large violations can move TNS more than many small fixes, and why a TNS target must always be read together with WNS and the failing-endpoint count.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8 on timing closure; the OpenSTA commands report_wns and report_tns.
