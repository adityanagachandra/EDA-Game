"""
Extract lesson content from Discrete Mathematics: An Open Introduction, 3rd ed (CC BY-SA 4.0)

PreTeXt puts teaching prose in <div class="para">, not <p>. The first version of
this script only converted p, li, headings and tables, so every explanatory
paragraph was dropped and the lessons kept their headings and exercises but no
teaching text. html2text reads the whole subtree, which is why this now uses it.
Math is moved out of the tree before conversion because html2text doubles each
backslash, and a doubled backslash breaks every LaTeX delimiter.
"""
import json
import re
import time
from pathlib import Path

import html2text
import requests
from bs4 import BeautifulSoup

BASE_URL = "https://discrete.openmathbooks.org/dmoi3"
OUTPUT_DIR = Path("data/lessons/teaching")

SLUG_MAP = {
    "discrete.logic.connectives": ("sec_propositional", "Propositional Logic"),
    "discrete.logic.truth_tables": ("sec_propositional", "Truth Tables"),
    "discrete.logic.quantifiers": ("sec_intro-statements", "Predicates and Quantifiers"),
    "discrete.sets.venn": ("sec_intro-sets", "Venn Diagrams"),
    "discrete.graphs.basics": ("sec_gt-intro", "Graph Theory Basics"),
    # discrete.graphs.paths is deliberately absent. Its lesson is hand-written
    # around walks, components and distance, which is what its generators and
    # its knowledge-point shard cover. Levin's sec_paths page is Euler circuits,
    # so running this script over it would replace a matching lesson with one
    # whose headings the shard cannot find.
    "discrete.graphs.trees": ("sec_trees", "Trees"),
    "discrete.recurrence": ("sec_recurrence", "Solving Recurrence Relations"),
}

MATH_TOKEN = "MATHUAPLACEHOLDER{}END"


def fetch_page(slug):
    """Fetch a section page and return its cleaned main content."""
    url = f"{BASE_URL}/{slug}.html"
    print(f"  Fetching {url}")
    resp = requests.get(url, timeout=30)
    resp.raise_for_status()
    soup = BeautifulSoup(resp.content, "html.parser")

    main = soup.select_one(".ptx-content") or soup.find("main") or soup.body
    if main is None:
        raise ValueError("No content container found")

    for tag in main.find_all(["script", "style", "nav", "aside", "noscript"]):
        tag.decompose()
    # Permalink anchors hold a link emoji, and the exercises are practice, not
    # teaching text. Both are removed before conversion so neither reaches a lesson.
    for cls in ("autopermalink", "exercises", "solution", "solutions", "hint", "answer", "colophon"):
        for tag in main.find_all(class_=cls):
            tag.decompose()
    return main, url


# PreTeXt defines these in its book preamble, so they never reach KaTeX. The
# renderer knows \\imp, \\iff and the blackboard letters, but not these three.
PRETEXT_MACROS = (
    (re.compile(r"\\Iff(?![A-Za-z])"), r"\\iff"),
    (re.compile(r"\\isom(?![A-Za-z])"), r"\\cong"),
    (re.compile(r"\\U(?![A-Za-z])"), r"\\mathcal{U}"),
)

DISPLAY_ENV = re.compile(
    r"^\\begin\{(equation\*?|displaymath|align\*?|gather\*?|multline\*?)\}([\s\S]*)\\end\{\1\}$"
)


def normalize_math(raw):
    """Return TeX the lesson renderer can parse.

    A display element carries its own environment, for example
    \\begin{equation*}...\\end{equation*}, which the renderer treats as raw prose.
    Its body is valid display math, so it becomes a $$ block instead.
    """
    tex = raw.strip()
    match = DISPLAY_ENV.match(tex)
    if match:
        tex = "$$\n" + match.group(2).strip() + "\n$$"
    elif tex.startswith("\\[") and tex.endswith("\\]"):
        tex = "$$\n" + tex[2:-2].strip() + "\n$$"
    else:
        # html2text hard-wraps inside a span, and a newline in an inline math
        # span breaks the delimiter pass in the renderer. TeX treats it as
        # whitespace anyway.
        tex = re.sub(r"\s*\n\s*", " ", tex)
    for pattern, replacement in PRETEXT_MACROS:
        tex = pattern.sub(replacement, tex)
    return tex


