package squire

import (
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

func (s *Story) AddChapter(chapter Chapter, lineNumber int, combinedErrors *CombinedStoryErrors) {
	chapter.Text = strings.TrimSpace(chapter.Text)
	chapter.OriginalOrder = len(s.Chapters)

	if len(chapter.Choices) == 0 {
		if chapter.Text == "" {
			combinedErrors.Append(lineNumber+1, "Missing chapter text")
		}

		combinedErrors.Append(lineNumber+1, "Dead end")
	}

	s.Chapters[chapter.Id] = chapter
}

type Chapter struct {
	Title         string
	Id            string
	Text          string
	Choices       []Choice
	Line          int
	OriginalOrder int
}

type Choice struct {
	Text      string
	ChapterId string
	Line      int
}

func (c *Chapter) AddChoice(choice Choice, lineNumber int, combinedErrors *CombinedStoryErrors) {
	if choice.Text == "" || choice.ChapterId == "" {
		combinedErrors.Append(lineNumber+1, "Invalid choice")
	} else {
		if len(c.Choices) == 0 && strings.TrimSpace(c.Text) == "" {
			combinedErrors.Append(lineNumber, "Missing chapter text")
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

func parseChoice(line string, lineNumber int) Choice {
	matches := choiceRegex.FindStringSubmatch(line)

	if len(matches) != 3 {
		return Choice{}
	}

	return Choice{
		Text:      matches[1],
		ChapterId: matches[2],
		Line:      lineNumber + 1,
	}
}

func ParseStory(contents string) (Story, error) {
	story := Story{
		Chapters: make(map[string]Chapter),
	}

	var combinedErrors CombinedStoryErrors
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

			chapter = Chapter{Line: lineNumber + 1}

			if len(matches) == 3 {
				chapter.Title = matches[1]
				chapter.Id = matches[2]
			} else {
				chapter.Title = "INVALID TITLE"
				chapter.Id = "INVALID ID"
				combinedErrors.Append(lineNumber+1, "Invalid chapter title")
			}

		case isChoice(line):
			choice := parseChoice(line, lineNumber)
			chapter.AddChoice(choice, lineNumber, &combinedErrors)

		default:
			chapter.Text += line + "\n"
		}
	}

	story.AddChapter(chapter, len(lines), &combinedErrors)

	error := combinedErrors.Finalize(story)

	return story, error
}
