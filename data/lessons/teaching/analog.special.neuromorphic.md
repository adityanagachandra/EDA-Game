# Neuromorphic and analog in-memory computing

Neuromorphic hardware computes with circuits inspired by neurons and synapses, and analog in-memory computing performs matrix-vector multiplication inside memory arrays using Ohm's and Kirchhoff's laws. Both trade digital exactness for energy efficiency, which turns analog non-idealities (conductance variation, IR drop, ADC resolution) into accuracy limits of the computation. They are an active area where circuit design, device physics and algorithms meet.

## Crossbar multiply-accumulate

In a resistive crossbar, each cell at row $i$ and column $j$ has conductance $G_{ij}$ encoding a weight. Driving row voltages $V_i$ while holding the columns at virtual ground makes each cell pass current $G_{ij}V_i$ by Ohm's law, and Kirchhoff's current law sums them on the column: $I_j = \sum_i G_{ij}V_i$. Rows at 0.1, 0.2 and 0.3 V through 50, 20 and 10 microsiemens give $5 + 4 + 3 = 12$ microamperes. A whole matrix-vector product happens in one step, with no data movement between memory and arithmetic units, which is the source of the energy advantage. The cells can be resistive RAM, phase-change memory, flash or SRAM-based capacitive structures.

## Signed weights and non-idealities

Conductances are positive, so signed weights use pairs of cells on two columns, $w = G^+ - G^-$, and the periphery subtracts the column currents. Real arrays depart from the ideal: device-to-device and cycle-to-cycle conductance variation acts as weight noise, conductance drift changes weights over time, wire resistance causes IR drop so that far cells see less voltage, and the column ADCs quantize the result. Each column needs an ADC, which often dominates energy and area, so the ADC resolution is chosen as low as accuracy allows. A design tool for these systems must therefore model the array as an analog circuit with statistics, much as mismatch is modelled with Pelgrom's law, and evaluate network accuracy under those models.

## Spiking neurons

Neuromorphic circuits often use spiking neurons. The leaky integrate-and-fire model integrates input current on a membrane capacitance $C$ with leak resistance $R$: $V(t) = IR(1 - e^{-t/\tau})$, $\tau = RC$. A spike fires when $V$ reaches the threshold $V_{th}$, after $t = \tau\ln(IR/(IR - V_{th}))$, and the membrane resets. With $R = 100$ megohms, $C = 10$ pF ($\tau = 1$ ms), $I = 10$ nA ($IR = 1$ V) and $V_{th} = 0.5$ V, the first spike comes after $\ln 2 = 0.693$ ms. If $IR$ is below threshold the neuron never fires. Information is carried in spike timing or rate, and circuits are event-driven, so simulating large networks of them is closer to digital event simulation than to SPICE.

## Going deeper

Sebastian et al., "Memory devices and applications for in-memory computing" (Nature Nanotechnology, 2020); Indiveri et al., "Neuromorphic silicon neuron circuits" (Frontiers in Neuroscience, 2011).
