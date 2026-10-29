# Continuous on-site learning and the IP boundary

Learned EDA models degrade under shift, and the data that would fix them, a customer's designs and flow results, cannot leave the customer's site. Continuous on-site learning retrains models where the data lives, using the accurate engines' results as labels. The IP boundary shapes the whole design: what trains where, what may leave, and how a model that changes every night stays trustworthy.

## Why learn on site

A pretrained congestion or timing model meets a customer's node, library, flow and design style, all different from its training data. Every accurate run the customer's flow already performs, every global route and every signoff STA, is a grounded label for the fast model, free of extra tool cost. Retraining on them closes the distribution gap. The base model that becomes an on-site domain model must have open weights that can be trained and run locally; a model reachable only through a remote API cannot be fine-tuned inside the boundary and can only drive the harness.

## The IP boundary

Design data, and anything that can reveal it, stays on site. Model weights trained on customer designs are themselves sensitive, because models can memorize training examples. Anonymized layout-derived feature maps can still encode the floorplan and macro placement. The architecture enforces the boundary: training, inference and data stay in the customer's sandbox with no network path for design data, and only reviewed, design-independent artifacts, such as a validated skill or a parameter rule, or aggregate metrics, may leave. That is how the skill flywheel shares know-how across customers without sharing designs.

## Keeping a learning loop safe

A model retrained nightly needs the same discipline as code. Version every model; run golden regression checks against the accurate engine before promoting a new version; keep a rollback to the last good one. Watch for catastrophic forgetting: fine-tuning on a memory-heavy project can degrade logic-heavy blocks the model handled before, which regression tests on older blocks catch. Record the tool version and settings with every label and never mix incompatible ones. Before relying on the model for a new design, measure its error against the accurate engine on a sample of that design's runs.

## Going deeper

Kirkpatrick et al., "Overcoming Catastrophic Forgetting in Neural Networks," PNAS (2017); Carlini et al., "Extracting Training Data from Large Language Models" (USENIX Security 2021), on memorization.
