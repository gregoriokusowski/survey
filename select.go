package survey

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

/*
Select is a prompt that presents a list of various options to the user
for them to select using the arrow keys and enter. Response type is a string.

	color := ""
	prompt := &survey.Select{
		Message: "Choose a color:",
		Options: []string{"red", "blue", "green"},
	}
	survey.AskOne(prompt, &color, nil)
*/
type Select struct {
	core.Renderer
	Message         string
	Options         []string
	Default         string
	StringerOptions []fmt.Stringer
	StringerDefault fmt.Stringer
	Help            string
	PageSize        int
	selectedIndex   int
	useDefault      bool
	showingHelp     bool
}

// the data available to the templates when processing
type SelectTemplateData struct {
	Select
	PageEntries   []string
	SelectedIndex int
	Answer        string
	ShowAnswer    bool
	ShowHelp      bool
}

var SelectQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else}}
  {{- if and .Help (not .ShowHelp)}} {{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}}{{end}}
  {{- "\n"}}
  {{- range $ix, $choice := .PageEntries}}
    {{- if eq $ix $.SelectedIndex}}{{color "cyan+b"}}{{ SelectFocusIcon }} {{else}}{{color "default+hb"}}  {{end}}
    {{- $choice}}
    {{- color "reset"}}{{"\n"}}
  {{- end}}
{{- end}}`

// OnChange is called on every keypress.
func (s *Select) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	// if the user pressed the enter key
	if key == terminal.KeyEnter {
		return []rune(s.StringerOptions[s.selectedIndex].String()), 0, true
		// if the user pressed the up arrow
	} else if key == terminal.KeyArrowUp {
		s.useDefault = false

		// if we are at the top of the list
		if s.selectedIndex == 0 {
			// start from the button
			s.selectedIndex = len(s.StringerOptions) - 1
		} else {
			// otherwise we are not at the top of the list so decrement the selected index
			s.selectedIndex--
		}
		// if the user pressed down and there is room to move
	} else if key == terminal.KeyArrowDown {
		s.useDefault = false
		// if we are at the bottom of the list
		if s.selectedIndex == len(s.StringerOptions)-1 {
			// start from the top
			s.selectedIndex = 0
		} else {
			// increment the selected index
			s.selectedIndex++
		}
		// only show the help message if we have one
	} else if key == core.HelpInputRune && s.Help != "" {
		s.showingHelp = true
	}

	// figure out the options and index to render
	opts, idx := paginate(s.PageSize, s.StringerOptions, s.selectedIndex)

	// render the options
	s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:        *s,
			SelectedIndex: idx,
			ShowHelp:      s.showingHelp,
			PageEntries:   opts,
		},
	)

	// if we are not pressing ent
	return []rune(s.StringerOptions[s.selectedIndex].String()), 0, true
}

type option struct {
	value string
}

func (o option) String() string {
	return o.value
}

func (s *Select) Prompt() (interface{}, error) {
	// Adapt old options
	if len(s.StringerOptions) == 0 && len(s.Options) != 0 {
		if s.StringerDefault == nil && s.Default != "" {
			s.StringerDefault = &option{value: s.Default}
		}
		for _, o := range s.Options {
			s.StringerOptions = append(s.StringerOptions, &option{value: o})
		}
	}

	// if there are no options to render
	if len(s.StringerOptions) == 0 {
		// we failed
		return "", errors.New("please provide options to select from")
	}

	// start off with the first option selected
	sel := 0
	// if there is a default
	if s.StringerDefault != nil {
		// find the choice
		for i, opt := range s.StringerOptions {
			// if the option correponds to the default
			if opt.String() == s.StringerDefault.String() {
				// we found our initial value
				sel = i
				// stop looking
				break
			}
		}
	}
	// save the selected index
	s.selectedIndex = sel

	// figure out the options and index to render
	opts, idx := paginate(s.PageSize, s.StringerOptions, sel)

	// ask the question
	err := s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:        *s,
			PageEntries:   opts,
			SelectedIndex: idx,
		},
	)
	if err != nil {
		return "", err
	}

	// hide the cursor
	terminal.CursorHide()
	// show the cursor when we're done
	defer terminal.CursorShow()

	// by default, use the default value
	s.useDefault = true

	rr := terminal.NewRuneReader(os.Stdin)
	rr.SetTermMode()
	defer rr.RestoreTermMode()
	// start waiting for input
	for {
		r, _, err := rr.ReadRune()
		if err != nil {
			return "", err
		}
		if r == '\r' || r == '\n' {
			break
		}
		if r == terminal.KeyInterrupt {
			return "", terminal.InterruptErr
		}
		if r == terminal.KeyEndTransmission {
			break
		}
		s.OnChange(nil, 0, r)
	}

	var val fmt.Stringer
	// if we are supposed to use the default value
	if s.useDefault {
		// if there is a default value
		if s.StringerDefault != nil {
			// use the default value
			val = s.StringerDefault
		} else {
			// there is no default value so use the first
			val = s.StringerOptions[0]
		}
		// otherwise the selected index points to the value
	} else {
		// the
		val = s.StringerOptions[s.selectedIndex]
	}
	if v, ok := val.(*option); ok {
		return v.String(), err
	}

	return val, err
}

func (s *Select) Cleanup(val interface{}) error {
	var ans string
	if stringerVal, ok := val.(fmt.Stringer); ok {
		ans = stringerVal.String()
	} else {
		ans = val.(string)
	}
	return s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:     *s,
			Answer:     ans,
			ShowAnswer: true,
		},
	)
}
