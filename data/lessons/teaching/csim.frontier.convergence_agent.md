# An agent that debugs convergence

Convergence debugging is the simulator task most likely to be handed to an agent: it is frequent, tedious and full of patterns. An agent reads the simulator's log, classifies the failure, proposes the smallest change, reruns the deck and repeats in a harness loop. The hard part is not making the deck run, which loosening tolerances always achieves, but proving the fix did not change the answer.

## Reading the failure

A useful agent starts from the simulator's own diagnosis. A singular matrix naming a node points to topology: the node has no DC path, touching only capacitors or gates. A DC iteration limit with oscillating iterates points to nonlinearity or multiple solutions. Timestep-too-small at 3.2 ns points to a discontinuity, a model kink or a stiff event near that time, and the agent must extract which nodes were changing fastest and which source edge occurs there. Partial convergence, such as gmin stepping that stopped at $10^{-9}$ S, is a failure, not a result. Classification determines which knobs are relevant and keeps the agent from trying every option blindly.

## Proposing the smallest fix

The agent edits at the smallest safe tier, preferring changes to the path over changes to the equations. A .nodeset steers Newton without altering the circuit; a realistic rise time on an ideal PWL corner removes an unphysical discontinuity; a targeted initial condition kicks an oscillator. Raising reltol or setting a large global gmin also makes decks run, but they change accuracy or the circuit itself, so the harness should reject them unless the request explicitly trades accuracy. Editing a device model's equations is a core change and never a convergence fix. Each attempt's new failure message goes back to the model as feedback for the next proposal.

## Proving the fix did not change the answer

A converged run is evidence only that the simulator finished. The agent must compare key results, such as operating-point voltages, delays or gains, against a trusted reference or golden values within tolerance. To show a glitch was numerical, rerun at tighter tolerances and smaller steps: numerical artifacts change, physical behavior converges. When a circuit has several valid DC solutions, a nodeset that selects one must be reported as a choice, not hidden. A reward of 'the run finished' is easy to game by loosening tolerances or adding gmin until anything runs, so a tool-grounded reward must include agreement with golden results and penalize changes that alter the circuit.

## Going deeper

The ngspice manual, sections on convergence and options; Najm, "Circuit Simulation" (Wiley, 2010), chapters 4 and 5.
