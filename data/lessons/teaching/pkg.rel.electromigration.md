# Electromigration and Black's equation

Electromigration is the slow transport of metal atoms by the momentum of conducting electrons. Over years it opens voids that break wires or raise their resistance, and pushes metal into hillocks that short to neighbours. Black's equation relates the median time to failure to current density and temperature, and it is the basis of the current-density rules that power grids, bumps and TSVs are signed off against.

## Current density

Electromigration depends on current density $J = I/A$, not current. A copper line 1 um wide and 0.5 um thick carrying 10 mA has $J = 0.01/(0.5 \times 10^{-8}$ cm^2$) = 2 \times 10^6$ A/cm^2, or 2 MA/cm^2. Units matter: one square micrometre is $10^{-8}$ cm^2. In packages the cross-sections are large but so are currents: a 20 um microbump at 0.2 A carries about 64000 A/cm^2, and solder fails at much lower current density than copper because tin diffuses fast. Current crowding where a narrow line enters a wide bump raises the local density several times above the average, which is why bump EM is analysed with detailed models.

## Black's equation

Black's equation states $MTTF = A J^{-n} e^{E_a/(kT)}$, where $n$ is the current exponent, $E_a$ the activation energy and $k$ Boltzmann's constant. At fixed temperature, changing current density from $J_1$ to $J_2$ scales life by $(J_1/J_2)^n$. With $n = 2$, a design that lasts 15 years at 1.5 MA/cm^2 lasts $15 \times (1.5/2)^2 = 8.44$ years at 2 MA/cm^2. The exponent is close to 2 when void nucleation limits life and close to 1 when void growth does, so the choice matters when extrapolating. EM signoff turns a target lifetime and temperature into a maximum allowed $J$ per wire width and layer, which the power-grid tool checks for every segment.

## Allowed current for a lifetime

Inverting at fixed temperature gives the current density allowed for a target life: $J_2 = J_1 (MTTF_1/MTTF_2)^{1/n}$. If a line reaches 5 years at 3 MA/cm^2 and must last 10 years with $n = 2$, the allowed density is $3/\sqrt{2} = 2.121$ MA/cm^2. Multiplying by the cross-section gives the allowed current, and so the wire width a power route needs. EM limits are statistical: Black's equation gives a median, and failure distributions are lognormal, so signoff uses a failure fraction far below 50%, which lowers the allowed density further. The limits differ for DC power current and for bidirectional signal current, which partly heals.

## Going deeper

Black, "Electromigration: a brief survey and some recent results", IEEE Transactions on Electron Devices (1969); Lienig and Thiele, "Fundamentals of Electromigration-Aware Integrated Circuit Design" (2018).
