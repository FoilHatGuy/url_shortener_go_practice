package cfg

import (
	"flag"
	"github.com/sakirsensoy/genv"
)

var (
	serverAdress    string
	baseURL         string
	fileStoragePath string
)

func Init() {
	flag.StringVar(&serverAdress, "a",
		genv.Key("SERVER_HOST").Default("localhost").String(),
		"help message for flagname")
	flag.StringVar(&baseURL, "b",
		genv.Key("BASE_URL").Default("/").String(),
		"help message for flagname")
	flag.StringVar(&fileStoragePath, "f",
		genv.Key("FILE_STORAGE_PATH").Default("./data").String(),
		"help message for flagname")
	flag.Parse()
}

var Shortener = &shortCfg{
	URLLength: genv.Key("SHORT_URL_LENGTH").Default(10).Int(),
}

var Server = &serverCfg{
	Host: serverAdress,
	Port: genv.Key("SERVER_PORT").Default("8080").String(),
}
var Router = &routerCfg{
	BaseURL: baseURL,
}
var Storage = &storageCfg{
	AutosaveInterval: genv.Key("STORAGE_AUTOSAVE_INTERVAL").Default(60).Int(),
	SavePath:         fileStoragePath,
	StorageType:      genv.Key("STORAGE_TYPE").Default("file").String(),
}
