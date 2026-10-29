# Three roles for a base model

Asking whether a model is good enough for agent-first EDA is the wrong question, because a model can play up to three different roles, with different requirements. It can drive the harness, it can run the continuous skill-selection loop, and it can become the on-site domain model. A closed API model can fill only the first. Matching models to roles is a design decision, made per role and sometimes per customer.

## Driving the harness

The first role is the agent that edits skill code and drives the edit, build, test and fix loop. It needs precise code editing (correct, minimal diffs against existing C++ source, not just fresh generation), reliable tool calling over long multi-turn loops without drifting, and a context window large enough for the relevant headers, the file being edited and the build log. Any strong tool-calling model, open-weights or closed, can fill this role, provided its prompts and the code it sees are allowed to reach wherever it runs. For a hosted API that means outside the site, which some deployments forbid.

## Running continuous selection

The second role runs constantly inside the fast modeling loop: reading metrics and their variance, choosing the next skill, deciding when to customize and when to recurse. It runs close to the data, often, and with low latency, so it is usually smaller than the harness model and fine-tuned on the selection records of the skill library. Because it must be fine-tuned and must read design-derived metrics continuously inside the customer's environment, it needs weights the deployment controls and can run locally.

## Becoming the on-site domain model

The third role is the domain model itself: a base model distilled and fine-tuned into the brain that knows the skill library, with a per-customer version trained on-site on that customer's records. This requires open weights, a license that permits commercial use, fine-tuning and shipping a fine-tuned derivative into a customer environment, and a size the training budget can afford. A closed model available only as a remote API cannot be fine-tuned and hosted on-site, so it is ruled out of roles two and three. It can still orchestrate in role one where the deployment allows its use.

## Going deeper

The agent-first EDA notes in this curriculum, section on models and learning; the license texts of the open-weights model families under consideration.
