package infra

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func init() { registrars = append(registrars, registerAPI) }

func registerAPI(reg *generator.Registry) {
	reg.Register("infra.api.tcl_evaluation", edakit.Func(tclEvaluation))
	reg.Register("infra.api.tcl_python", &edakit.Bank{Items: tclPythonItems})
	reg.Register("infra.api.dsl_design", &edakit.Bank{Items: dslDesignItems})
}

var apiWords = []string{"u1", "u2", "clk", "rst", "n12", "core", "a", "b", "c", "x", "y", "met1", "met2", "in0", "out"}

// apiList builds a random Tcl list literal (without the outer braces) and returns it with its
// top-level element count.
func apiList(ctx generator.GeneratorContext) (string, int) {
	n := edakit.Scaled(ctx, 3, 6)
	var parts []string
	for i := 0; i < n; i++ {
		switch edakit.Int(0, 4) {
		case 0, 1:
			parts = append(parts, edakit.Pick(apiWords))
		case 2:
			k := edakit.Int(2, 3)
			ws := edakit.Shuffle(apiWords)[:k]
			parts = append(parts, "{"+strings.Join(ws, " ")+"}")
		case 3:
			ws := edakit.Shuffle(apiWords)[:2]
			parts = append(parts, "{"+ws[0]+" {"+ws[1]+" "+edakit.Pick(apiWords)+"}}")
		default:
			parts = append(parts, "{}")
		}
	}
	return strings.Join(parts, " "), n
}

func tclEvaluation(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		l, n := apiList(ctx)
		q := "In Tcl, which value does the command `llength {" + l + "}` return?"
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("llength counts top-level elements. A braced group, however much it contains and even when empty, is one element, so the list {%s} has %d elements.", l, n),
			map[string]string{"llength": edakit.IntStr(n)})
	case 2:
		v := edakit.Pick([]string{"U42", "core_0", "met3", "clk_div", "n1077", "reg_a"})
		pre := edakit.Pick([]string{"inst_", "top/", "", "pin:"})
		suf := edakit.Pick([]string{"/Z", "_q", "", ".A1"})
		quoted := edakit.Int(0, 1) == 0
		var body string
		var length int
		if quoted {
			body = "\"" + pre + "[set w]" + suf + "\""
			length = len(pre) + len(v) + len(suf)
		} else {
			body = "{" + pre + "[set w]" + suf + "}"
			length = len(pre) + len("[set w]") + len(suf)
		}
		q := edakit.Sprintf("A Tcl script runs `set w %s` and then `set s %s`. Which value does `string length [set s]` return?", v, body)
		var why string
		if quoted {
			why = edakit.Sprintf("Inside double quotes Tcl performs command substitution, so [set w] is replaced by %s and s becomes %s%s%s, which has %d characters.", v, pre, v, suf, length)
		} else {
			why = edakit.Sprintf("Inside braces Tcl performs no substitution, so s holds the literal text %s[set w]%s, which has %d characters.", pre, suf, length)
		}
		return edakit.Numeric(q, edakit.IntStr(length), why, map[string]string{"length": edakit.IntStr(length)})
	}
	a := edakit.Int(7, 60)
	b := edakit.Int(2, 9)
	if a%b == 0 {
		a++
	}
	if edakit.Int(0, 1) == 0 {
		a = -a
	} else if edakit.Int(0, 2) == 0 {
		b = -b
	}
	op := edakit.Pick([]string{"/", "%"})
	// Tcl integer division rounds toward negative infinity; the remainder takes the divisor's sign.
	qt := a / b
	if (a%b != 0) && ((a < 0) != (b < 0)) {
		qt--
	}
	r := a - b*qt
	ans := qt
	if op == "%" {
		ans = r
	}
	q := edakit.Sprintf("In Tcl, which value does `expr {%d %s %d}` return?", a, op, b)
	return edakit.Numeric(q, edakit.IntStr(ans),
		edakit.Sprintf("Tcl integer division rounds toward negative infinity, so %d / %d = %d, and the remainder %d %% %d = %d has the sign of the divisor, keeping a = b x quotient + remainder: %d x %d + %d = %d. C and C++ truncate toward zero instead, a classic porting bug.", a, b, qt, a, b, r, b, qt, r, a),
		map[string]string{"quotient": edakit.IntStr(qt), "remainder": edakit.IntStr(r)})
}

