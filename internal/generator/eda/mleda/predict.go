package mleda

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// overlap1 is the overlap length of [a0, a1] and [b0, b1].
func overlap1(a0, a1, b0, b1 float64) float64 {
	return math.Max(0, math.Min(a1, b1)-math.Max(a0, b0))
}

// congestion: RUDY density of one net (KP 1), RUDY of a tile from overlapping nets (KP 2), and
// the NRMSE of a predicted congestion map (KP 3).
func congestion(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		w := float64(edakit.Int(2, 40))
		h := float64(edakit.Int(2, 40))
		d := (w + h) / (w * h)
		ans := edakit.Num(d, 4)
		return edakit.Numeric(edakit.Sprintf("A placed net has a bounding box %s um wide and %s um tall. RUDY spreads the net's estimated wirelength (its half-perimeter) uniformly over the box. What wire density does the net add at every point inside the box, in um of wire per square um? Round to 4 decimal places.", f(w, 0), f(h, 0)),
			ans, edakit.Sprintf("$D = (w + h)/(w h) = %s / %s = %s$ per um. Small boxes give high density, which is how RUDY marks local pin clusters as congested.", f(w+h, 0), f(w*h, 0), ans),
			map[string]string{"area": f(w*h, 0)})
	case 2:
		t := float64(edakit.Pick([]int{10, 20}))
		n := edakit.Scaled(ctx, 2, 4)
		total := 0.0
		var desc, terms []string
		for i := 0; i < n; i++ {
			var x0, y0, w, h, ov float64
			for tries := 0; ov <= 0; tries++ {
				if tries >= generator.MaxDrawAttempts {
					x0, y0, w, h = 0, 0, t, t
					ov = t * t
					break
				}
				x0 = float64(edakit.Int(-10, int(t)-2))
				y0 = float64(edakit.Int(-10, int(t)-2))
				w = float64(edakit.Int(4, 24))
				h = float64(edakit.Int(4, 24))
				ov = overlap1(x0, x0+w, 0, t) * overlap1(y0, y0+h, 0, t)
			}
			d := (w + h) / (w * h)
			total += d * ov
			desc = append(desc, edakit.Sprintf("net %d spans x from %s to %s and y from %s to %s", i+1, f(x0, 0), f(x0+w, 0), f(y0, 0), f(y0+h, 0)))
			terms = append(terms, edakit.Sprintf("%s \\cdot %s", f(d, 4), f(ov, 0)))
		}
		r := total / (t * t)
		ans := edakit.Num(r, 4)
		return edakit.Numeric(edakit.Sprintf("A congestion tile covers x and y from 0 to %s um. Bounding boxes, in um: %s. Using RUDY (each net's density $(w+h)/(wh)$ times its overlap area with the tile), what is the average wire density in the tile, in um per square um? Round to 4 decimal places.", f(t, 0), strings.Join(desc, "; ")),
			ans, edakit.Sprintf("Sum of density times overlap area: $%s = %s$ um of wire. Divide by the tile area %s: %s per um.", strings.Join(terms, " + "), f(total, 4), f(t*t, 0), ans),
			map[string]string{"wire": f(total, 4)})
	default:
		n := 9
		act := make([]float64, n)
		pred := make([]float64, n)
		for i := range act {
			act[i] = float64(edakit.Int(0, 12))
			pred[i] = math.Max(0, act[i]+float64(edakit.Int(-3, 3)))
		}
		act[edakit.Int(0, n-1)] = float64(edakit.Int(13, 20))
		mx, mn := act[0], act[0]
		se := 0.0
		for i := range act {
			mx = math.Max(mx, act[i])
			mn = math.Min(mn, act[i])
			se += (pred[i] - act[i]) * (pred[i] - act[i])
		}
		rmse := math.Sqrt(se / float64(n))
		r := rmse / (mx - mn)
		ans := edakit.Num(r, 3)
		return edakit.Numeric(edakit.Sprintf("A congestion model predicts the overflow of a 3 by 3 tile window. Read row by row, the router's actual overflows are %s and the predictions are %s. Report NRMSE, the root-mean-square error divided by the range (max minus min) of the actual values. Round to 3 decimal places.", floats(act, 0), floats(pred, 0)),
			ans, edakit.Sprintf("Squared errors sum to %s over 9 tiles, so RMSE $= \\sqrt{%s/9} = %s$. Range $= %s - %s = %s$. NRMSE $= %s$. A low NRMSE can still miss the one hot tile that matters.", f(se, 0), f(se, 0), f(rmse, 4), f(mx, 0), f(mn, 0), f(mx-mn, 0), ans),
			map[string]string{"rmse": f(rmse, 4)})
	}
}

