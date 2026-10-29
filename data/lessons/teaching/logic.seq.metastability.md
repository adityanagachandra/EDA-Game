# Metastability

When a flip-flop's input changes inside its setup-hold window, the internal bistable element can be left balanced between 0 and 1, and it takes an unbounded, random time to fall to one side. This metastability cannot be eliminated, only made improbable. It matters whenever a signal is asynchronous to the clock that samples it, which happens at every clock-domain crossing and every external input.

## How often a flip-flop goes metastable

Let the sampling clock run at $f_c$ and the asynchronous data toggle at rate $f_d$. Each data transition lands at a uniformly random phase of the clock. A flip-flop has an effective aperture $T_w$ (often written $T_0$), a few picoseconds around the edge, inside which a transition causes metastability. The rate of metastable events is approximately

$$R = f_c \, f_d \, T_w.$$

With $f_c = 500$ MHz, $f_d = 10$ MHz and $T_w = 20$ ps: $R = 5\times10^{8} \cdot 10^{7} \cdot 2\times10^{-11} = 10^{5}$ events per second. Metastability is not rare; a hundred thousand events a second is ordinary. What must be rare is an event that lasts long enough to cause a failure.

## Exponential resolution

Near the balance point the cross-coupled inverters behave as a linear system with positive feedback, so a small initial imbalance $v_0$ grows as $v_0 e^{t/\tau}$, where $\tau$ is the regeneration time constant (a few picoseconds in modern processes). Resolution takes longer the closer the start is to balance, and with a uniform starting point the probability that the flip-flop is still unresolved after time $t$ is

$$P(\text{unresolved at } t) = e^{-t/\tau}.$$

With $\tau = 15$ ps, after 60 ps the probability is $e^{-4} \approx 0.0183$; after 150 ps it is $e^{-10} \approx 4.5 \times 10^{-5}$. The tail is memoryless: an element that has already waited does not become more likely to resolve in the next instant, which is why the exponential model fits so well.

## Buying reliability with time

Because the failure probability is $e^{-t/\tau}$, improving reliability by a factor $X$ costs extra resolution time

$$\Delta t = \tau \ln X.$$

With $\tau = 20$ ps, a factor of $10^{6}$ costs $20 \cdot \ln 10^{6} \approx 276.3$ ps. That is the logarithmic bargain that makes synchronizers work: each extra flip-flop stage adds nearly a whole clock period of resolution time and multiplies the mean time between failures by an enormous factor.

The parameters $\tau$ and $T_w$ are not in a standard Liberty timing table; they are characterized separately by circuit simulation, and $\tau$ degrades at low supply voltage, where the inverters are weak. A synchronizer that is safe at nominal voltage can become unsafe in a low-voltage mode.

## Going deeper

Ginosar, "Metastability and Synchronizers: A Tutorial" (IEEE Design and Test of Computers, 2011); Kleeman and Cantoni, "Metastable Behavior in Digital Systems" (IEEE Design and Test, 1987).
