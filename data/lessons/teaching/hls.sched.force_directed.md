# Force-directed scheduling

Force-directed scheduling, introduced by Paulin and Knight, attacks latency-constrained scheduling by balancing the expected demand on each unit type across control steps. It treats each unscheduled operation as spread uniformly over its time frame, sums those probabilities into distribution graphs, and repeatedly fixes the operation whose placement best flattens them. It is the classic example of a scheduler that reasons about all operations at once rather than greedily.

## Time frames as probabilities

Fix a latency bound $L$ and compute each operation's frame $[\mathrm{ASAP}_i, \mathrm{ALAP}_i]$, which has width $w_i = \mu_i + 1$. Before any decision, assume operation $i$ is equally likely to start in any step of its frame: $p_i(s) = 1 / w_i$ for $s$ inside the frame and 0 outside. An operation with frame $[4, 6]$ has probability $1/3 \approx 0.33$ in each of steps 4, 5 and 6, and 0 in step 3. A critical operation has width 1, so its probability is 1 in its only step: it is already decided. The model ignores dependences between the probabilities of connected operations, which is a simplification the algorithm corrects by recomputing frames after every decision.

## Distribution graphs

For each unit type $k$ and step $s$, the distribution graph sums the probabilities of the type-$k$ operations: $DG_k(s) = \sum_{i \in k} p_i(s)$. It estimates how many type-$k$ units step $s$ will need. With multiplications m1 in $[1, 3]$, m2 in $[4, 4]$ and m3 in $[1, 4]$ under $L = 4$, step 2 receives $1/3$ from m1 and $1/4$ from m3, so $DG(2) = 0.58$; step 4 receives 1 from m2 and $1/4$ from m3, so $DG(4) = 1.25$. A peak in the distribution graph signals a step where multiplications will collide and force an extra multiplier. The goal of the algorithm is a schedule whose distribution graphs are as flat as possible, since the peak of the final, fully decided graph is the number of units needed.

## Self-force and the iteration

Fixing operation $i$ in step $j$ changes its probability from $1/w_i$ everywhere in its frame to 1 at $j$ and 0 elsewhere. The self-force of that choice is $F_i(j) = \sum_s DG_k(s)\,\Delta p_i(s) = DG_k(j) - \frac{1}{w_i} \sum_{s \in \text{frame}} DG_k(s)$: the distribution value at $j$ minus the average over the frame. A negative force means $j$ is a quieter step than average, so placing the operation there flattens the graph. The full algorithm adds predecessor and successor forces, because fixing $i$ shrinks the frames of its neighbours, then fixes the operation and step with the lowest total force, recomputes frames and distribution graphs, and repeats until everything is placed. Each iteration costs a full recomputation, which makes the method roughly cubic in the number of operations.

## Going deeper

Paulin and Knight, "Force-directed scheduling for the behavioral synthesis of ASICs" (IEEE TCAD, 1989); De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5.
