package infra

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func init() { registrars = append(registrars, registerArch) }

func registerArch(reg *generator.Registry) {
	reg.Register("infra.arch.shared_database_tools", &edakit.Bank{Items: sharedDatabaseToolsItems})
	reg.Register("infra.arch.plugin", &edakit.Bank{Items: pluginItems})
	reg.Register("infra.arch.callbacks", edakit.Func(callbacks))
}

// archEvents are the database edit kinds an observer may subscribe to.
var archEvents = []string{"move", "swap", "connect", "create"}

var archEventText = map[string]string{
	"move":    "instance moved",
	"swap":    "master swapped",
	"connect": "pin connected or disconnected",
	"create":  "instance created",
}

type archObserver struct {
	name string
	subs map[string]bool
}

func archObservers(ctx generator.GeneratorContext) []archObserver {
	all := []string{"timer", "density map", "global router", "parasitic estimator", "power analyzer"}
	n := edakit.Scaled(ctx, 3, 5)
	pick := edakit.Shuffle(all)[:n]
	var out []archObserver
	for _, nm := range pick {
		subs := map[string]bool{}
		k := edakit.Int(1, 3)
		for _, e := range edakit.Shuffle(archEvents)[:k] {
			subs[e] = true
		}
		out = append(out, archObserver{nm, subs})
	}
	return out
}

func archDescribeObservers(obs []archObserver) string {
	var parts []string
	for _, o := range obs {
		var ev []string
		for _, e := range archEvents {
			if o.subs[e] {
				ev = append(ev, archEventText[e])
			}
		}
		parts = append(parts, "the "+o.name+" subscribes to "+names(ev))
	}
	return strings.Join(parts, "; ")
}

func callbacks(ctx generator.GeneratorContext) generator.Problem {
	obs := archObservers(ctx)
	m := edakit.Scaled(ctx, 4, 9)
	edits := make([]string, m)
	for i := range edits {
		edits[i] = edakit.Pick(archEvents)
	}
	count := func(e string) int {
		c := 0
		for _, o := range obs {
			if o.subs[e] {
				c++
			}
		}
		return c
	}
	var seq []string
	for _, e := range edits {
		seq = append(seq, archEventText[e])
	}
	switch variant(ctx) {
	case 1:
		total := 0
		var work []string
		for _, e := range edits {
			total += count(e)
			work = append(work, edakit.IntStr(count(e)))
		}
		q := edakit.Sprintf("A design database notifies registered observers on every edit: %s. An optimizer then makes these edits in order: %s. How many callback invocations does the database make in total?", archDescribeObservers(obs), strings.Join(seq, ", "))
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("Each edit calls every observer subscribed to its kind: %s, which sums to %d.", strings.Join(work, " + "), total),
			map[string]string{"callbacks": edakit.IntStr(total)})
	case 2:
		subs := map[string]bool{}
		for _, e := range edakit.Shuffle(archEvents)[:edakit.Int(2, 3)] {
			subs[e] = true
		}
		var tokens []string
		eager, lazy := 0, 0
		dirty := false
		var subText []string
		for _, e := range archEvents {
			if subs[e] {
				subText = append(subText, archEventText[e])
			}
		}
		n := edakit.Scaled(ctx, 6, 12)
		for try := 0; try < 60; try++ {
			tokens, eager, lazy, dirty = nil, 0, 0, false
			hadQuery := false
			for i := 0; i < n; i++ {
				if i > 0 && edakit.Int(0, 2) == 0 || (i == n-1 && !hadQuery) {
					tokens = append(tokens, "QUERY")
					hadQuery = true
					if dirty {
						lazy++
						dirty = false
					}
					continue
				}
				e := edakit.Pick(archEvents)
				tokens = append(tokens, archEventText[e])
				if subs[e] {
					eager++
					dirty = true
				}
			}
			if lazy >= 2 && eager > lazy {
				break
			}
		}
		q := edakit.Sprintf("A timer subscribes to %s. Instead of recomputing in each callback, it sets a dirty flag and recomputes once at the next timing query, and only if the flag is set. The session runs: %s. How many times does the timer recompute?", names(subText), strings.Join(tokens, ", "))
		return edakit.Numeric(q, edakit.IntStr(lazy),
			edakit.Sprintf("An eager timer would recompute on each of the %d relevant callbacks. The lazy one recomputes only at queries that find the flag set, which happens %d times.", eager, lazy),
			map[string]string{"eager": edakit.IntStr(eager), "lazy": edakit.IntStr(lazy)})
	}
	// KP 3: batch per transaction.
	nt := edakit.Scaled(ctx, 2, 4)
	var txs [][]string
	var txText []string
	for t := 0; t < nt; t++ {
		k := edakit.Int(2, 4)
		var tx []string
		var tt []string
		for i := 0; i < k; i++ {
			e := edakit.Pick(archEvents)
			tx = append(tx, e)
			tt = append(tt, archEventText[e])
		}
		txs = append(txs, tx)
		txText = append(txText, edakit.Sprintf("transaction %d: %s", t+1, strings.Join(tt, ", ")))
	}
	batched, perEdit := 0, 0
	var work []string
	for _, tx := range txs {
		got := 0
		for _, o := range obs {
			hit := false
			for _, e := range tx {
				if o.subs[e] {
					hit = true
					perEdit++
				}
			}
			if hit {
				got++
			}
		}
		batched += got
		work = append(work, edakit.IntStr(got))
	}
	q := edakit.Sprintf("A database batches notifications: at each commit, every observer subscribed to at least one edit kind in the transaction gets exactly one call carrying the list of edits. Observers: %s. The optimizer commits %s. How many observer calls does the database make?", archDescribeObservers(obs), strings.Join(txText, "; "))
	return edakit.Numeric(q, edakit.IntStr(batched),
		edakit.Sprintf("Per transaction, count the observers touched at least once: %s, total %d. Per-edit notification would have made %d calls.", strings.Join(work, " + "), batched, perEdit),
		map[string]string{"batched": edakit.IntStr(batched), "per_edit": edakit.IntStr(perEdit)})
}

