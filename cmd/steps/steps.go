package steps

type StepSchema struct {
	Name    string
	Header  string
	Options []Item
}

type Item struct {
	Title       string
	Description string
}

type Steps struct {
	Steps map[string]StepSchema
}

func InitSteps() *Steps {
	return &Steps{
		map[string]StepSchema{
			"dateOptions": {
				Name: "Date Options",
				Options: []Item{
					{
						Title:       "Year",
						Description: "All the months for a calendar year",
					},
					{
						Title:       "Month",
						Description: "A specific month",
					},
					{
						Title:       "Range",
						Description: "All the months inside a range",
					},
				},
				Header: "What would you like to generate?",
			},
		},
	}
}
