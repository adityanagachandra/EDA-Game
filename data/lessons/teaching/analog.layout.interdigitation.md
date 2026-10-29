# Interdigitated layout and gradient cancellation

When matched devices must sit in a single row, as fingers of multi-finger transistors usually do, the layout choice is the order of fingers. Interdigitation alternates the devices; mirrored patterns such as ABBA make the order symmetric. Which patterns cancel which gradients is a small computation on finger positions, and it is exactly what a primitive generator in an analog layout tool must get right.

## Finger positions and centroids in one row

Number the fingers along the row and measure each finger's position from the row centre, so a row of $n$ fingers has positions $-(n-1)/2, \dots, (n-1)/2$. A device's centroid is the mean position of its fingers. For "AABB" the positions are $-1.5, -0.5, 0.5, 1.5$: A's centroid is -1, B's is $+1$, a separation of 2 pitches. For "ABAB", A sits at -1.5 and 0.5, centroid -0.5; B at $+0.5$; separation 1. For "ABBA", A at $\pm 1.5$ and B at $\pm 0.5$: both centroids are 0. Plain interdigitation halves the separation of the side-by-side layout, and the mirrored pattern removes it. Dummy fingers at both ends do not change centroids; they equalize the surroundings of the outer fingers.

## Linear gradients

Under a linear gradient of $g$ millivolts per pitch along the row, each finger's threshold shifts by $g$ times its position, and a device's mean shift is $g$ times its centroid. The systematic mismatch between A and B is therefore $g\,|\bar{x}_A - \bar{x}_B|$. With $g = 0.4$ mV per pitch, "AABB" gives 0.8 mV, "ABAB" gives 0.4 mV, and "ABBA" gives zero. Compare this with random mismatch of perhaps 1 mV for the same devices: an unlucky order can create a systematic error as large as the random one, and unlike random mismatch it does not shrink with area. The calculation is cheap enough that a placer can evaluate it inside its inner loop for every candidate order.

## Higher-order gradients

Real gradients are not perfectly linear. A quadratic component $c\,x^2$, for example from stress near a well edge or a hot spot centred on the array, shifts each device by $c$ times the mean of $x^2$ over its fingers. In "ABBA", A's fingers at $\pm 1.5$ give mean $x^2 = 2.25$ and B's at $\pm 0.5$ give 0.25, so with $c = 0.1$ mV per pitch squared the mismatch is 0.2 mV, even though the linear term cancels. Patterns such as "ABBAABBA" or "ABABBABA" spread both devices more evenly and reduce this. The general lesson for a tool designer: a matching pattern is characterized by which moments of the position distribution it equalizes, and a generator should choose the pattern by the gradients it must cancel, not by habit.

## Going deeper

Hastings, "The Art of Analog Layout" (2nd ed., 2005), chapters on device matching.
