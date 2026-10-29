# Scattering parameters

Scattering parameters describe a network by how it reflects and transmits waves at its ports, measured against a reference impedance. They are what vector network analyzers measure, what package and connector vendors deliver, and what channel simulators consume. Unlike impedance or admittance parameters, they stay well defined at high frequency, where opens and shorts are hard to realize.

## Reflection at a port

At each port $k$ define incident and outgoing waves $a_k$ and $b_k$, normalized to a reference impedance $Z_{ref}$, usually 50 ohm. The S-matrix relates them, $\mathbf{b} = S\mathbf{a}$. With all other ports terminated in $Z_{ref}$, $S_{11} = b_1/a_1$ is the reflection coefficient looking into port 1:
$$S_{11} = \frac{Z_{in} - Z_{ref}}{Z_{in} + Z_{ref}}.$$

Example: a port presenting 100 ohm in a 50 ohm system has $|S_{11}| = 1/3$, or $20\log_{10}(1/3) = -9.54$ dB, a return loss of 9.54 dB. A port presenting 60 ohm gives -20.8 dB. S-parameters depend on the reference impedance: the same device has different S-parameters against 50 and 100 ohm, so a tool must always carry $Z_{ref}$ with the data.

## Insertion loss in decibels

$S_{21} = b_2/a_1$ with port 2 matched is the transmission from port 1 to port 2. Because S-parameters are ratios of wave amplitudes, $|S_{21}|$ in dB is $20\log_{10}|S_{21}|$, and since wave power is proportional to $|a|^2$, $|S_{21}|^2$ is the ratio of power delivered to a matched load to incident power. So $|S_{21}|_{dB} = 10\log_{10}(P_{out}/P_{in})$ as well.

Example: $|S_{21}| = 0.8$ is -1.94 dB; delivering half the incident power is -3.01 dB. Insertion loss is usually quoted as the positive number, 3.01 dB. Channel specifications set a maximum insertion loss at the Nyquist frequency because equalizers can only recover so much; a channel tool sweeps $|S_{21}|$ and compares it with that mask.

## Power balance, passivity and reciprocity

With port 2 matched, incident power at port 1 is reflected ($|S_{11}|^2$), transmitted ($|S_{21}|^2$) or absorbed. So the absorbed fraction is $1 - |S_{11}|^2 - |S_{21}|^2$. A passive network can only absorb, so this must never be negative; a lossless network has a unitary S-matrix, giving $|S_{11}|^2 + |S_{21}|^2 = 1$. A reciprocal network, one without magnetized ferrites or active devices, has $S_{21} = S_{12}$.

Example: $|S_{11}| = 0.2$ and $|S_{21}| = 0.85$ give an absorbed fraction of $1 - 0.04 - 0.7225 = 0.2375$, so 23.75% of the power is lost in the network. A lossless two-port with $|S_{11}| = -20$ dB (0.1) has $|S_{21}| = \sqrt{0.99}$, -0.044 dB. Extracted S-parameter models that violate passivity or reciprocity cause unstable or wrong time-domain simulation, so tools check and enforce these properties.

## Going deeper

Pozar, "Microwave Engineering" (4th ed.), chapter 4; Hall and Heck, "Advanced Signal Integrity for High-Speed Digital Designs", chapter 9.
