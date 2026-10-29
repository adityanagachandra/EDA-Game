# Optical proximity correction

The image of a drawn rectangle is not a rectangle: corners round, line ends pull back, and widths shift with their neighbours. Optical proximity correction pre-distorts the mask so the printed shape lands on target. It is one of the largest compute jobs in the whole design-to-silicon chain and the main reason mask data looks nothing like the layout.

## Why printed shapes differ from drawn ones

The projection lens passes only low spatial frequencies, so sharp features blur. Corners lose intensity because they have neighbours on fewer sides, so they round. Line ends shorten, because the end of a line receives light from only one side. A line's printed width depends on what lies within the optical interaction range of a few hundred nanometres, so identical drawn lines print differently in dense and sparse contexts. Resist and etch add their own biases. Without correction, a gate drawn at 20 nm could print several nanometres off, and a line end could pull back enough to miss its via.

## Rule-based and model-based OPC

Rule-based OPC applies table-driven fixes: bias each edge by an amount that depends on its width and its space to the neighbour, add hammerheads to line ends and serifs to corners. It is fast but cannot capture complex two-dimensional context. Model-based OPC fragments each edge into short segments, simulates the aerial image and resist response with a calibrated compact model at a control site on each fragment, measures the edge placement error (EPE, printed minus target position), and moves the fragments to reduce it. It iterates until EPE is within tolerance. After correction, a verification run simulates the result across dose and focus to find hotspots that pinch or bridge.

## Inverse lithography and cost

Inverse lithography technology (ILT) poses the mask as the unknown of an optimization: find the mask pixels whose simulated print best matches the target across a process window. The solution includes assist features naturally and is often curvilinear. It gives better process window than fragment-moving OPC but costs far more compute, and curvilinear masks were slow on shaped-beam writers. GPUs and multi-beam writers have moved ILT from hotspots and memory arrays toward wider use. For design, OPC means context dependence: the same cell prints differently in different placements, and regular layouts reduce that spread.

## Going deeper

Cobb, "Fast Optical and Process Proximity Correction Algorithms for Integrated Circuit Manufacturing" (PhD thesis, UC Berkeley, 1998); Ma and Arce, "Computational Lithography" (2010).
