# RFC: locale-independent learning with versioned presentation

Status: design-only proposal for [upstream issue #260](https://github.com/chuma-beep/mathua/issues/260), revised in response to [the maintainer's review](https://github.com/chuma-beep/mathua/issues/260#issuecomment-6450486380). No shared interface is approved yet. This PR changes documentation only: no translated content files, locale bundles, parser code or runtime behavior.

Code baseline: [upstream a87386b9](https://github.com/chuma-beep/mathua/tree/a87386b9a0548e2dd544a204bac0c35fc5822b1f), inspected on 2026-10-28. Prototype baseline: [fork 13d096ad](https://github.com/sl1mus/mathua/tree/13d096ad50f1832b43f136dc30f3fa1b28cb3a47). The document-only commit is based directly on that upstream revision, which descends from the requested eeccd223 baseline. The prototype is evidence only and is not the merge base. Proposed names and paths below describe interface seams for review, not approved API contracts.

## 1. Scope and invariants

One concept graph supports every locale. Concept IDs, prerequisites, KP ordering, difficulty, grading semantics, diagnostic evidence, XP, FIRe propagation, mastery and review state remain language-independent. A translation is neither another concept nor another distinct mathematical problem.

The first complete localized course is arithmetic. Completeness includes interface copy, all reachable exercise variants, explanations, teaching slices, reference lessons, diagnostic feedback and accessible media descriptions. Until that gate passes, label Russian coverage as partial and disclose English fallback.

No live machine translation is required to serve a lesson. Core learning remains available without an AI provider. This RFC does not redesign FIRe, mastery, scheduling or AI routing.

## 2. Existing components and proposed boundaries

| Existing component | Observation at baseline | Proposed change |
| --- | --- | --- |
| `internal/generator/generator.go` | `Problem` combines prose, answer, facts, difficulty and scaffolds | Add optional stable template/variant metadata for migrated generators. Keep the legacy English fields during migration. |
| `internal/generator/registry.go` | Single and batch paths generate problems. `explain` interpolates solution schemas using facts | Render presentations from the same generated instance. Cover both generation paths. Do not redraw to translate. |
| `internal/solutions/solutions.go`, `data/lessons/solutions/` | Authored schemas assemble explanations from instance facts | Add versioned locale schemas and explicit branch coverage. Preserve the English loader as fallback. |
| `internal/lessons/lesson.go`, `data/lessons/lessons.json`, `data/lessons/kp/` | Lessons map source documents to concepts. KP sections refer to lesson headings | Add a presentation overlay with stable teaching-block IDs. Do not resolve a translated KP by an English heading. |
| `internal/server/server.go` | Practice, study, quiz and diagnostic use different handlers | Apply one locale resolver and presentation envelope to every migrated delivery path. Keep answer material server-side. |
| `internal/engine/engine.go` | Durable study anchors are keyed by question text. `Anchor` holds expected answer, explanation and difficulty | Add instance-ID lookup behind existing anchoring. Retain legacy text lookup for old clients. |
| `web/next-app/lib/api.ts` | Study submissions carry `question` and a retry-stable `submission_id` | Add optional `problem_id`, `answer_locale` and `answer_format`. Preserve retry idempotency. |
| `LearnStepper.tsx`, study page, quiz hosts and `DiagnosticFlow.tsx` | Each surface owns answer state | Re-render copy without remounting the question host, resetting its timer or fetching another problem. |
| `internal/grader/router.go`, numeric, tuple, choice and symbolic graders | Grading dispatch already distinguishes answer types | Normalize only an explicitly declared input format before the existing grader. Keep the legacy English path unchanged. |
| `internal/diagnostic/report.go` and result components | Evidence states and reporting are evolving | Translate stable feedback codes and parameters. Do not infer evidence from translated labels. |
| `internal/scheduler/`, `internal/fsrs/`, mastery and FIRe code | Upstream already has FSRS fields and an FSRS package alongside SM-2 state | Add no locale fields or localization imports to learning algorithms. Preserve the upstream migration independently. |
| `internal/ai/provider.go` | A provider interface and separate failover infrastructure already exist | Reuse that boundary for optional assistance. Do not create a competing provider abstraction. |

A production merge needs a fresh upstream comparison. The fork prototype is not a substitute for the recent upstream engine, auth, diagnostic or idempotency changes.

## 3. Frontend internationalization and preference

Use semantic message IDs such as `learning.checkAnswer`, with separate English and Russian catalogs. Do not use English prose as the durable translation key. Catalog validation checks key parity, parameter types and plural forms. Use locale-aware formatting for interface numbers and dates, independently from mathematical input parsing.

Keep the current route structure initially. The selector changes presentation in place. A React locale provider/store must not key the application or learning host by locale.

Preference order after hydration: explicit saved browser preference, supported browser language, deployment default. Persist an explicit selection and synchronize browser tabs. For static export, render and hydrate using the same deployment default, then resolve the browser preference. Update the document language and mark fallback blocks with their actual language. No server-wide mutable locale is allowed.

Account preference synchronization is a separate optional API change. If added later, a deliberate selection in the current session wins over an older account value. Guest-to-account migration must not move learning evidence.

The prototype demonstrates in-place switching and browser persistence. Its English-key dictionary is temporary and will be replaced incrementally with semantic IDs.

## 4. Versioned curriculum overlays

Keep prerequisites, grading types and concept identity in the existing canonical definitions. Store translated presentation outside those definitions. Proposed layout:

```text
data/locales/en/concepts/<concept_id>.json
data/locales/ru/concepts/<concept_id>.json
data/locales/<locale>/templates/<template_id>.json
data/locales/<locale>/lessons/<document_id>.md
data/locales/<locale>/manifest.json
web/next-app/messages/<locale>.json
```

Each content unit declares `schema_version`, `content_id`, `locale`, `source_revision` and `translation_revision`. A concept overlay contains labels and references to stable KP/block IDs, not a second concept definition. A template variant has its own stable ID and declared placeholder contract.

A manifest records review status, validation status and provenance per unit. These are separate axes: a human review does not imply machine validation. Availability is derived from both, the source revision and the required-content inventory.

Use a content digest of normalized source content and its dependency contract rather than the entire repository commit as `source_revision`. A template depends on parameter names/types, branch IDs, option IDs, grading contract revision and relevant solution schema. A lesson depends on its blocks, captions and referenced assets. An unrelated code change does not invalidate every translation.

Source changes mark dependent translations stale. Adding a variant or teaching block makes affected coverage incomplete. A deliberate compatibility decision can preserve review status, but must be recorded and revalidated. Never silently serve stale content as validated Russian.

The proposed template contract contains `schema_version`, `content_id`, `concept_id`, a validated `locale`, `source_revision`, `translation_revision`, statement/explanation content and a typed parameter declaration. This documentation PR supplies no executable schema, translated content or locale bundle. A later schema PR supplies validated fixtures before any delivery code consumes an overlay.

Template rendering interpolates already computed facts. It does not evaluate translator-authored expressions or recompute answers. Branch-sensitive examples such as rational inequalities require distinct variants or explicit branch selectors. Upstream's correction for issue #267 illustrates why placeholder presence alone is insufficient.

Imported lessons retain source URLs, attribution, applicable license identifiers, modification notices and asset provenance. Translations inherit the applicable source-content obligations. Shared source documents are translated once and referenced by multiple concepts. Source-license review is a release gate, not an assumption that the repository's code license covers every lesson.

## 5. Mathematical identity and localized presentation

Distinguish three identities:

- `template_id` plus `variant_id` and a generator contract revision identify the authored mathematical form.
- `problem_id` is an opaque server-issued identifier for one served instance, scoped to its learner/session and canonical grading anchor.
- `variety_key` represents the same mathematical content across presentations. It excludes locale, translated strings, presentation revision and the random delivery ID.

The variety key uses canonical typed parameters and stable option semantics under the existing variety policy. Re-serving identical mathematics with a different delivery ID must not inflate diversity. Difficulty changes count only where the existing mathematical variety policy supports them.

The server anchor retains expected answer, grading type, difficulty, KP reference, canonical choice mapping and parameter/contract revision. It also retains the facts and content revision needed to render the same served instance later. A public DTO does not expose the expected answer or internal answer-bearing facts for unanswered tasks.

Conceptual interfaces, subject to review:

```text
Generate(context) -> CanonicalProblem
Present(CanonicalProblem, locale, contentRevision) -> Presentation
ResolveServedInstance(student, problemID) -> CanonicalProblem + grading anchor
NormalizeAnswer(raw, answerLocale, declaredFormat, answerContract) -> CanonicalAnswer | InputError
Grade(anchor, CanonicalAnswer) -> existing grading result
```

A presentation carries requested/effective locale, content revision, statement, stable-ID choice labels and fallback reason. Feedback is delivered after submission. A presentation lookup never calls the generator or scheduler.

Choice IDs and stored option order remain fixed during a language switch. Submit the selected option ID, not its translated label. Map it to the existing canonical expected value server-side. Legacy text-embedded choices remain English until their identity contract is migrated.

Transport choice: an authenticated, read-only per-instance presentation lookup, proposed as `GET /api/problems/{problem_id}/presentation?locale={locale}`. Generation returns the initial presentation with the opaque problem ID. Subsequent switches request a presentation for that existing instance. The endpoint resolves learner and session ownership on every read, performs no generation and writes no learning evidence. Clients ignore responses for an older selector request if a later switch has already occurred.

Before an answer is submitted, this endpoint returns statement, choices and allowed input guidance only. It never returns expected answers, solution prose, raw facts or scaffold answers that disclose the solution. After submission, a separately authorized feedback lookup can return a translated verdict and authored explanation from the stored result. That read reuses the same ownership checks and does not submit or grade again. This preserves the upstream fix that removed answers and explanations from unanswered `/api/lessons/` payloads.

The alternative locale bundle is deferred. Per-instance lookup avoids shipping unused locales and makes authorization and revision pinning explicit. Section 12 fixes its security and caching contract.

## 6. Switching with an answer already entered

Preserve the exact draft, selected option ID, host identity, elapsed-time origin, problem ID and submission ID. Do not submit, redraw or record an attempt because the selector changed.

Record the draft's input locale and representation when entry begins. Switching the interface alone does not reinterpret an existing `1,5` as another number or rewrite the draft. Keep that draft's parsing locale until it is cleared. Explain this beside the input where needed. MathLive canonical LaTeX uses its declared format, not the interface locale. Server-side parsing remains authoritative.

Switching while a submission is in flight preserves the request and its idempotency key. Render the eventual verdict in the current interface language using stable feedback codes or a read-only feedback lookup after the result exists. Do not submit again to obtain a translated verdict. Existing grading-unavailable outcomes remain infrastructure errors, not learner mistakes.

Expired anchors return the established expiry behavior. Do not create a new question as a side effect of a presentation lookup. The learner explicitly continues to a newly served instance.

## 7. Locale-aware answer parsing without weaker grading

Add an input-boundary adapter rather than locale-dependent engine rules. Proposed request fields are additive and validated against the served answer contract:

```text
problem_id: opaque instance identifier
answer: raw input
answer_locale: en | ru
answer_format: plain_numeric | mathlive_latex | option_id | typed_text
submission_id: the same value for a retry of one submission
```

For legacy clients with no new fields, retain existing English behavior and text anchors. A migrated server rejects a conflicting `question`/`problem_id` pair. Never accept an expected answer supplied by the client.

For Russian plain numeric input, accept a decimal comma only in a scalar grammar that declares it. Interpret `1,5` as a decimal there. Reject grouping and mixed separators in the initial grammar. In tuples, sets, intervals and matrices, commas already carry structural meaning. Use format-specific parsers and explicit separators rather than global comma replacement. Test canonical LaTeX fractions and mixed numbers separately.

Units use reviewed aliases mapped to canonical unit IDs. Text answers use per-template accepted aliases or stable semantic IDs, never arbitrary machine translation. Unsupported textual formats retain an explicit fallback rather than guessed equivalence.

Normalization preserves existing precision, tolerances and mathematical equivalence. Invalid or ambiguous syntax returns a localized input error before grading or evidence mutation. A valid mathematical answer still follows the existing correct/incorrect path.

## 8. Learning engine and optional AI

No locale enters FIRe graph traversal, mastery thresholds, scheduler selection, review intervals or XP policy. Locale is optional presentation metadata, never part of a learning-state key. Do not store parallel Russian and English progress rows. This holds for current SM-2 paths, upstream FSRS work and future scheduler changes.

Optional AI explanations reuse `internal/ai.Provider` and existing failover policy. The application-level request can include requested language and the minimal authorized mathematical context. Provider IDs, models and routing configuration stay out of curriculum assets and learning records. Section 12 defines secret handling, prompt-input trust and context minimization.

Validate structure, required mathematical facts and output safety before display. These checks do not prove arbitrary generated prose correct. Label optional generated assistance distinctly and fall back to the authored explanation on timeout, invalid output or unavailable providers. Do not use optional AI output for grading, evidence, coverage certification or automatic publication of a translation. Offline AI-assisted drafts require the same human review and validation as other translations.

## 9. Coverage and validation gates

Build the denominator from registered canonical concepts, reachable generator variants, lesson/KP mappings, feedback codes and asset dependencies. Do not hardcode the curriculum size. This permits the expansion target discussed in issue #260 without duplicating concepts.

A locale/course manifest enumerates the scope it claims to cover. Report missing, translated, reviewed, validated and stale units separately. Fully available means every required unit is current, reviewed and validated, including accessibility copy and the course's diagnostic paths. English fallback never counts as coverage.

Proposed automated checks:

- Schema validity, known IDs, duplicate IDs and placeholder/type parity.
- Source/dependency digests, removed targets, missing branches and stale translations.
- Math delimiter parsing, referenced assets, stable choice-ID parity and provenance completeness.
- Generated fixtures covering each declared branch, checking localized operands and solutions against the same canonical answer.
- Counterexamples for sign changes, undefined denominators, negative values and ambiguous answer syntax.
- EN/RU answer normalization with correct, incorrect and invalid examples for each supported answer format.
- UI tests for draft preservation, pending submissions, tab synchronization, hydration and language attributes.
- Server tests for anchor ownership, expiry, restarts, retries and equivalent input producing identical learning updates across locales.
- Learning invariance tests: language changes create no attempt, XP, diversity, mastery, FIRe or scheduling update.

Automated tests cannot certify linguistic quality or prove all prose mathematically correct. Require a named reviewer and mathematical review for each release unit. Preserve existing English regression suites and upstream invariants.

## 10. Incremental migration and reviewable PRs

1. **Design only:** agree on identity, content contracts, draft parsing behavior, fallback policy and availability gates.
2. **Interface infrastructure:** semantic English catalogs, Russian UI catalog, preference persistence and in-place switching. No curriculum or engine rewrite.
3. **Content schema and checks:** overlays, source digests, provenance and a coverage report. Existing English loaders remain authoritative.
4. **Instance identity:** additive anchor metadata and presentation lookup, first for numeric addition. Preserve legacy clients and retry idempotency. Include storage migrations and both SQLite/Postgres paths where persistence changes.
5. **Answer boundary:** explicit formats and locale parsing with parity tests. Stable choice IDs precede translated choices.
6. **Arithmetic vertical slice:** statements, all variants, authored solutions, KP slices, reference lessons, media and diagnostic feedback. Include learn, study, review and quiz delivery. Mark the course fully available only after its gate passes.
7. **Expansion:** migrate additional concepts and locales by the same inventory. AI assistance and account preference synchronization can ship independently.

Each implementation PR is rebased on the agreed upstream baseline and contains its own English regression evidence. Do not merge the fork's deployment configuration or unrelated grading fixes as part of localization. Keep the fork available as an exploratory demonstration.

Migrate additive anchor/storage fields first, then clients. Old in-flight English sessions remain gradeable until normal expiry. Rollback disables the new presentation path while preserving canonical attempts and review state. Revision retention must keep content needed by live anchors available until expiry. Newly stale content cannot receive a fresh validated-course claim.

## 11. Proposed decisions for maintainer confirmation

These are explicit recommendations, not claims of approval.

| Decision | Recommendation and ownership |
| --- | --- |
| Identity | Add stable template/variant contracts and opaque per-delivery instance IDs. Scope instances to the receiving learner and learning session. The underlying mathematics and variety key remain independent of language and delivery ID. |
| Transport | Use the authenticated per-instance presentation lookup in Section 5, with a separate read-only post-submission feedback lookup. No locale bundle in the first implementation. |
| Frontend and content | Use semantic message catalogs and versioned presentation overlays, consistent with the maintainer's acceptance of this direction. No locale-specific concept graph. |
| Draft behavior | Preserve the exact draft and its input locale until the learner clears it. Interface switching alone never changes parsing, option IDs, elapsed time or evidence. |
| Review ownership | `sl1mus` coordinates the Russian contribution, inventory and proposed translation PRs. Each release unit requires a named Russian-language reviewer and a named mathematical reviewer, recorded in review metadata. The maintainer designates or accepts those reviewers and retains merge/release approval. Reviewers are not presumed appointed by this RFC. If either review role is unfilled, content remains unavailable as a validated Russian course. |
| Source revisions | Use SHA-256 of a deterministic serialization of source content and its dependency contract. Normalize UTF-8 line endings to LF and sort object keys. Preserve mathematical tokens, whitespace inside math, branch IDs and option IDs. Do not normalize away semantic changes. Record the digest algorithm/schema version and translation revision. Any digest mismatch is stale until explicit re-review and validation. |
| First implementation slice | Migrate numeric addition for `arith.add.single`, `arith.add.double`, `arith.add.carry` and `arith.add.triple` as a technical pilot. This pilot is not a complete course. |
| First complete Russian course | Cover every canonical concept with `domain == "arithmetic"` in `data/concepts/arithmetic.json` at the approved release baseline. Appendix A pins the inspected inventory. Include every reachable variant, KP slice, reference lesson, solution, feedback code and required media description, plus arithmetic diagnostic, learn, review and quiz paths. |
| Scope growth | A source addition changes the derived required-content inventory. It cannot silently inherit availability. Regenerate and review the course manifest before claiming coverage of the expanded scope. |

The arithmetic course includes its word problems, choice, ordering, tuple and expression formats. It is not limited to the four pilot addition generators. Fraction-domain concepts and other domains are outside this first course, even when linked as references. Inventory and localize every dependency needed to complete an arithmetic lesson or diagnostic path before the complete-course gate passes. Otherwise the course remains explicitly partial.

The source digest includes referenced solution schemas, declared parameters, generator/variant contract revision, canonical option identities and the grading-contract revision. Lesson digests include stable teaching-block IDs, text, captions, accessibility copy and asset revisions. A corpus snapshot revision pins the manifest and the exact expected concept set. New source content is unreviewed by default. There is no automatic reviewed/validated status inheritance.

## 12. Security boundaries and interface acceptance tests

This section fixes design constraints for the later interface PR. The tests below are required implementation acceptance tests, not tests implemented or passed by this documentation PR.

### Instance ownership on every read and submit

Resolve the authenticated account or guest principal from the existing credential mechanism. A raw `student_id` query parameter is not proof of ownership. Store learner identity and the receiving learning-session identity in the served-instance record. The client cannot broaden either by supplying an ID.

Every statement presentation, feedback, anchor read and answer submission checks both identities before returning data or mutating state. Knowing an opaque problem ID grants no access by itself. A different account, a different guest or a different learning session under the same account cannot read or submit the instance. Return the same not-found response for unknown and unauthorized IDs to avoid disclosing their existence. No answer-bearing data is returned on denial.

A session restored after a server restart must reestablish the same valid authenticated session binding. Do not bypass checks for restored anchors. Existing expiry behavior applies only after ownership succeeds. Opening a new learning session does not transfer an old instance. This session-bound proposal needs maintainer confirmation before adding fields to persisted anchors.

### Locale allowlist before path resolution

Map request input to a compiled `Locale` enum with explicit entries for English and Russian before any file access. Validate initial generation, presentation, feedback and answer-input locale fields consistently. Do not derive an allowlist from directory names or permit a client to supply a file path.

An absent preference uses the documented default resolver. An explicit unknown or malformed locale returns a stable invalid-locale response without filesystem access. Browser-language negotiation can map a supported language family to an enum value, but that negotiation is separate from accepting explicit API parameters.

Construct paths only from enum-to-directory constants and validated content IDs obtained from the canonical registry/manifest. Reject traversal, encoded separators, absolute paths, NUL characters and malformed IDs before IO. Client-supplied content IDs are not joined into paths. Bound lookups to trusted immutable assets and reject links escaping that root during build validation. Errors never expose host paths.

### Overlays are untrusted content

Placeholders are typed named facts declared by the template contract. They cannot contain arithmetic evaluation, function calls, property access, imports, arbitrary code or executable directives. Unknown placeholders and type mismatches fail validation. Escape interpolated values for their output context rather than treating values as raw HTML.

Use the existing sanitized Markdown rendering boundary in `KatexContent.tsx` and verify all new rendering surfaces follow an equivalent policy. Do not bypass it with raw inner HTML. Disable KaTeX trust features and validate link/media protocols and asset references. Strip scripts, event handlers, unsafe URLs, embedded frames and active SVG content. Cover both translation-authored text and substituted values. Content validation at build/load time complements rendering sanitization and never replaces it.

### Cache privacy and answer visibility

All personalized generation, presentation, feedback and submission responses, including errors, send `Cache-Control: private, no-store`. CDN/reverse-proxy rules explicitly bypass public caching for those routes. Service workers and browser application caches do not store those payloads for offline or cross-session reuse. Clear in-memory personalized entries at logout, guest replacement or session change.

Where a private in-memory presentation cache exists, its key includes principal, learning session, problem ID, effective locale, content revision and feedback state. Authentication and ownership run before reading that cache. The client checks its active session before applying an asynchronous response. Public immutable catalogs can use revision-addressed caching only if they contain no learner data or answer-bearing problem payloads.

Unanswered DTOs expose no expected answer, canonical answer mapping, solution explanation or raw answer-bearing facts. Review every locale and fallback path, not only the English DTO. Feedback becomes available only from the stored authorized submission result. Never rely on hiding a client-side answer field.

### AI provider and prompt-input boundary

Provider credentials remain in server-side secret configuration. They never enter catalogs, overlays, client bundles, DTOs, logs or model prompts. Optional assistance sends only the minimum context authorized for that learner and feature. Do not send another learner's history or a full profile by default.

Treat translated statements, lesson prose and retrieved content as untrusted prompt data, not instructions. Keep them separate from server-controlled instructions. The assistance feature grants no tools or privileged actions because translated text requests them. Restrict any feature tools by a server-side allowlist independent of content. Validate outputs and retain the authored fallback.

AI never grades, modifies canonical answers, records mastery evidence, approves translations or automatically publishes content. AI-assisted drafts remain subject to named review. This design reuses the existing provider interface without changing routing or scheduler behavior.

### Required interface test matrix

| Test | Required result |
| --- | --- |
| Account A's instance requested by account B | Uniform denial for presentation, feedback and submit. No instance metadata, answer or evidence leaks. |
| Two guests, or two learning sessions under one account | Cross-guest and cross-session reads/submits fail. A valid same-session read still works. |
| Restart, expiry and conflicting legacy question/ID | Valid persisted ownership survives restart. Expiry and conflicting identifiers follow explicit error paths without redraw or duplicate evidence. |
| Unknown locale, traversal and encoded traversal | Malformed input is rejected before filesystem IO. Include unknown tags, slash/backslash, percent-encoded traversal, absolute paths and malicious content IDs. |
| Public-cache probe | Warm an actual configured edge/proxy with A's request, then request as B, another session and anonymously. No public cache hit or A payload is served. Assert no-store headers and service-worker exclusion. Header-only unit checks are insufficient. |
| Hostile translation content | Scripts, HTML handlers, unsafe links/media, active SVG, malicious math directives and expression-like placeholders do not execute or escape the content root. Typed fact interpolation remains inert. |
| AI prompt/secret probe | Hostile translated instructions cannot reveal a server secret, invoke privileged tools, change a grade or publish content. Provider errors return the authored fallback. |
| Language switch with a draft | Preserve problem identity, parameters, typed answer, option selection, timer origin and draft input locale. No attempt, XP, variety, FIRe, mastery or review-state write occurs. |
| Switch during submit and retry | Reuse the submission ID. One attempt and one learning update at most. The final feedback uses the selected presentation without resubmission. |
| Unanswered and post-submission DTOs | Inspect every supported locale and fallback for answer leakage. Explanation/answer visibility is authorized by the stored result, not a client flag. |

Run interface acceptance tests before merging delivery code. Pure documentation checks cannot establish these runtime guarantees.

## Appendix A. First arithmetic course inventory

Scope rule: all records whose canonical `domain` is `arithmetic` in `data/concepts/arithmetic.json` at baseline `a87386b9a0548e2dd544a204bac0c35fc5822b1f`. The list below is a review snapshot, not a new source of concept identity. A future manifest generator must read the canonical file and fail on drift from a declared release snapshot.

- `arith.add.single`
- `arith.add.double`
- `arith.add.carry`
- `arith.add.triple`
- `arith.add.word`
- `arith.sub.single`
- `arith.sub.double`
- `arith.sub.borrow`
- `arith.sub.word`
- `arith.place.tens`
- `arith.place.hundreds`
- `arith.place.thousands`
- `arith.round.tens`
- `arith.round.hundreds`
- `arith.round.thousands`
- `arith.mult.concept`
- `arith.mult.2_5_10`
- `arith.mult.tables`
- `arith.mult.double`
- `arith.mult.triple`
- `arith.mult.word`
- `arith.div.concept`
- `arith.div.basic`
- `arith.div.remainder`
- `arith.div.long`
- `arith.div.word`
- `arith.factor.find`
- `arith.factor.prime`
- `arith.factor.prime_fact`
- `arith.factor.gcf`
- `arith.factor.lcm`
- `arith.factor.composite`
- `arith.exp.concept`
- `arith.exp.evaluate`
- `arith.exp.product_rule`
- `arith.exp.quotient_rule`
- `arith.exp.power_rule`
- `arith.sqrt.perfect`
- `arith.sqrt.simplify`
- `arith.neg.number_line`
- `arith.neg.add_sub`
- `arith.neg.mult_div`
- `arith.order_ops.basic`
- `arith.order_ops.full`
- `arith.order_ops.nested`
- `arith.dec.intro`
- `arith.compare.whole`
- `arith.place.expanded`
- `arith.place.powers_ten`
- `arith.estimate.sum`
- `arith.estimate.product`
- `arith.word.multi_step`
- `arith.div.interpret_remainder`
- `arith.div.zero`
- `arith.neg.word`
- `arith.exp.negative_base`
- `arith.exp.power_product`
- `arith.exp.power_quotient`
- `arith.sqrt.cube`
- `arith.sqrt.estimate`

## Code references at the inspected baseline

- [Canonical problem and generator contract](https://github.com/chuma-beep/mathua/blob/a87386b9a0548e2dd544a204bac0c35fc5822b1f/internal/generator/generator.go).
- [Registry explanation assembly](https://github.com/chuma-beep/mathua/blob/a87386b9a0548e2dd544a204bac0c35fc5822b1f/internal/generator/registry.go).
- [Durable grading anchors](https://github.com/chuma-beep/mathua/blob/a87386b9a0548e2dd544a204bac0c35fc5822b1f/internal/engine/engine.go).
- [Frontend submission identity and retries](https://github.com/chuma-beep/mathua/blob/a87386b9a0548e2dd544a204bac0c35fc5822b1f/web/next-app/lib/api.ts).
- [Lesson and KP loader](https://github.com/chuma-beep/mathua/blob/a87386b9a0548e2dd544a204bac0c35fc5822b1f/internal/lessons/lesson.go).
- [Provider interface](https://github.com/chuma-beep/mathua/blob/a87386b9a0548e2dd544a204bac0c35fc5822b1f/internal/ai/provider.go).
- [Scheduler state including FSRS fields](https://github.com/chuma-beep/mathua/blob/a87386b9a0548e2dd544a204bac0c35fc5822b1f/internal/scheduler/scheduler.go).
