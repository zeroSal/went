package template

type Renderer interface {
	Render(data map[string]string) error
	Read(path string) ([]byte, bool)
	GetFiles() []string
}
