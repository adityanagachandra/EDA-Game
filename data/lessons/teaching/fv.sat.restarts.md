# Restarts and the Luby sequence

A CDCL solver restarts by clearing its trail back to decision level 0 while keeping everything it has learned. Restarts protect it against early bad decisions that would otherwise trap it in a hard part of the search. When to restart is a policy question, and the classic answers are the Luby sequence and geometric schedules. They are simple enough to compute by hand and important enough that every solver exposes them as a tunable parameter.

## The Luby sequence

The Luby sequence is 1, 1, 2, 1, 1, 2, 4, 1, 1, 2, 1, 1, 2, 4, 8, ... It is defined by $t_i = 2^{k-1}$ if $i = 2^k - 1$, and otherwise $t_i = t_{i - 2^{k-1} + 1}$ for $2^{k-1} \le i < 2^k - 1$. In words: the first $2^k - 1$ terms are two copies of the first $2^{k-1} - 1$ terms followed by $2^{k-1}$. So term 7 is 4, terms 8 to 14 repeat terms 1 to 7, and term 15 is 8. To find term 12, note $8 \le 12 < 15$, so $t_{12} = t_{12 - 8 + 1} = t_5 = 1$. Luby, Sinclair and Zuckerman showed that this schedule is within a logarithmic factor of the best possible restart strategy for any randomized algorithm whose run-time distribution is unknown, which is why it became a default. Each term is a power of two, and the sum of the first $2^k - 1$ terms is $k \cdot 2^{k-1}$, so the total work grows only slightly faster than the number of runs.

## Restart budgets

A solver multiplies the Luby sequence by a unit $u$: run $j$ lasts $u \cdot t_j$ conflicts, and then the solver restarts. The total number of conflicts by the end of run $r$ is $u \sum_{j=1}^{r} t_j$. With $u = 100$ and $r = 6$, the terms are 1, 1, 2, 1, 1, 2, summing to 8, so 800 conflicts have been performed. The sequence keeps returning to short runs while occasionally allowing long ones: by term 15 a run of 8 units occurs, and the longest run grows without bound. That combination is what a restart policy needs: frequent restarts to escape bad early decisions, plus ever longer runs so that hard instances that need long uninterrupted search can still be finished.

## Geometric schedules and what survives a restart

A geometric schedule starts with $c$ conflicts and multiplies by a factor $g > 1$: run $j$ lasts $\lfloor c \cdot g^{j-1} \rfloor$ conflicts. With $c = 100$ and $g = 1.5$, the first three runs last 100, 150 and 225 conflicts, 475 in total. Geometric schedules restart less and less often, the opposite of Luby's recurring short runs. Whatever the schedule, a restart clears only the trail: learned clauses stay, VSIDS activities stay, and saved phases stay, so the solver usually re-descends quickly into a similar region and then diverges where the new clauses and scores say it should. Modern solvers often replace fixed schedules by dynamic ones that restart when recent learned clauses look worse (higher LBD) than the long-run average.

## Going deeper

Luby, Sinclair and Zuckerman, "Optimal Speedup of Las Vegas Algorithms" (Information Processing Letters, 1993); Huang, "The Effect of Restarts on the Efficiency of Clause Learning" (IJCAI 2007); Audemard and Simon, "Refining Restarts Strategies for SAT and UNSAT" (CP 2012).
