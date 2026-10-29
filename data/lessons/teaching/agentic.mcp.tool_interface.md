# Engines as agent tools over MCP

Once engines are callable, they need a standard way to be offered to agents, including agents built by other teams or other companies. The Model Context Protocol (MCP) is an open protocol for this: a server exposes tools with names, descriptions and typed input schemas, and an agent's host calls them over JSON-RPC. Exposing an EDA engine as an MCP server makes it usable by any compliant agent, which raises questions of interface design and security that the engine designer must answer.

## Tools with typed inputs and outputs

An MCP server advertises a list of tools. Each has a name, a plain-language description and an input schema written in JSON Schema, so the caller knows exactly which arguments exist, their types, their ranges and which are required. The host lists the tools to the model, the model asks to call one, and the host sends the call to the server and returns the result. For a timer, a tool `report_worst_paths` might take a design handle, a corner name, a path group and a count from 1 to 1000, and return structured records of endpoint, slack and path. Servers can also expose resources (readable data such as a report) and prompts, but tools are the core of engine integration.

## Designing an engine's tool surface

Good EDA tools are coarse-grained and stateless in their arguments. They take a handle to a design state in the design database rather than the design itself, since a netlist of millions of cells does not belong in a message. Long operations such as global routing are started as jobs that return a job identifier, with separate tools to poll status and fetch results, so a call never blocks for hours. Results are summaries with handles to the full data, such as the worst 20 paths and a handle to the full report. Calls are idempotent where possible and record provenance. A tool per engine operation, not per internal function, keeps the surface small enough for a model to use well.

## Security at the tool boundary

The server, not the model, enforces what is allowed. Every argument is validated against the schema and against policy: the design handle must belong to this session's workspace, file paths must lie inside allowed directories, and tier rules apply, so a tool that changes engine code goes through the harness rather than writing files directly. Tool descriptions from third-party servers are themselves untrusted text an agent reads, so a host should allowlist servers and tools per task. The engine still runs in its sandbox, on-site, with the same limits whether called by a human script or by an agent.

## Going deeper

The Model Context Protocol specification, on tools and their input schemas; the JSON Schema specification.