// timingPrediction: MAE (KP 1), R squared (KP 2), and violation-status flips (KP 3) of a
// pre-route endpoint-slack predictor.
func timingPrediction(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	act := make([]float64, n)
	pred := make([]float64, n)
	for i := range act {
		act[i] = float64(edakit.Int(-40, 60))
		if act[i] == 0 {
			act[i] = 3
		}
		pred[i] = act[i] + float64(edakit.Int(-18, 18))
		if pred[i] == 0 {
			pred[i] = -2
		}
	}
	stem := edakit.Sprintf("A GNN predicts post-route setup slack at %d endpoints from the placed netlist. Signoff STA after routing gives actual slacks %s ps; the model predicted %s ps, in the same order.", n, floats(act, 0), floats(pred, 0))
	switch kp(ctx) {
	case 1:
		s := 0.0
		for i := range act {
			s += math.Abs(pred[i] - act[i])
		}
		ans := edakit.Num(s/float64(n), 2)
		return edakit.Numeric(stem+" What is the mean absolute error, in ps? Round to 2 decimal places.",
			ans, edakit.Sprintf("Absolute errors sum to %s ps over %d endpoints: MAE $= %s$ ps.", f(s, 0), n, ans), map[string]string{"sum_abs": f(s, 0)})
	case 2:
		mean := 0.0
		for _, a := range act {
			mean += a
		}
		mean /= float64(n)
		ssr, sst := 0.0, 0.0
		for i := range act {
			ssr += (act[i] - pred[i]) * (act[i] - pred[i])
			sst += (act[i] - mean) * (act[i] - mean)
		}
		r2 := 1 - ssr/sst
		ans := edakit.Num(r2, 3)
		return edakit.Numeric(stem+" What is the coefficient of determination $R^2 = 1 - \\sum(y - \\hat{y})^2 / \\sum(y - \\bar{y})^2$? Round to 3 decimal places.",
			ans, edakit.Sprintf("$\\bar{y} = %s$. Residual sum of squares %s, total sum of squares %s. $R^2 = 1 - %s/%s = %s$. A high $R^2$ across all endpoints says little about accuracy near zero slack, where decisions are made.", f(mean, 3), f(ssr, 0), f(sst, 3), f(ssr, 0), f(sst, 3), ans),
			map[string]string{"ssr": f(ssr, 0), "sst": f(sst, 3)})
	default:
		flips := 0
		for i := range act {
			if (act[i] < 0) != (pred[i] < 0) {
				flips++
			}
		}
		return edakit.Numeric(stem+" At how many endpoints does the prediction get the violation status wrong (predicts passing when the endpoint fails, or failing when it passes)?",
			edakit.IntStr(flips), edakit.Sprintf("Compare signs endpoint by endpoint: %d endpoints flip. These are the errors that change what an optimizer does, regardless of the average error.", flips),
			map[string]string{"flips": edakit.IntStr(flips)})
	}
}

