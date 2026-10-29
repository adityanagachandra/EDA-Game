# Transition density propagation

Signal probability says how often a net is 1, not how often it changes. Transition density, the expected number of transitions per unit time, is the quantity power actually depends on. Najm's transition-density method propagates it through logic using Boolean differences, without assuming one transition per cycle. It handles asynchronous activity and multiple transitions per cycle, and it is a standard model inside vectorless power estimators.

## The propagation rule

For a gate $y = f(x_1, \ldots, x_n)$ with independent inputs, an input transition propagates to the output exactly when the Boolean difference $\partial y/\partial x_i = f|_{x_i=1} \oplus f|_{x_i=0}$ is 1. Najm's rule is
$$D(y) = \sum_i P\!\left(\frac{\partial y}{\partial x_i}\right) D(x_i).$$
For AND, $\partial y/\partial a = b$, so $D(y) = p_b D_a + p_a D_b$. For OR, $\partial y/\partial a = \bar b$, so $D(y) = (1 - p_b)D_a + (1 - p_a)D_b$; with $p_a = 0.2$, $p_b = 0.35$, $D_a = 0.4$ and $D_b = 1$ per ns, $D(y) = 0.26 + 0.8 = 1.06$ per ns. For XOR every input transition propagates, so densities add.

## Wider gates

The rule extends to any number of inputs. For a 3-input AND, $\partial y/\partial a = bc$, so $D(y) = p_b p_c D_a + p_a p_c D_b + p_a p_b D_c$; with $p = (0.25, 0.25, 0.75)$ and $D = (0.8, 1.6, 1.1)$ per ns this is $0.15 + 0.3 + 0.069 = 0.519$ per ns. For a 3-input OR, an input matters only when the other two are 0. The pattern shows why wide AND and OR gates filter activity: a transition needs every other input at its non-controlling value, which becomes unlikely as the gate widens.

## Through a network

Propagate in topological order, carrying both $p$ and $D$ for every node. For $n = a \wedge b$ and $y = n \vee c$, with $p = (0.35, 0.7, 0.45)$ and $D = (0.1, 1, 1)$ per ns: $p_n = 0.245$ and $D(n) = 0.7 \cdot 0.1 + 0.35 \cdot 1 = 0.42$, then $D(y) = 0.55 \cdot 0.42 + 0.755 \cdot 1 = 0.986$ per ns. Power is then $\tfrac{1}{2} C V^2 D$ per node, since each transition costs half of $CV^2$ on average. The method assumes independent inputs and treats transitions as instantaneous, so it ignores inertial filtering of short glitches and correlation from reconvergence.

## Going deeper

Najm, "Transition density: a new measure of activity in digital circuits", IEEE Transactions on Computer-Aided Design (1993); Najm, "A survey of power estimation techniques in VLSI circuits", IEEE Transactions on VLSI Systems (1994).
