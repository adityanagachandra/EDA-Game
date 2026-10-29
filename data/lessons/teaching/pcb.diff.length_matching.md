# Length matching from a skew budget

Signals that must arrive together, the two traces of a differential pair or the data lines of a memory byte lane and their strobe, are given a skew budget in picoseconds. The layout tool turns that budget into a length tolerance using the layer's delay per inch, then tunes short traces with serpentines until every member of the group is inside the window. Matching is a timing constraint expressed in geometry, and the conversion is where errors creep in.

## Skew from a length mismatch

Skew is the arrival-time difference. For two traces on the same layer, skew equals the length difference times the delay per inch: $\Delta t = \Delta L \cdot t_{pd}$. Example: a 50 mil mismatch on a layer with 170 ps per inch is $0.05 \cdot 170 = 8.5$ ps. Within a differential pair, **intra-pair skew** turns part of the differential signal into common mode during the skewed interval, which shrinks the eye and raises emission; for multi-gigabit links the budget is a small fraction of the unit interval. Mismatch accumulates at every bend and pin: the inner trace of a turn is shorter than the outer, so routers compensate near the mismatch rather than at one end, keeping the two halves aligned along the whole route.

## Length tolerance from a budget

Inverting gives the tolerance: $\Delta L_{max} = \Delta t_{max} / t_{pd}$. Example: a 5 ps budget on stripline with $\varepsilon_r = 3.9$ has $t_{pd} = 1000\sqrt{3.9}/11.8 \approx 167.4$ ps per inch, so $\Delta L_{max} = 5/167.4 \approx 0.0299$ in, which rounds down to 29 mil; rounding up would exceed the budget. The budget itself comes from a timing analysis of the interface: setup and hold windows at the receiver minus driver uncertainty, jitter and other terms, so the length tolerance is a downstream product of a timing budget. When nets change layers the conversion must be per segment, which is why tools match delay rather than length.

## Matching a group to a strobe

In a memory byte lane every DQ must arrive within a window around its strobe, say $\pm 5$ ps. The tool computes each member's delay, compares it with the strobe's, and flags those outside the window. Example: a 2000 mil strobe at 160 ps per inch and a 5 ps window allow $\pm 31.25$ mil, so a 2040 mil DQ is outside the window: either it is shortened, or the strobe and the other members are lengthened toward it. Tuning adds **serpentines** (accordion meanders) to the short members. Meanders couple to themselves: parallel segments too close together let the signal shortcut through crosstalk, so the real delay added is less than the added length. Rules set a minimum gap between meander segments, often a few trace widths, and SI simulation confirms the result.

## Going deeper

JEDEC DDR4 and DDR5 design guides from memory vendors, for byte-lane matching rules. Bogatin, "Signal and Power Integrity Simplified" (3rd ed., 2018), on skew and mode conversion.
