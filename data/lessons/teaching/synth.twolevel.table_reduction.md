# Row and column dominance

Between essential-prime selection and search, covering algorithms shrink the table with dominance rules. A column that is easier to satisfy than another makes the harder one redundant, and a row that covers no more than another (at no lower cost) can be dropped. Alternating these reductions with essential selection often solves the table outright and always shrinks the search.

## Essential rows in a covering table

In a general covering table, row $P_i$ covers a set of columns and every column must be covered. A row is **essential** if it is the only row marking some column. For the table $P_1 = \{c_1, c_2, c_4, c_5\}$, $P_2 = \{c_1, c_2, c_4, c_5\}$, $P_3 = \{c_5\}$, $P_4 = \{c_1, c_3\}$, $P_5 = \{c_2\}$, $P_6 = \{c_4\}$, only column $c_3$ has a single mark, so $P_4$ is the one essential row. Selecting it deletes columns $c_1$ and $c_3$, and the remaining table is smaller. Essential selection is safe: any cover must contain that row, so committing to it loses nothing.

## Column dominance

Let $R(c)$ be the set of rows that mark column $c$. If $R(c_j) \subseteq R(c_i)$, any row that covers $c_j$ also covers $c_i$, so once $c_j$ is covered $c_i$ is covered for free: delete $c_i$, the column with the larger row set. Two identical columns keep one. For example, with $R(c_1) = \{P_2, P_4\}$ and $R(c_5) = \{P_1, P_2, P_4\}$, column $c_5$ is deleted. The rule looks backwards at first: the dominating column is the one removed. Deleting columns never changes the set of valid covers, so column dominance is exact.

## Row dominance and the cyclic core

If row $P_q$'s remaining columns are a subset of row $P_p$'s and $P_p$ costs no more, then some minimum cover avoids $P_q$: replace it by $P_p$. So delete the dominated row (with unit cost, ties keep the lower index). Row dominance is safe for finding a minimum but may discard some alternative optima. Alternate essential selection, row dominance and column dominance until nothing changes. Either the table vanishes, giving the minimum cover from the selected rows, or a **cyclic core** remains in which every column is marked at least twice and no dominance applies. Only the cyclic core needs branch and bound, so these reductions decide how much search an exact minimizer must do.

## Going deeper

Coudert, "On solving covering problems", DAC 1996; Hachtel and Somenzi (1996), on unate covering reductions.
