# Thread pools and partitioned work

Parallel EDA engines rarely spawn a thread per task. They split the work into pieces (regions for routing, nets for extraction, partitions for placement) and feed them to a fixed pool of worker threads. How the pieces are ordered and assigned decides how close the run gets to the ideal: with uneven pieces, one long task can leave the other workers idle at the end. The makespan, the time until the last worker finishes, is the number to reason about.

## Greedy list scheduling

In list scheduling, tasks wait in a queue and each idle worker takes the next one. Equivalently, each task in order goes to the worker that becomes free first. The makespan is the largest worker load.

Example: tasks of 5, 3, 8, 2, 4 and 6 ms on 2 workers, ties going to the lower-numbered worker. Worker 1 takes 5 and worker 2 takes 3. Worker 2 is free first (at 3), so it takes 8 and is busy until 11. Worker 1 (at 5) takes 2, then at 7 takes 4, ending at 11. Both are at 11; the tie goes to worker 1, which takes 6 and finishes at 17. Makespan: 17 ms, while worker 2 sits idle from 11. Greedy scheduling is never worse than twice the optimum, but the order of the queue matters.

## Longest task first

Sorting the queue from longest to shortest before scheduling, the LPT rule, avoids leaving a long task for the end, when it cannot be balanced against anything. For the same tasks, the order is 8, 6, 5, 4, 3, 2. Worker 1 takes 8, worker 2 takes 6. Worker 2 is free first and takes 5 (ending at 11), worker 1 at 8 takes 4 (12), worker 2 at 11 takes 3 (14), worker 1 at 12 takes 2 (14). Makespan: 14 ms instead of 17.

LPT is guaranteed to be within a factor $4/3 - 1/(3p)$ of optimal. Engines apply the same idea by estimating task cost in advance: routing regions by pin count, nets by fanout. When costs cannot be predicted, smaller pieces and work stealing (an idle worker takes queued work from a busy one) keep the imbalance small at the end.

## Lower bounds on makespan

No schedule on $p$ workers can finish before
$$T^* \ge \max\left(\frac{\sum_i t_i}{p},\ \max_i t_i\right):$$
the total work must be shared among $p$ workers, and the longest task runs on one worker whatever happens. With dependencies between tasks, the longest dependency chain, the critical path, is a further bound.

In the example the total is 28 ms, so 2 workers need at least 14 ms, and the longest task is 8 ms. The bound is 14 ms, which LPT achieves, so LPT is optimal here. Comparing a schedule with this bound tells an engineer whether to fix the scheduler (a large gap) or to split the longest task (when $\max_i t_i$ dominates, as with one giant net or one congested region). A single 8 ms task on 8 workers with 16 ms total work gives a bound of 8 ms, not 2.

## Going deeper

Graham, "Bounds on Multiprocessing Timing Anomalies" (SIAM J. Applied Mathematics, 1969); the Taskflow documentation, used by OpenTimer, for task-graph scheduling with work stealing.
