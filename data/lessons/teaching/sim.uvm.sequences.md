# Sequences, sequencers and drivers

UVM separates what traffic to generate from how to put it on the pins. Sequences describe scenarios as streams of randomized transactions; drivers implement the protocol. The two meet through a sequencer and a precise handshake. Understanding that handshake, arbitration among sequences, and virtual sequences explains how one environment supports hundreds of tests without changing a line of driver code.

## The sequence, sequencer and driver handshake

A sequence's `body()` task produces items. For each item it calls `start_item(req)`, which blocks until the sequencer grants it access to the driver; it randomizes the item; then `finish_item(req)` hands the item over and blocks until the driver is done. On the other side, the driver loops on `seq_item_port.get_next_item(req)`, drives the item onto the interface over as many cycles as the protocol needs, and calls `seq_item_port.item_done()`, which releases the sequence. Randomizing between start_item and finish_item is **late randomization**: constraints can depend on the state at the moment the driver is ready. If a driver never calls item_done, the sequence blocks in finish_item forever.

## Arbitration, responses and virtual sequences

Several sequences can run on one sequencer at once, for example background traffic plus an error-injection sequence. The sequencer arbitrates among their start_item requests according to its arbitration mode, such as FIFO order, priority or random. A sequence that needs uninterrupted access for a multi-item operation takes `lock()` or `grab()`. Responses travel back the same way: the driver returns a response item and the sequence retrieves it with `get_response`. A **virtual sequence** generates no items itself; it starts sequences on several agents' sequencers in a coordinated order, for example configuring registers through one interface and then starting data traffic on two others. Virtual sequences are where multi-interface scenarios live.

## Layering and reuse

Keeping stimulus in sequences means tests vary traffic without editing the environment: a test selects which sequence to start, sets a default sequence for a sequencer, or uses a factory override to replace the item type with a derived class carrying extra constraints. Sequences can be **layered**: a packet-level sequence produces packets, and a translation layer converts each packet into the byte-level items a lower sequencer and driver understand, so each layer randomizes and collects coverage at its own abstraction. Libraries of sequences, chosen at random in regression, mix scenarios in random order so their interactions are exercised. The driver stays small, protocol-only and stable, which also makes it the right place for protocol assertions.

## Going deeper

Accellera, "UVM 1.2 User's Guide" (2015), chapters on sequences and sequencers; Salemi, "The UVM Primer" (2013).
