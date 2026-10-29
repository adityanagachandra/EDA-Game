# Dirty flags and lazy invalidation

A dirty flag is one bit per cached result that says "my inputs changed since I was computed." Setting flags is cheap, recomputing is not, and the gap between the two is where incremental tools win. An eager tool recomputes after every edit; a lazy one marks results dirty and recomputes only when a value is read. Combined with early cutoff, which stops propagation when a recomputed value comes out unchanged, dirty flags let an optimizer make thousands of edits between expensive queries.

## Eager recomputation after every edit

In the eager scheme, each edit immediately recomputes every result downstream of the edited node. The cost of a sequence of edits is the sum of the downstream set sizes, edit by edit, and results shared by several edits are paid for every time.

Example: edges $A \to C$, $B \to C$, $C \to D$, $D \to E$. A script edits $A$ and then $B$. Editing $A$ recomputes $C$, $D$ and $E$: three evaluations. Editing $B$ recomputes $C$, $D$ and $E$ again: three more. Total six. If the script makes $k$ edits that all feed the same cone of size $m$, eager evaluation costs $km$. Eager evaluation is simple and always leaves values valid, which is why it survives in small tools, but it scales badly when an optimizer makes many edits per query.

## Mark dirty, recompute on demand

In the lazy scheme, an edit walks forward and sets the dirty flag of every downstream result, stopping at results already dirty, since their descendants were flagged earlier. Nothing is recomputed until a read. A read of a dirty result first refreshes its dirty inputs, recursively, then clears the flag. If every result is eventually read once, each dirty result is recomputed exactly once.

For the same example, edits to $A$ and $B$ flag $C$, $D$ and $E$, the union of the two downstream sets. One query pass recomputes three results instead of six. In general lazy cost is the size of the union, eager cost the sum of the parts, and the saving grows with the overlap between edits. The price is that every read must check a flag, and code that reads cached fields directly, bypassing the check, silently returns stale data.

## Early cutoff on unchanged values

Dirty means "might have changed," not "has changed." When a result is recomputed and its new value equals the old one, its dependents see identical inputs and need no work. Early cutoff exploits this: re-evaluate a node only if at least one of its inputs actually changed value.

Example: edges $A \to B$, $B \to C$, $C \to D$, $A \to E$. An edit changes $A$, so $B$ and $E$ are re-evaluated. Suppose $B$ comes out unchanged, which happens in timing when the arrival at a pin is set by another, slower fanin. Then $C$ and $D$ are not re-evaluated, so only $B$ and $E$ are: two evaluations instead of four. A node reachable by another path that avoids $B$ is still re-evaluated. Cutoff needs a cheap equality test, which is one reason timers compare arrivals with a small tolerance rather than recomputing bit-exact values.

## Going deeper

Mokhov, Mitchell and Peyton Jones, "Build Systems à la Carte" (ICFP 2018), for dirty bits and early cutoff in build systems; the OpenSTA source for how it marks and refreshes invalid arrivals.
