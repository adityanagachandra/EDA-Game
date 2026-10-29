# Allocation of functional units

Allocation chooses how many units of each type the datapath contains and which unit types to use at all. It sets the area budget that scheduling works within, and it is where the first latency-area trade-off is made. Allocation decisions also include choosing between specialized units and multi-function units, such as separate adders and subtractors or combined add/subtract units.

## Datapath area from an allocation

Each unit type has a cost from the operator library, measured in LUTs and DSP blocks on an FPGA or in cell area on an ASIC. The functional-unit and register area of an allocation is $\sum_k a_k c_k$ over unit counts $a_k$ and costs $c_k$. Four adders of 80 units, two multipliers of 300 and six registers of 60 cost $320 + 600 + 360 = 1280$ units. Multipliers usually dominate, which is why allocation focuses on them, but the count of registers and the multiplexers that sharing introduces can matter as much in heavily shared designs. Interconnect is decided later by binding, so allocation-time area is a lower estimate; a design with fewer units but many shared multiplexers can end up larger than one with more units.

## Minimum allocation for a latency target

Under a latency bound of $L$ steps, $n_k$ single-cycle operations of type $k$ need at least $\lceil n_k / L \rceil$ units, so the smallest functional-unit area any schedule could use is $\sum_k \lceil n_k / L \rceil c_k$. Eleven additions and nine multiplications in 8 steps need at least 2 adders and 2 multipliers; at 60 and 700 units that is 1520. Whether a schedule reaches this allocation depends on the dependences, but it is the natural first allocation to try, increasing a type's count only when the scheduler cannot meet the latency. Sweeping $L$ traces the latency-area curve: halving the latency roughly doubles the units of every type that is resource-bound, while types with few operations stay at one unit.

## Specialized or multi-function units

A combined unit can execute several operation types, at a cost higher than any one specialized unit but lower than two. With five additions and three subtractions in 2 steps, separate units need $\lceil 5/2 \rceil = 3$ adders and $\lceil 3/2 \rceil = 2$ subtractors; at 30 and 50 units that is 190. Combined add/subtract units of 40 units each need $\lceil 8/2 \rceil = 4$, costing 160. Pooling operation types smooths the per-step demand, since a step heavy in additions and light in subtractions no longer wastes a subtractor. The cost is a function-select signal from the controller, slightly more delay through the unit, and sometimes worse scheduling flexibility. Choosing unit types is part of allocation, and HLS tools expose it through resource directives.

## Going deeper

Gajski, Dutt, Wu and Lin, "High-Level Synthesis" (1992), chapter on allocation; De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 6.