var sharedDatabaseToolsItems = []edakit.Item{
	{KP: 1, Q: "In a flow of point tools, the placer writes DEF, the router reads it, and the timer reads the routed DEF plus SPEF. What is the main cost of this handoff?",
		A: "Every step re-parses and rebuilds the whole design, so small loops are slow", D: []string{"DEF cannot store placement", "Point tools cannot share a Liberty file", "Files force every tool to use the same algorithm"},
		Why: "Writing and parsing gigabytes of text at every hand-off makes tight optimize-and-analyze loops impractical."},
	{KP: 1, Q: "A placer wants timing feedback after every few thousand moves. Why is that hard when the timer is a separate program reading files?",
		A: "Each query would require writing the design out and starting the timer again", D: []string{"Timers cannot read placed designs", "Placers do not know which cells exist", "File-based timers produce different slacks"},
		Why: "Interactive feedback needs the timer to see the placer's state directly. File hand-off turns every query into a full export and import."},
	{KP: 1, Q: "Point tools from different vendors exchange LEF, DEF, Verilog and SPEF. What does that architecture get right?",
		A: "Any tool can be swapped, because the formats are the contract between them", D: []string{"It removes all parsing work", "It guarantees the same results from every vendor", "It lets tools share memory directly"},
		Why: "Standard formats decouple tools. The cost is the serialization at every boundary."},
	{KP: 1, Q: "A team hands designs between tools as DEF and finds that a custom per-instance attribute they rely on disappears. Why?",
		A: "The receiving format or reader has no field for it, so it is dropped in the round trip", D: []string{"DEF files are limited to 1000 instances", "Attributes are stored only in GDS", "The router deletes all attributes on purpose"},
		Why: "A file boundary keeps only what both sides represent. Anything outside the format is lost unless both tools agree on an extension."},
	{KP: 2, Q: "OpenROAD runs floorplanning, placement, CTS, routing and timing in one process. What do these engines share?",
		A: "One in-memory OpenDB database holding the design", D: []string{"A directory of DEF files rewritten after each step", "A network socket per engine", "A copy of the design per engine, merged at the end"},
		Why: "The tools are modules of one application that read and edit the same database, so results are visible to the next engine immediately."},
	{KP: 2, Q: "In a shared-database tool, the resizer swaps a cell's master. What must happen for the timer to see the change?",
		A: "The timer is notified through the database (for example a callback) and invalidates the affected timing", D: []string{"The resizer writes a new Verilog file", "The user reloads the design", "Nothing; the timer re-times from scratch on every query"},
		Why: "Sharing data is only half the job. Engines that cache derived results must hear about edits so they can update incrementally."},
	{KP: 2, Q: "Which engine pairing benefits most from sharing one database rather than files?",
		A: "A repair loop where the resizer edits and the timer re-reports thousands of times", D: []string{"Writing the final GDS once", "Reading the Liberty files at start-up", "Printing the final area report"},
		Why: "The benefit scales with the number of round trips. Tight edit-and-analyze loops are where file hand-off hurts most."},
	{KP: 2, Q: "A shared database stores the netlist, placement and routing. Which data does a timer typically keep in its own structures instead?",
		A: "Its timing graph with arrival and required times", D: []string{"The instance names", "The cell masters", "The die area"},
		Why: "Engine-specific derived data lives in the engine, built from the database and kept in sync with it."},
	{KP: 3, Q: "Two engines edit the same database at once from different threads without coordination. What is the risk?",
		A: "One engine can read a half-applied edit, or both can corrupt shared lists", D: []string{"The database doubles its memory", "The engines produce identical results", "The Tcl shell stops accepting commands"},
		Why: "Shared mutable state needs a discipline: one writer at a time, or partitioned regions, or transactions."},
	{KP: 3, Q: "A detailed router deletes a net's wires while the timer holds a pointer to one of the wire objects. What design avoids a crash?",
		A: "Engines hold stable ids or handles and hear about deletions through callbacks", D: []string{"Never deleting anything from the database", "Copying the database before each engine runs", "Running the timer only at the end of the flow"},
		Why: "Ids can be checked for validity, and deletion callbacks let an engine drop its references before the object goes away."},
	{KP: 3, Q: "A user runs global placement, then edits the netlist by hand in Tcl, then asks for timing. Why can the report be stale in a careless implementation?",
		A: "The timer's cached graph was not invalidated by the manual edit", D: []string{"Tcl edits never reach the database", "Global placement locks the netlist", "Timing reports always ignore manual edits"},
		Why: "Every path into the database, including scripts, must trigger the same notifications that engine edits do."},
	{KP: 3, Q: "Why do shared-database tools still write DEF, Verilog and SPEF at stage boundaries?",
		A: "For checkpoints, debugging and hand-off to other tools such as sign-off", D: []string{"Because engines cannot run without reading files", "Because the database is erased between commands", "To make the flow run faster"},
		Why: "Files remain the interchange and checkpoint mechanism; the shared database removes them only from the inner loops."},
}

