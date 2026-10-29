# License management and tool economics

Commercial EDA tools are licensed, usually as a pool of floating licenses served by a license server: a job checks out a license when it starts a licensed feature and checks it back in when done. Licenses, not machines, are often the scarce resource, since a site may have thousands of cores but a few dozen licenses of a signoff tool. Planning around licenses is ordinary scheduling arithmetic, and tool and flow builders influence it directly by deciding when a license is held.

## Licenses cap concurrency

If every job needs one license for its whole run and machines are plentiful, at most $L$ jobs run at once. $J$ jobs of $t$ hours then run in $\lceil J/L \rceil$ waves, and the batch takes $\lceil J/L \rceil\, t$ hours.

Example: 23 signoff timing runs of 3 h each with 5 licenses run in $\lceil 23/5 \rceil = 5$ waves: 15 h. Adding machines does nothing; adding a sixth license gives $\lceil 23/6 \rceil = 4$ waves, 12 h, and an eighth gives 3 waves, 9 h. This is why license counts are negotiated against schedules: a tape-out crunch that needs many signoff runs in parallel is a license question first. Open-source tools remove this constraint entirely, one of their practical attractions for regressions and sweeps, where thousands of runs would otherwise queue for licenses.

## License utilization

Over a batch, the licenses offer $L \times T_{wall}$ license-hours, and the jobs use $J t$. Utilization is
$$U = \frac{J t}{L\, \lceil J/L \rceil\, t} = \frac{J}{L \lceil J/L \rceil}.$$
In the example, $U = 23/(5 \times 5) = 92\%$: the last wave has only 3 jobs, so 2 licenses sit idle for 3 h. With 20 jobs, utilization would be 100%; with 21, a fifth wave holds a single job and $U = 21/25 = 84\%$.

Licenses are expensive, so sites measure utilization from license-server logs and fill idle capacity with lower-priority jobs. Queueing systems that are aware of licenses hold a job until its license is free, instead of starting it and letting it wait while it occupies a machine.

## Check in licenses early

A job holds a license for as long as it keeps it checked out, which is not necessarily its whole run. Suppose each job has a licensed phase of $b$ hours (the timing analysis itself) and then $c$ hours of unlicensed work (parsing reports, plotting, archiving). Holding the license until the end makes each job occupy it for $b + c$ hours: the batch takes $\lceil J/L \rceil (b + c)$. Checking in after the licensed phase lets the next job start: licensed phases run back to back in $\lceil J/L \rceil$ waves of $b$, and only the last wave's post-processing adds $c$, for $\lceil J/L \rceil\, b + c$.

Example: 12 jobs, 4 licenses, $b = 2$ h, $c = 1$ h. Holding: $3 \times 3 = 9$ h. Checking in: $3 \times 2 + 1 = 7$ h, saving 2 h with the same licenses. Flow builders apply the same idea by splitting jobs so that licensed steps do only licensed work.

## Going deeper

Vendor documentation for floating license servers (for example the FlexNet Publisher administration guide); Kahng et al., "VLSI Physical Design" (2nd ed., 2022), chapter 1, on design cost and tools.
