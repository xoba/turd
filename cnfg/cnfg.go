package cnfg

type Config struct {
	Mode           string
	AWSProfile     string
	Port           int
	PublicKeyFile  string
	PrivateKeyFile string
	Seed           int
	Delete         bool
}
