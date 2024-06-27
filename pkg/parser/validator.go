package parser

import (
	"fmt"
	"slices"
	"strings"
)

type StoryError struct {
	Msg  string
	Line int
}

func (e StoryError) Error() string {
	return fmt.Sprintf("%d: %s", e.Line, e.Msg)
}

// CombinedStoryError represents a collection of story errors.
type CombinedStoryError struct {
	Errors []StoryError
}

func (e CombinedStoryError) Error() string {
	msg := ""

	for _, err := range e.Errors {
		msg += err.Error() + "\n"
	}

	return strings.TrimSpace(msg)
}

func validateTitle(errs CombinedStoryError, story Story) CombinedStoryError {
	if story.Title == "" {
		errs.Errors = append(errs.Errors, StoryError{Line: 1, Msg: "Missing title"})
	}

	return errs
}

func validateAuthor(errs CombinedStoryError, story Story) CombinedStoryError {
	if story.Author == "" {
		line := 2

		if story.Title == "" {
			line = 1
		}

		errs.Errors = append(errs.Errors, StoryError{Line: line, Msg: "Missing author"})
	}

	return errs
}

func validateChapterID(errs CombinedStoryError, chapter Chapter) CombinedStoryError {
	if chapter.ID == "" {
		errs.Errors = append(errs.Errors, StoryError{Line: chapter.StartLine, Msg: "Missing chapter id"})
	}

	return errs
}

func validateChapterTitle(errs CombinedStoryError, chapter Chapter) CombinedStoryError {
	if chapter.Title == "" {
		errs.Errors = append(errs.Errors, StoryError{Line: chapter.StartLine, Msg: "Missing chapter title"})
	}

	return errs
}

func validateChapterText(errs CombinedStoryError, chapter Chapter) CombinedStoryError {
	if chapter.Body == "" {
		errs.Errors = append(errs.Errors, StoryError{Line: chapter.StartLine + 1, Msg: "Missing chapter text"})
	}

	return errs
}

func validateChoiceID(errs CombinedStoryError, choice Choice, validChapterIDs []string) CombinedStoryError {
	if choice.ChapterID == "" {
		errs.Errors = append(errs.Errors, StoryError{Line: choice.Line, Msg: "Missing choice id"})
	}

	if !slices.Contains(validChapterIDs, choice.ChapterID) {
		errs.Errors = append(errs.Errors, StoryError{Line: choice.Line, Msg: "Invalid chapter id for choice"})
	}

	return errs
}

func validateChoiceText(errs CombinedStoryError, choice Choice) CombinedStoryError {
	if choice.Text == "" {
		errs.Errors = append(errs.Errors, StoryError{Line: choice.Line, Msg: "Missing choice text"})
	}

	return errs
}

func validateChapterIsReachable(errs CombinedStoryError, chapter Chapter, choiceChapterIDs []string) CombinedStoryError {
	if !slices.Contains(choiceChapterIDs, chapter.ID) {
		errs.Errors = append(errs.Errors, StoryError{Line: chapter.StartLine, Msg: "Unreachable chapter"})
	}

	return errs
}

func validateNotDeadEnd(errs CombinedStoryError, chapter Chapter) CombinedStoryError {
	if !chapter.IntentionalDeadEnd && len(chapter.Choices) == 0 {
		errs.Errors = append(errs.Errors, StoryError{Line: chapter.EndLine, Msg: "Dead end"})
	}

	return errs
}

func validate(story Story) error {
	errs := CombinedStoryError{}
	validChapterIDs := []string{}
	choiceChapterIDs := []string{}

	for _, chapter := range story.Chapters {
		validChapterIDs = append(validChapterIDs, chapter.ID)

		for _, choice := range chapter.Choices {
			choiceChapterIDs = append(choiceChapterIDs, choice.ChapterID)
		}
	}

	errs = validateTitle(errs, story)
	errs = validateAuthor(errs, story)

	for _, chapter := range story.Chapters {
		errs = validateChapterID(errs, chapter)
		errs = validateChapterTitle(errs, chapter)
		errs = validateChapterIsReachable(errs, chapter, choiceChapterIDs)
		errs = validateChapterText(errs, chapter)
		errs = validateNotDeadEnd(errs, chapter)

		for _, choice := range chapter.Choices {
			errs = validateChoiceID(errs, choice, validChapterIDs)
			errs = validateChoiceText(errs, choice)
		}
	}

	if len(errs.Errors) > 0 {
		return errs
	}

	return nil
}
