# Boolean resubstitution

Resubstitution re-expresses a node using other nodes already in the network, so its own logic can be deleted. The Boolean version asks a precise question: is the node's function, on its care set, a function of a few chosen divisors? If so, the node can be rebuilt from them with a small gate. ABC's resub command does this with truth tables over windows, and it is one of the most effective AIG optimizations.

## Which divisor values occur

Simulate the candidate divisors $g_1, \dots, g_k$ on every input vector of a window (or on many random vectors) and record which value combinations of the divisors appear. For $g_1 = a + b$ and $g_2 = \overline{ad}$ over $a, b, c, d$: $g_1 = 0$ requires $a = 0$, which forces $g_2 = 1$, so the pair $(0, 0)$ never occurs. Three of the four pairs appear. Missing combinations are controllability don't cares for any new node built from these divisors, and they are exactly the freedom that makes a small replacement possible.

## The functional-dependency test

$f$ can be written as $f = h(g_{i_1}, \dots, g_{i_s})$ for some function $h$ if and only if no two input vectors agree on all the chosen divisors but disagree on $f$. Check it by grouping vectors by their divisor values: each group must have a single $f$ value, and that value defines $h$ on the group. To find the smallest divisor set, test subsets in order of size. For $f = \sum m(2, 3, 14, 15)$ with candidates $bc$, $c\bar{a}$, $\overline{d + a}$ and $a + c$, no single divisor separates the ON and OFF vectors, but a pair does, so the answer is 2, and $h$ is read off the groups.

## Using don't cares and accepting the move

Vectors in the node's don't-care set (unreachable or unobservable) need not be separated, so they are dropped before the test. That can let fewer divisors suffice: with four don't-care vectors removed, a function that needed two divisors may depend on just one. The move is accepted if the new node, typically one to three AND gates over the divisors, is smaller than the logic freed (the node's MFFC). ABC computes truth tables over windows with up to a dozen or so leaves, uses simulation to filter candidates, and verifies the dependency exactly on the window. Resubstitution is Boolean, so it finds rewrites algebraic substitution cannot.

## Going deeper

Mishchenko and Brayton, "Scalable logic synthesis using a simple circuit structure", IWLS (2006); the ABC documentation for the resub command.
