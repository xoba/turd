package cnfg

type Config struct {
	Mode           string
	AWSProfile     string
	Port           int
	PublicKeyFile  string
	PrivateKeyFile string
	Seed           int
	Delete         bool
	Debug          bool
	Lisp           string
	File           string
	Profile        string // name of profile output, if any
}
