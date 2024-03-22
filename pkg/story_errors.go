package squire

import (
	"errors"
	"fmt"
	"strings"
)

type CombinedStoryErrors struct {
	Errors []StoryError
}

type StoryError struct {
	Line int
	Msg  string
}

func (e *CombinedStoryErrors) Error() string {
	var error error

	if len(e.Errors) > 0 {
		var sb strings.Builder
		for _, e := range e.Errors {
			sb.WriteString(fmt.Sprintf("%d: %s\n", e.Line, e.Msg))
		}
		error = errors.New(strings.TrimRight(sb.String(), "\n"))
	}

	return error.Error()
}

func (e *CombinedStoryErrors) Prepend(lineNumber int, msg string) {
	e.Errors = append([]StoryError{{lineNumber, msg}}, e.Errors...)
}

func (e *CombinedStoryErrors) Append(lineNumber int, msg string) {
	e.Errors = append(e.Errors, StoryError{lineNumber, msg})
}

func (c *CombinedStoryErrors) AppendAfterOthersAtLine(line int, msg string) {
	// combinedErrors is ordered by line number. We want to insert the new error
	// at the correct position based on its line number
	for i, e := range c.Errors {
		if e.Line > line {
			c.Errors = append(c.Errors, StoryError{})
			copy(c.Errors[i+1:], c.Errors[i:])
			c.Errors[i] = StoryError{line, msg}
			return
		}
	}

	c.Append(line, msg)
}

func (e *CombinedStoryErrors) findChoicesForNonExistentChapters(story Story) {
	for _, c := range story.Chapters {
		for _, choice := range c.Choices {
			if _, ok := story.Chapters[choice.ChapterId]; !ok {
				e.AppendAfterOthersAtLine(choice.Line, "Invalid chapter id for choice")
			}
		}
	}
}

func (e *CombinedStoryErrors) findUnreachableChapters(story Story) {
	reachableChapters := make(map[string]bool)
	reachableChapters[story.StartChapterId] = true

	for _, c := range story.Chapters {
		for _, choice := range c.Choices {
			reachableChapters[choice.ChapterId] = true
		}
	}

	for id := range story.Chapters {
		if _, ok := reachableChapters[id]; !ok {
			e.AppendAfterOthersAtLine(story.Chapters[id].Line, "Unreachable chapter")
		}
	}
}

func (e *CombinedStoryErrors) addFrontMatterErrors(story Story) {
	if story.Author == "" {
		e.Prepend(1, "Missing author")
	}

	if story.Title == "" {
		e.Prepend(1, "Missing title")
	}
}

func (e *CombinedStoryErrors) Finalize(story Story) error {
	e.findChoicesForNonExistentChapters(story)
	e.findUnreachableChapters(story)
	e.addFrontMatterErrors(story)

	if (len(e.Errors)) == 0 {
		return nil
	}

	return e
}
