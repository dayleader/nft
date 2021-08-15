package domain

// AppContext - application main context.
type AppContext struct {
	GeneratorParams GeneratorParams
}

// NewAppContext - constructs a new application context.
func NewAppContext() *AppContext {
	return &AppContext{
		GeneratorParams: GeneratorParams{},
	}
}
