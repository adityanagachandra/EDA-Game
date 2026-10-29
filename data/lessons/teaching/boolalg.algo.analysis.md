# Algorithm analysis for Boolean tools

Boolean tools live on the edge between polynomial and exponential cost. A tool designer must predict, before writing code, whether an approach scales: how much memory an explicit representation takes, how many calls a recursion makes, and how a measured runtime will extrapolate to the next design size. These estimates come from three habits: count the size of the data, write the recurrence of the algorithm, and know what each growth rate does when the input grows.

## Memory of an explicit truth table

A truth table of $n$ variables has $2^n$ bits. Packed into $w$-bit words it needs $2^n / w$ words, $2^{n-6}$ for 64-bit words when $n \ge 6$. A 16-input table is 1,024 words (8 KiB), a 20-input table 16,384 words (128 KiB), and a 30-input table $2^{24}$ words, 128 MiB.

Each added variable doubles the memory, which is why truth tables are used for cuts of up to about 16 inputs and why everything larger is represented implicitly by AIGs, BDDs or CNF. When a design flow computes a truth table per node, multiply by the number of nodes too: a million 6-input cuts at one word each is 8 MB, but a million 12-input cuts is 512 MB.

## Recurrences of recursive Boolean algorithms

Recursive algorithms that split on a variable obey recurrences. A tautology check that splits on one variable per level without pruning makes $T(k) = 2T(k-1) + 1$ calls with $T(0) = 1$, which unrolls to $T(k) = 2^{k+1} - 1$: 31 calls for 4 variables, about two million for 20. A procedure that does $c$ units of work and recurses on both cofactors satisfies $T(k) = 2T(k-1) + c$, still $\Theta(2^k)$.

Pruning changes the base, not the shape: if one branch always has one fewer variable and the other two fewer, $T(k) = T(k-1) + T(k-2) + 1$ grows like $1.618^k$. Memoization changes the shape: a recursion that caches results per distinct subfunction does work proportional to the number of distinct subfunctions, which is the BDD size, not $2^k$.

## Extrapolating exponential running time

Measured runtimes extrapolate by growth rate. For $\Theta(2^n)$, adding $k$ variables multiplies time by $2^k$: an exact method taking 2 s at $n = 40$ takes 8 s at $n = 42$ and about 34 minutes at $n = 50$. For $\Theta(n^3)$, scaling $n$ by a factor $r$ multiplies time by $r^3$: doubling the size costs 8 times as long.

The practical test is what happens on the next technology node or the next design generation. Polynomial algorithms degrade gracefully; exponential ones hit a wall. That is why EDA engines that are exponential in the worst case (SAT, exact minimization) are wrapped in resource limits and fall back to heuristics, and why an engine's documentation should state its complexity in the size of the actual input, which for a truth table is already $2^n$.

## Going deeper

Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), chapters 3 and 4. Knuth, "The Art of Computer Programming", volume 4A, on Boolean function representation costs.