def protect_math(node):
    """Replace every math element with a placeholder and return its raw text.

    html2text escapes a backslash as two, which turns \\( into \\\\( and breaks the
    delimiters the renderer looks for. Holding the math out of the conversion and
    putting it back verbatim is simpler and safer than undoing the escape later,
    because real LaTeX line breaks are also two backslashes.
    """
    saved = []
    for el in list(node.select(".process-math, .displaymath")):
        if el.parent is None:
            continue
        if el.find_parent(class_="process-math") or el.find_parent(class_="displaymath"):
            continue
        token = MATH_TOKEN.format(len(saved))
        saved.append(normalize_math(el.get_text()))
        el.replace_with(token)
    return saved


def restore_math(text, saved):
    for i, raw in enumerate(saved):
        text = text.replace(MATH_TOKEN.format(i), raw)
    return text


def find_section(main, section_title):
    """Return the section that holds the named heading, or None.

    Matching is by the heading's own .title span because PreTeXt splits the label
    ("Section 4.1" and "Definitions") across spans. A title that does not exist
    on the page falls back to the whole page, which is correct for a concept that
    names a chapter rather than a subsection.
    """
    for heading in main.find_all(["h1", "h2", "h3", "h4", "h5", "h6"]):
        title = heading.find(class_="title")
        if title and title.get_text(" ", strip=True) == section_title:
            return heading.find_parent("section"), heading
    return None, None


def to_markdown(node):
    converter = html2text.HTML2Text()
    converter.body_width = 0
    converter.ignore_links = True
    converter.ignore_images = True
    converter.ignore_emphasis = False

    lines = converter.handle(str(node)).split("\n")
    cleaned = []
    blank = False
    for line in lines:
        line = line.rstrip()
        if not line:
            if not blank:
                cleaned.append("")
            blank = True
        else:
            cleaned.append(line)
            blank = False
    return "\n".join(cleaned).strip()


def extract_section(main, section_title):
    # Two concepts can share one page, so every mutation happens on a copy.
    # Mutating main in place left the first concept's math placeholders in the
    # tree for the second concept, which then found no math to restore.
    #
    # The whole page is kept, not just the named section. Knowledge-point shards
    # cite headings throughout the page ("Subsection Truth Tables", "Section 0.2
    # Mathematical Statements"), so a scoped lesson orphans those shards.
    node = BeautifulSoup(str(main), "html.parser")

    saved = protect_math(node)
    body = restore_math(to_markdown(node), saved)
    # A bare numeral set as inline math is re-escaped by the fixed-point pass
    # in the API and reaches the reader as literal $6$. Numerals are upright
    # text, so unwrap them here.
    body = re.sub(r"\\\((\d+)\\\)", r"\1", body)
    body = re.sub(r"\n{3,}", "\n\n", body).strip()
    return f"# {section_title}\n\n{body}\n"


def main():
    with open("data/lessons/mapping.json") as f:
        mapping = json.load(f)

    dm_concepts = {k: v for k, v in mapping.items() if v.get("source") == "discrete_math"}
    print(f"Discrete Math concepts: {len(dm_concepts)}")

    slug_groups = {}
    for cid in sorted(dm_concepts.keys()):
        if cid not in SLUG_MAP:
            print(f"  SKIP {cid}: no slug mapping")
            continue
        url_slug, section_name = SLUG_MAP[cid]
        slug_groups.setdefault(url_slug, []).append((cid, section_name))

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    for url_slug, concepts in slug_groups.items():
        print(f"\nFetching {url_slug}...")
        try:
            main, url = fetch_page(url_slug)
        except Exception as e:
            print(f"  FAILED: {e}")
            continue

        for cid, section_name in concepts:
            print(f"  Extracting {cid} ({section_name})...")
            md = extract_section(main, section_name)
            attribution = (
                f"> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed]({url})"
                " by Oscar Levin, CC BY-SA 4.0\n\n"
            )
            out_path = OUTPUT_DIR / f"{cid}.md"
            out_path.write_text(attribution + md)
            print(f"    -> Wrote {out_path} ({len(md)} chars)")
        time.sleep(1)


if __name__ == "__main__":
    main()
