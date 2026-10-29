# The PCB tool landscape

Board design tools range from open-source suites to enterprise platforms. For a tool designer the useful questions are not which buttons they have, but what data model each exposes, how open its file formats and scripting are, how it represents constraints, and how designs move between tools. This lesson maps the main families: KiCad, LibrePCB and gEDA/Lepton in open source; Altium Designer, Cadence Allegro and OrCAD, Siemens Xpedition and PADS, Zuken, and Autodesk's EAGLE lineage commercially.

## Open-source tools

**KiCad** is the most widely used open-source suite: schematic and PCB editors, an interactive push-and-shove router, a 3D viewer and a Python scripting interface. Its schematics and boards are readable S-expression text files, which makes them diffable and parseable by external tools. **LibrePCB** is built around a versioned library in which every symbol, package and device is identified by a UUID, so libraries can be shared without name clashes. **gEDA** was an early open toolkit (a schematic editor, the netlister gnetlist and the pcb layout tool); **Lepton EDA** is a community fork of the gEDA schematic and netlisting tools, scripted in Guile Scheme, with netlist back ends written as small Scheme programs. **Freerouting** is an open-source autorouter that exchanges boards through Specctra DSN and SES files.

## Commercial families

**Altium Designer** is a unified schematic, layout and library environment, paired with Altium's cloud collaboration service. **Cadence** offers **OrCAD** (OrCAD Capture for schematics and mainstream layout) and **Allegro**, an enterprise layout platform with a central constraint manager. **Siemens** offers **Xpedition** for enterprise design and **PADS** for mainstream design, both from its Mentor Graphics acquisition. **Zuken** offers **CR-8000** and **eCADSTAR**. **EAGLE** was acquired by Autodesk from CadSoft, and Autodesk moved its electronics capability into the Fusion electronics workspace. Enterprise tools emphasize what large teams need: a shared constraint system, concurrent multi-user editing, library governance and packaging co-design, rather than drawing features.

## Data models and interoperability

For automation, a tool's data model matters more than its editor. Text formats with a scripting API (KiCad) let a team run custom checks across hundreds of designs; binary formats confine automation to the vendor's API. Moving a design between tools preserves connectivity and geometry best, because every tool has nets, footprints, traces and vias. **Constraints** translate worst, because each tool models net classes, matching groups and rule scopes differently, so they are usually rebuilt by hand. Neutral formats help at the edges: Specctra DSN hands a placed board to an external router, IPC-2581 and ODB++ carry a finished board to manufacturing, and the gEDA netlister's pluggable back ends show the seam pattern, one stable connectivity model with many small format writers.

## Going deeper

The KiCad developer documentation on the S-expression file formats and the Python API. The Lepton EDA documentation on netlister back ends.
