package cfg

import (
	"flag"
	"github.com/sakirsensoy/genv"
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
)

var (
	serverAdress    string
	baseURL         string
	fileStoragePath string

	Shortener shortCfg
	Server    serverCfg
	Storage   storageCfg
	Router    routerCfg
)

func Initialize() {
	flag.StringVar(&serverAdress, "a",
		genv.Key("SEVER_ADDRESS").Default("localhost").String(),
		"help message for flagname")
	flag.StringVar(&baseURL, "b",
		genv.Key("BASE_URL").Default("http://localhost:8080").String(),
		"help message for flagname")
	flag.StringVar(&fileStoragePath, "f",
		genv.Key("FILE_STORAGE_PATH").Default("./data").String(),
		"help message for flagname")
	flag.Parse()

	Shortener = shortCfg{
		URLLength: genv.Key("SHORT_URL_LENGTH").Default(10).Int(),
	}

	Server = serverCfg{
		Host: serverAdress,
		Port: genv.Key("SERVER_PORT").Default("8080").String(),
	}
	Router = routerCfg{
		BaseURL: baseURL,
	}
	Storage = storageCfg{
		AutosaveInterval: genv.Key("STORAGE_AUTOSAVE_INTERVAL").Default(30).Int(),
		SavePath:         fileStoragePath,
		StorageType:      genv.Key("STORAGE_TYPE").Default("file").String(),
	}

}
