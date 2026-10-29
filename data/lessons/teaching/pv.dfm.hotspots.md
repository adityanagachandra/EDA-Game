# Lithography hotspot detection

Design rules are a simplified contract with lithography. Some rule-clean patterns still print badly: lines pinch, spaces bridge, line ends pull back, especially at the edges of the process window. Hotspot detection finds these patterns before tape-out, using pattern libraries, learned classifiers and lithography simulation. It is the clearest case in physical verification of a fast screen backed by an accurate engine.

## Precision and recall of a detector

A detector is judged against ground truth from lithography simulation. If it flags 56 true hotspots and 286 clean clips and misses 28 hotspots, its recall is $56 / 84 = 66.7$ percent and its precision $56 / 342 = 16.4$ percent. Recall measures how many real hotspots are caught; every miss is potential yield loss. Precision measures how many flags are real; every false alarm costs simulation time or an unnecessary layout change. Signoff flows require very high recall and buy precision with confirmatory simulation.

## What makes a hotspot

A hotspot is a pattern that fails to print within tolerance somewhere in the process window, the range of focus and exposure dose the fab expects. Typical failures are pinching (a line narrows), bridging (a space closes) and line-end shortening. Printing depends on all geometry within the optical interaction range, so the same pattern can be safe in one neighbourhood and a hotspot in another: a hotspot is a pattern together with its context. OPC corrects most proximity effects; hotspots are what remains difficult after correction.

## Detection methods

Full-chip lithography simulation is the most accurate and the slowest: it computes the aerial image and resist response everywhere. Pattern matching looks up known bad patterns exactly or with tolerances; it is fast and precise but blind to new patterns. Machine-learning classifiers generalize to unseen clips but make false alarms and misses, so their flags are confirmed by simulation. Fixes change geometry where there is room: wire spreading, widening, or rerouting with a hotspot cost. Pattern libraries from foundries are confidential, which shapes how detectors are deployed.

## Going deeper

The ICCAD 2012 CAD contest on fuzzy pattern matching for physical verification; the "Handbook of Algorithms for Physical Design Automation" (Alpert, Mehta and Sapatnekar, eds.), chapters on lithography and DFM.
