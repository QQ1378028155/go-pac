package conf

type Config struct {
	Repository *string
	Android    *AndroidConfig
	Xcode      *XcodeConfig
}

type AndroidConfig struct {
	Store         *string
	StorePassword *string
	Alias         *string
	AliasPassword *string
}

type XcodeConfig struct {
	Sign      *string
	Provision *string
}
