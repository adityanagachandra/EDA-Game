# Leakage power and multi-Vt

Leakage is the power a circuit burns while doing nothing: current through transistors that are nominally off. In modern processes it can be a large fraction of total power, especially in idle blocks. Libraries offer several threshold-voltage (Vt) flavours of every cell precisely to trade leakage against speed, and power and optimization tools need the exponential law that governs that trade.

## Subthreshold current is exponential in Vt

Below threshold a transistor is not perfectly off: diffusion current flows that depends exponentially on the gate overdrive,
$$I_{sub} \propto 10^{(V_{GS} - V_t)/S},$$
where the subthreshold swing $S$ is the gate voltage change for a tenfold change in current. $S$ is at least about 60 mV per decade at room temperature and typically 70 to 100 mV per decade in real devices. For an off device $V_{GS} = 0$, so its leakage falls by a decade for every $S$ of threshold. Leakage also rises steeply with temperature, both through $S$ and through a falling $V_t$, which is why leakage is analyzed at the hot corner. Other components, gate tunnelling and junction leakage, add to it.

## Leakage ratios between Vt flavours

Two flavours whose thresholds differ by $\Delta V_t$ differ in off-current by
$$\frac{I_{LVT}}{I_{HVT}} = 10^{\Delta V_t / S}.$$
With $\Delta V_t = 100$ mV and $S = 90$ mV per decade the ratio is $10^{1.11} = 12.9$. An SVT device leaking 20 nA becomes about 1.6 nA in a flavour 100 mV higher. A cell's leakage power is its leakage current times $V_{DD}$, summed over its off devices, so a block of a million cells averaging 5 nW each leaks 5 mW. Libraries typically provide several flavours, for example ultra-low, low, standard and high Vt, with identical footprints and pins so that they can be swapped in place.

## DIBL and the cost of speed

Leakage also depends on the drain voltage. Drain-induced barrier lowering reduces the threshold as the drain-source voltage grows, $V_t = V_{t0} - \eta V_{DS}$, so an off device with its full supply across it leaks more. Lowering $V_{DD}$ from 1.0 V to 0.7 V with $\eta = 0.1$ raises the effective threshold by 30 mV, cutting leakage by $10^{30/80} = 2.4$ times at $S = 80$ mV per decade, on top of the lower supply in $P = I V_{DD}$. The other side of the trade is speed: a higher threshold reduces gate overdrive $V_{DD} - V_t$ and so drive current. Every Vt choice is therefore a leakage-versus-delay decision, which is what multi-Vt optimization in synthesis and physical design automates.

## Going deeper

Roy, Mukhopadhyay and Mahmoodi-Meimand, "Leakage current mechanisms and leakage reduction techniques in deep-submicrometer CMOS circuits" (Proceedings of the IEEE, 2003); Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), section 5.5.
