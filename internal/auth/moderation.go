package auth

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// ErrNameNotAllowed is deliberately vague. Naming the matched word turns the
// filter into a guide for evading it and tells a learner which slur the system
// found in their name.
var ErrNameNotAllowed = errors.New("that name is not allowed")

// blockedTerm is one curated entry. `exact` means the whole compacted token
// must match, not merely contain the pattern. Short words such as ass and dick
// are exact because they sit inside ordinary names: Dickens and Cassandra must
// keep working.
type blockedTerm struct {
	pattern *regexp.Regexp
	exact   bool
}

func repeatedLetterPattern(term string, exact bool) *regexp.Regexp {
	var b strings.Builder
	if exact {
		b.WriteString("^")
	}
	for _, r := range term {
		b.WriteString(regexp.QuoteMeta(string(r)))
		b.WriteString("+")
	}
	if exact {
		b.WriteString("$")
	}
	return regexp.MustCompile(b.String())
}

func blocked(term string, exact bool) blockedTerm {
	return blockedTerm{pattern: repeatedLetterPattern(term, exact), exact: exact}
}

// The list is intentionally short and visible. A large list produces silent
// false positives that nobody can explain. Repeated letters are matched by the
// pattern itself, so "fuuuck" and "niiigger" are caught without collapsing
// legitimate repeated letters in ordinary names.
var blockedTerms = []blockedTerm{
	blocked("fuck", false),
	blocked("shit", true),
	blocked("bitch", false),
	blocked("cunt", false),
	blocked("pussy", false),
	blocked("whore", false),
	blocked("slut", false),
	blocked("bastard", false),
	blocked("asshole", false),
	blocked("motherfucker", false),
	blocked("nigger", false),
	blocked("nigga", false),
	blocked("faggot", false),
	blocked("retard", false),
	blocked("kike", false),
	blocked("spic", false),
	blocked("chink", false),
	blocked("tranny", false),
	blocked("wetback", false),
	blocked("tits", false),

	// Exact-only short terms. These are slurs or profanity on their own but
	// occur inside legitimate words, so substring matching would block names
	// such as Dickens, Cassandra and title.
	blocked("ass", true),
	blocked("dick", true),
	blocked("fag", true),
	blocked("tit", true),
	blocked("coon", true),
	blocked("dyke", true),
}

var reservedNames = map[string]bool{
	"admin":         true,
	"administrator": true,
	"help":          true,
	"mathua":        true,
	"moderator":     true,
	"official":      true,
	"owner":         true,
	"root":          true,
	"security":      true,
	"staff":         true,
	"support":       true,
	"system":        true,
	"teacher":       true,
}

// compactName lowercases, strips accents, maps the usual leetspeak and keeps
// only letters and digits. The result is what the patterns run against.
func compactName(s string) string {
	s = norm.NFD.String(strings.ToLower(strings.TrimSpace(s)))
	var b strings.Builder
	for _, r := range s {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		switch r {
		case '0':
			r = 'o'
		case '1', '!', '|':
			r = 'i'
		case '3':
			r = 'e'
		case '4', '@':
			r = 'a'
		case '5', '$':
			r = 's'
		case '7':
			r = 't'
		case '8':
			r = 'b'
		case '9':
			r = 'g'
		}
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func nameIsAllowed(compact string) bool {
	if compact == "" {
		return true
	}
	if reservedNames[compact] {
		return false
	}
	for _, term := range blockedTerms {
		if term.pattern.MatchString(compact) {
			return false
		}
	}
	return true
}

// exactTermInTokens catches an exact-only term when it is a separate word.
// The compacted whole name cannot do that: "what the shit" becomes
// "whattheshit" and the exact pattern no longer matches, while "Kshitij" and
// "Shittu" must stay legal because those are real names that contain the
// letters by accident.
func exactTermInTokens(name string) bool {
	tokens := strings.FieldsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	for _, token := range tokens {
		compact := compactName(token)
		if compact == "" {
			continue
		}
		if reservedNames[compact] {
			return true
		}
		for _, term := range blockedTerms {
			if term.exact && term.pattern.MatchString(compact) {
				return true
			}
		}
	}
	return false
}

// ValidateDisplayName applies the same moderation to a display name that a
// username receives. Display names stay repeatable, because the username is
// the unique handle. The empty and length messages are preserved for callers
// that already depend on them.
func ValidateDisplayName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name is required")
	}
	if len([]rune(name)) > 50 {
		return errors.New("name must be 1-50 characters")
	}
	if !nameIsAllowed(compactName(name)) || exactTermInTokens(name) {
		return ErrNameNotAllowed
	}
	return nil
}

// SafeDisplayName is the OAuth entry point. A provider-supplied name is not
// user input this service accepted, so a blocked or malformed name must not
// fail the sign-in. Keep the provider name when it passes, fall back to the
// email local part when that is acceptable, then use a neutral label.
func SafeDisplayName(name, email string) string {
	name = strings.TrimSpace(name)
	if ValidateDisplayName(name) == nil {
		return name
	}
	local := strings.TrimSpace(strings.SplitN(email, "@", 2)[0])
	if ValidateDisplayName(local) == nil {
		return local
	}
	return "Learner"
}

// validateUsernameModeration is called after the character rules, so the
// compact form only contains letters, digits and separators the username
// policy already allows.
func validateUsernameModeration(username string) error {
	if !nameIsAllowed(compactName(username)) {
		return ErrNameNotAllowed
	}
	return nil
}
