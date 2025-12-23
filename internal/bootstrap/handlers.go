package bootstrap

type Handlers struct{}

func GetHandlers(repos *Repositories) *Handlers {
	return &Handlers{}
}
