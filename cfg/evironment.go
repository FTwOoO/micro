package cfg

const (
	TestingEnv = Environment("test")
	LocalEnv   = Environment("local")
	DevEnv     = Environment("dev")
	StagingEnv = Environment("staging")
	ProdEnv    = Environment("prod")
)

type Environment string
