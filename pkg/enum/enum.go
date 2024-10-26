package enum

type Enum interface {
	EnumDescriptions() []string
}

func Hello() string {
	return "test"
}
