# Sandbox resource budgets

Every sandbox carries CPU, memory and time limits. Those limits are a budget that must be sized: too tight and correct runs are killed, too loose and a runaway skill starves everything else. Sizing them is concrete arithmetic on host capacity and on the measured runtimes of the golden suite, and it decides how many candidate changes the harness can evaluate at once.

## How many sandboxes fit on a host

A host with $C$ cores and $M$ GB of memory, running sandboxes that each reserve $c$ cores and $m$ GB, fits
$$\min\left(\left\lfloor C / c \right\rfloor, \left\lfloor M / m \right\rfloor\right)$$
sandboxes at once. A 64-core host with 512 GB running detailed-routing sandboxes of 4 cores and 48 GB fits $\min(16, 10) = 10$: memory binds, and 24 cores sit idle. Reshaping the sandbox to 6 cores and 48 GB still fits 10 and makes each run faster. Matching the sandbox shape to the host is the cheapest throughput gain a harness has.

## Setting a timeout from golden runtimes

A timeout should kill hung runs without killing slow correct ones. If golden-suite runtimes for a skill are roughly normal with mean $\mu$ and standard deviation $\sigma$, a timeout of $\mu + k\sigma$ kills a correct run with probability $1 - \Phi(k)$. With $\mu = 180$ s, $\sigma = 25$ s and $k = 3$, the timeout is 255 s and the false-kill rate is about 0.13%. Runtimes are often skewed to the right, so in practice the timeout is set from a high percentile of measured runtimes, with a margin, and re-measured when the skill changes.

## Wasted compute from hangs

A hung run burns its whole timeout. If a harness evaluates $N$ candidates per hour, a fraction $h$ hang, each sandbox holds $c$ cores and the timeout is $T$ seconds, the wasted compute is $N h T c / 3600$ core-hours per hour. With 120 candidates, 5% hanging, a 600 s timeout and 4 cores, that is $120 \times 0.05 \times 600 \times 4 / 3600 = 4$ core-hours every hour. Halving the timeout halves the waste but raises false kills. Hang detection that watches progress, such as a router's iteration counter not advancing, can end hung runs far earlier than the timeout.

## Going deeper

The Linux kernel documentation on cgroup v2 CPU and memory controllers.
