# Fast versus accurate litho simulation

Computational lithography has a natural pair of engines. Fast ones, such as calibrated compact optical and resist models, pattern matching and machine-learned hotspot detectors, run at full-chip scale. Accurate ones, such as rigorous electromagnetic mask simulation and physical resist models, run on small clips. An agent-driven OPC or DTCO loop uses the fast engines to search and the accurate ones to decide, and the craft is knowing when to recurse.

## Two classes of litho engine

A compact model approximates the aerial image with a few precomputed optical kernels and a calibrated resist threshold or blur, so it can evaluate billions of OPC fragments. Pattern matching and ML hotspot detectors are faster still: they classify clips as risky from geometry alone. Rigorous simulation solves Maxwell's equations through the three-dimensional mask, including EUV shadowing, and models resist chemistry, at a cost that limits it to small areas. The fast engines are calibrated against the accurate ones and against wafer data, and are trustworthy only inside the range of patterns they were calibrated on.

## When to recurse to the accurate engine

The agent should spend accurate-engine time where the fast answer is least reliable or matters most. Clear triggers: a pattern outside the fast model's calibration domain, such as a new pitch from a rule change; a decision between candidates whose fast-score difference is smaller than the fast model's known error; a pattern class where audits show the two engines disagree, for example line ends; and every candidate about to be accepted. The variance between fast and accurate results on a fixed audit sample is the signal that a surrogate is drifting and needs recalibration.

## Proxy gaming and honest reporting

If a learner is rewarded by a fast litho score, it will find the model's blind spots, for example layouts just outside the calibrated range that score perfectly but print badly. Guards include penalties for leaving the calibrated domain, accurate-engine audits of the top candidates, and acceptance based only on rigorous verification and, eventually, wafer data. An ML hotspot detector trained on past layouts will miss new hotspot types introduced by a rule change, so sampling for rigorous checks must increase after any change. Reports should state which engine produced each number and the audit agreement for that pattern class.

## Going deeper

Ma and Arce, "Computational Lithography" (2010); Yang et al., "Layout hotspot detection with feature tensor generation and deep biased learning", DAC (2017).
