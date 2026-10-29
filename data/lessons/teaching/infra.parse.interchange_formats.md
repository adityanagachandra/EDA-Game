# EDIF, IP-XACT, UPF and Si2 standards

Beyond Verilog, Liberty, LEF/DEF, SPEF and GDS, a tool meets a set of interchange standards for netlists, IP metadata and power intent. Each comes with its own syntax and its own parsing strategy: S-expressions, XML with a schema, or Tcl. Consortia such as Si2 steward several of them. A tool designer should recognize each format, know what it carries, and know what reader it needs.

## EDIF netlist interchange

EDIF (Electronic Design Interchange Format) describes netlists as nested S-expressions, LISP-style parenthesized lists: libraries contain cells, cells contain views, views have an interface (ports) and contents (instances and nets). One recursive reader can parse any EDIF file, because every construct is a list whose first element is a keyword. The price is verbosity: keywords repeat for every object. EDIF remains common in FPGA flows and older tool hand-offs. The parser's job is the same as for a Verilog netlist: build the same instances, nets and terminals in the object model. Importing the same design from EDIF and from Verilog should give identical connectivity, which is a cheap and powerful regression test for both readers.

## IP-XACT component metadata

IP-XACT (IEEE 1685) is an XML schema for describing IP blocks: components, bus interfaces with their port mappings, memory maps, address blocks and registers with their fields. It does not describe logic. It is the single source from which tools generate register RTL, C header files and documentation, and from which SoC integration tools wire blocks together by matching bus interfaces. A reader uses a standard XML parser and validates against the schema, rather than a hand-written grammar. A register of 32 bits with fields at bit offsets 0, 8 and 16 generates both the decode logic and the software mask constants, so hardware and software cannot drift apart.

## UPF power intent and the Si2 coalitions

UPF (Unified Power Format, IEEE 1801) captures power intent outside the RTL: `create_power_domain`, supply nets and ports, isolation, level shifting, retention and power states. Like SDC it is Tcl, so a tool reads it with a Tcl interpreter that defines these commands, and keeping it separate lets one RTL be implemented under different power strategies. Si2 (Silicon Integration Initiative) is an industry consortium whose member coalitions maintain standards, among them OpenAccess, a design database API, and the reference LEF/DEF parsers on which OpenDB's readers are based. Shared reference parsers matter because a format is only as standard as its corner cases: when tools share a parser, the same file means the same design to each of them.

## Going deeper

IEEE Std 1685 (IP-XACT); IEEE Std 1801 (UPF); the EDIF 2 0 0 specification; the Si2 OpenAccess documentation.