var tclPythonItems = []edakit.Item{
	{KP: 1, Q: "OpenROAD exposes its engines through a Tcl shell. What does a typical Tcl command such as a global-placement command do underneath?",
		A: "Calls into the C++ engine, which reads and edits the shared database", D: []string{"Writes a DEF file and launches a separate placer binary", "Interprets a placement algorithm written in Tcl", "Sends the design to a remote server"},
		Why: "The Tcl layer is a thin command interface. The compiled engine does the work on the in-memory database, so scripting costs almost nothing at run time."},
	{KP: 1, Q: "Why have EDA tools kept Tcl as their command language for decades?",
		A: "It is small, embeddable and SDC is already written in it", D: []string{"It is the fastest language for numerical kernels", "It has the strictest type system", "Foundries require it for signoff"},
		Why: "Tcl is easy to embed in a C or C++ application, and SDC, the constraint language, is Tcl. Speed comes from the compiled engines, not from Tcl."},
	{KP: 1, Q: "Besides Tcl, what scripting interface does OpenROAD offer for working with its database and engines?",
		A: "A Python interface to the same engines and database objects", D: []string{"A JavaScript interface running in a browser", "A Perl interface only for reports", "No other interface; Tcl is the only way in"},
		Why: "OpenROAD can run with a Python interpreter, giving access to the database objects and engine calls, which suits data analysis and machine-learning work."},
	{KP: 1, Q: "How are C++ database classes usually made callable from Tcl and Python in open EDA tools?",
		A: "A wrapper generator produces bindings from the C++ headers or interface files", D: []string{"Each class is rewritten by hand in both languages", "The C++ code is compiled to Tcl bytecode", "The scripts parse the C++ source at run time"},
		Why: "Tools such as SWIG generate glue code that maps C++ classes and methods to script objects, so one C++ API serves several languages."},
	{KP: 2, Q: "A script loops over all 2 million instances in Python, calling one database getter per instance for its location, and it takes minutes. What is the usual fix?",
		A: "Move the loop into a compiled command or fetch data in bulk", D: []string{"Rewrite the same loop in Tcl", "Add more threads to the Python loop", "Save the design to DEF first"},
		Why: "Each script-to-C++ call has overhead far above the work of one getter. Bulk calls or a compiled command amortize the crossing."},
	{KP: 2, Q: "A Tcl procedure stores a database object handle in a variable, then a command deletes that net. The procedure later uses the handle. What should a well-designed API do?",
		A: "Detect that the object is gone and raise a script error", D: []string{"Silently return data from the deleted net", "Recreate the net", "Crash the tool as intended"},
		Why: "Script handles must be validated, or a stale handle becomes a use-after-free inside the C++ engine."},
	{KP: 2, Q: "A flow engineer wants a custom report: all nets with more than 50 sinks and their driver cells. What is the right layer to implement it in?",
		A: "A script over the database API", D: []string{"A change to the router's core", "A new Liberty attribute", "A modification of the DEF grammar"},
		Why: "Queries and reports that only read the database belong in scripts. Core edits are for algorithms that must run at compiled speed."},
	{KP: 2, Q: "Which task needs a compiled C++ extension rather than a script?",
		A: "A new detailed-placement move evaluated millions of times", D: []string{"Reading a configuration file of 20 parameters", "Printing a summary of utilization", "Selecting which corner to analyze"},
		Why: "Inner-loop work runs millions of times and must be compiled. Scripts orchestrate and configure; they should not sit in hot loops."},
	{KP: 3, Q: "A Tcl command returns a list of 100,000 object names as one string, and the script then looks each name up again. What design would avoid the waste?",
		A: "Return object handles or a collection that commands accept directly", D: []string{"Return a longer string", "Write the names to a file", "Require the user to type the names"},
		Why: "Name round-trips cost string building and hashing. Collections of handles, as in SDC-style get commands, pass objects without re-lookup."},
	{KP: 3, Q: "Two scripts set the same global Tcl variable to configure an engine, and results depend on which ran last. What API change fixes this?",
		A: "Pass configuration as explicit command options or a config object", D: []string{"Rename the variable", "Run the scripts in a random order", "Store the variable in DEF"},
		Why: "Hidden global state makes runs order-dependent. Explicit parameters make each call reproducible and easy to log."},
	{KP: 3, Q: "A tool exposes a Python API and an agent will drive it. Which property of the API matters most for the agent?",
		A: "Typed, documented calls that return structured results", D: []string{"Short command names", "Commands that print human-only text", "As many commands as possible"},
		Why: "An agent needs to know argument types and to parse the result. Free-text output must be scraped and breaks when wording changes."},
	{KP: 3, Q: "A Python script calls an engine, which throws a C++ exception. What should the binding layer do?",
		A: "Translate it into a Python exception with the message", D: []string{"Let the process abort", "Ignore it and return zero", "Print it and continue with stale data"},
		Why: "Exceptions must not cross the language boundary raw. The binding catches them and raises a script-level error that the caller can handle."},
}

