// EDA-Game for macOS: the static reader in site/ inside a native window.
//
// The page is served from the app bundle through a custom eda:// scheme rather than file://,
// because the reader fetches its lesson data with fetch(), which WebKit refuses for file URLs.
// Progress and notes live in the page's localStorage, which WebKit keeps per app.
import Cocoa
import UniformTypeIdentifiers
import WebKit

let scheme = "eda"
let siteRoot = Bundle.main.resourceURL!.appendingPathComponent("site")

final class BundleSchemeHandler: NSObject, WKURLSchemeHandler {
    func webView(_ webView: WKWebView, start task: WKURLSchemeTask) {
        guard let url = task.request.url else { return }
        var path = url.path
        if path.isEmpty || path == "/" { path = "/index.html" }
        let file = siteRoot.appendingPathComponent(path).standardizedFileURL
        guard file.path.hasPrefix(siteRoot.path), let data = try? Data(contentsOf: file) else {
            task.didFailWithError(URLError(.fileDoesNotExist))
            return
        }
        let type = UTType(filenameExtension: file.pathExtension)?.preferredMIMEType ?? "application/octet-stream"
        let response = HTTPURLResponse(url: url, statusCode: 200, httpVersion: "HTTP/1.1",
                                       headerFields: ["Content-Type": type, "Content-Length": "\(data.count)"])!
        task.didReceive(response)
        task.didReceive(data)
        task.didFinish()
    }

    func webView(_ webView: WKWebView, stop task: WKURLSchemeTask) {}
}

final class AppDelegate: NSObject, NSApplicationDelegate, WKNavigationDelegate, WKUIDelegate {
    var window: NSWindow!
    var web: WKWebView!

    func applicationDidFinishLaunching(_ note: Notification) {
        let config = WKWebViewConfiguration()
        config.setURLSchemeHandler(BundleSchemeHandler(), forURLScheme: scheme)
        config.websiteDataStore = .default()
        web = WKWebView(frame: .zero, configuration: config)
        web.navigationDelegate = self
        web.uiDelegate = self
        web.allowsBackForwardNavigationGestures = true

        window = NSWindow(contentRect: NSRect(x: 0, y: 0, width: 1200, height: 820),
                          styleMask: [.titled, .closable, .miniaturizable, .resizable, .fullSizeContentView],
                          backing: .buffered, defer: false)
        window.title = "EDA-Game"
        window.minSize = NSSize(width: 420, height: 480)
        window.contentView = web
        window.center()
        window.setFrameAutosaveName("EDA-Game main window")
        window.makeKeyAndOrderFront(nil)

        buildMenu()
        web.load(URLRequest(url: URL(string: "\(scheme)://app/index.html")!))
        NSApp.activate()
    }

    func applicationShouldTerminateAfterLastWindowClosed(_ app: NSApplication) -> Bool { true }

    // Links that leave the reader open in the default browser.
    func webView(_ webView: WKWebView, decidePolicyFor action: WKNavigationAction,
                 decisionHandler: @escaping (WKNavigationActionPolicy) -> Void) {
        if let url = action.request.url, url.scheme != scheme, url.scheme != "about" {
            NSWorkspace.shared.open(url)
            decisionHandler(.cancel)
        } else {
            decisionHandler(.allow)
        }
    }

    // target="_blank" links.
    func webView(_ webView: WKWebView, createWebViewWith config: WKWebViewConfiguration,
                 for action: WKNavigationAction, windowFeatures: WKWindowFeatures) -> WKWebView? {
        if let url = action.request.url { NSWorkspace.shared.open(url) }
        return nil
    }

    func webView(_ webView: WKWebView, didFinish navigation: WKNavigation!) {
        // Keep the window title in step with the page.
        if let t = webView.title, !t.isEmpty { window.title = t }
    }

    @objc func goHome(_ sender: Any?) { web.evaluateJavaScript("location.hash = '#/'") }
    @objc func goGlossary(_ sender: Any?) { web.evaluateJavaScript("location.hash = '#/glossary'") }
    @objc func goNotes(_ sender: Any?) { web.evaluateJavaScript("location.hash = '#/notes'") }
    @objc func back(_ sender: Any?) { web.goBack() }
    @objc func forward(_ sender: Any?) { web.goForward() }
    @objc func reload(_ sender: Any?) { web.reload() }
    @objc func zoomIn(_ sender: Any?) { web.pageZoom = min(web.pageZoom + 0.1, 2.5) }
    @objc func zoomOut(_ sender: Any?) { web.pageZoom = max(web.pageZoom - 0.1, 0.5) }
    @objc func zoomReset(_ sender: Any?) { web.pageZoom = 1 }

    func buildMenu() {
        let main = NSMenu()
        func submenu(_ title: String, _ items: [NSMenuItem]) {
            let holder = NSMenuItem()
            let menu = NSMenu(title: title)
            items.forEach(menu.addItem)
            holder.submenu = menu
            main.addItem(holder)
        }
        func item(_ title: String, _ action: Selector, _ key: String, _ mods: NSEvent.ModifierFlags = .command) -> NSMenuItem {
            let i = NSMenuItem(title: title, action: action, keyEquivalent: key)
            i.keyEquivalentModifierMask = mods
            return i
        }
        submenu("EDA-Game", [
            item("About EDA-Game", #selector(NSApplication.orderFrontStandardAboutPanel(_:)), ""),
            .separator(),
            item("Hide EDA-Game", #selector(NSApplication.hide(_:)), "h"),
            item("Hide Others", #selector(NSApplication.hideOtherApplications(_:)), "h", [.command, .option]),
            item("Show All", #selector(NSApplication.unhideAllApplications(_:)), ""),
            .separator(),
            item("Quit EDA-Game", #selector(NSApplication.terminate(_:)), "q"),
        ])
        submenu("Edit", [
            item("Undo", Selector(("undo:")), "z"),
            item("Redo", Selector(("redo:")), "z", [.command, .shift]),
            .separator(),
            item("Cut", #selector(NSText.cut(_:)), "x"),
            item("Copy", #selector(NSText.copy(_:)), "c"),
            item("Paste", #selector(NSText.paste(_:)), "v"),
            item("Select All", #selector(NSText.selectAll(_:)), "a"),
        ])
        submenu("View", [
            item("Actual Size", #selector(zoomReset(_:)), "0"),
            item("Zoom In", #selector(zoomIn(_:)), "+"),
            item("Zoom Out", #selector(zoomOut(_:)), "-"),
            .separator(),
            item("Reload", #selector(reload(_:)), "r"),
            .separator(),
            item("Enter Full Screen", #selector(NSWindow.toggleFullScreen(_:)), "f", [.command, .control]),
        ])
        submenu("Go", [
            item("Back", #selector(back(_:)), "["),
            item("Forward", #selector(forward(_:)), "]"),
            .separator(),
            item("Domains", #selector(goHome(_:)), "1"),
            item("Glossary", #selector(goGlossary(_:)), "2"),
            item("Notes", #selector(goNotes(_:)), "3"),
        ])
        submenu("Window", [
            item("Minimize", #selector(NSWindow.performMiniaturize(_:)), "m"),
            item("Close", #selector(NSWindow.performClose(_:)), "w"),
        ])
        NSApp.mainMenu = main
    }
}

let app = NSApplication.shared
let delegate = AppDelegate()
app.delegate = delegate
app.setActivationPolicy(.regular)
app.run()
