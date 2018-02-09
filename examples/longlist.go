package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

type S struct {
	s string
}

func (s *S) String() string {
	return s.s
}

func s(s string) *S {
	return &S{s: s}
}

// the questions to ask
var simpleQs = []*survey.Question{
	{
		Name: "struct",
		Prompt: &survey.Select{
			Message: "Choose a letter for struct:",
			StringerOptions: []fmt.Stringer{
				s("a"),
				s("b"),
				s("c"),
				s("d"),
				s("e"),
				s("f"),
				s("g"),
				s("h"),
				s("i"),
				s("j"),
			},
		},
		Validate: survey.Required,
	},
	{
		Name: "letter",
		Prompt: &survey.Select{
			Message: "Choose a letter for string:",
			Options: []string{
				"a",
				"b",
				"c",
				"d",
				"e",
				"f",
				"g",
				"h",
				"i",
				"j",
			},
		},
	},
	{
		Name: "letters",
		Prompt: &survey.MultiSelect{
			Message: "Choose some letters:",
			Options: []string{
				"a",
				"b",
				"c",
				"d",
				"e",
				"f",
				"g",
				"h",
				"i",
				"j",
			},
		},
	},
}

func main() {
	answers := struct {
		Struct  *S
		Letter  string
		Letters []string
	}{}

	// ask the question
	err := survey.Ask(simpleQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("you chose %s.\n", answers.Struct)
	fmt.Printf("you chose %s.\n", answers.Letter)
	fmt.Printf("you chose %s.\n", answers.Letters)
}
