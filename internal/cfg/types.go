package cfg

type ConfigT struct {
	Shortener ShortCfg
	Server    ServerCfg
	Storage   StorageCfg
}
type ShortCfg struct {
	Secret    string
	URLLength int
}

type ServerCfg struct {
	Address        string
	Port           string
	BaseURL        string
	CookieLifetime int
}
type StorageCfg struct {
	AutosaveInterval int
	SavePath         string
	StorageType      string
	DatabaseDSN      string
}