var dslDesignItems = []edakit.Item{
	{KP: 1, Q: "SDC constraints are commands executed by a Tcl interpreter inside the tool. What kind of domain-specific language is SDC?",
		A: "An embedded DSL hosted in Tcl", D: []string{"A standalone compiled language", "A markup language like XML", "A binary format"},
		Why: "SDC adds domain commands such as create_clock to Tcl and inherits Tcl syntax, variables and control flow."},
	{KP: 1, Q: "Liberty has its own grammar of groups and attributes and no loops or variables. What kind of language is it?",
		A: "An external declarative data language", D: []string{"An embedded scripting language", "A general-purpose programming language", "A hardware description language"},
		Why: "Liberty is declarative data with a dedicated parser. It describes cells; it does not compute."},
	{KP: 1, Q: "What is the main advantage of an embedded DSL over an external one for tool configuration?",
		A: "It reuses the host language's parser, variables and control flow", D: []string{"It is always faster to execute", "It cannot have syntax errors", "It needs no documentation"},
		Why: "Embedding gives loops, procedures and expressions for free. The cost is that users can write arbitrary code, which is harder to analyze."},
	{KP: 1, Q: "What is the main advantage of a declarative external DSL such as a YAML or JSON flow configuration?",
		A: "A tool can validate and diff it without executing anything", D: []string{"It can express any computation", "It runs faster than compiled code", "It removes the need for a schema"},
		Why: "Declarative data can be checked against a schema, compared between runs and generated by programs safely."},
	{KP: 2, Q: "A user writes create_clock with a period of -2 in an SDC file. When should a well-designed reader report it?",
		A: "When the command executes, with the file and line", D: []string{"Only after timing analysis finishes", "Never; negative periods are legal", "Only at tape-out"},
		Why: "Semantic checks belong at the command, where the context is known, so the user can fix the source line."},
	{KP: 2, Q: "A flow config has the key PLACE_DENSTY (misspelled). The flow silently uses the default density. What DSL design choice would catch it?",
		A: "Reject unknown keys against a declared schema", D: []string{"Accept any key for flexibility", "Make all keys optional", "Store config in environment variables"},
		Why: "Silently ignoring unknown keys turns typos into wrong results. A schema with known keys and types fails fast."},
	{KP: 2, Q: "A team wants agents to propose flow configurations. Which representation is safest to accept from an agent?",
		A: "Declarative parameters validated against a schema", D: []string{"Arbitrary Tcl evaluated in the tool", "Shell commands run as the user", "Patches to the engine source"},
		Why: "Validated data limits what an agent can do to the intended knobs. Arbitrary code can do anything."},
	{KP: 2, Q: "A Yosys synthesis script is a sequence of named passes with options. What does that make it?",
		A: "A command language over a pass registry", D: []string{"A netlist format", "A timing constraint language", "A layout description"},
		Why: "Each command names a registered pass that transforms the design in RTLIL. The script orders the passes."},
	{KP: 3, Q: "An SDC file uses Tcl loops to generate 10,000 set_false_path commands. A static checker without a Tcl interpreter must report all exceptions. What does it need?",
		A: "To execute the Tcl, because the commands exist only after evaluation", D: []string{"Only a regular-expression scan of the file", "Only the first line", "The Liberty file"},
		Why: "In an embedded DSL the effective commands are the result of running the program. A pure text scan cannot see generated commands."},
	{KP: 3, Q: "A DSL for design rules lets authors write each rule as a declarative expression. What does that let the engine do that imperative scripts would not?",
		A: "Analyze, reorder and batch rules for efficient evaluation", D: []string{"Skip checking rules", "Avoid storing geometry", "Run without a layout"},
		Why: "Declarative rules expose their structure, so the engine can share work across rules and run them in parallel."},
	{KP: 3, Q: "Two versions of a tool interpret the same configuration key differently. What practice prevents silent changes in results?",
		A: "Version the schema and record it with every run", D: []string{"Never change the tool", "Remove the key", "Let each user decide"},
		Why: "A versioned schema makes changes explicit, and recording it with results lets a regression explain a difference."},
	{KP: 3, Q: "A DSL interpreter is embedded in a tool that agents can call. What is the main safety risk of a Tcl-based DSL there?",
		A: "Tcl can run shell commands and touch files unless restricted", D: []string{"Tcl cannot express numbers", "Tcl has no comments", "Tcl is too slow to parse"},
		Why: "Tcl has exec, file and socket commands. A safe interpreter or a restricted command set is needed when input is not trusted."},
}
