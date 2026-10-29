# Minimum area and via rules

Area and via rules are where routers most often produce violations: short wire stubs are too small, and wide wires need more than one via cut. These rules also tie verification to yield, because a single via cut is one of the most failure-prone features in a chip. The computations are small, but they recur on millions of connections.

## Minimum area

Minimum area bounds the area of every merged polygon on a layer. A metal stub 50 nm wide and 122 nm long has area 6100 square nm; with an 8000 square nm rule it needs a length of $\lceil 8000 / 50 \rceil = 160$ nm, so 38 nm more. Routers fix area violations by extending a wire end along its direction, which can in turn create a new spacing violation in front of it. Minimum enclosed area (for holes) is the dual rule. Because area is a property of the merged polygon, two short stubs that touch can pass together while each would fail alone.

## Fitting a via array

Where two wide wires overlap, a via array connects them. With cuts of size $v$, cut spacing $s$ and enclosure $e$, the usable length along x is $W - 2e$, and $n$ cuts need $nv + (n - 1)s$, so $n = \lfloor (W - 2e + s) / (v + s) \rfloor$. For a 313 by 480 nm overlap with $v = 50$, $s = 50$ and $e = 30$: $\lfloor 303 / 100 \rfloor = 3$ columns and $\lfloor 470 / 100 \rfloor = 4$ rows, 12 cuts. Real decks add rules for cut arrays (larger spacing inside big arrays), so generators must follow the deck, not this formula alone.

## Minimum-cut rules

A wide wire carries more current and its via is a larger reliability risk, so decks require multiple cuts when a via lands on a wide wire: for example 2 cuts at widths of at least 250 nm and 4 cuts at widths of at least 740 nm. A connection with width 600 nm and 1 cut fails; one with width 690 nm and 3 cuts passes. Separately, DFM flows add redundant cuts wherever space allows, even where the rule does not demand them. The rule check needs the via's connectivity to both wires and their widths, which is why it is implemented on derived layers that combine cut and metal information.
