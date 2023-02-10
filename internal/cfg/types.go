package cfg

type shortCfg struct {
	UrlLength int
}

type serverCfg struct {
	Host string
	Port string
}
type routerCfg struct {
	BaseURL string
}
type storageCfg struct {
	AutosaveInterval int
	SavePath         string
	StorageType      string
}
