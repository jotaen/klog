package cli

type Command struct {
	Main        func(Environment, []string) int
	Name        string
	Alias       []string
	Description string
}
