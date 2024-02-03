package assistants

type Assistant interface {
	Query(message string) (response []string, err error)
	Explain(command string) (response []string, err error)
}
