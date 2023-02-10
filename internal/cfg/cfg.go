package cfg

import (
	"flag"
	"github.com/sakirsensoy/genv"
)

var (
	server_adress     string
	base_url          string
	file_storage_path string
)

func Init() {
	flag.StringVar(&server_adress, "a",
		genv.Key("SERVER_HOST").Default("localhost").String(),
		"help message for flagname")
	flag.StringVar(&base_url, "b",
		genv.Key("BASE_URL").Default("/").String(),
		"help message for flagname")
	flag.StringVar(&file_storage_path, "f",
		genv.Key("FILE_STORAGE_PATH").Default("./data").String(),
		"help message for flagname")
	flag.Parse()
}

var Shortener = &shortCfg{
	UrlLength: genv.Key("SHORT_URL_LENGTH").Default(10).Int(),
}

var Server = &serverCfg{
	Host: server_adress,
	Port: genv.Key("SERVER_PORT").Default("8080").String(),
}
var Router = &routerCfg{
	BaseURL: base_url,
}
var Storage = &storageCfg{
	AutosaveInterval: genv.Key("STORAGE_AUTOSAVE_INTERVAL").Default(60).Int(),
	SavePath:         file_storage_path,
	StorageType:      genv.Key("STORAGE_TYPE").Default("file").String(),
}
