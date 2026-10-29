# Rule decks and the IP boundary

Physical verification touches the most sensitive data in the chain: the foundry's rule decks, device recognition and hotspot libraries, and the customer's complete layout. An agent working on DRC or LVS must operate inside a boundary that keeps all of it on site. The boundary is part of the architecture: what the agent may read, where it runs, what it may log, and what can flow back as shared know-how.

## Whose IP is whose

Rule decks, device models and pattern libraries are foundry IP, licensed under NDA to customers and tool vendors; they encode process capability. Some decks are delivered encrypted so that only the checking tool can read them. The layout and netlists are the customer's IP. An agent can work from rule names, marker locations and measured values in the results, plus the design-rule documentation available under the same agreement. Reading protected deck text by other means would breach the license.

## Architecture of the boundary

The agent drives the engines through typed tool interfaces running in a sandbox with no network access, on site. Logs sent anywhere are redacted of coordinates, cell and net names and geometry. Models fine-tuned on one customer's layouts may memorize them, so they must not serve another customer. Foundry pattern libraries cannot be bundled into skills shared across customers. These rules are enforced by the harness, not left to the model's discretion.

## Boundary decisions in practice

Concrete requests test the boundary. Sending failing clips to a cloud hotspot service is a transfer of design data, so the agent uses the on-site predictor or simulation instead. A fix strategy learned at one customer can be promoted to the shared library only as generic logic that passes golden and differential tests on reference designs, without that customer's data. Public benchmarks for agentic DRC use open PDKs such as SKY130 and GF180MCU or synthetic layouts. Quoting a confidential rule to a licensed user inside the deployment is acceptable; letting it reach a shared log or training set is not.