// drcHotspot: precision (KP 1), F1 (KP 2), and precision at a deployment prevalence (KP 3).
func drcHotspot(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1, 2:
		tp := edakit.Scaled(ctx, 10, 120)
		fp := edakit.Int(5, 150)
		fn := edakit.Int(5, 100)
		tiles := edakit.Pick([]int{10000, 40000, 250000})
		stem := edakit.Sprintf("A DRC hotspot predictor scores %d gcells of a routed block. Against the DRC engine's report, it flags %d true hotspots correctly, raises %d false alarms and misses %d real hotspots.", tiles, tp, fp, fn)
		if k == 1 {
			p := 100 * float64(tp) / float64(tp+fp)
			ans := edakit.Num(p, 2)
			return edakit.Numeric(stem+" What is its precision, in percent? Round to 2 decimal places.",
				ans, edakit.Sprintf("Precision $= TP/(TP+FP) = %d/%d = %s$ percent: the share of flagged gcells that are real.", tp, tp+fp, ans), map[string]string{"flagged": edakit.IntStr(tp + fp)})
		}
		pr := float64(tp) / float64(tp+fp)
		rc := float64(tp) / float64(tp+fn)
		f1 := 2 * pr * rc / (pr + rc)
		ans := edakit.Num(f1, 3)
		return edakit.Numeric(stem+" What is its F1 score (as a fraction between 0 and 1)? Round to 3 decimal places.",
			ans, edakit.Sprintf("Precision $= %d/%d = %s$, recall $= %d/%d = %s$. $F_1 = 2PR/(P+R) = 2TP/(2TP+FP+FN) = %d/%d = %s$. True negatives never enter, which is why F1 suits rare hotspots.", tp, tp+fp, f(pr, 4), tp, tp+fn, f(rc, 4), 2*tp, 2*tp+fp+fn, ans),
			map[string]string{"precision": f(pr, 4), "recall": f(rc, 4)})
	default:
		r := float64(edakit.Int(70, 95)) / 100
		fpr := float64(edakit.Int(1, 50)) / 1000
		prev := edakit.Pick([]float64{0.001, 0.002, 0.005, 0.01, 0.02, 0.05})
		p := r * prev / (r*prev + fpr*(1-prev))
		ans := edakit.Num(100*p, 2)
		return edakit.Numeric(edakit.Sprintf("A hotspot classifier has recall %s and a false-positive rate of %s on non-hotspot gcells. On a new design, %s percent of gcells are real hotspots. What precision should the team expect, in percent? Round to 2 decimal places.", f(r, 2), f(fpr, 3), f(100*prev, 1)),
			ans, edakit.Sprintf("Per gcell: true flags $= %s \\times %s = %s$, false flags $= %s \\times %s = %s$. Precision $= %s/(%s + %s) = %s$ percent. At low prevalence even a small false-positive rate swamps the true hotspots.", f(r, 2), f(prev, 3), f(r*prev, 6), f(fpr, 3), f(1-prev, 3), f(fpr*(1-prev), 6), f(r*prev, 6), f(r*prev, 6), f(fpr*(1-prev), 6), ans),
			map[string]string{"true_flags": f(r*prev, 6)})
	}
}