var pluginItems = []edakit.Item{
	{KP: 1, Q: "A Yosys pass is a C++ class whose constructor passes the name \"my_opt\" to the Pass base class, and one static instance is declared. How does `my_opt` become a command?",
		A: "The static instance registers itself by name in a global registry at start-up", D: []string{"Yosys parses the C++ source at run time", "The user edits Yosys's main function", "The pass is found by scanning the file system"},
		Why: "Self-registering objects are the classic plug-in pattern: constructing the instance adds it to a name-to-object map that the shell looks up."},
	{KP: 1, Q: "A placer chooses a density penalty by name from a config file: \"bell\", \"exponential\" or \"custom\". Which pattern lets new penalties be added without editing the placer?",
		A: "An abstract penalty interface plus a factory registry keyed by name", D: []string{"A switch statement in the placer's core loop", "A global variable per penalty", "Copying the placer for each penalty"},
		Why: "The core calls the interface; the registry maps names to constructors. Adding a penalty adds a registration, not a core edit."},
	{KP: 1, Q: "What does a plug-in registry store per entry?",
		A: "A name and a way to construct the implementation (a factory)", D: []string{"A copy of the design database", "The plug-in's source code", "The plug-in's test results"},
		Why: "The registry maps names to constructors; the host creates instances on demand."},
	{KP: 1, Q: "Two plug-ins both register under the name \"route_cost\". What should the registry do?",
		A: "Refuse the duplicate or report it, rather than silently replacing one", D: []string{"Keep both and pick randomly", "Merge their code", "Rename the second one at random"},
		Why: "Silent replacement makes behavior depend on load order, which is a nondeterminism bug."},
	{KP: 2, Q: "Yosys can load an extension with `yosys -m my_plugin.so`. What kind of plug-in is that?",
		A: "A shared library loaded at run time that registers its passes", D: []string{"A Tcl script interpreted by Yosys", "A Verilog module synthesized by Yosys", "A pass compiled into the Yosys binary"},
		Why: "Dynamic loading lets users add passes without rebuilding Yosys; the loaded library's static objects register themselves."},
	{KP: 2, Q: "OpenROAD's tools (placement, CTS, routing) are compiled into one binary rather than loaded at run time. What does that choice buy?",
		A: "One consistent build where all modules agree on the database and compiler", D: []string{"The ability to add tools without recompiling", "Smaller memory use than any other design", "Independence from the C++ ABI"},
		Why: "Compiled-in modules avoid ABI drift between separately built pieces; the price is that adding a module means rebuilding."},
	{KP: 2, Q: "A plug-in built with one compiler version crashes inside a host built with another, though both compile cleanly. What is the likely cause?",
		A: "The C++ ABI (layout, name mangling, standard-library types) differs across the boundary", D: []string{"The plug-in has a syntax error", "The host's Tcl version is too new", "The design is too large"},
		Why: "C++ objects crossing a dynamic-library boundary must agree on layout and runtime. A C interface or a pinned toolchain avoids this."},
	{KP: 2, Q: "Which plug-in boundary is the most robust to compiler differences between host and plug-in?",
		A: "A C-callable interface passing plain structs and opaque handles", D: []string{"Passing std::map objects by reference", "Passing C++ exceptions across the boundary", "Sharing class hierarchies with virtual inheritance"},
		Why: "The C ABI is stable across compilers; standard-library types and exceptions are not."},
	{KP: 3, Q: "A host adds a field to the struct it passes to plug-ins. Old plug-ins now read the wrong offsets. What prevents this?",
		A: "A version or size field in the struct, and adding fields only at the end", D: []string{"Recompiling the host only", "Renaming the struct", "Making all fields private"},
		Why: "Versioned, append-only structs let old plug-ins keep working and let the host detect what each plug-in understands."},
	{KP: 3, Q: "A plug-in asks the host for an API version it does not provide. When should the host find out?",
		A: "At load time, by checking the plug-in's declared interface version", D: []string{"When the plug-in first crashes", "At the end of the flow", "Never; versions are advisory"},
		Why: "Checking a declared version at load turns a late, confusing crash into an immediate, clear error."},
	{KP: 3, Q: "A plug-in keeps hidden global state between calls. Why does that hurt a host that runs it on many designs in parallel?",
		A: "Calls interfere through the shared state, so results depend on scheduling", D: []string{"Global state is forbidden by C++", "It makes the plug-in compile slower", "It prevents the host from loading the plug-in"},
		Why: "Plug-ins should take their state as explicit arguments so instances are independent and results reproducible."},
	{KP: 3, Q: "Why do tool vendors keep the plug-in interface small, even when users ask for access to internals?",
		A: "Every exposed function becomes a promise that constrains future changes to the core", D: []string{"Small interfaces compile to faster code in every case", "Large interfaces cannot be documented", "Users cannot call more than ten functions"},
		Why: "A wide interface freezes internals. A narrow, stable one lets the core evolve behind it."},
}
