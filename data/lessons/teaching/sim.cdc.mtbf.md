# Synchronizer MTBF

A synchronizer cannot eliminate metastability, so its quality is a failure rate. The mean time between failures (MTBF) model turns flop parameters, clock rates and the time allowed for resolution into a number of seconds or years. Because the resolution time enters an exponential, small design changes move MTBF by many orders of magnitude, which is why synchronizer design is a quantitative engineering decision rather than a rule of thumb.

## The MTBF formula

A flop that samples an asynchronous input goes metastable at a rate proportional to how often the input changes and how often it is sampled: $T_0 f_{clk} f_{data}$ events per second, where $T_0$ is a window parameter of the flop (in seconds), $f_{clk}$ the sampling clock frequency and $f_{data}$ the data toggle rate. A metastable flop resolves exponentially: the probability that it is still unresolved after time $t$ is about $e^{-t/\tau}$, where $\tau$ is the resolution time constant of the flop, a few to tens of picoseconds in modern processes. A failure happens when metastability lasts longer than the available resolution time $t_r$. So
$$\text{MTBF} = \frac{e^{t_r/\tau}}{T_0 \, f_{clk} \, f_{data}}.$$
Both $T_0$ and $\tau$ are characterized per flop and per process corner.

## Computing MTBF

Take $\tau = 20$ ps, $T_0 = 50$ ps, $f_{clk} = 500$ MHz, $f_{data} = 50$ MHz and $t_r = 400$ ps. The denominator is $50 \times 10^{-12} \cdot 5 \times 10^8 \cdot 5 \times 10^7 = 1.25 \times 10^6$ per second. The exponent is $400/20 = 20$, and $e^{20} \approx 4.85 \times 10^8$. So MTBF $\approx 4.85 \times 10^8 / 1.25 \times 10^6 \approx 388$ seconds: a failure every six and a half minutes, useless. Now allow nearly a full 2 ns clock period for resolution, $t_r = 1800$ ps: the exponent is 90 and MTBF is about $10^{33}$ seconds. The answer depends almost entirely on $t_r/\tau$; the frequencies matter only linearly.

## Adding stages and designing to a target

Adding a third synchronizer flop gives roughly one more destination clock period $T$ of resolution time, multiplying MTBF by $e^{T/\tau}$. At 1 GHz with $\tau = 20$ ps that factor is $e^{50}$, about $10^{21.7}$: almost 22 orders of magnitude. Conversely, the resolution time needed for a target MTBF is $t_r = \tau \ln(\text{MTBF} \cdot T_0 f_{clk} f_{data})$, and every extra factor of 10 in MTBF costs only $\tau \ln 10 \approx 2.3\tau$ of time. A chip contains many synchronizers, and the system failure rate is their sum, so a product with 1,000 synchronizers needs each to be 1,000 times better than the system target. High clock rates, slow corners (larger $\tau$) and logic between the synchronizer flops all shrink $t_r$ and must be included.

## Going deeper

Ginosar, "Metastability and synchronizers: a tutorial", IEEE Design and Test of Computers (2011); Kleeman and Cantoni, "Metastable behavior in digital systems", IEEE Design and Test of Computers (1987).
