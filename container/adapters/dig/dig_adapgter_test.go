package dig_test

type Server struct {
	ID     string
	Logger *Logger
}

type Logger struct {
	Name string
}

func provideServerTypeA(logger *Logger) *Server {
	return &Server{
		ID:     "A",
		Logger: logger,
	}
}
