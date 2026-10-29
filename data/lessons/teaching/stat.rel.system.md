# System reliability and fault tolerance

A chip is a system of many parts, each with its own lifetime, and its reliability follows from how they are connected: in series, any failure kills the system; with redundancy, some failures are tolerated. The block-diagram algebra of series, parallel and voting structures is the quantitative basis of fault-tolerant design, from triple modular redundancy in safety-critical logic to spare lanes in links.

## Series systems

If a system fails when any of $n$ independent parts fails, its reliability is the product $R_{sys}(t) = \prod_i R_i(t)$. For identical Weibull parts this is $e^{-n(t/\eta)^\beta}$, itself Weibull with the same shape and scale $\eta\,n^{-1/\beta}$. A chip with 10,000 identical electromigration-critical wire segments, each with $\eta = 1000$ years and $\beta = 2$, behaves like one part with $\eta = 1000/100 = 10$ years. This is the weakest-link effect: the more sites, the earlier the first failure, and lifetimes that are generous per site can be marginal for a whole chip. Hazard rates add in series, which is the same statement in rate form.

## Triple modular redundancy

Triple modular redundancy runs three copies of a module and takes a majority vote, so it survives any single failure. With module reliability $R$ and a perfect voter, the system works if at least two copies do:
$$R_{TMR} = R^3 + 3R^2(1 - R) = 3R^2 - 2R^3.$$
With $R = 0.95$, $R_{TMR} = 0.99275$. TMR helps only when $R > 0.5$; below that, voting on mostly failed copies is worse than one copy. Over a long mission with exponential modules, TMR's reliability eventually drops below a single module's, so it suits short missions or systems that repair failed copies. The voter is a series element and must be far more reliable than the modules. For soft errors, TMR with frequent resynchronization is the standard defence in radiation-exposed logic.

## Parallel and k-of-n redundancy

A parallel system works if any of $m$ independent paths works, so $R = 1 - (1 - R_1)^m$. Three paths each at 0.9 give $1 - 0.001 = 0.999$. The general $k$-of-$n$ structure, which needs $k$ of $n$ units, is a binomial tail and contains series ($k = n$), parallel ($k = 1$) and TMR ($k = 2$, $n = 3$) as special cases. Independence is the critical assumption: a shared supply, a shared clock or a common-mode fault such as a design bug present in every copy defeats redundancy entirely, which is why safety standards require analysis of common-cause failures and sometimes diverse implementations.

## Going deeper

Siewiorek and Swarz, "Reliable Computer Systems: Design and Evaluation" (3rd ed., 1998); Koren and Krishna, "Fault-Tolerant Systems" (2007), chapters 1 and 2.
