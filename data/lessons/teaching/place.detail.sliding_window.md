# Sliding-window reordering

Local reordering takes a small window of consecutive cells in a row, tries every order of them packed into the same span and keeps the best. With three or four cells the enumeration is tiny, and sliding the window along every row catches improvements that pairwise swaps miss. It is a standard pass in FastDP-style and GPU detailed placers.

## Packing a window in a new order

Take $k$ adjacent cells in a row occupying a span from $x_0$. For an order $\pi$, pack the cells left to right from $x_0$ with no gaps (or preserving the window's total whitespace in fixed positions): the first cell sits at $x_0$, the next at $x_0$ plus the first width, and so on. Pin positions follow from cell positions; for a pin at a cell's centre it is $x + w/2$. With widths A = 2, B = 4 and C = 8 packed from 18 in order A B C, the centres are 19, 22 and 30. The packing stays legal because the cells fill exactly the same span, so only wirelength changes, and only on nets that touch the window.

## Trying every order

With $k$ cells there are $k!$ orders: 6 for three cells, 24 for four, 120 for five. For each order, compute the x span of every net touching the window and sum. The y spans do not change, because the cells stay in the row. Keep the order with the smallest sum, and apply it if it beats the current order. Example: A to fixed pins at 28 and 55, B to 39, C to 13, and a net from B to C, with widths 4, 2, 8 from x = 15. The current order costs some total; trying all six orders finds that C B A costs only 54 um. Nets between two window cells must be evaluated with both cells' new positions.

## Window size and sweeping

The window slides across each row with some overlap, for example by one or two cells at a time, so an improvement found in one window can enable another in the next. Window size is a trade-off: $k!$ grows fast, so $k = 3$ or 4 is typical, and dynamic programming over positions can handle somewhat larger windows when cells must stay in their row segment. Because windows in different rows, or far apart in one row, do not share cells, many can be optimised in parallel if they also share no nets, which is how GPU detailed placers batch them. The pass is cheap and complements swaps, which change positions of distant cells, and independent-set matching, which reassigns several cells at once.

## Going deeper

Pan, Viswanathan and Chu, "An efficient and effective detailed placement algorithm" (ICCAD 2005); Lin et al., "ABCDPlace: accelerated batch-based concurrent detailed placement on multithreaded CPUs and GPUs" (IEEE TCAD, 2020).
