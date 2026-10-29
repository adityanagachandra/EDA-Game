# Power electronics and power semiconductor devices

Power electronics converts electrical energy between voltages and forms with switching circuits: buck and boost converters, inverters, chargers. Its components are power semiconductors (silicon MOSFETs and IGBTs, silicon carbide and gallium nitride devices), inductors and capacitors. Design is dominated by efficiency, thermal limits and electromagnetic effects, and its simulation and modelling problems differ from those of small-signal analog design.

## The buck converter

A buck converter switches its input onto an inductor with duty cycle $D$ and filters the result with an LC network. In steady state the inductor's average voltage must be zero (volt-second balance): during the on-time it sees $V_{in} - V_{out}$ and during the off-time $-V_{out}$, so $D(V_{in} - V_{out}) = (1 - D)V_{out}$, which gives $V_{out} = DV_{in}$ for an ideal converter in continuous conduction. Stepping 12 V to 3.3 V needs $D = 0.275$. Real converters need a slightly larger duty cycle to cover switch and inductor losses, and a feedback controller adjusts $D$ to regulate the output. Boost converters obey $V_{out} = V_{in}/(1 - D)$ by the same balance argument.

## Inductor ripple

During the on-time $DT$, the inductor current rises with slope $(V_{in} - V_{out})/L$, so the peak-to-peak ripple is $\Delta I = (V_{in} - V_{out})D/(Lf_s)$. For 12 V to 3.3 V with $L = 4.7\ \mu$H at 1 MHz: $\Delta I = 8.7 \cdot 0.275/(4.7 \times 10^{-6} \cdot 10^6) = 0.509$ A. Ripple sets the inductor's peak current and core loss, the output voltage ripple through the output capacitor, and the boundary of continuous conduction (the converter enters discontinuous mode when the average load current falls below $\Delta I/2$). Raising the switching frequency shrinks the inductor and capacitor for the same ripple, which is the main reason to switch faster, limited by switching losses.

## Losses in the switch and wide-bandgap devices

A switch dissipates conduction loss $I^2R_{on}D$ while on and switching loss during each transition, when voltage and current overlap; a common estimate is $\frac{1}{2}V_{in}I(t_r + t_f)f_s$. With 5 A, 10 milliohms, $D = 0.275$, 12 V, 10 ns rise and fall and 1 MHz: conduction 0.069 W, switching 0.6 W. Switching dominates at high frequency, which is where wide-bandgap devices help: silicon carbide and gallium nitride have higher critical fields, so for a given blocking voltage they achieve lower on-resistance with smaller devices and smaller capacitances, which allows faster transitions. GaN dominates lower-voltage, high-frequency converters and SiC higher-voltage applications such as traction inverters. Design tools for power electronics combine circuit simulation with device models that capture these transitions, thermal networks and parasitic inductances of packages and layouts.

## Going deeper

Erickson and Maksimovic, "Fundamentals of Power Electronics" (3rd ed., 2020); Baliga, "Fundamentals of Power Semiconductor Devices" (2nd ed., 2019).
