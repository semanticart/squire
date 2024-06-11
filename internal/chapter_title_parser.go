// Package internal provides things we want to test and use internally but not expose.
// This file contains the implementation of the chapter title parser.
package internal

import (
	"regexp"
	"strings"
)

var (
	// # Something isn't right here. {#intro}
	// # Going to school {#backpack} !!
	newChapterRegex = regexp.MustCompile(`^# (.+) {#(.+)}(\s+!{2})?$`)

	// # {#phone}
	newChapterMissingTextRegex = regexp.MustCompile(`^# {#(.+)}(\s+!{2})?$`)

	// # Something isn't right here.
	newChapterMissingIDRegex = regexp.MustCompile(`^# (.+?)(\s+!{2})?$`)
)

func isIntentionalDeadEnd(marker string) bool {
	return strings.TrimSpace(marker) == "!!"
}

// ParseChapterTitleAndID function parses the chapter title, ID, and isIntentionalDeadEnd from the given line.
func ParseChapterTitleAndID(line string) (string, string, bool) {
	matches := newChapterRegex.FindStringSubmatch(line)

	if len(matches) != 4 {
		matches = newChapterMissingTextRegex.FindStringSubmatch(line)

		if len(matches) != 3 {
			matches = newChapterMissingIDRegex.FindStringSubmatch(line)

			if len(matches) != 3 {
				return "", "", false
			}

			return matches[1], "", isIntentionalDeadEnd(matches[2])
		}

		return "", matches[1], isIntentionalDeadEnd(matches[2])
	}

	return matches[1], matches[2], isIntentionalDeadEnd(matches[3])
}
