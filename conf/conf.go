package conf

type Config struct {
	Repository *string
	Android    *AndroidConfig
}

type AndroidConfig struct {
	Store         *string
	StorePassword *string
	Alias         *string
	AliasPassword *string
}
