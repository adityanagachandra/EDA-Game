# Power density and hot spots

Power density is power per unit area, and its spatial variation, not the total power, decides the peak temperature on a die. Modern chips mix low-density caches with high-density arithmetic blocks, so the hottest spot can be many kelvin above the die average. Floorplanners, thermal sensors and throttling policies are all designed around hot spots, so computing them is a core package-chip co-design skill.

## Average power density

Average power density is $q'' = P/A$. A 300 W die of 6 cm^2 averages 50 W/cm^2. For scale, a kitchen hot plate is a few tens of W/cm^2, and high-performance processors exceed that on average and by several times locally. The unit conversion matters: W/mm^2 times 100 gives W/cm^2. Average density sets the heat-sink and package design, because the total heat must leave; but the cooling solution sized for the average does not guarantee the peak, because heat entering a small area must spread laterally before it reaches the large sink, and spreading has its own resistance.

## Hot-spot density ratio

A block with power $P_b$ over area $A_b$ has local density $P_b/A_b$, and its ratio to the die average measures how concentrated the heat is. A 40 W execution cluster on 25 mm^2 has 1.6 W/mm^2, or 160 W/cm^2, against a 50 W/cm^2 average, a ratio of 3.2. Hot-spot ratios of several times are common because activity is concentrated in arithmetic and register-file circuits. The ratio guides floorplanning: spreading hot blocks apart, interleaving them with cooler caches, and placing them under the best part of the heat-sink contact all lower the peak without changing total power.

## Temperature above the average

A simple vertical model gives each region a temperature rise of $q'' R''$, where $R''$ is an area-specific thermal resistance in K mm^2/W from that region to the coolant. With $R'' = 20$ K mm^2/W, a 0.5 W/mm^2 region rises 10 K and a 1.6 W/mm^2 hot spot rises 32 K, a 22 K difference. The vertical model overstates the hot spot, because lateral spreading in silicon and copper carries heat sideways; the smaller the hot spot, the more it spreads, so a tiny hot spot is cooler than this estimate. Thermal simulation captures spreading, and thinned dies in 3D stacks spread much less, which is why hot spots are worse in 3D.

## Going deeper

Skadron and others, "Temperature-aware microarchitecture", ISCA (2003); Incropera and others, "Fundamentals of Heat and Mass Transfer", on spreading resistance.
