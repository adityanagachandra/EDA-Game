# The IP boundary around process data

Foundry process data is among the most protected information in the industry: device models, rule decks, process recipes and yield data are covered by non-disclosure agreements. An agentic flow that reads PDKs must keep that data where it belongs. The boundary is enforced by architecture: where models run, what tools can read, what the network allows, and what may flow back to a shared library.

## What is confidential and why it is architecture

Confidential process data includes SPICE model parameters, rule decks and their values, extraction stacks, process recipes, and measured yield and defect data. General physics, textbook models and open-source algorithms are not. Because agents copy text freely, the boundary cannot rest on instructions alone. It is upheld by running models on-site or in the customer's environment, by executing skills in sandboxes with no outbound network, and by access control on PDK files. Logs inherit the confidentiality of what they contain: an agent log that quotes rule values can only be stored inside the same boundary.

## What may cross the boundary

The skill flywheel moves technique, not data. A DFM skill improved at one customer can be generalized and validated on non-confidential test cases and promoted to the shared library, without the customer's layouts, PDK or yield data. A vendor can improve a model serving many foundry customers by training centrally on public or synthetic data and adapting locally inside each boundary. To debug an engine with an outside team, a foundry can share a minimal synthetic clip and a rule rewritten with dummy values that reproduce the bug. Fine-tuned on-site model weights may memorize process values, so they inherit the same protection.

## Narrow interfaces instead of shared text

The strongest pattern is to give an agent a narrow query interface instead of the data: a tool that answers "is this layout clean under the deck" or "what is the spacing for these widths" without exposing the full deck. The agent gets feedback for optimization while the rule text never enters its context. Simple obfuscation, such as renaming layers or scaling all dimensions, does not protect the process, because ratios, rule structure and patterns still reveal it. Enforcement belongs in the sandbox and network policy, independent of the agent's instructions, because instructions can be ignored or manipulated.

## Going deeper

The agent-first EDA anchor on sandboxing and the IP boundary; the open SKY130 and GF180MCU PDKs, as a contrast in what open process data allows.
