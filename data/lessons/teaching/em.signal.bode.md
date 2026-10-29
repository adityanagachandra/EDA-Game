# Bode magnitude and phase

A Bode plot shows a transfer function's magnitude in decibels and its phase against log frequency. Because poles and zeros contribute straight-line asymptotes, a designer can sketch a response in seconds and read bandwidth, loop gain and phase margin off it. Frequency-response engines report these numbers, and model-fitting tools are judged by how well they reproduce them.

## Magnitude in decibels

Magnitude in decibels is $20\log_{10}|H|$ for voltage or current ratios (and $10\log_{10}$ for power ratios). Products of factors become sums in dB, which is why Bode plots add contributions. A single real pole $K/(1 + jf/f_p)$ has
$$|H|_{dB} = 20\log_{10}K - 10\log_{10}\left(1 + (f/f_p)^2\right).$$

Example: $K = 100$ and $f_p = 10$ MHz. The DC gain is 40 dB. At 10 MHz the exact gain is $40 - 10\log_{10}2 = 36.99$ dB, the familiar 3 dB point. At 30 MHz it is $40 - 10\log_{10}10 = 30$ dB. At 1 MHz it is $40 - 0.04 = 39.96$ dB. Each factor of 10 in voltage is 20 dB; a factor of 2 is about 6 dB. Power ratios use $10\log_{10}$, so a 3 dB power drop is the same point as a 3 dB voltage drop.

## Asymptotes of minus 20 dB per decade per pole

Far below a pole its factor is about 1 (0 dB); far above it falls as $f_p/f$, a straight line of -20 dB per decade. The asymptotes meet at $f_p$, where the true curve is 3 dB lower. For several real poles, above all of them the asymptotic magnitude is
$$|H|_{dB} \approx K_{dB} - \sum_k 20\log_{10}(f/f_{p,k}).$$
A zero contributes $+20$ dB per decade above its frequency.

Example: 80 dB DC gain with poles at 1 MHz and 100 MHz. At 1 GHz the first pole has taken $20\log_{10}1000 = 60$ dB and the second $20\log_{10}10 = 20$ dB, leaving 0 dB. Between the poles the slope is -20 dB per decade; beyond both it is -40. The asymptotic error is largest near each corner and at most 3 dB per isolated pole.

## Phase of poles and zeros

A real left-half-plane pole contributes phase $-\arctan(f/f_p)$: 0 well below it, $-45^\circ$ at the pole and $-90^\circ$ well above. A left-half-plane zero contributes $+\arctan(f/f_z)$. The phase of a product is the sum of the phases, so
$$\angle H = -\sum_k \arctan(f/f_{p,k}) + \sum_m \arctan(f/f_{z,m}).$$

Example: poles at 1 MHz and 10 MHz, evaluated at 3 MHz: $-\arctan 3 - \arctan 0.3 = -71.6^\circ - 16.7^\circ = -88.3^\circ$. Phase changes over about two decades around each pole, so it starts moving well before the magnitude does. That spread is what limits feedback: by the time a loop's gain has fallen to 1, the extra poles have already eaten into the phase.

## Going deeper

Sedra and Smith, "Microelectronic Circuits" (8th ed.), appendix F on Bode plots; Franklin, Powell and Emami-Naeini, "Feedback Control of Dynamic Systems", chapter 6.
