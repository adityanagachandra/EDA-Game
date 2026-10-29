# Diagnosing LVS errors

An LVS report is a list of symptoms: unmatched nets, unmatched devices, property errors. Each class of physical mistake leaves a typical signature, and one mistake can produce hundreds of symptoms. Tool designers shape reports around these signatures, and designers and agents fix the root cause, not the symptoms.

## Shorts and opens

A short merges two nets: two different label names land on one extracted net, and devices of both schematic nets attach to it. An open splits a net: one schematic net matches two layout nets, each with part of the connections, typically because of a missing via, a gap, or a crossing without a via. A label placed on a text layer that does not attach to the pin's metal names nothing, so the port appears missing. These are connectivity errors, and they change the graph that refinement works on.

## Device and property errors

When the topology matches but devices differ, the report shows device-level errors. A W of 400 nm against 450 nm with a 5 percent tolerance is a property error. A missing implant changes the recognized device type. Fingers that are not merged, because parallel reduction is off or the fingers differ in L, give a device-count mismatch. A PMOS whose well is tied to the wrong supply mismatches on its bulk terminal. These errors are local to the device and do not cascade, unlike shorts and opens.

## Debugging order

Fix shorts first, then opens, then device and property errors, rerunning after each step. A single short perturbs refinement around it and can produce dozens of unmatched nets that disappear once it is fixed. Short finders locate a short by a shortest path through the merged net's shapes between a point labelled with one name and a point labelled with the other. When a clean block suddenly fails with hundreds of errors after a small edit, suspect one connectivity change. Hard IP without a netlist is black-boxed: only its pins are compared.
