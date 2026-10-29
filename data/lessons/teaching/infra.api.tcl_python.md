# Tcl and Python tool APIs

EDA engines are written in C++ for speed, but users and flows drive them from scripts. OpenROAD exposes its engines and its OpenDB database through a Tcl shell and through Python. The scripting layer is how flows are built, how reports are customized and, increasingly, how agents call tools. Designing that layer well decides what can be automated cheaply, what must be compiled, and how safe and reproducible scripted runs are.

## Scripting shells over compiled engines

In a modern EDA tool the command language is a thin layer. A global-placement command typed in the OpenROAD Tcl shell parses its options and calls into the C++ placer, which reads and edits the shared database; no placement arithmetic happens in Tcl. Tcl has stayed the industry's command language because it is small and easy to embed in a C or C++ program, and because SDC, the constraint language, is Tcl. Python is the second interface: OpenROAD can run with a Python interpreter and expose the database objects and engine calls, which suits data analysis and machine-learning experiments. The bindings are usually generated rather than handwritten: a wrapper generator such as SWIG reads interface descriptions of the C++ classes and emits glue code for each language, so one C++ API serves both.

## What to script and what to compile

Every call from a script into C++ crosses a boundary that costs far more than a simple getter. A Python loop that calls one location getter on each of 2,000,000 instances pays that overhead 2,000,000 times and can take minutes, while the same loop in C++ takes milliseconds. The rule follows: scripts orchestrate, configure and report; inner loops are compiled. A custom report of nets with more than 50 sinks and their drivers is a script over the database API. A new detailed-placement move evaluated millions of times is a compiled extension. Between the two sit bulk calls: a command that returns all locations at once, or that accepts a whole list of objects, moves the loop into C++ while keeping the script in control. Good APIs offer such bulk operations for the queries users run most.

## Designing a script API

Several rules separate a robust script API from a fragile one. Commands should accept and return object handles or collections rather than names, so a script does not convert 100,000 objects to strings and look each up again. Handles must be validated: if a script holds a handle to a net that another command deleted, the API should raise a script error rather than read freed memory. Configuration should travel as explicit command options or configuration objects, not as global variables whose last writer wins, so that each call is reproducible and loggable. C++ exceptions must be caught at the boundary and turned into script errors. When the caller is an agent rather than a person, the same rules matter more: typed, documented commands that return structured results can be called and parsed reliably, while free-text reports must be scraped and break when wording changes.

## Going deeper

The OpenROAD documentation on its Tcl and Python interfaces; the SWIG documentation; John Ousterhout, "Scripting: Higher-Level Programming for the 21st Century", IEEE Computer, 1998.
