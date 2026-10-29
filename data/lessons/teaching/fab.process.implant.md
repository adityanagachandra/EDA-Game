# Ion implantation profiles

Ion implantation shoots ionized dopant atoms into the wafer at a chosen energy and dose. It sets wells, threshold adjustments, halos, extensions and source-drain regions, each through its own mask layer. To first order the implanted atoms form a Gaussian in depth. That one model gives peak concentration, the doping at any depth and the junction depth, which are the numbers device and process simulators start from.

## Dose, range and the Gaussian profile

An implant is specified by species, energy and dose $Q$, the number of ions per square centimetre. Energy sets the projected range $R_p$, the mean depth, and the straggle $\Delta R_p$, the standard deviation of depth. The profile is $N(x) = N_p \exp\left(-(x - R_p)^2 / 2\Delta R_p^2\right)$. Integrating it over depth must return the dose, so the peak is $N_p = Q/(\sqrt{2\pi}\,\Delta R_p)$. For $Q = 1 \times 10^{15}$ cm^-2 and $\Delta R_p = 20$ nm $= 2 \times 10^{-6}$ cm, $N_p = 10^{15}/(2.507 \times 2 \times 10^{-6}) = 1.99 \times 10^{20}$ cm^-3.

## Concentration at a given depth

Once $N_p$ is known, the concentration at depth $x$ follows from the distance to the peak in units of straggle. At one straggle from the peak the factor is $e^{-1/2} = 0.607$; at two straggles, $e^{-2} = 0.135$; at three, $e^{-4.5} = 0.011$. Take $N_p = 2 \times 10^{20}$ cm^-3, $R_p = 50$ nm and $\Delta R_p = 20$ nm. At $x = 80$ nm the depth is 1.5 straggles from the peak, so $N = 2 \times 10^{20} \times e^{-1.125} = 6.49 \times 10^{19}$ cm^-3. Real profiles are skewed and have channelling tails, so production simulators use Pearson distributions or Monte Carlo, but the Gaussian is the shape to reason with.

## Junction depth

The metallurgical junction sits where the implanted concentration equals the background doping $N_B$ of the opposite type. Solving $N(x_j) = N_B$ beyond the peak gives $x_j = R_p + \Delta R_p \sqrt{2 \ln(N_p/N_B)}$. With $N_p = 2 \times 10^{20}$, $N_B = 1 \times 10^{17}$ cm^-3, $R_p = 50$ nm and $\Delta R_p = 20$ nm: $\ln(2000) = 7.60$, $\sqrt{15.2} = 3.90$, so $x_j = 50 + 78 = 128$ nm. The logarithm makes junction depth insensitive to dose but linear in straggle, which is why shallow junctions need low energies, heavy ions and minimal later diffusion.

## Going deeper

Plummer, Deal and Griffin, "Silicon VLSI Technology", chapter 8; Ziegler's SRIM documentation for range tables.
