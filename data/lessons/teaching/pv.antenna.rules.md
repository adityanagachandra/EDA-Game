# Antenna rules

During manufacturing, plasma etching of each metal layer charges the exposed metal. If that metal is connected to a transistor gate and to no diffusion that can drain the charge, the charge flows through the thin gate oxide and can damage it. Antenna rules bound the ratio of collecting area to gate area. The check needs connectivity at an intermediate process step, which makes it unlike any purely geometric rule.

## The area antenna ratio

The area ratio is the area of metal connected to a gate, divided by the total gate area on that net. All gates on the net share the collected charge, so their areas add. A net with metal 2 pieces of 24 by 0.14 um and 40 by 0.2 um (3.36 and 8 square um) connected to one gate with W = 0.42 um and L = 0.13 um (0.0546 square um) has a ratio of $11.36 / 0.0546 = 208.1$. Gate area is W times L of the channel. Metal that also touches a source or drain diffusion is protected, because the junction leaks the charge away, which is the basis of the diode fix.

## The longest safe wire

Inverting the ratio gives the longest wire a router may run on one layer before it needs a fix. For a ratio limit $R$, gate area $G$ and wire width $w$, the length must satisfy $w L / G \le R$, so $L \le R G / w$. With $R = 850$, $G = 0.495$ square um and $w = 0.4$ um, $L \le 1051.9$ um, so 1051 um. Routers use this bound to plan layer changes before they create a violation, which is cheaper than repairing afterwards. A larger fanout raises $G$ and relaxes the bound.

## Side-wall ratios

For thick metal layers, charge collects through the side walls exposed during etch, so decks measure the side-wall area: the perimeter of the connected metal times the layer thickness. A 44 by 0.2 um wire, 0.5 um thick, has perimeter $2(44 + 0.2) = 88.4$ um and side-wall area 44.2 square um; on a gate of 0.26 square um the ratio is 170. Decks may also bound via-cut area ratios separately. Each layer can carry its own limit and measurement mode, so an antenna engine evaluates a small formula per layer per net rather than one global number.

## Going deeper

The LEF/DEF Language Reference, sections on antenna properties; the OpenROAD documentation for antenna checking and repair.
