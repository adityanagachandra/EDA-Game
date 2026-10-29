# Interpolating NLDM tables

Liberty's non-linear delay model stores each timing arc as small tables of delay and output slew, indexed by the input slew and the output load. A real pin almost never lands on a table point, so every delay a timer computes is an interpolation, and some are extrapolations. This lookup runs for every arc evaluation in the design, which makes it one of the hottest loops in any timer.

## Linear interpolation along one axis

A typical table has index_1 for input slew and index_2 for output load, often seven points each, with the delay at every grid point. If the query lies exactly on a grid line, for instance the slew equals an index_1 value, only the load needs interpolating: $v = v_1 + t\,(v_2 - v_1)$ with $t = (c - c_1)/(c_2 - c_1)$. Take load points 8 and 16 fF with delays 67 and 97 ps at slew 40 ps. At 14 fF, $t = 0.75$ and $v = 67 + 0.75 \times 30 = 89.5$ ps. Linear interpolation is exact for a function that is linear in load, which delay roughly is over a small window, and that is why tables need only a few points per axis.

## Bilinear interpolation in a window

In general the query falls inside a $2 \times 2$ window of slews $s_1 < s_2$ and loads $c_1 < c_2$. Interpolate along load at both slews, then along slew between the two results. In one formula,
$$v = \frac{v_{11}(s_2-s)(c_2-c) + v_{12}(s_2-s)(c-c_1) + v_{21}(s-s_1)(c_2-c) + v_{22}(s-s_1)(c-c_1)}{(s_2-s_1)(c_2-c_1)},$$
where $v_{ij}$ is the value at slew $s_i$ and load $c_j$. With $s$ at fraction $t_s$ and $c$ at fraction $t_c$, this is $(1-t_s)(1-t_c)v_{11} + (1-t_s)t_c v_{12} + t_s(1-t_c) v_{21} + t_s t_c v_{22}$. For corners 60, 68, 67 and 82 ps, $t_s = 0.75$ and $t_c = 0.375$, the result is 70.22 ps. The same lookup gives the output slew, which feeds the next stage.

## Extrapolation beyond the table

When the load or slew lies outside the table, timers such as OpenSTA extrapolate with the same formula using the nearest edge window, so a fraction greater than 1 or less than 0 appears. With load index points 1 and 2 fF and a query at 2.25 fF, $t_c = 1.25$. Extrapolation is only as good as the linear trend of the last window, and delay curves bend: a heavily loaded cell's real delay often grows faster than the edge slope predicts. Characterization ranges are chosen to cover the loads and slews the design should see, and max_capacitance and max_transition design rules keep the design inside them. A timer that silently extrapolates far beyond the table produces numbers no one characterized, so reports flag it.

## Going deeper

The Liberty reference manual on lu_table_template and NLDM; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 3; the OpenSTA source, TableModel.
