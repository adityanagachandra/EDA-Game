# E-classes and congruence closure

The machinery that keeps an e-graph consistent is congruence closure: when two classes are declared equal, every pair of e-nodes that now apply the same operator to the same classes must be merged as well, and so on until a fixpoint. Counting e-classes before and after unions is the clearest way to see hashconsing, union-find and congruence at work.

## Hashconsing builds the initial classes

Insert terms bottom up. Each leaf is an e-node in its own class. For a composite term, canonicalize its children to their class ids and look up the e-node (operator, children) in a hash table: if it exists, reuse its class; otherwise create a new e-node and class. So writing $t_1 = a + b$ twice creates one class, not two. For the leaves $a$, $b$, $c$ and terms $t_1 = b + a$, $t_2 = b \cdot b$, $t_3 = c \cdot b$, $t_4 = t_1 + t_1$, there are 7 distinct e-nodes and hence 7 e-classes. Note that $a + b$ and $b + a$ are different e-nodes until a commutativity rule merges them.

## Unions and congruence

Asserting $x \equiv y$ merges their classes in the union-find structure. Then congruence must be restored: any two e-nodes $f(c_1, \dots)$ and $f(c_1', \dots)$ whose children are now pairwise equal must merge. With $t_1 = a \cdot c$, $t_2 = b \cdot c$, $t_3 = a \cdot b$, $t_4 = t_3 \cdot c$, asserting $a \equiv b$ leaves 6 classes from the union alone, and then $a \cdot c$ and $b \cdot c$ become congruent and merge, giving 5. Each merge can trigger more, up the parent chains, so closure iterates to a fixpoint. The procedure is the classic congruence closure of Nelson and Oppen and of Downey, Sethi and Tarjan.

## Rebuilding efficiently

Restoring congruence after every single merge is expensive when saturation performs thousands of merges per iteration. The egg library defers it: merges during an iteration go into a worklist, and a single **rebuild** pass at the end re-canonicalizes the affected parents and repairs the hash table. The result is identical, but the amortized cost is much lower. Merges caused by congruence alone, beyond the asserted unions, measure how much a single equality propagates through shared structure; when the asserted classes have no parents in common, congruence adds nothing. Correctness depends on the invariant holding before every search phase.

## Going deeper

Nelson and Oppen, "Fast decision procedures based on congruence closure", Journal of the ACM (1980); Downey, Sethi and Tarjan, "Variations on the common subexpression problem", Journal of the ACM (1980); Willsey et al., POPL (2021).
