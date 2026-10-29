# Incremental STA inside a tool

Inside an optimizer, a static timer is not a batch program that reads files and prints a report. It is a service living in the same process as the database, answering slack queries while the resizer, placer and router edit the netlist under it. Turning the timing-graph mathematics of incremental update into working software takes three pieces of infrastructure: a way to hear about edits, a way to record and refresh what they invalidate, and an understanding of which edits and queries are cheap.

## Edits reach the timer through callbacks

Tools that share one database learn about each other's edits through observers. The timer registers callbacks with the database; when an engine creates or deletes an instance, swaps its master, connects or disconnects a pin, or moves a cell, the database calls the timer before or after the change. The callback does not recompute timing. It updates the timer's own graph to match the database, adding or removing vertices and arcs, and it marks what is now stale: the delays of the affected arcs, the arrivals downstream and the required times upstream.

Callbacks must be fast, because a placement pass may produce tens of thousands of them. Recording "these pins are invalid" is constant work per edit; recomputing on each callback would multiply the cost of every edit by the size of a cone. Edits that change no electrical fact, such as renaming a net, should not invalidate anything.

## Invalidate, then propagate in level order

Invalidation splits into two directions. An arc change at $u \to v$ can change arrivals at $v$ and in its fanout cone, and required times at $u$ and in its fanin cone. A timer keeps separate invalid sets for the two, plus a set of arcs whose delays must be recalculated.

When a query arrives, the timer recomputes delays for invalid arcs, then processes invalid arrivals in increasing level order, so that every fanin is final before a pin's max is taken; required times are processed in decreasing level order. Two refinements matter. First, demand: a query for one endpoint's slack refreshes only its fanin cone and its own required time. Second, early termination: if a pin's new arrival equals its old one, its fanout need not be visited. OpenSTA implements this levelized lazy scheme; OpenTimer instead builds a task graph for the invalid region and runs it on a parallel task scheduler.

## What makes an incremental update expensive

Incrementality pays only when the changed region is small and the questions are local. Three patterns break it. A constraint edit, such as a new clock period or a changed input delay, alters required times at many endpoints, so the backward update can cover most of the design. A global query after every edit, such as asking for worst negative slack following each buffer insertion, forces a near-global refresh each time, so 10,000 edits cost 10,000 global updates. And structural edits that lengthen paths, such as inserting a buffer, can raise the levels of downstream pins, so level bookkeeping must itself be updated incrementally.

Good optimizer code batches edits, asks local questions (the slack of the endpoints it touched), and saves global summaries for checkpoints. A useful mental model: the cost of an update is the size of the region a query forces valid, not the size of the edit.

## Going deeper

The OpenSTA source and its documentation on network edits and incremental timing; Huang and Wong, "OpenTimer: A High-Performance Timing Analysis Tool" (ICCAD 2015) and the later OpenTimer 2 paper on task-based incremental timing.
