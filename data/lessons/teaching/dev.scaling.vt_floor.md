# The threshold-voltage floor

The threshold voltage has a floor set by leakage: the off-current at $V_{GS} = 0$ is the current at threshold divided by ten for every swing's worth of threshold. With swing stuck near 60 to 80 mV/decade, a leakage budget fixes a minimum threshold, and with it a minimum useful supply voltage. This is the arithmetic that ended voltage scaling and that every low-power design trade-off still runs into.

## Minimum threshold for an off-current budget

Define threshold by a constant current, $I_{V_t}$ per micron at $V_{GS} = V_t$. Below threshold the current falls a decade per swing $S$, so
$$I_{off} = I_{V_t}\cdot 10^{-V_t/S} \quad\Rightarrow\quad V_t \ge S\log_{10}\frac{I_{V_t}}{I_{off,max}}.$$

Example: $I_{V_t} = 300$ nA/µm, a budget of 100 pA/µm and $S = 75$ mV/decade. The ratio is 3000, so $V_t \ge 75 \times 3.477 = 261$ mV.

Each decade of leakage budget costs one swing of threshold. A high-performance device with a loose budget might sit near 200 mV; a low-power mobile device with a strict budget near 400 mV or more. With supply around three to four thresholds for reasonable overdrive and noise margins, the floor on $V_t$ becomes a floor on $V_{DD}$. Steeper-swing devices (tunnel FETs, negative-capacitance FETs) are research attempts to break this floor.

## The hot corner sets the floor

Swing is proportional to absolute temperature for a fixed slope factor: $S(T) = S_{300}\cdot T/300$. Leakage budgets must hold at the hottest operating temperature, where swing is worst and threshold is also lower.

Example: $S_{300} = 70$ mV/decade, $I_{V_t} = 200$ nA/µm and a budget of 50 pA/µm. At 300 K, $V_t \ge 70 \times \log_{10}4000 = 252$ mV. At 398 K (125 °C) the swing is 92.9 mV/decade and the floor rises to 335 mV. This simple version ignores the drop of threshold with temperature itself, which raises the requirement further.

So the threshold must be set for hot silicon, but delay often matters most at the slow corner, which at low voltage may be cold because of temperature inversion. That tension, leakage set hot, speed set by whichever corner is slowest, is why multi-corner multi-mode signoff exists and why threshold flavor choice is corner-dependent.

## The on-off ratio

The ratio of on-current to off-current measures how good a switch a transistor is. With $I_{off} = I_{V_t}10^{-V_t/S}$,
$$\log_{10}\frac{I_{on}}{I_{off}} = \log_{10}\frac{I_{on}}{I_{V_t}} + \frac{V_t}{S}.$$

Example: $I_{on} = 1000$ µA/µm, $I_{V_t} = 0.3$ µA/µm, $V_t = 300$ mV and $S = 75$ mV/decade. The ratio is $\log_{10}(3333) + 4 = 7.52$ decades.

Logic needs several decades of on-off ratio so that millions of idle devices do not swamp the active ones. Raising threshold increases the off-side decades but lowers $I_{on}$ through reduced overdrive; lowering supply cuts $I_{on}$ directly. Every point on the power-performance curve that a process or library offers is a choice of these terms. Better swing, from FinFET and nanosheet electrostatics, buys extra decades at the same threshold, which processes spent on lower threshold and supply.

## Going deeper

Frank, Dennard and others, "Device scaling limits of Si MOSFETs and their application dependencies", Proceedings of the IEEE (2001); Ionescu and Riel, "Tunnel field-effect transistors as energy-efficient electronic switches", Nature (2011).
