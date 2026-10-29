# Genetic algorithms and evolutionary search

Genetic algorithms evolve a population of candidate solutions by selection, crossover and mutation. They have been applied to partitioning, floorplanning, placement and analog sizing, and evolutionary strategies such as CMA-ES and NSGA-II remain standard for black-box and multi-objective tuning. For EDA the key design questions are the encoding, whether crossover of two good parents tends to produce a good child, and how strongly selection favours the fittest.

## Fitness-proportional selection

Roulette-wheel selection picks individual i with probability proportional to its fitness, $p_i = f_i / \sum_j f_j$. With fitness values 6, 37, 25, 26 and 9 (total 103), individual 4 is chosen with probability $26/103 = 0.252$ per draw, and the fittest, individual 2, with $37/103 = 0.359$.

Fitness must be positive and should be scaled sensibly. If the cost is wirelength, fitness might be a baseline minus wirelength, or its reciprocal. Raw proportional selection has two failure modes: one super-fit individual takes over the population early (premature convergence), or late in the run all fitnesses are similar and selection becomes nearly random. Rank-based selection and fitness scaling fix both by making the selection pressure independent of the raw cost scale.

## Crossover on an encoding

One-point crossover cuts two parent strings at the same position and joins the head of one to the tail of the other. Encode a bipartition as a bit per cell. Parent A = 100011 and parent B = 000010 with a cut after gene 4 give the child 1000 from A and 10 from B, that is 100010. With cell areas 3, 6, 6, 9, 6, 8 the child puts cells 1 and 5 in partition 1, total area 9.

Crossover only helps if good partial solutions, building blocks, survive the cut. For a bipartition the bit order should follow netlist locality, so connected cells are near each other in the string; otherwise crossover scrambles good clusters. Permutation encodings, used for placement orders or Polish expressions, need special operators (order or cycle crossover) because cutting and joining two permutations usually duplicates some elements. Mutation, flipping a few bits, maintains diversity.

## Tournament selection pressure

Tournament selection draws k individuals at random and keeps the fittest. Drawing with replacement from a population of N, the single best individual wins whenever it is drawn at least once:
$$P = 1 - \left(1 - \tfrac{1}{N}\right)^k.$$
With $N = 93$ and $k = 5$: $1 - (92/93)^5 = 0.0526$, about $k/N$ for small $k/N$. Raising k raises selection pressure: the best individual gets more offspring, so the population converges faster but loses diversity.

Tournament selection depends only on fitness order, not scale, and is easy to parallelize, which is why it is popular in practice. For multi-objective EDA problems, NSGA-II ranks individuals by Pareto dominance level and crowding distance and then uses tournaments on that ranking. Compared with annealing, a population costs more evaluations per generation but explores several regions at once, which pays off when evaluations can run in parallel.

## Going deeper

Goldberg, "Genetic Algorithms in Search, Optimization and Machine Learning" (1989); Deb, Pratap, Agarwal and Meyarivan, "A fast and elitist multiobjective genetic algorithm: NSGA-II" (IEEE Transactions on Evolutionary Computation, 2002).
