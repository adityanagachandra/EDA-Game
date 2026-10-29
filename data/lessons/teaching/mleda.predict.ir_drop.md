# IR-drop prediction

IR drop is the voltage lost across the resistance of the power grid as cells draw current. Computing it means solving a large resistive network, and dynamic analysis repeats that over many time steps and vectors. Learned IR-drop predictors map power and grid feature maps to a drop map, trained on solver results. The metric that matters is whether the worst drops and the hotspots are found.

## What the solver computes

The label comes from a power-grid solver, which applies Kirchhoff's laws to the grid. On a single strap fed from a pad, each segment carries the current of every tap beyond it. With segment resistance 0.1 ohm and taps drawing 10, 20 and 30 mA from the pad outward, the segments carry 60, 50 and 30 mA, and the drop at the last tap is $0.1(60 + 50 + 30) = 14$ mV (ohm times mA is mV). A real grid is a mesh solved as a sparse linear system $G v = i$ with millions of nodes. The predictor learns how current maps and grid structure turn into drops, the mapping the solver computes exactly.

## Errors in the worst drop

Signoff checks the worst drop against a budget, so the error in the predicted maximum matters more than the average error. If the true worst drop is 57 mV and the predicted worst is 68 mV, the error is $11/57 = 19.3$ percent. The two maxima need not occur in the same tile. A model can predict the right worst value in the wrong place, which passes this metric and still misleads a designer about where to add straps. Evaluations therefore report the maximum error, the location of the predicted hotspot and the per-tile error, rather than one aggregate number.

## Hotspot recall

Designers act on tiles above a drop threshold, so IR-drop predictors are also judged as hotspot detectors. Recall is the fraction of tiles truly above the threshold that the model also places above it. If three tiles exceed 45 mV in the solver result and the model puts two of them above 45 mV, recall is 67 percent. A missed hotspot is the costly error: no one re-simulates a region the model called safe. In a flow, the predictor is a fast mode for early power-grid planning and for screening many vectors, and the solver still produces the signoff result.

## Going deeper

Huang et al., "Machine Learning for Electronic Design Automation: A Survey," ACM TODAES (2021), section on IR-drop analysis; the CircuitNet dataset papers (Chai et al.), which include IR-drop labels.
