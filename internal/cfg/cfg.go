package cfg

import "github.com/sakirsensoy/genv"

var Shortener = &shortCfg{
	UrlLength: genv.Key("SHORT_URL_LENGTH").Default(10).Int(),
}

var Server = &serverCfg{
	Host: genv.Key("SERVER_HOST").Default("localhost").String(),
	Port: genv.Key("SERVER_PORT").Default("8080").String(),
}
var Router = &routerCfg{
	BaseURL: genv.Key("BASE_URL").Default("/").String(),
}
var Storage = &storageCfg{
	AutosaveInterval: genv.Key("STORAGE_AUTOSAVE_INTERVAL").Default(60).Int(),
	SavePath:         genv.Key("FILE_STORAGE_PATH").Default("./data").String(),
	StorageType:      genv.Key("STORAGE_TYPE").Default("file").String(),
}
