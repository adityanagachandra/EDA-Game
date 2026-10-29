# Chemical-mechanical polishing and density

Chemical-mechanical polishing (CMP) flattens the wafer after deposition so that the next lithography step stays in focus, and it defines damascene copper lines by removing excess metal. Its removal rate depends on the local pattern density, so a layout's density map is a process input. This is why design rules impose minimum and maximum density in every window, and why fill tools insert dummy shapes.

## Preston's law

A rotating pad with an abrasive, chemically active slurry presses on the wafer. Removal rate follows Preston's empirical law, $RR = K_p P v$, proportional to applied pressure $P$ and relative velocity $v$, with $K_p$ capturing slurry chemistry and pad properties. Scaling from a known operating point needs no $K_p$: a step that removes 200 nm/min at 20 kPa and 1.0 m/s should remove $200 \times (25 \times 1.2)/(20 \times 1.0) = 300$ nm/min at 25 kPa and 1.2 m/s. Real processes deviate at low pressure and when chemistry dominates, but Preston's law is the model a CMP simulator starts from.

## Pattern density sets local removal

On a patterned wafer the pad touches raised features first. If a region has pattern density $\rho$, the fraction of area that is raised, the pressure on the raised parts is about $P/\rho$, so they polish at the blanket rate divided by $\rho$. The time to remove a step of height $h$ is $t = h \rho / K$, where $K$ is the blanket rate. At 200 nm/min blanket, a 400 nm step clears in 0.6 min at density 0.3 but in 1.6 min at density 0.8. Regions of different density therefore finish at different times and end up at different heights. In copper CMP the same physics produces dishing of wide lines and erosion of dense arrays.

## Density rules and dummy fill

To keep the post-CMP surface within the lithography focus budget, foundries require metal density between a minimum and maximum in every window, often tens of micrometres on a side and stepped across the die. If a 50 um by 50 um window must reach 25% density and holds 400 um^2 of metal, it needs $0.25 \times 2500 - 400 = 225$ um^2 of fill. Fill tools add floating or tied dummy shapes, which adds coupling capacitance to nearby signals, so timing-aware fill keeps dummies away from critical nets. Extraction must see the fill, or timing will be optimistic.

## Going deeper

Preston (1927) on glass polishing; Stine et al., "Rapid Characterization and Modeling of Pattern-Dependent Variation in Chemical-Mechanical Polishing", IEEE Transactions on Semiconductor Manufacturing (1998).
