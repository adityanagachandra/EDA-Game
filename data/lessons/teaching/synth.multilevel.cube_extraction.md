# Single-cube extraction

The simplest shared divisor is a single cube, such as $ab$ appearing inside many product terms. Extracting it as a new node trades several copies of its literals for one node and one literal per use. The value formula is short, and it shows exactly when sharing pays: often enough to use, never so often that extraction is free.

## Counting occurrences

A cube $c$ divides a cube $q$ when every literal of $c$ appears in $q$. Count how many cubes across all functions contain $c$; call it $m$. In $F_1 = aef + bce + def$ and $F_2 = abcf + abdf + aef + be$, the cube $df$ divides $def$ and $abdf$, so $m = 2$. The same count over all pairs of literals, kept in a matrix of cube-literal incidences, is how single-cube extraction finds candidates. Cubes of $m = 1$ are useless for sharing, since a node used once saves nothing; the interesting candidates divide two or more cubes, possibly in different functions.

## Literals after extracting a cube

Create $X = c$ with $k$ literals and replace $c$ by $X$ in each of the $m$ cubes it divides. Each such cube loses $k$ literals and gains one, and $X$ itself adds $k$. The new total is $L - mk + m + k$. For $F_1 = abd + cde + cdf$ and $F_2 = ae + af + bd + cd$, extracting $X = bd$ ($k = 2$, $m = 2$, $L = 17$) gives $17 - 4 + 2 + 2 = 17$: no gain. Extracting $X = cd$ instead ($k = 2$, $m = 3$) gives $17 - 6 + 3 + 2 = 16$. The new node also adds a level of logic to each use, which is the delay cost of sharing.

## The value formula

The saving is $mk - m - k = (m - 1)(k - 1) - 1$. It is positive only when $m \ge 2$, $k \ge 2$ and not both equal 2: a two-literal cube used twice saves nothing, used three times saves 1, and a three-literal cube used three times saves 3. For $F_1 = abcf + cd + ce$ and $F_2 = abde + bcef + abf$, a best candidate is $ab$, used 3 times, saving $3 \cdot 2 - 3 - 2 = 1$ (the cube $abf$, used twice, ties with $2 \cdot 3 - 2 - 3 = 1$). Greedy extraction picks the best value, updates the counts and repeats. The formula assumes literal count is area; with a real library a node also costs a gate and an extra level, so tools require a margin.

## Going deeper

Rajski and Vasudevamurthy, IEEE TCAD (1992), on single-cube and double-cube extraction; Rudell, "Logic Synthesis for VLSI Design", PhD thesis, UC Berkeley (1989), on rectangle covering.
