# SPEF name maps and file size

Full-chip SPEF files run to tens of gigabytes, and most of their bytes are names: long hierarchical instance and net names repeated in every entry. The name map writes each name once and refers to it by a short index token. The arithmetic of that saving explains why name maps are universal, and the same reasoning about repeated strings guides the design of every large EDA exchange format.

## Bytes saved per net

With a name map, each reference to a name of length $\ell$ is replaced by a token of length $d$ (an asterisk and digits, such as *4721), saving $\ell - d$ bytes per reference. The map itself costs one line per name: the token, a space, the name and a newline, $d + \ell + 2$ bytes. A net referenced $r$ times saves $r(\ell - d) - (d + \ell + 2)$. For $\ell = 57$, $d = 7$ and $r = 13$, the saving is $13 \times 50 - 66 = 584$ bytes. References are frequent: the D_NET line, every internal node name such as net:3, and every coupling entry on another net that touches this one all repeat the name.

## Scaling to a whole chip

Multiply the per-net saving by the net count. With 2 million nets of 31-character names, 7-character tokens and 16 references each, every net saves $16 \times 24 - 40 = 344$ bytes, so the map removes 688 MB. Parsing benefits as much as disk: a reader hashes each long name once while loading the map and afterwards works with integer indices, so lookup and memory are both cheaper. The same idea appears elsewhere, in string tables of netlist databases and in DEF and Verilog readers that intern names, and a tool designer applies it whenever a format repeats large keys.

## The break-even point

A map entry pays for itself once $r(\ell - d) > \ell + d + 2$, that is $r > (\ell + d + 2)/(\ell - d)$. For $\ell = 74$ and $d = 4$, $r > 80/70 = 1.14$, so two references suffice. Because every net in a detailed SPEF is named on its D_NET line and in each internal node name, almost every net clears the bar, and writers simply map every name. The exceptions are short names, where $\ell$ is close to $d$; once indices reach many digits, very short names can cost more mapped than unmapped, and a writer can choose to leave them inline. This is the same per-entry cost-benefit test as any dictionary compression scheme.

## Going deeper

IEEE Std 1481 (SPEF), the NAME_MAP section; Salomon, "Data Compression: The Complete Reference" (4th ed., 2007), on dictionary methods.
