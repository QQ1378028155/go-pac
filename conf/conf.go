package conf

//Config is used as a json protocol
type Config struct {
	Repository *string
	Android    *AndroidConfig
	Xcode      *XcodeConfig
}

//AndroidConfig contains some fields for Android
type AndroidConfig struct {
	Store         *string
	StorePassword *string
	Alias         *string
	AliasPassword *string
}

//XcodeConfig contains some fields for Xcode
type XcodeConfig struct {
	Sign      *string
	Provision *string
}
