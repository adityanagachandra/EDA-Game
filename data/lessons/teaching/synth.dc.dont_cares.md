# Don't-care conditions

A don't care is a point where the circuit's value does not matter, and every don't care is freedom to simplify. Some come from the specification (external don't cares); others arise inside a network because some value combinations can never occur or can never be observed. Understanding where don't cares come from and when using them is safe is central to node simplification and resubstitution.

## External don't cares

External don't cares are given by the environment. A BCD decoder never sees input codes 10 to 15, so its outputs on those codes are free. A five-state FSM encoded in three bits has three unused codes, which are don't cares for the next-state and output logic if the machine can never reach them. An output that is ignored under some condition (a valid bit is low) is a don't care on those inputs. In a K-map these cells may join any group or none. The risk: a don't care is a promise. If the impossible input does occur, after a soft error in an FSM for instance, the circuit outputs whatever the minimizer chose.

## Don't cares inside a network

Inside a Boolean network two more kinds appear. **Satisfiability don't cares**: for a node $y = f(x)$, combinations where $y \ne f(x)$ can never happen, so the set $y \oplus f(x)$ is a don't care for any node that sees both $y$ and $x$. If $y = ab$ feeds $z$ alongside $a$, then $y = 1$ with $a = 0$ never occurs. **Observability don't cares**: an input of a node is unobservable when the outputs do not depend on it; for $z = y + c$, the value of $y$ does not matter whenever $c = 1$. Formally the ODC of $y$ with respect to $z$ is the complement of the Boolean difference $z|_{y=1} \oplus z|_{y=0}$.

## Using don't cares safely

A node may be replaced by any function that agrees with it wherever the node is both controllable (its local input pattern can occur) and observable (its value can reach an output). Inside that care set it must match; outside it is free. Exact global don't cares can be huge, so tools compute them over a window of nearby logic, by BDDs on small windows or by simulation and SAT. Two cautions: don't cares of different nodes interact, so using full ODCs at several nodes at once can be wrong; and removing logic covered by don't cares also removes redundancy that would make stuck-at faults untestable, which is why don't-care optimization improves testability.

## Going deeper

Bartlett et al., "Multilevel logic minimization using implicit don't cares", IEEE TCAD (1988); Savoj and Brayton, "The use of observability and external don't cares for the simplification of multi-level networks", DAC (1990).
