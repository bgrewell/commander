package mutations

type Mutation interface {
	Apply(input string) (output string)
}
