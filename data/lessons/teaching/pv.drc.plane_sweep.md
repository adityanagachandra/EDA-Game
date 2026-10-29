# Sweep-line spacing checks

Batch DRC engines process a layer with a sweep: a vertical line moves across the die, and only shapes near the line are compared. The sweep turns a quadratic all-pairs problem into one whose cost grows with the number of shapes and the number of near pairs. Its data structures, an event queue in x and an interval structure in y, are the heart of many checkers and Boolean engines.

## The active set

In a spacing sweep for rule $S$, a shape enters the active set at its left edge $x_1$ and leaves at $x_2 + S$, so it is still present when any shape starting within the rule distance arrives. The active set at a sweep position is the set of intervals $[x_1, x_2 + S]$ containing it. Its maximum size bounds memory and the work per event. Shapes $[46, 78]$ and $[28, 68]$ with $S = 27$ are active on $[46, 105]$ and $[28, 95]$, so both are active at $x = 60$. The maximum is found by scanning event points in order, adding at left events and removing at right events.

## Querying the active set

When a shape arrives, the engine asks which active shapes come closer than $S$ in y, that is, which y-intervals $[a, b]$ satisfy $a < y_2 + S$ and $b > y_1 - S$. An interval tree or a balanced tree of y-intervals answers this in $O(\log n + k)$ time for $k$ results. For a new shape on $[133, 154]$ with $S = 14$, the active interval $[105, 119]$ fails ($119 < 119$ is false) while $[94, 131]$ is returned. Only returned shapes get the exact distance test. The query window is the shape's y-extent grown by the rule distance.

## What the sweep saves

Brute force compares $n(n-1)/2$ pairs. The sweep compares only pairs near in both x and y, which in real layouts is roughly linear in $n$ because each shape has a bounded number of neighbours. For six scattered shapes, all 15 pairs may be far apart, and the sweep then makes no exact tests at all. The total cost is $O(n \log n + k)$ for $n$ shapes and $k$ near pairs. The same sweep structure computes Boolean operations and union areas, so one well-engineered sweep kernel serves many checks.

## Going deeper

de Berg, Cheong, van Kreveld and Overmars, "Computational Geometry: Algorithms and Applications" (3rd ed.), chapters on line-segment intersection and interval trees.
