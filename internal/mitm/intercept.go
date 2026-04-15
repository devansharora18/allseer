package mitm

type Interceptor struct {
	ca *CA
}

func NewInterceptor(ca *CA) *Interceptor {
	return &Interceptor{ca: ca}
}
