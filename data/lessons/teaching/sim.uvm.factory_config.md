# The UVM factory and configuration database

A reusable environment must let each test change behavior without editing shared code. UVM provides two mechanisms for this. The **factory** substitutes derived classes for base classes at creation time. The **configuration database** passes values and handles, including virtual interfaces, down the hierarchy by path. Both are lookups resolved at run time, so they are flexible, and their failure modes are silent unless code checks for them.

## Factory creation and overrides

Components and objects registered with the factory (through the `uvm_component_utils` or `uvm_object_utils` macros) are created with `my_driver::type_id::create("drv", this)` rather than `new()`. The factory checks its override table first. A **type override**, such as `set_type_override_by_type(base_item::get_type(), err_item::get_type())`, makes every factory creation of `base_item` return an `err_item`, a derived class that might add error-injection constraints. An **instance override** applies only at matching hierarchical paths, for example to the driver of one agent out of two. Anything built with `new()` bypasses the factory, so overrides silently fail to apply to it, which is the most common factory bug.

## The configuration database

`uvm_config_db#(T)::set(context, "path", "field", value)` stores a value for components whose hierarchical names match the path, which can contain wildcards. A component retrieves it with `uvm_config_db#(T)::get(this, "", "field", var)`, normally in build_phase. The typical use is passing a virtual interface from the top-level module, which instantiates the real interface, into the driver and monitor, which are classes and cannot reach module instances directly. Precedence matters: during build, a set made from higher in the hierarchy wins over one made lower, so a test can override an environment's defaults. After build, the last set wins. Configuration objects bundle many fields into one class to keep this manageable.

## Silent failures and defensive use

Both mechanisms resolve strings at run time, so mistakes do not fail compilation. A path typo, such as setting "env.agnt*" for a component named "env.agent0", matches nothing; the get returns 0, leaves the variable unchanged, and the component runs with its default. A get that runs before the corresponding set, because of phase ordering, has the same effect. Defensive code checks the return value of every get and reports a fatal error for required items such as virtual interfaces, and tests print the factory's override table and the config_db's contents in debug runs. Changing configuration during the run phase invites races between set and get and is best avoided: configuration shapes construction and belongs in build.

## Going deeper

Accellera, "UVM 1.2 User's Guide" (2015), chapters on the factory and the configuration database; IEEE Std 1800.2-2020, uvm_factory and uvm_config_db.
