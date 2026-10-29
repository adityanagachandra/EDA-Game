# Database callbacks and observers

When engines share one design database, an engine that caches derived data, a timer's arrival times or a placer's density map, must learn when the design changes. The observer pattern does this: engines register callbacks with the database, and every edit notifies the subscribers. OpenDB provides this through callback objects attached to a block. The design questions are how many notifications an edit stream generates and how to keep their cost low.

## Observers subscribe to edit events

An observer registers for the kinds of edit it cares about. A timer cares about master swaps and pin connections, because both change delays and topology; a density map cares about moves and new instances; a global router cares about moves and connections. On each edit the database calls every observer subscribed to that kind. With a timer on {swap, connect}, a density map on {move, create} and a router on {move, connect}, the edit stream move, swap, connect, move costs $2 + 1 + 2 + 2 = 7$ callbacks. Subscribing only to relevant kinds matters: an engine that subscribes to everything pays for every edit in the flow, and a million-move detailed placement then makes millions of useless calls.

## Mark dirty, recompute on query

A callback should be cheap. Recomputing timing inside every callback is ruinous when an optimizer makes thousands of edits between questions. Instead the observer records what changed, often just a dirty flag or a set of dirty objects, and recomputes when someone asks. For the session swap, move, swap, QUERY, connect, QUERY, QUERY with a timer subscribed to swap and connect, an eager timer recomputes 3 times, once per relevant edit. A lazy one recomputes at the first query (dirty from the two swaps), at the second (dirty from the connect), and not at the third (still clean): 2 times. The edits between queries are merged into one update, which is how incremental timers stay fast inside optimization loops.

## Batching notifications per transaction

Even cheap callbacks cost a function call and often a lock per edit. A database can batch: edits inside a transaction are collected, and at commit each observer that subscribed to at least one of them gets a single call carrying the list. With three transactions, if the timer is touched in all three and the density map in two, the batched cost is 5 calls, however many edits each transaction held. Batching also gives observers a consistent view: they never see a half-applied change such as an instance moved but not yet reconnected. The trade-off is latency, since an observer learns of an edit only at commit, and memory for the pending edit list.

## Going deeper

The OpenDB documentation on dbBlockCallBackObj; Gamma, Helm, Johnson and Vlissides, "Design Patterns" (1994), the Observer pattern.
