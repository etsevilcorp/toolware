package configagnostic

type Config[I any, O any] struct {
	Name         string
	Description  string
	InputSchema  I
	OutputSchema O
}
