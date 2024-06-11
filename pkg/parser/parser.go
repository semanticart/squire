// Package parser provides a parser and validation for Squire story content.
package parser

import (
	"regexp"
	"strings"

	"github.com/semanticart/squire/internal"
)

var (
	// - [do not answer](#ignore-phone)
	// - [jump in a nearby lion's mouth](#lion)
	choiceRegex = regexp.MustCompile(`^- \[(.+)\]\(#(.+)\)$`)

	// - [pick up phone]
	// - [do not answer](#)
	choiceWithoutIDRegex = regexp.MustCompile(`^- \[(.+)\]`)

	// - (#ignore-phone)
	choiceWithoutTextRegex = regexp.MustCompile(`^- \(#(.+)\)`)
)

func isFrontMatter(line string) bool {
	return strings.HasPrefix(line, "% ")
}

func isNewChapter(line string) bool {
	return strings.HasPrefix(line, "# ")
}

func isChoice(line string) bool {
	return strings.HasPrefix(line, "- [") || strings.HasPrefix(line, "- (")
}

func parseChoice(line string, lineNumber int) Choice {
	matches := choiceRegex.FindStringSubmatch(line)

	if len(matches) != 3 {
		matches = choiceWithoutIDRegex.FindStringSubmatch(line)

		if len(matches) != 2 {
			matches = choiceWithoutTextRegex.FindStringSubmatch(line)

			if len(matches) != 2 {
				return Choice{Line: lineNumber}
			}

			return Choice{
				ChapterID: matches[1],
				Line:      lineNumber,
			}
		}

		return Choice{
			Text: matches[1],
			Line: lineNumber,
		}
	}

	return Choice{
		Text:      matches[1],
		ChapterID: matches[2],
		Line:      lineNumber,
	}
}

// Parse function parses Squire story content and returns a Story struct and an
// error of combined validation issues.
func Parse(content string) (Story, error) {
	story := Story{}
	chapter := Chapter{}

	lines := strings.Split(content, "\n")

	for lineNumber, line := range lines {
		switch {
		case isFrontMatter(line):
			if story.Title == "" {
				story.Title = strings.TrimPrefix(line, "% ")
			} else {
				story.Author = strings.TrimPrefix(line, "% ")
			}
		case isNewChapter(line):
			if chapter.ID != "" || chapter.Title != "" {
				chapter.EndLine = lineNumber
				story.appendChapter(chapter)
			}

			title, ID, intentionalDeadEnd := internal.ParseChapterTitleAndID(line)

			chapter = Chapter{
				StartLine:          lineNumber + 1,
				Title:              title,
				ID:                 ID,
				IntentionalDeadEnd: intentionalDeadEnd,
			}
		case isChoice(line):
			choice := parseChoice(line, lineNumber+1)

			chapter.Choices = append(chapter.Choices, choice)
		default:
			chapter.Body += line + "\n"
		}
	}

	if chapter.ID != "" || chapter.Title != "" {
		chapter.EndLine = len(lines)
		story.appendChapter(chapter)
	}

	return story, validate(story)
}
