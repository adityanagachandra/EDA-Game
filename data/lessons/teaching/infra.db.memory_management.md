# Handles, arenas and memory management

An EDA database creates and deletes millions of small objects while several engines hold references to them. How memory is allocated and how references are represented decide speed, footprint and, above all, whether a deleted object can crash another engine. The standard toolkit is object tables with stable handles, generation counters, arenas and pools for short-lived data, and compaction at safe points. This lesson covers each and the failure each one prevents.

## Handles, tables and arenas

Instead of handing out raw pointers, a database stores objects of one kind in a table and hands out handles: the index of the slot, often with a generation number. A handle is half the size of a 64-bit pointer when it is 32 bits, and stays valid when the table grows and moves. When an object is deleted its slot goes on a free list and is reused; the slot's generation counter increments, so a stale handle carrying the old generation is detected instead of silently reading the new occupant. Short-lived data follows a different pattern. A router's search state for one net or a parser's temporary tokens can come from an arena: allocation is a pointer increment, objects lie contiguously, and the whole arena is released at once when the pass ends. General heap allocation per object is slower and scatters objects across memory.

## Dangling references and leaks

Two failures dominate. A dangling reference arises when an engine keeps a raw pointer to an object that is later moved or freed. If a legalizer caches pointers into a vector of cells and the database then grows that vector, every cached pointer names the old storage, and reads become undefined behaviour that may pass tests and crash in production. Handles with validity checks, or storage that never moves objects, prevent it. A leak is memory never returned. Leaks look harmless per call but accumulate in long loops: 200 bytes per ECO iteration over 2,000,000 iterations is 400 MB. Agent-driven flows that run thousands of incremental edits in one session make both failures more likely, so long-running engines need leak tests and sanitizers in their regression suites, not only correct results on single runs.

## Compaction, side arrays and sharing

Deletion leaves holes. After deleting 30% of nets, a loop over the net table skips many empty slots, and density suffers. Compaction moves live objects together and remaps handles, which is only safe at a point where no engine holds old handles, or with a published old-to-new mapping. Engines often attach their own data through side arrays indexed by object number, such as a per-net criticality; when a net is created, each owner must grow its side array, usually in response to a creation callback. Sharing across threads is simplest when the database is frozen during a parallel phase: reading immutable data needs no locks, while reference counting every object makes threads contend on atomic counters. Finally, expect freed memory to stay in the process: allocators keep freed pages for reuse, so resident size may not fall when a design is closed.

## Going deeper

The OpenDB source in the OpenROAD repository (its dbTable and object id scheme); Richard Fabian, "Data-Oriented Design" (2018), the chapter on existence-based processing.
