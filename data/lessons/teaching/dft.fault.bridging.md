# Bridging faults

Shorts between neighbouring wires are among the most common real defects, and they do not behave like a stuck line: the shorted nets fight, and the result depends on their values and on the drivers' strengths. Bridging fault models describe that fight in logic terms so that ATPG and fault simulation can target layout-extracted shorts. A tool designer must handle their quadratic candidate count and their excitation condition, which differs from stuck-at activation.

## Wired and dominant models

A bridge joins two nets $x$ and $y$. In the **wired-AND** model both nets take $x \wedge y$, as when a strong pull-down wins; in **wired-OR** both take $x \vee y$. In a **dominant** bridge, one driver is much stronger, so the victim net copies the aggressor: "$x$ dominates $y$" means $y$ takes $x$'s value while $x$ is unaffected. Each model is a different fault, so a pair of nets can carry up to four. A bridge is a **non-feedback** bridge if neither net lies in the other's fan-in; a feedback bridge can create a loop and oscillation, and is modeled separately.

## Excitation and candidate counts

A bridge has an effect only when the two nets carry **opposite** values in the fault-free circuit; if they agree, shorting them changes nothing. So the excitation condition is $x \neq y$, and counting exciting vectors means enumerating where the two functions differ. Detection additionally requires the corrupted value to reach an output. For a wired-AND, only the net that was 1 is corrupted (it becomes 0), so propagation must start from that net. Without layout, every pair of $N$ nets is a candidate: $\binom{N}{2}$ pairs, so 1000 nets give 499,500 pairs per model. Production flows extract only physically adjacent pairs, with their critical area, from the layout, which cuts the list to roughly linear size and makes bridging ATPG practical.

## Detection in practice

To check a bridge by simulation, compute fault-free values, replace both bridged nets with the bridged value (AND, OR, or the dominant net's value), resimulate their fan-out, and compare outputs. For a dominant bridge where $g_3$ dominates $g_2$, only $g_2$ changes, and only on vectors where $g_3 \neq g_2$. Stuck-at test sets detect many bridges by accident, because a stuck-at test for $y$ stuck-at-0 that also sets $x = 0$ detects the wired-AND bridge. Real shorts are often resistive and behave in between the logical models, which is why IDDQ testing (measuring quiescent supply current, which a short raises) was historically used alongside logic tests.

## Going deeper

Bushnell and Agrawal (2000), chapter 4; Mei, "Bridging and stuck-at faults," IEEE Transactions on Computers, 1974; Wang, Wu and Wen (2006), chapter 2.
