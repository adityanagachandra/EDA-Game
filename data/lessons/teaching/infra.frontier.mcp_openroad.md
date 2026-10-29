# MCP tool wrappers around OpenROAD

The Model Context Protocol (MCP) lets an engine be exposed as a set of typed tools that any agent can discover and call. Wrapping OpenROAD this way turns its Tcl commands into a contract: named tools with input and output schemas, structured results, sessions that hold a design, and a sandbox around the process. Designing that wrapper well is a small exercise in API design with large consequences for every agent that uses it.

## Typed tools over an engine

An MCP server advertises tools, each with a name, a plain-language description and a JSON schema for its inputs and outputs. For OpenROAD, a tool such as report_timing might take a path group and a count and return structured fields: WNS, TNS and a list of worst paths with their slacks and endpoints, rather than the raw log text. Input schemas validate before the engine runs: a utilization declared as a number between 0.1 and 0.9 rejects 45 (a percent typed as a fraction) with a clear error instead of an absurd floorplan. The description states preconditions, such as "requires a placed design", because the agent reads descriptions to decide which tool to call and in which order, exactly as it reads skill manifests.

## Choosing tool granularity

Granularity decides how usable the server is. A single eval_tcl tool that runs any string is powerful but has no typed contract, no validation, unstructured output and no way to audit or restrict what runs. The opposite extreme, one tool per database accessor, floods the agent with round trips. Coarse actions work best: initialize the floorplan, run global placement with these parameters, repair timing, report timing. Per-object operations should be batched, so moving 10,000 cells is one call with a list, not 10,000 calls, the same reason a plug-in is called per batch rather than per cell. Long operations such as detailed routing, which can take an hour, should start a job and return an id, with a separate tool to poll or cancel it.

## Sessions, sandboxing and structured results

Loading a large ODB is expensive, so the server keeps a session holding the current design between calls, and each tool acts on that state. When the design is customer IP, the OpenROAD process and its data stay on-site, in a sandbox with resource limits and no outbound network; only structured results cross to the agent. Tools should be idempotent where possible: calling set_place_density twice with the same value leaves the same state, so an agent that retries after a timeout does no harm. Every call is logged with its arguments, the engine version and a checksum of the resulting database, so a session can be audited, replayed from scratch and reverted, which keeps agent actions reviewable in the same way a code change is.

## Going deeper

The Model Context Protocol specification; the OpenROAD Tcl and Python API documentation; the agent-first notes in this curriculum on engines as agent tools.
