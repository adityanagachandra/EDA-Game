# Timing closure

Timing closure is the iterative process of bringing a design from its first timing estimate to a netlist and layout that pass signoff in every scenario. It is not one algorithm but a loop of analysis and optimization across synthesis, placement, clock tree synthesis and routing, in which each stage's estimate is optimistic relative to the next and every fix can disturb something else.

## Reading the shape of the problem

The first decision is diagnostic. WNS gives the worst violation, TNS and the count of failing endpoints give its extent. A block with WNS -40 ps and TNS -9000 ps over 2000 endpoints has a systematic problem, such as a poor floorplan, congestion, a missing exception or an over-tight constraint, and patching the worst path will barely move TNS. A slack wall, many endpoints near the worst value, means fixing one path just exposes the next. A design that met timing at placement and fails after routing usually suffered from optimistic wire estimates: detours, layer assignment and coupling appear only once wires exist. Teams therefore over-constrain early stages, for instance with a tighter clock, so that later losses are covered.

## Choosing the fix

Each lever fits a cause. Gate sizing helps when a cell drives too much load, with diminishing returns once gates are large. When wire delay dominates, buffering, layer promotion or moving cells closer attack the wire. Swapping high-Vt cells for low-Vt cells speeds a path at the cost of leakage power. Useful skew moves clock edges to lend time between stages. Logic restructuring shortens deep paths when local fixes are exhausted. Every fix must be evaluated in all active scenarios, since a change that helps the slow corner's setup can break the fast corner's hold. Fixes are ordered too: setup first, because sizing and buffering reshape many paths, then hold, by inserting delay that costs little setup slack.

## Late-stage closure

As the design stabilizes, the clocks switch from ideal to propagated after CTS, and the skew estimate leaves clock uncertainty so skew is not counted twice. Signoff timing uses the accurate extraction and timer; implementation tools have built-in fast timers whose numbers correlate closely but not exactly, so flows keep a little margin. Late fixes are ECOs: small, local netlist changes applied incrementally with the rest of the layout frozen. After the base layers are fixed, only a metal-only ECO is possible, using spare cells already placed in the layout. Closure ends when every scenario passes at signoff with the agreed margins, and the fixes are verified for function as well as timing.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8 on timing closure; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009).
