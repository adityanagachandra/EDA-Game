# IBIS-AMI and SerDes channel simulation

Multi-gigabit serial links depend on equalization: transmit FFE, receive CTLE and DFE, and clock-data recovery. These are algorithms, often adaptive, that IBIS tables cannot describe. IBIS-AMI (Algorithmic Modeling Interface) extends IBIS with executable models that a channel simulator calls through a standard API. Combined with S-parameter channel models, they let a board designer predict eye opening and bit error rate before building the board.

## The AMI model structure

An IBIS-AMI model has two parts: the usual IBIS buffer for the analog front end, and an **algorithmic** part delivered as a compiled shared library plus an .ami parameter file that declares the model's settings (tap values, CTLE curves, adaptation modes). The simulator loads the library and calls three entry points: **AMI_Init**, which receives the channel's impulse response and may return an equalized one; **AMI_GetWave**, which processes blocks of waveform samples in the time domain; and **AMI_Close**, which frees the model's state. Because the interface is standard, any compliant simulator can run any vendor's model, and the vendor's equalization algorithms stay compiled and private. It is an engine seam in exactly the four-layer sense.

## Transmit FFE as a filter

A transmit feed-forward equalizer is a short FIR filter on the symbol stream: $y[n] = c_{-1} x[n+1] + c_0 x[n] + c_1 x[n-1]$, with a pre-cursor tap $c_{-1}$, a main tap $c_0$ and a post-cursor tap $c_1$, usually normalized so $|c_{-1}| + c_0 + |c_1| = 1$. Example with $c_{-1} = -0.1$, $c_0 = 0.7$, $c_1 = -0.2$: a transition, $x[n-1] = -1$, $x[n] = +1$, $x[n+1] = +1$, gives $-0.1 + 0.7 + 0.2 = 0.8$, while a repeated $+1$ after $+1$ gives $-0.1 + 0.7 - 0.2 = 0.4$. Transitions are emphasized relative to runs, a high-pass shaping that pre-compensates the channel's low-pass loss. Tap values are typical AMI parameters swept in channel simulation.

## Statistical and time-domain flows

A channel simulation cascades transmitter, Touchstone channel and receiver. In the **statistical** flow, if the equalization is linear and time-invariant, AMI_Init transforms the channel impulse response, and the simulator builds the eye and BER contours analytically from it, fast enough to sweep many settings. When the equalizer adapts or acts nonlinearly, as a **DFE** does by subtracting weighted past decisions to cancel post-cursor ISI, or as clock recovery does, the **time-domain** flow runs millions of bits through AMI_GetWave. Results can differ between simulators even with identical models, because jitter injection, sampling and BER extrapolation are each tool's own. A DFE cannot cancel pre-cursor ISI, since it only sees past bits; that is left to FFE and CTLE.

## Going deeper

The IBIS specification, chapters on the Algorithmic Modeling Interface. Hall and Heck, "Advanced Signal Integrity for High-Speed Digital Designs" (2009), chapter on equalization.
