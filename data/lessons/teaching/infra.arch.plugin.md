# Plug-in architectures

Every long-lived EDA tool grows by extension: new synthesis passes, new placement penalties, new reports, new readers. A plug-in architecture lets that growth happen without editing the core. The patterns are old and simple, a registry, a factory and an interface, but the hard parts are the boundary itself: how plug-ins are built and loaded, and how the contract between host and plug-in survives change.

## Registries and factories

A registry maps a name to a way of constructing an implementation. Yosys shows the pattern: a pass is a C++ class derived from `Pass`, its constructor passes a name such as "my_opt" to the base class, and one static instance of the class is declared. Constructing that instance at start-up adds it to a global name-to-pass map, so typing `my_opt` in the shell finds and runs it. The same pattern serves policies inside an engine: a placer defines an abstract density-penalty interface, implementations register under names such as "bell", and a config file selects one by name. The core calls only the interface. Adding a penalty adds a registration, not a core edit. Duplicate names should be rejected, since silent replacement makes behavior depend on load order.

## Compiled-in modules versus loaded libraries

A plug-in can be compiled into the binary or loaded at run time. Yosys supports both: built-in passes, and shared libraries loaded with `yosys -m my_plugin.so`, whose static objects register their passes when the library loads. OpenROAD compiles its tools into one application, each module adding its Tcl commands at start-up. Compiled-in modules are built together, so they agree on the database layout, compiler and standard library; the price is a rebuild to add one. Loaded libraries can be added by users, but a C++ object crossing a library boundary depends on the C++ ABI: class layout, name mangling, standard-library types and exceptions. A plug-in built with another compiler version can compile cleanly and crash. A C-callable boundary passing plain structs and opaque handles is the robust choice.

## Versioning the plug-in contract

Once published, a plug-in interface is a promise. If the host adds a field in the middle of a struct, old plug-ins read wrong offsets. The defenses are standard: every struct carries a version or size field, fields are only appended, the plug-in declares the interface version it was built against, and the host checks it at load time so a mismatch is a clear error rather than a later crash. Plug-ins should keep no hidden global state, so instances run independently on many designs in parallel and results do not depend on scheduling. And the interface stays narrow: every function exposed is a function the core can never change freely again, so hosts expose a small, stable seam rather than their internals.

## Going deeper

The Yosys manual chapter on writing extensions; the OpenROAD developer guide on adding a tool; Gamma, Helm, Johnson and Vlissides, "Design Patterns" (1994), on Factory Method and Abstract Factory.
