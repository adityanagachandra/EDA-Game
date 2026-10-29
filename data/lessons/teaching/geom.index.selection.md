# Choosing a spatial index

Interval trees, segment trees, R-trees, k-d trees, quadtrees, bin grids and bounding volume hierarchies all answer "what is near here", but each is tuned to a different query, data shape and update pattern. Choosing well is a design decision that sets the cost of every check and every move in the engine. This concept is about matching structure to workload, and about why hierarchy and spatial locality make full-chip checking tractable at all.

## Matching the query to the structure

Name the query first. A stabbing query on intervals (which wires cross a cut line) is an interval tree. A window query over boxes (which shapes lie near this rectangle) is an R-tree or a bin grid, and a distance query becomes a window query once the window is grown by the distance. A nearest-neighbour query over points (the nearest legal site, the closest pin) is a k-d tree or, for a regular site array, plain arithmetic. An aggregate over coverage (covered length at a sweep stop) is a segment tree. A bounding volume hierarchy, a tree of nested boxes built top-down, serves the same queries as an R-tree when objects are known in advance. The wrong match still works but pays a large factor.

## Static versus dynamic data

Some workloads load data once and query it many times (batch DRC on a finished layout); others interleave updates and queries (an incremental router, a placer's legalizer, an interactive editor). Static workloads favour structures built in bulk: a k-d tree split at medians, or an R-tree bulk-loaded by sorting along a space-filling curve or by tiling, giving tight, barely overlapping nodes. Dynamic workloads favour structures with cheap local updates: bin grids and R-trees. R-trees stay height-balanced under updates, but their MBRs bloat and overlap over time, so query cost creeps up; R*-tree reinsertion or periodic rebuilding restores it. A k-d tree degrades under updates because its splits no longer match the data.

## Hierarchy and why DRC stays tractable

Two facts keep checking a billion shapes feasible. First, rules have a small interaction distance, so each shape needs comparing only with nearby shapes, found by a sweep or a spatial index; the work is near-linear instead of quadratic. Second, layouts are hierarchical: a cell stored once and placed thousands of times can be checked once, and only its interaction region, the band within the rule distance of other geometry, needs checking in each context. Hierarchical engines therefore index cell bounding boxes as well as shapes and flatten selectively where instances overlap heavily. When a design has little repetition, the bookkeeping of hierarchy costs more than it saves, and flat processing with good spatial indexing wins.

## Going deeper

Samet, "Foundations of Multidimensional and Metric Data Structures" (2006); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on physical verification; the KLayout documentation on deep (hierarchical) mode in DRC.
