# Latch regeneration and synchronizer MTBF

When a flop samples an asynchronous signal, data sometimes changes within the tiny window around the clock edge, and the flop's internal latch is left balanced between 0 and 1. How quickly it resolves depends on the latch's regeneration time constant, a transistor-level property. Synchronizer design, and the MTBF numbers that justify it, follow directly from that exponential resolution.

## Regeneration in a cross-coupled latch

A latch's storage loop is two inverters in a ring. Near its balance point, where both nodes sit at the switching threshold, each inverter is a high-gain amplifier. A small voltage difference $\Delta V$ between the nodes is amplified around the loop and grows exponentially, $\Delta V(t) = \Delta V_0 e^{t/\tau}$, where $\tau$ is set by the inverters' transconductance and node capacitance, roughly $C/g_m$. If data arrives so close to the edge that $\Delta V_0$ is tiny, the latch takes a long time to reach a valid logic level. The probability of still being unresolved after time $t$ therefore decays as $e^{-t/\tau}$. A faster process and a larger, well-sized latch both reduce $\tau$.

## The MTBF formula

A synchronizer fails if its first flop is still unresolved when the next stage samples it. With resolution time $t_r$, clock frequency $f_c$, data toggle rate $f_d$ and an effective metastability window $T_w$, the mean time between failures is
$$\text{MTBF} = \frac{e^{t_r/\tau}}{T_w f_c f_d}.$$
Example: $\tau = 15$ ps, $T_w = 20$ ps, $f_c = 1$ GHz, $f_d = 100$ MHz and $t_r = 600$ ps. The rate term is $T_w f_c f_d = 2 \times 10^{6}$ per second, and $t_r/\tau = 40$, so $\log_{10}\text{MTBF} = 40/\ln 10 - \log_{10}(2 \times 10^6) = 17.37 - 6.30 = 11.07$, about $10^{11}$ s, a few thousand years. Taking logs is the practical way to handle these numbers.

## Designing for a target MTBF

Solving for the resolution time gives $t_r = \tau \ln(\text{MTBF} \cdot T_w f_c f_d)$. Because MTBF is exponential in $t_r$, each extra $\tau \ln 10$ of resolution time adds a factor of ten. Adding a third flop to a two-flop synchronizer adds roughly one clock period minus flop overhead to $t_r$: at 1 GHz with 80 ps of overhead and $\tau = 15$ ps, that is $920/(15 \ln 10) = 26.6$ decades. This is why two or three flops suffice for almost any rate, and why synchronizer flops are special library cells with small $\tau$ and are kept close together with no logic between them. Static timing does not analyze metastability; clock-domain-crossing checkers verify the synchronizer structure instead, and $\tau$ and $T_w$ come from transistor-level characterization.

## Going deeper

Ginosar, "Metastability and synchronizers: a tutorial" (IEEE Design and Test, 2011); Weste and Harris, "CMOS VLSI Design" (4th ed.), section 10.6.
