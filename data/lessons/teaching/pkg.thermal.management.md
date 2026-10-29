# Thermal management: TIMs and heat sinks

Thermal management is the set of materials and devices that move heat from the package to the environment: thermal interface materials (TIMs), lids and spreaders, heat sinks, fans, heat pipes, vapour chambers and liquid cooling. Each contributes a resistance to the path. The designer's job is to size these elements so that the junction stays below its limit at the worst-case power and ambient.

## Interface materials

Two solid surfaces touch only at asperities, leaving air gaps that conduct poorly. A TIM fills the gaps. Its resistance is $R = BLT/(kA)$ plus two contact resistances, where BLT is the bond-line thickness. A 40 um grease at 5 W/(m K) over 300 mm^2 gives $40/(5 \times 300) = 0.0267$ K/W, about 2.7 K at 100 W. Indium and liquid-metal TIMs conduct far better than greases; solder TIMs are used between a die and a lid in high-power parts. Bond-line thickness depends on surface flatness and clamping pressure, so warpage, which bends the die and lid apart, directly raises TIM resistance.

## Required heat-sink resistance

Given the power, the junction limit and the other resistances, the heat sink must satisfy $\theta_{sa} \le (T_{j,max} - T_a)/P - \theta_{jc} - \theta_{cs}$. For 150 W, a 95 °C limit, 40 °C air, $\theta_{jc} = 0.1$ and $\theta_{cs} = 0.05$ K/W, the budget is $55/150 - 0.15 = 0.217$ K/W. If no affordable sink meets it, the options are a better TIM, a lower $\theta_{jc}$ by a soldered lid or direct die contact, a lower ambient, or less power. This subtraction is how thermal requirements flow from the system to the package team, and why a package's $\theta_{jc}$ is a competitive number.

## Convection and fin area

A heat sink finally transfers heat to air by convection, with resistance $R = 1/(h A)$, where $h$ is the heat-transfer coefficient and $A$ the wetted area. Natural convection gives $h$ of a few to about 10 W/(m^2 K); forced air tens to about 100 or more; liquid cooling far higher. To reach 0.2 K/W with forced air at $h = 50$ W/(m^2 K) needs $A = 1/(50 \times 0.2) = 0.1$ m^2, 1000 cm^2, which is why sinks have many fins. Fin efficiency below one raises the needed area further. Liquid cold plates and immersion reach much lower resistances, which is why high-power accelerators move to liquid cooling.

## Going deeper

Incropera and others, "Fundamentals of Heat and Mass Transfer", chapter 3 on extended surfaces; Prasher, "Thermal interface materials: historical perspective, status, and future directions", Proceedings of the IEEE (2006).
