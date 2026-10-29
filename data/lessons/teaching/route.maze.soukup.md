# Soukup's fast maze router

Soukup's router (1978) takes a different trade: it gives up the guarantee of a shortest route to gain speed. It walks greedily toward the target and falls back to breadth-first expansion only when it is blocked. On typical layouts with long open stretches it touches far fewer cells than Lee, and it still finds a route whenever one exists.

## Greedy first, breadth-first when blocked

From the current cell, Soukup's router moves in the direction of the target, depth-first, without branching. When the straight move is blocked, it switches to breadth-first expansion from the cells reached so far until it finds a cell from which it can again move toward the target, and then resumes the greedy walk. The search alternates between these two modes until the target is reached.

On an open grid between pins at Manhattan distance $d$, the greedy phase alone reaches the target after about $d$ cells. Lee would label on the order of $d^2$ cells in the same situation. The breadth-first fallback is what keeps the method complete: if a route exists, the expansion eventually finds a way around every wall.

## What is guaranteed and what is not

Soukup's router always finds a route if one exists, but not necessarily a shortest one. The greedy walk can commit to the wrong side of an obstacle and then detour around it, producing a route longer than Lee's or Hadlock's. Hadlock, by contrast, also speeds up Lee but preserves optimality by ordering search by detours.

A ranking of search effort on a lightly blocked grid is usually Soukup fewest, then Hadlock, then Lee. In a dense maze of obstacles the greedy moves keep failing and Soukup's advantage shrinks toward Lee's behaviour. For non-critical nets a somewhat longer route costs little, and the runtime saved across millions of nets is the point.

## Where the idea went

The core idea, biasing the search toward the target so open space is crossed with little expansion, survives in every modern detailed router. What replaced Soukup's specific procedure is A*: a priority search with an admissible, consistent heuristic gets most of the same speedup while keeping optimality, and it handles general edge costs (vias, bends, congestion) that a pure greedy walk cannot weigh.

So a reviewer should object if a router uses Soukup-style search for timing-critical nets, where extra length costs delay. For bulk nets in open regions, a fast incomplete-optimality search is a legitimate engineering choice, as long as the router can rip up and reroute what it gets wrong.

## Going deeper

J. Soukup, "Fast Maze Router" (1978); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.).
