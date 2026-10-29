# Model abstraction layers and swappability

Models improve every few months, and the best model for editing C++ may not be the best for skill selection or the right one for a particular customer's security posture. An agent-first harness is therefore built model-agnostic: one tool-call interface for the harness, and a thin adapter per model. Swapping or mixing models becomes configuration plus a benchmark run, not a rewrite.

## One tool-call interface

The harness talks to the model through a single internal interface: send a conversation of messages and tool results, receive either text or structured tool calls. The tools themselves (read a manifest, propose an edit, run a build, run tests, read a report) are defined once, with typed arguments, independent of any model. The harness's logic, its tier policy, its sandbox and its test gates never mention a model by name. This separation means a new model affects only how the conversation is encoded and decoded, and every safety property of the harness survives the swap.

## Thin adapters per model

Models differ in surface details that the adapter absorbs. Each has its own chat template (how system, user, assistant and tool messages are marked up), its own tool-call format (a JSON object, a tagged block, a function-call field), its own stop conditions, context limit and tokenizer. Self-hosted models are served by an inference runtime such as vLLM or SGLang, which may expose an API compatible with a common format; hosted models are reached through their vendor's API. The adapter maps the harness's messages and tools into the model's format and parses its output back, rejecting malformed tool calls with an error the harness can feed back.

## Mixing and swapping models

With adapters in place, mixing is a design goal, not a compromise. One model may make the C++ edits, a smaller fine-tuned model may run selection, and a different base may serve a customer with stricter requirements. A swap is decided empirically: rerun the skill-editing benchmark (compile-pass rate, skill-edit correctness, QoR deltas on real engines) with the candidate behind the same harness, and compare. A model that leads generic coding leaderboards may not lead on this task. A swap is also versioned and reversible like any other change, since a new model can alter the behavior of every agent that uses it.

## Going deeper

The vLLM and SGLang documentation on serving and tool-call parsing; the agent-first EDA notes in this curriculum, section on models.
