# Heat conduction and Fourier's law

Heat leaves a chip mainly by conduction through solids: silicon, solder, thermal interface material, lid and heat sink, before convection carries it away. Conduction obeys Fourier's law, the heat-flow analogue of Ohm's law, and its one-dimensional form gives the thermal resistance of every layer in a package. That resistance is the building block of all package thermal models.

## Fourier's law

Fourier's law states that heat flux is proportional to the temperature gradient: $q'' = -k\,\nabla T$, with thermal conductivity $k$ in W/(m K). In one dimension through a slab of thickness $L$ with faces at $T_1$ and $T_2$, the heat flow is $Q = k A (T_1 - T_2)/L$. The steady heat equation $\nabla \cdot (k \nabla T) + q''' = 0$ follows from energy conservation. Conductivities span four orders of magnitude: copper about 400, silicon about 150 at room temperature, solder around 50, thermal interface materials a few, mold compound and underfill below one, and silicon dioxide about 1.4. Low-conductivity layers, even thin ones, therefore dominate the heat path.

## Resistance of a layer

Rearranging gives the thermal resistance $R = \Delta T/Q = L/(k A)$, in K/W. With $L$ in micrometres and $A$ in mm^2 the factors of $10^{-6}$ cancel, so $R = L/(kA)$ directly. A 750 um silicon die of 100 mm^2 has $R = 750/(150 \times 100) = 0.05$ K/W; a 50 um interface layer at 4 W/(m K) over the same area has $50/400 = 0.125$ K/W, two and a half times the die's resistance despite being fifteen times thinner. At 100 W, the interface layer alone drops 12.5 K. This is why thermal interface materials and their bond-line thickness receive so much engineering attention.

## Sizing a material for a budget

The formula inverts to size a material: if a layer of thickness $L$ and area $A$ may drop at most $\Delta T$ while carrying $Q$, it needs $k \ge Q L/(A \Delta T)$. For 150 W through a 30 um bond line over 400 mm^2 with a 3 K budget, $k \ge 150 \times 30/(400 \times 3) = 3.75$ W/(m K). Two caveats. One-dimensional resistance assumes heat flows straight through; if the heat source is smaller than the layer, heat spreads, and the effective resistance is larger near the source and smaller far from it. And conductivity depends on temperature: silicon's falls as it heats, so hot dies conduct worse, a nonlinearity that full simulators include.

## Going deeper

Incropera, DeWitt, Bergman and Lavine, "Fundamentals of Heat and Mass Transfer" (7th ed., 2011), chapters 2 and 3; Lasance and Poppe (eds.), "Thermal Management for LED Applications" (2014), on thermal resistance.
