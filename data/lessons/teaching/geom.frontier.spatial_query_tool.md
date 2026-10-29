# Spatial queries as agent tools

An agent fixing violations or exploring a layout cannot read a billion shapes; it asks questions. Exposing a geometry engine's spatial queries as typed tools, over an interface such as MCP, lets a model ask what lies in a window, how far two nets are apart or where violations cluster, and get small, exact answers. The tool boundary is where units, bounds, determinism, cost and the IP boundary must all be designed explicitly.

## Typed spatial queries

A good spatial tool takes typed, unit-explicit inputs and returns structured records. A region query names the layer, the window corners and the database unit; a distance query names two nets or shapes; a violation query names a rule and a region. Free-text queries push parsing onto the model and invite misunderstandings, and whole-layer dumps overflow its context. Outputs should be records an agent can compute with: for a spacing violation, the two shapes, the measured and required distances and the marker box, so the agent can work out how far an edge must move and the harness can check the fix. Internally, each tool maps to an index query: a window query on an R-tree or grid, grown by a margin when the question is about distance.

## Bounded, deterministic results

A window over a dense region can match millions of shapes. The tool should return a bounded page, the total count, and a continuation token, so the agent knows its view is partial and can refine the window or continue. Results must be deterministic: sorted by a stable key such as layer, then coordinates, then id, so that two identical queries give identical answers and harness diffs stay meaningful. Caches must be invalidated when edits touch their region, or the agent will chase violations that no longer exist. Example: a query for metal-2 shapes in a $10 \times 10$ um window returns the first 500 of 12,400 records sorted by lower-left corner, with a token for the next page.

## Cost and safety at the tool boundary

Geometry queries range from microseconds (a small window) to hours (a full-chip XOR). Tools should be region-scoped and report a cost estimate, so a planner can check only the edited area plus its interaction halo instead of rerunning everything after every change. The engine runs next to the data, inside the design site's sandbox, so the layout never leaves the site and only queries and small results cross the boundary. Editing tools need more: an edit should go through a transaction that runs the relevant checks and a differential comparison and produces a reviewable diff before commit. A fix is complete only when the edited region and its halo show the intended change and nothing new.

## Going deeper

The Model Context Protocol specification, on tool schemas and structured results; Samet, "Foundations of Multidimensional and Metric Data Structures" (2006), on window and distance queries.
