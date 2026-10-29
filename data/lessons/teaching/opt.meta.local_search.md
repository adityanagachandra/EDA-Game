# Local search, neighborhoods and tabu search

Local search starts from a solution and repeatedly moves to a neighbour chosen by a rule. Its behaviour is defined entirely by the move set, which decides what counts as a neighbour and hence what counts as a local optimum. Fiduccia-Mattheyses partitioning, detailed-placement swaps, rip-up and reroute, and gate resizing passes are all local search. Tabu search adds memory so the walk can leave local optima without cycling.

## Neighbourhoods and local minima

A local minimum is a solution no neighbour improves on, so it depends on the move set. Take ten partitioner configurations in a line with costs 50, 47, 25, 24, 32, 30, 93, 72, 60, 53 and moves of plus or minus 1. Configuration 3 (cost 24) is lower than both neighbours (25 and 32); configuration 5 (30) beats 32 and 93; configuration 9 (53), an endpoint, beats its one neighbour 60. That is three local minima, of which only configuration 3 is global.

Enlarging the move set to plus or minus 2 removes the local minimum at 5, since configuration 3 is then a neighbour and is cheaper. Larger neighbourhoods have fewer local minima but cost more to scan per step. This trade-off is central to engine design: FM moves single cells, while multilevel partitioners effectively move whole clusters at coarse levels.

## Steepest descent and where it stops

Steepest descent moves to the best neighbour as long as that improves the cost, and stops at the first local minimum. From configuration 1 in costs 71, 20, 10, 63, 45, 15, 60, 44, 46, the neighbours are 0 (cost 71) and 2 (cost 10); it moves to 2. There the neighbours cost 20 and 63, both worse, so it stops at cost 10. Starting at 4 instead, it would move to 5 and stop at 15, a different local minimum.

Descent is fast and deterministic, and it is the inner loop of many engines: a greedy detailed placer accepts only improving swaps, and a sizer upsizes the gate with the best slack gain per area until nothing helps. Its weakness is that the start fixes the answer. Remedies are random restarts, perturbation followed by descent (iterated local search), allowing worse moves under control (simulated annealing), or forbidding recent moves (tabu search).

## Tabu search

Tabu search always moves to the best non-tabu neighbour, even when it is worse, and keeps a tabu list of recently visited configurations, or recently reversed moves, that may not be revisited for a fixed tenure. Costs 84, 50, 52, 59, 15, 37, 82, 26 with moves of plus or minus 2, start 5, and the current plus last 2 visited configurations tabu. From 5 the walk takes 4 (cost 15), then the best non-tabu neighbour of 4 is 2 (52), then 1 (50), then 3 (59), then 4 again (15), then 5 (37). The best cost seen is 15.

The forced uphill moves let it leave a local minimum, and the tabu list stops it stepping straight back. Aspiration criteria override the list when a tabu move would beat the best solution found. FM partitioning has the same flavour: each cell moves once per pass, is then locked, and the pass keeps the best prefix of moves, so it can climb out of local minima within a pass.

## Going deeper

Fiduccia and Mattheyses, "A linear-time heuristic for improving network partitions" (DAC, 1982); Glover, "Tabu search, part I" (ORSA Journal on Computing, 1989); Aarts and Lenstra (eds.), "Local Search in Combinatorial Optimization" (1997).
