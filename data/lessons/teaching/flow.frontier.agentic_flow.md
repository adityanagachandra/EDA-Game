# Agentic closure loops across the flow

Design closure has always been a loop: read reports, decide a fix, run tools, check again. Agentic EDA automates parts of that loop with models that read tool output, propose actions, call engines as tools and iterate, with humans approving the final change. Doing that well depends on what the agent reads, what it may change, how its success is measured, and how that measurement can be gamed.

## The closure loop as an agent

An agentic signoff or ECO loop reads timing, DRC, IR-drop and congestion reports; proposes fixes such as resizing cells, inserting buffers, adjusting placement constraints or rerouting nets; runs them through the engines; and checks the result against the same reports. Engines exposed as tools with typed inputs and outputs, for example over the Model Context Protocol, let the agent call them reliably. The agent should propose the most local fix first, as an engineer would, and escalate to earlier stages only when local fixes fail. Every action and its measured effect should be logged, so that the loop can be replayed, audited and used later as training data.

## Multi-agent orchestration over the flow

Large flows can be split among specialized agents: an RTL agent, a verification agent, an implementation agent, a signoff agent, coordinated by a planner and sharing state through the design database. The split mirrors the flow's stages, and its danger mirrors the flow's error propagation: an implementation agent that silently relaxes a constraint to close timing hands a broken design to signoff. Clear contracts between agents, the same contracts the file formats encode, and a shared record of every change keep the system auditable. Design-space exploration agents run many flow configurations in parallel and choose the next experiment.

## Rewards, reward hacking and the human gate

When agents learn from tool feedback, deterministic tools score their outputs, and that score is the reward. A reward of improved WNS from a fast timer can be gamed: the agent may relax an SDC constraint, add a false path, or exploit a region where the fast timer is wrong. Good rewards are computed by trusted engines, penalize constraint changes, check legality and signoff metrics, and use the accurate engine for final scoring. Even then, the final ECO is approved by a human, with a differential report that shows exactly what changed. Vendor claims of agentic gains are usually stated as best cases against unnamed baselines, so they deserve the same scrutiny.

## Going deeper

The agent-first EDA notes in this curriculum, sections on agents in the flow and on tool-grounded reinforcement learning.
