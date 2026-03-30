package errors

import (
	"fmt"
	"strconv"
	"strings"
)

func SyntaxError(message string, line int, location int,
	text string, filename string) error {
	println(message, line)

	arr := strings.Split(text, "\n")
	col := 0
	inc := 0
	line_inc := 0

	for inc < location {
		if line_inc < len(arr) && inc+len(arr[line_inc]) < location {
			inc += len(arr[line_inc])
			line_inc++
		} else {
			col = location - inc
			break
		}

	}

	var spacer = strings.Repeat(" ", len(strconv.Itoa(line)))

	filename_text := spacer + blue_text("--> ") + filename +
		fmt.Sprintf(":%d:%d", line, col+1)

	// TODO: Fix the fucked math
	var final_line = "  " + spacer + strings.Repeat(" ", max((col-line)+1, 0))

	final_line += red_text("^ here")
	return fmt.Errorf("%s: %s\n%s\n%s\n%s %s\n %s",
		red_text("Error"), message, filename_text, blue_text(spacer+" |"),
		blue_text(fmt.Sprintf("%d |", line)), strings.Split(text, "\n")[line-1],
		final_line)
}

func RuntimeError(message string, line int, filename string) error {
	return fmt.Errorf("\nRuntime Error: %s\n    at %s:%d",
		message, filename, line)
}

type ReturnErrorType struct {
	Value    any
	Line     int
	Filename string
}

func (self ReturnErrorType) Error() string {
	return RuntimeError("'return' outside of function",
		self.Line, self.Filename).Error()
}

type BreakErrorType struct {
	Line int
	Filename string
}

func (self BreakErrorType) Error() string {
	return RuntimeError("'break' outside of loop",
		self.Line, self.Filename).Error()
}

type ContinueErrorType struct {
	Line int
	Filename string
}

func (self ContinueErrorType) Error() string {
	return RuntimeError("'break' outside of loop",
		self.Line, self.Filename).Error()
}

func red_text(text string) string {
	return "\x1b[31m" + text + "\x1b[0m"
}

func bold_text(text string) string {
	return "\033[1m" + text + "\033[0m"
}

func blue_text(text string) string {
	return "\x1b[34m" + text + "\x1b[0m"
}
