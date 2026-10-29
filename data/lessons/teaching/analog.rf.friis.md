# The Friis cascade formula

A receiver is a cascade: filter, LNA, mixer, IF amplifiers, ADC. The Friis formula combines their noise factors and gains into the noise factor of the chain. It explains why the first stage dominates, why lossy components in front are so costly, and it is the formula every RF budget spreadsheet and system-level design tool evaluates.

## Two stages

Refer each stage's added noise to the chain's input. The first stage adds $(F_1 - 1)$ times the source noise. The second adds $(F_2 - 1)$ times the source noise at its own input, but its input already carries the signal amplified by $G_1$, so referred to the chain input its contribution is divided by $G_1$. The total is $F = F_1 + (F_2 - 1)/G_1$, with all quantities linear. For an LNA with NF 2 dB ($F_1 = 1.585$) and gain 15 dB ($G_1 = 31.6$) followed by a mixer with NF 10 dB ($F_2 = 10$), $F = 1.585 + 9/31.6 = 1.870$, so NF = 2.72 dB. The mixer's 10 dB costs only 0.72 dB because the LNA's gain hides it.

## Longer chains

The formula extends stage by stage: $F = F_1 + (F_2 - 1)/G_1 + (F_3 - 1)/(G_1G_2) + \dots$, each stage's excess noise divided by all the available gain in front of it. With a third stage of NF 15 dB ($F_3 = 31.6$) after a mixer with 5 dB conversion gain in the example above, its term is $30.6/(31.6 \cdot 3.16) = 0.306$, raising the total to $F = 2.176$ and NF = 3.38 dB. The general lesson: enough gain early makes later stages almost irrelevant to noise, but too much gain early hurts linearity, because later stages then see larger signals. Balancing the noise cascade against the linearity cascade is the central trade-off of receiver planning.

## Lossy stages in front

A passive matched component with loss $L$ (a filter, a switch, a cable) at 290 K has $F = L$ and gain $1/L$. Put it first and Friis gives $F = L + (F_2 - 1)L = L\,F_2$: in decibels the loss adds directly to the following noise figure. A 1.5 dB filter in front of a 1.5 dB LNA gives a 3 dB receiver. That is why antenna switches and preselect filters are designed for minimum insertion loss, and why some systems put the LNA before the filter when interference allows. The formula assumes each noise figure is specified for the source impedance the stage actually sees; noise figure depends on source impedance, so a cascade analysis with mismatched interfaces needs the full noise parameters, which RF simulators handle and spreadsheets do not.

## Going deeper

Friis, "Noise figures of radio receivers" (Proceedings of the IRE, 1944); Razavi, "RF Microelectronics" (2nd ed., 2011), chapter 2.
