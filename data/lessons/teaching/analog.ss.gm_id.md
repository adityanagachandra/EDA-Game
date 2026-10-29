# The gm/Id methodology

The gm/Id methodology sizes transistors by choosing their transconductance efficiency $g_m/I_D$ first. That one number sets the inversion level, and with it the trade-off between current, speed, gain and signal swing. It works from weak to strong inversion where the square law does not, and it maps naturally onto tables of simulated device data, which is why it underlies most modern hand sizing and many automated sizers.

## Transconductance efficiency

$g_m/I_D$ measures how much transconductance each unit of current buys, in $\mathrm{V}^{-1}$. Most analog specs ask for a transconductance: GBW sets $g_m$ of an input pair, noise sets $g_m$ of critical devices, a pole position sets $g_m$ of a driver. Given the required $g_m$ and a chosen $g_m/I_D$, the current follows immediately: $I_D = g_m/(g_m/I_D)$. A stage needing $g_m = 1$ mS at $g_m/I_D = 20\ \mathrm{V}^{-1}$ needs $I_D = 50\ \mu$A; at $g_m/I_D = 10$ it needs 100 microamperes. Power is current times supply, so $g_m/I_D$ is a direct power knob. The width then follows from the current density at that $g_m/I_D$, which is the job of a lookup table.

## Inversion level and its limits

$g_m/I_D$ is a monotone function of the inversion level. In strong inversion the square law gives $g_m/I_D = 2/V_{OV}$: $g_m/I_D = 10$ corresponds to $V_{OV} = 200$ mV. As the gate voltage drops toward threshold the device enters moderate and then weak inversion, where the current is exponential in $V_{GS}$ and $g_m/I_D$ approaches its ceiling $1/(nU_T)$, with $U_T = kT/q \approx 25.9$ mV at 300 K and slope factor $n$ around 1.2 to 1.5. For $n = 1.3$ the ceiling is about 29.7 per volt. Typical design values span 5 (strong inversion, fast, large swing) to about 25 (weak inversion, efficient, slow). Unlike the square law, the methodology stays valid across the whole range because it does not assume a current equation.

## The speed and efficiency trade-off

Efficiency is not free. At high $g_m/I_D$ the device is wide for its current, so its gate capacitance is large and its transit frequency $f_T \approx g_m/(2\pi C_{gg})$ falls. At low $g_m/I_D$ it is small and fast but burns current. A device at $g_m/I_D = 15$ with $I_D = 40\ \mu$A has $g_m = 600\ \mu$S; with $C_{gg} = 50$ fF, $f_T = 600/(2\pi \cdot 50) = 1.91$ GHz (microsiemens over femtofarads is $10^9$ per second). The product $(g_m/I_D) f_T$ peaks in moderate inversion, which is why many optimized designs sit there. A sizer that searches over $g_m/I_D$ and $L$ instead of $W$ and $V_{GS}$ works in variables that map directly onto these trade-offs and stay well-scaled across processes.

## Going deeper

Jespers and Murmann, "Systematic Design of Analog CMOS Circuits Using Pre-Computed Lookup Tables" (2017); Silveira, Flandre and Jespers, "A gm/ID based methodology for the design of CMOS analog circuits" (IEEE JSSC, 1996).
