package cfg

type shortCfg struct {
	Secret    string
	URLLength int
}

type serverCfg struct {
	Address string
	Port    string
	BaseURL string
}
type storageCfg struct {
	AutosaveInterval int
	SavePath         string
	StorageType      string
}
