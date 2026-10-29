# Fringe capacitance

The field from a wire does not stop at its bottom face. Field lines leave its sidewalls and top and curve down to the conductor below, adding fringe capacitance along each edge. For narrow wires this fringe term is larger than the area term, so any extraction model that is accurate at advanced nodes is mainly a model of fringing. Extractors separate the two because they scale differently: area capacitance with width, fringe with length alone.

## Area plus edges

A practical model writes a wire's capacitance to ground as an area term plus a fringe term per edge: $C = c_a W L + 2 c_f L$, where $c_a$ is in aF/um$^2$ and $c_f$ in aF/um per edge. The fringe term does not depend on width, because the fringing field near an edge looks the same whether the wire is wide or narrow, once it is wider than a few dielectric thicknesses. A wire 0.05 um wide and 180 um long with $c_a = 80$ aF/um$^2$ and $c_f = 75$ aF/um has area capacitance $80 \times 0.05 \times 180 = 720$ aF and fringe $2 \times 75 \times 180 = 27000$ aF: 27.72 fF in total, almost all of it fringe. Such coefficients come from field-solver runs, fitted per layer.

## The fringe share

The fraction of ground capacitance that is fringe is $2c_f/(c_a W + 2c_f)$, independent of length. For $W = 0.05$ um, $c_a = 80$ aF/um$^2$ and $c_f = 30$ aF/um, the area term is 4 aF/um and the fringe 60 aF/um: 93.8% fringe. Width sizing therefore changes capacitance much less than the parallel-plate formula suggests: doubling the width of this wire adds only 4 aF/um, under 7%, while halving its resistance. This asymmetry is why wide wires are effective for long nets and why sizing a wire to cut delay does not cost proportionally more capacitance. It also explains why extraction errors in fringe models matter more than errors in area models.

## The crossover width

Fringe exceeds area capacitance when $c_a W < 2 c_f$, that is below the width $W^* = 2c_f/c_a$. With $c_a = 20$ aF/um$^2$ and $c_f = 80$ aF/um per edge, $W^* = 8$ um. Every signal wire at an advanced node is far below such a crossover, so its ground capacitance is dominated by edges. The model's assumption that $c_f$ is width-independent fails for very narrow wires, where the two edges' fringe fields interact, and in dense wiring, where neighbours intercept fringe lines and turn them into coupling. Production extractors therefore tabulate fringe against neighbour spacing too, rather than using one constant per layer.

## Going deeper

Sakurai and Tamaru, "Simple formulas for two- and three-dimensional capacitances", IEEE Transactions on Electron Devices (1983); Cheng, Lillis, Lin and Chang, "Interconnect Analysis and Synthesis" (2000).