// irDrop: the ground-truth drop along a strap (KP 1), the error in predicted worst drop (KP 2),
// and the recall of predicted hotspots (KP 3).
func irDrop(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1:
		n := edakit.Scaled(ctx, 3, 6)
		r := float64(edakit.Int(1, 5)) / 10
		cur := make([]float64, n)
		for i := range cur {
			cur[i] = float64(edakit.Int(2, 30))
		}
		drop := 0.0
		var terms []string
		for j := 0; j < n; j++ {
			down := 0.0
			for i := j; i < n; i++ {
				down += cur[i]
			}
			drop += r * down
			terms = append(terms, edakit.Sprintf("%s \\cdot %s", f(r, 1), f(down, 0)))
		}
		ans := edakit.Num(drop, 2)
		return edakit.Numeric(edakit.Sprintf("To label an IR-drop dataset, a solver computes the drop along a power strap fed from a pad at one end. The strap has %d taps; each segment (pad to tap 1, and between consecutive taps) has resistance %s ohm. The taps draw %s mA, from the pad end outward. What is the IR drop at the last tap, in mV? Round to 2 decimal places.", n, f(r, 1), floats(cur, 0)),
			ans, edakit.Sprintf("Each segment carries the current of every tap beyond it. Drop $= \\sum_j R \\cdot I_{\\ge j} = %s = %s$ mV (ohm times mA is mV).", strings.Join(terms, " + "), ans),
			map[string]string{"total_current": f(sumF(cur), 0)})
	default:
		n := edakit.Scaled(ctx, 6, 10)
		act := make([]float64, n)
		pred := make([]float64, n)
		for i := range act {
			act[i] = float64(edakit.Int(10, 60))
			pred[i] = math.Max(1, act[i]+float64(edakit.Int(-12, 12)))
		}
		thr := float64(edakit.Pick([]int{35, 40, 45}))
		act[edakit.Int(0, n-1)] = thr + float64(edakit.Int(3, 20))
		stem := edakit.Sprintf("A CNN predicts static IR drop per tile from power and power-grid feature maps. For %d tiles, the solver's drops are %s mV and the model's predictions are %s mV, in the same order.", n, floats(act, 0), floats(pred, 0))
		if k == 2 {
			ma, mp := maxF(act), maxF(pred)
			e := 100 * math.Abs(mp-ma) / ma
			ans := edakit.Num(e, 2)
			return edakit.Numeric(stem+" What is the percentage error of the predicted worst-case drop relative to the true worst-case drop, regardless of which tile each occurs in? Round to 2 decimal places.",
				ans, edakit.Sprintf("True worst %s mV, predicted worst %s mV. Error $= |%s - %s|/%s = %s$ percent. The two maxima need not sit on the same tile, so a small error here can hide a misplaced hotspot.", f(ma, 0), f(mp, 0), f(mp, 0), f(ma, 0), f(ma, 0), ans),
				map[string]string{"true_max": f(ma, 0)})
		}
		hot, caught := 0, 0
		for i := range act {
			if act[i] > thr {
				hot++
				if pred[i] > thr {
					caught++
				}
			}
		}
		ans := edakit.Num(100*float64(caught)/float64(hot), 2)
		return edakit.Numeric(stem+edakit.Sprintf(" A tile is a hotspot when its drop exceeds %s mV. What is the model's hotspot recall, in percent? Round to 2 decimal places.", f(thr, 0)),
			ans, edakit.Sprintf("True hotspots: %d tiles above %s mV. The model also puts %d of them above the threshold. Recall $= %d/%d = %s$ percent. A missed hotspot is the costly error: it is the one no one re-simulates.", hot, f(thr, 0), caught, caught, hot, ans),
			map[string]string{"hot": edakit.IntStr(hot)})
	}
}

func sumF(xs []float64) float64 {
	s := 0.0
	for _, x := range xs {
		s += x
	}
	return s
}

func maxF(xs []float64) float64 {
	m := xs[0]
	for _, x := range xs {
		m = math.Max(m, x)
	}
	return m
}

