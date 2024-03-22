package squire

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// # Something isn't right here. {#intro}
	chapterTitleRegex = regexp.MustCompile(`# (.+) {#(.+)}`)
	// - [jump in a nearby lion's mouth](#lion)
	choiceRegex = regexp.MustCompile(`- \[(.+)\]\(#(.+)\)`)
	// - [ or - (
	choiceStartRegex = regexp.MustCompile(`- [\[\(]`)
)

type Story struct {
	Title          string
	Author         string
	StartChapterId string
	Chapters       map[string]Chapter
}

func (s *Story) AddChapter(chapter Chapter, lineNumber int, combinedErrors *error) {
	chapter.Text = strings.TrimSpace(chapter.Text)

	if len(chapter.Choices) == 0 {
		if chapter.Text == "" {
			*combinedErrors = errors.Join(*combinedErrors, fmt.Errorf("%d: Missing chapter text", lineNumber+1))
		}

		*combinedErrors = errors.Join(*combinedErrors, fmt.Errorf("%d: Dead end", lineNumber+1))
	}

	s.Chapters[chapter.Id] = chapter
}

type Chapter struct {
	Title   string
	Id      string
	Text    string
	Choices []Choice
}

type Choice struct {
	Text      string
	ChapterId string
}

func (c *Chapter) AddChoice(choice Choice, lineNumber int, combinedErrors *error) {
	if choice.Text == "" || choice.ChapterId == "" {
		*combinedErrors = errors.Join(*combinedErrors, fmt.Errorf("%d: Invalid choice", lineNumber+1))
	} else {
		if len(c.Choices) == 0 && strings.TrimSpace(c.Text) == "" {
			*combinedErrors = errors.Join(*combinedErrors, fmt.Errorf("%d: Missing chapter text", lineNumber))
		}
		c.Choices = append(c.Choices, choice)
	}
}

func isFrontMatter(line string) bool {
	return strings.HasPrefix(line, "% ")
}

func isNewChapter(line string) bool {
	return strings.HasPrefix(line, "# ")
}

func isChoice(line string) bool {
	return choiceStartRegex.MatchString(line)
}

func parseChoice(line string) Choice {
	matches := choiceRegex.FindStringSubmatch(line)

	if len(matches) != 3 {
		return Choice{}
	}

	return Choice{
		Text:      matches[1],
		ChapterId: matches[2],
	}
}

func ParseStory(contents string) (Story, error) {
	story := Story{
		Chapters: make(map[string]Chapter),
	}

	var combinedErrors error
	var chapter Chapter

	lines := strings.Split(contents, "\n")

	for lineNumber, line := range lines {
		switch {
		case isFrontMatter(line):
			if story.Author == "" && story.Title != "" {
				story.Author = strings.TrimPrefix(line, "% ")
			}

			if story.Title == "" {
				story.Title = strings.TrimPrefix(line, "% ")
			}

		case isNewChapter(line):
			if chapter.Id != "" {
				story.AddChapter(chapter, lineNumber-1, &combinedErrors)

				if story.StartChapterId == "" {
					story.StartChapterId = chapter.Id
				}
			}

			matches := chapterTitleRegex.FindStringSubmatch(line)

			if len(matches) == 3 {
				chapter = Chapter{
					Title: matches[1],
					Id:    matches[2],
				}
			} else {
				chapter = Chapter{
					Title: "INVALID",
					Id:    "INVALID",
				}
				combinedErrors = errors.Join(combinedErrors, fmt.Errorf("%d: Invalid chapter title", lineNumber+1))
			}

		case isChoice(line):
			choice := parseChoice(line)
			chapter.AddChoice(choice, lineNumber, &combinedErrors)

		default:
			chapter.Text += line + "\n"
		}
	}

	story.AddChapter(chapter, len(lines), &combinedErrors)

	if story.Author == "" {
		combinedErrors = errors.Join(fmt.Errorf("1: Missing author"), combinedErrors)
	}

	if story.Title == "" {
		combinedErrors = errors.Join(fmt.Errorf("1: Missing title"), combinedErrors)
	}

	return story, combinedErrors
}
