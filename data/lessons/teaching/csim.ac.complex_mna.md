# The complex MNA system

AC analysis solves the linearized circuit at each frequency in a sweep. In the phasor domain a capacitor's current is $j\omega C$ times its voltage, so MNA becomes a complex linear system $(G + j\omega C)x = b$, with $G$ and $C$ the same matrices transient analysis uses. One complex sparse factorization per frequency point is the entire cost of an AC sweep, which is why AC is so much cheaper than transient.

## Admittance entries

A capacitor stamps like a resistor whose conductance is the admittance $j\omega C$: $+j\omega C$ on its two diagonals and $-j\omega C$ off the diagonal. So the diagonal entry of a node with a resistor and a capacitor to ground is $Y_{11} = G + j\omega C$. Units pair conveniently: $2\pi$ times megahertz times picofarads is microsiemens. A 1 pF capacitor at 10 MHz contributes $\omega C = 2\pi \times 10^{7} \times 10^{-12} = 62.8$ uS to the imaginary part, while a 2 k$\Omega$ resistor contributes 500 uS to the real part. Inductors in branch rows contribute $-j\omega L$ on their diagonal, and controlled sources and conductances stay real.

## Solving one complex node

A current source $I$ into a node with $R$ and $C$ to ground gives $V = I/(G + j\omega C)$, so $|V| = I/\sqrt{G^2 + (\omega C)^2}$. With $I = 5$ uA, $R = 2$ k$\Omega$ ($G = 500$ uS) and $C = 5$ pF at 10 MHz ($\omega C = 314$ uS), $|Y| = 590.5$ uS and $|V| = 8.47$ mV. Microamperes over microsiemens give volts. A simulator performs the same division in complex arithmetic inside an LU factorization. Since AC analysis is linear, the source amplitude is arbitrary: SPICE users typically set it to 1 so that node voltages read directly as transfer functions.

## Phase of the response

The phase of $V = I/(G + j\omega C)$ is the source phase minus the phase of the admittance: $\angle V = -\arctan(\omega C/G)$. For $R = 2$ k$\Omega$, $C = 0.5$ pF at 100 MHz, $\omega C = 314$ uS against $G = 500$ uS, so the phase is -32.1 degrees. At the corner frequency $\omega = G/C$ the phase is -45 degrees and the magnitude is down by $\sqrt{2}$. AC output reports real and imaginary parts, from which magnitude in dB and phase in degrees are derived. Phase is where errors in device capacitances show up first, which is why AC results are a sensitive check of a model's charge equations.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter on frequency-domain analysis; Vlach and Singhal, "Computer Methods for Circuit Analysis and Design" (2nd ed., 1994).
