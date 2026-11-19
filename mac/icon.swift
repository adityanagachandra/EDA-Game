// Draws the app icon: a dark rounded square with a small standard-cell grid and "EDA".
// Usage: swift mac/icon.swift out.png
import AppKit

let size = 1024.0
let image = NSImage(size: NSSize(width: size, height: size))
image.lockFocus()
let inset = 100.0, body = NSRect(x: inset, y: inset, width: size - 2 * inset, height: size - 2 * inset)
let tile = NSBezierPath(roundedRect: body, xRadius: 185, yRadius: 185)
NSGradient(starting: NSColor(calibratedRed: 0.10, green: 0.12, blue: 0.16, alpha: 1),
           ending: NSColor(calibratedRed: 0.03, green: 0.04, blue: 0.06, alpha: 1))!.draw(in: tile, angle: -90)

// Rows of cells, like a placed layout.
let accent = NSColor(calibratedRed: 0.29, green: 0.56, blue: 1.0, alpha: 1)
let widths: [[Double]] = [[90, 150, 60, 120, 110], [140, 70, 130, 90, 100], [60, 110, 150, 80, 120]]
for (r, row) in widths.enumerated() {
    var x = 210.0
    let y = 250.0 + Double(r) * 70
    for (i, w) in row.enumerated() {
        let cell = NSBezierPath(roundedRect: NSRect(x: x, y: y, width: w, height: 48), xRadius: 8, yRadius: 8)
        accent.withAlphaComponent((r + i) % 3 == 0 ? 0.9 : 0.35).setFill()
        cell.fill()
        x += w + 14
    }
}

let text = NSAttributedString(string: "EDA", attributes: [
    .font: NSFont.systemFont(ofSize: 250, weight: .bold),
    .foregroundColor: NSColor.white,
    .kern: 6,
])
let t = text.size()
text.draw(at: NSPoint(x: (size - t.width) / 2, y: 520))
image.unlockFocus()

let rep = NSBitmapImageRep(data: image.tiffRepresentation!)!
try! rep.representation(using: .png, properties: [:])!.write(to: URL(fileURLWithPath: CommandLine.arguments[1]))
