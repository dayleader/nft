package domain

// AppContext - application main context.
type AppContext struct {
	GeneratorParams GeneratorParams
	IpfsParams      IpfsParams
}

// NewAppContext - constructs a new application context.
func NewAppContext() *AppContext {
	return &AppContext{
		GeneratorParams: GeneratorParams{},
		IpfsParams:      IpfsParams{},
	}
}
