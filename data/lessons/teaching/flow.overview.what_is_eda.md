# What electronic design automation is

Electronic design automation (EDA) is the software that turns an idea for an electronic system into something a factory can build and a tester can check. It covers integrated circuits, printed circuit boards and packages. Older names overlap with it: ECAD (electronic computer-aided design) and CAE (computer-aided engineering). For a tool designer the key fact is that every EDA program is an algorithm over a precise model: a graph, a set of Boolean functions, a geometry database or a system of equations.

## Design, analysis and verification tools

EDA tools fall into three kinds by what they do to a design. Synthesis or implementation tools transform it: logic synthesis turns register-transfer code into gates, a placer assigns cells to locations, a router draws wires. Analysis tools compute a property without changing anything: a static timer computes arrival times and slack, an extractor computes resistance and capacitance, a power tool computes switching power. Verification tools check that two descriptions agree or that a property holds: a simulator runs test stimulus, an equivalence checker proves that a netlist matches its RTL, a design-rule checker proves that layout geometry obeys the foundry's rules. Most real tools mix the three: an optimizer calls an analyzer inside its loop.

## Every tool is an algorithm on a model

A tool never sees the chip. It sees a model. A gate netlist is a hypergraph whose vertices are cells and whose hyperedges are nets. A timing graph is a directed acyclic graph with delays on its edges. Logic is a Boolean network, often stored as an and-inverter graph. Layout is a set of polygons on numbered layers. Circuit simulation is a system of nonlinear differential-algebraic equations. Once the model is fixed, the tool's job is a classical problem: shortest paths, partitioning, satisfiability, sparse linear solves, computational geometry. That is why an applied mathematician can learn to write EDA tools.

## Why automation is unavoidable

A modern chip holds billions of transistors and is checked against thousands of geometric rules and many timing scenarios. No team can place, wire or check that by hand. Automation also makes design repeatable: the same inputs and settings give the same layout, so engineers can change one thing and measure its effect. The cost is that tool quality now bounds chip quality. A placer that wastes 5% of wirelength, or a timer that is pessimistic by 20 ps, costs power, area or frequency on every chip that uses it.

## Going deeper

Lavagno, Markov, Martin and Scheffer (eds.), "Electronic Design Automation for IC System Design, Verification, and Testing" and its companion volume on implementation (2nd ed., 2016).
