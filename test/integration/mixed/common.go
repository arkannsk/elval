package mixed

// NoAnnotations структура без аннотаций - .gen.go не должен генерироваться
type NoAnnotations struct {
	ID   int
	Name string
}

// Helper вспомогательная структура без аннотаций
type Helper struct {
	Code string
}
