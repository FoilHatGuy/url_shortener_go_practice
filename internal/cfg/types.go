package cfg

type ConfigT struct {
	Shortener shortCfg
	Server    serverCfg
	Storage   storageCfg
}
type shortCfg struct {
	Secret    string
	URLLength int
}

type serverCfg struct {
	Address        string
	Port           string
	BaseURL        string
	CookieLifetime int
}
type storageCfg struct {
	AutosaveInterval int
	SavePath         string
	StorageType      string
	DatabaseDSN      string
}
