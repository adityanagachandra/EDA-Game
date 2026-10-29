# Observability don't cares

A node's value matters only if a change in it can reach an output. Wherever it cannot, the node is unobservable, and those input conditions are its observability don't cares. They are computed with the Boolean difference, they accumulate along paths, and they are the main source of freedom for simplifying nodes deep inside a network.

## The Boolean difference

For an output $z$ that depends on node $y$, the Boolean difference $\partial z/\partial y = z|_{y=1} \oplus z|_{y=0}$ is 1 exactly where flipping $y$ flips $z$. The **observability don't care** of $y$ with respect to $z$ is its complement. Take $y = ab$ and $z = y d\bar{c}$: $z|_{y=1} = d\bar{c}$ and $z|_{y=0} = 0$, so $\partial z/\partial y = d\bar{c}$. $y$ is observable only when $d = 1$ and $c = 0$, one of the four $(c, d)$ values. Over all 16 combinations of $a, b, c, d$ the ODC has $3 \cdot 4 = 12$ minterms, because $a$ and $b$ do not affect observability here.

## Care sets and counting

The **observability care set** is where the Boolean difference is 1. For $z = y\bar{c} + d$: $z|_{y=1} = \bar{c} + d$ and $z|_{y=0} = d$, so $\partial z/\partial y = \bar{c}\bar{d}$, true for one of the four $(c, d)$ values. With $a$ and $b$ free that is 4 of the 16 input combinations in the care set and 12 in the ODC. A simplification of $y$ may change its value anywhere in the ODC. Combined with controllability don't cares, the full local freedom of a node is everything outside the set of inputs that are both reachable and observable.

## Several outputs and reconvergence

If $y$ fans out to several outputs, it is unobservable only where it is unobservable at every one: the ODC is the intersection of the per-output ODCs. With $z_1 = y(c \oplus d)$ and $z_2 = y\bar{c} \oplus d$, $\partial z_1/\partial y = c \oplus d$ and $\partial z_2/\partial y = \bar{c}$; both are 0 only when $c = 1$ and $d = 1$, so 4 of 16 combinations are ODC. Along a single fanout-free path, ODCs accumulate gate by gate. With reconvergent fanout, simple per-branch rules fail, because flipping $y$ can change two paths that recombine, and exact ODCs need the Boolean difference of the whole output function, usually computed on a window.

## Going deeper

Savoj and Brayton, "The use of observability and external don't cares for the simplification of multi-level networks", DAC (1990); Damiani and De Micheli, "Observability don't care sets and Boolean relations", ICCAD (1990).
