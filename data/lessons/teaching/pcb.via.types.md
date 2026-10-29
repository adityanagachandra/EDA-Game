# Via types

A via is a plated hole that connects copper on different layers. Its type is defined by which layers it spans and how it is made: through vias span the whole board, blind vias start at an outer layer and stop inside, buried vias join only inner layers, and microvias are small laser-drilled blind structures spanning one dielectric. Via type drives cost, density, routing capacity and signal integrity. A board router models each type as a layer span with its own cost and rules.

## Through, blind, buried and micro

A **through via** is mechanically drilled through the finished laminated stack and plated in one step; it is the cheapest via and connects to any layer it passes. A **blind via** is visible on one surface and ends at an inner layer. A **buried via** joins inner layers only; its sub-stack is drilled and plated before the outer layers are laminated over it, so every distinct buried span is an extra process sequence. A **microvia** is laser-drilled through a single thin dielectric, with an aspect ratio of at most 1:1 in current IPC usage. In a router, a via is an edge between layer graphs, annotated with the span it occupies and a cost. A through via blocks every layer; a blind or buried via blocks only the layers it spans, which is its routing value.

## Pads, antipads and fanout structures

Each via has a **pad** on the layers it connects and an **antipad**, a clearance opening, on the planes it passes without connecting. The antipad size sets the via's capacitance to the planes and hence its impedance. On fine-pitch BGAs, a via cannot fit between pads; **via-in-pad** places it in the land, and the hole must be filled and plated over so solder does not wick into it during reflow. Microvias can be **stacked**, each on a copper-filled one below, which is compact but concentrates thermal-cycling stress at the interfaces, or **staggered**, offset from each other, which uses more area and is generally more robust. These choices change the DRC rule set the router must obey per via type.

## Via stubs and resonance

A through via that carries a signal from L1 to L3 of a 12-layer board leaves the barrel below L3 unused: an open-ended **stub**. A stub resonates when it is a quarter wavelength long, $f_0 = c / (4 L \sqrt{\varepsilon_r})$, where it notches the channel's transmission. Example: a 60 mil stub in $\varepsilon_r = 4.0$ gives $11.8 / (4 \cdot 0.06 \cdot 2) \approx 24.6$ GHz, while a 120 mil stub puts the notch near 12.3 GHz, inside the band of a 25 Gbaud NRZ link whose Nyquist frequency is 12.5 GHz. Fixes are to route on layers that use most of the barrel, to **backdrill** the stub from the far side with a controlled-depth drill, or to use blind vias. The stub length is a direct output of layer assignment, so a router that knows the channel's band can choose layers to keep stubs short.

## Going deeper

IPC-2226, "Sectional Design Standard for High Density Interconnect Printed Boards". Bogatin, "Signal and Power Integrity Simplified" (3rd ed., 2018), chapter on via discontinuities.
