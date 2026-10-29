# Dopant diffusion and thermal budget

Dopants placed by implantation do not stay put. Every hot step afterwards lets them spread by diffusion, which moves junctions, widens doped regions and shortens effective channel length. The process engineer tracks a thermal budget: the accumulated product of diffusivity and time. A tool designer meets this through device models, short-channel behaviour and through why advanced flows use fast spike and laser anneals rather than long furnace bakes.

## Fick's law and the diffusion length

Dopant flux follows Fick's first law, $J = -D\, \partial C/\partial x$, and mass conservation gives the diffusion equation $\partial C/\partial t = D\, \partial^2 C/\partial x^2$, the same equation as heat flow. Its solutions spread over a characteristic length $2\sqrt{Dt}$. A fixed dose diffusing from a thin surface layer becomes a Gaussian with standard deviation $\sqrt{2Dt}$; a constant surface concentration gives a complementary error function profile. With $D = 1 \times 10^{-14}$ cm^2/s and a 60-minute anneal, $Dt = 3.6 \times 10^{-11}$ cm^2, $\sqrt{Dt} = 6 \times 10^{-6}$ cm, and $2\sqrt{Dt} = 120$ nm.

## Arrhenius temperature dependence

Diffusivity depends steeply on temperature: $D = D_0 e^{-E_a/kT}$, with activation energies of roughly 3 to 4 eV for common dopants in silicon and $k = 8.617 \times 10^{-5}$ eV/K. The ratio between two temperatures is $D_2/D_1 = \exp[(E_a/k)(1/T_1 - 1/T_2)]$. For $E_a = 3.5$ eV, going from 900 C (1173 K) to 1000 C (1273 K) gives $1/T_1 - 1/T_2 = 6.70 \times 10^{-5}$ per K, times $3.5/8.617 \times 10^{-5}$ is 2.72, and the ratio is $e^{2.72} \approx 15$. A hundred degrees changes diffusion by more than an order of magnitude.

## Adding up the thermal budget

Successive anneals add their $Dt$ products when the dopant and its diffusion mechanism are the same: total $Dt = \sum_i D_i t_i$, and the final spread is $2\sqrt{\sum D_i t_i}$. A 10 s spike at $5 \times 10^{-14}$ cm^2/s contributes $5 \times 10^{-13}$ cm^2; a 300 s anneal at $2 \times 10^{-15}$ cm^2/s contributes $6 \times 10^{-13}$ cm^2; together $2\sqrt{1.1 \times 10^{-12}} = 21$ nm. This is why later steps run cold: any heat after the source-drain implant moves junctions that set channel length. Real dopants also show transient enhanced diffusion after implant damage, which a simple constant-$D$ budget understates.

## Going deeper

Plummer, Deal and Griffin, "Silicon VLSI Technology", chapter 7; Crank, "The Mathematics of Diffusion".