var datasetItems = []edakit.Item{
	{KP: 1, Q: "What does an open dataset such as CircuitNet supply to an ML-for-EDA researcher?",
		A: "Layout-derived feature maps paired with tool-computed labels such as congestion, DRC violations and IR drop", D: []string{"Trained models that replace the router", "Foundry design rules for a production node", "Hand-labeled images of fabricated chips"},
		Why: "The scarce resource is labels from real tool runs. Such datasets package many flow runs as features and labels so that models can be compared on the same data."},
	{KP: 1, Q: "Why do most ML-for-EDA datasets come from running many parameter variants of a few designs through a flow?",
		A: "Labels come only from expensive tool runs, and the number of distinct open designs is small", D: []string{"Parameter variants are required by the file formats", "A single run produces no usable features", "Designs must be anonymized by changing parameters"},
		Why: "Each label costs a placement, route or analysis run, and open RTL is limited, so variety is bought by varying flow settings."},
	{KP: 1, Q: "A dataset for synthesis-recipe learning pairs and-inverter graphs with the QoR after many recipes. What is the label?",
		A: "The area or delay that a tool reports after running each recipe", D: []string{"The recipe a human expert preferred", "The number of lines of RTL", "The graph's node count before synthesis"},
		Why: "The label is a tool's measurement of the outcome, which is exactly what the model is trained to predict."},
	{KP: 1, Q: "A congestion dataset stores feature maps such as cell density, RUDY and macro regions per gcell. What makes the problem image-to-image?",
		A: "Inputs and labels are both 2-D grids over the die, so a convolutional encoder-decoder can map one to the other", D: []string{"The labels are photographs of the die", "The netlist is stored as a bitmap", "Routing tracks are pixels in the GDS file"},
		Why: "Rasterizing features and labels to the gcell grid turns congestion prediction into dense per-pixel regression."},
	{KP: 2, Q: "A random 80/20 split of tiles from 10 designs gives a high test score. Why is that score optimistic?",
		A: "Tiles from the same design appear in both splits, so the model is tested on designs it has already seen", D: []string{"Twenty percent is too small a test set", "Random splits always overfit the test set", "Tiles are independent, so the split is fine"},
		Why: "Tiles of one design share placement, netlist style and flow settings. Testing on held-out designs measures what a tool will face."},
	{KP: 2, Q: "Labels in a dataset came from one tool version, and the team now runs a newer router with different congestion behavior. What is the risk?",
		A: "The labels describe the old tool, so the model predicts the wrong tool's outcome", D: []string{"None, because congestion is a property of the netlist alone", "The features become unreadable", "The model's parameter count changes"},
		Why: "A label is a measurement by a specific tool and setting. Change the tool and the target distribution shifts."},
	{KP: 2, Q: "A dataset reports DRC violation counts per tile from a commercial flow. A team trains a model and applies it in an OpenROAD flow. What should they expect?",
		A: "Degraded accuracy, because the router, rule deck and node differ from those that produced the labels", D: []string{"Identical accuracy, because DRC rules are universal", "Better accuracy, because open tools are simpler", "No predictions at all"},
		Why: "Hotspot labels depend on the router's behavior and the rule deck. Transfer across flows is a distribution shift that must be measured."},
	{KP: 2, Q: "Why is the number of distinct designs, not the number of tiles, the better measure of a dataset's size for generalization?",
		A: "Generalization to a new design depends on how many independent designs the model has seen", D: []string{"Tiles are not stored in the dataset", "Every design has the same number of tiles", "Tile counts cannot be computed"},
		Why: "Millions of correlated tiles from a handful of designs carry far less information about unseen designs than the count suggests."},
	{KP: 3, Q: "A paper reports state-of-the-art congestion accuracy on a public dataset but uses a different train-test split from earlier work. What is the problem?",
		A: "The results are not comparable to earlier numbers on the same dataset", D: []string{"There is no problem if the dataset is the same", "The model must be retrained from scratch", "The labels change when the split changes"},
		Why: "A benchmark is the data plus the split plus the metric. Changing any one makes cross-paper comparisons meaningless."},
	{KP: 3, Q: "Most open ML-for-EDA datasets use a few open or academic process kits. What should a tool designer conclude about models trained on them?",
		A: "Their accuracy on an advanced production node is unknown until it is measured there", D: []string{"They transfer to every node because physics is shared", "They are useless for any purpose", "They are certified for signoff on any node"},
		Why: "Feature-label relations such as pin density versus DRC count depend on the node's rules. Transfer has to be tested, not assumed."},
	{KP: 3, Q: "A team wants a dataset of timing labels. Which source gives labels a model can be trusted to learn?",
		A: "Signoff-quality STA on routed and extracted designs with documented constraints and corners", D: []string{"Slack estimates from the placer's internal model", "Slacks copied from different designs' reports without their SDC", "Random slacks generated to balance the classes"},
		Why: "A model inherits the quality of its labels. Labels must come from the analysis the prediction stands in for, with its setup recorded."},
	{KP: 3, Q: "Why do dataset papers in this field publish the flow scripts and parameters used to produce each sample?",
		A: "So that others can regenerate the labels and audit how each one was produced", D: []string{"Because the scripts are the model", "To satisfy the netlist file format", "So that the dataset can be compressed"},
		Why: "Labels are tool outputs. Without the recipe that produced them, a sample cannot be reproduced or checked."},
}
