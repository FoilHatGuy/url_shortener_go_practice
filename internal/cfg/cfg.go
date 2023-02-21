package cfg

import (
	"flag"
	"fmt"
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
)

func Initialize() {
	flag.StringVar(&serverAdress, "a", "localhost:8080", "help message for flagname")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "help message for flagname")
	flag.StringVar(&fileStoragePath, "f", "./data/data", "help message for flagname")
	flag.Parse()

	fmt.Println(serverAdress, baseURL, fileStoragePath)
	fmt.Println(
		genv.Key("SERVER_ADDRESS").Default("NO SUCH FIELD").String(),
		genv.Key("BASE_URL").Default("NO SUCH FIELD").String(),
		genv.Key("FILE_STORAGE_PATH").Default("NO SUCH FIELD").String())

	Shortener = shortCfg{
		URLLength: genv.Key("SHORT_URL_LENGTH").Default(10).Int(),
	}

	Server = serverCfg{
		Address: genv.Key("SERVER_ADDRESS").Default(serverAdress).String(),
		Port:    genv.Key("SERVER_PORT").Default("8080").String(),
		BaseURL: genv.Key("BASE_URL").Default(baseURL).String(),
	}
	Storage = storageCfg{
		AutosaveInterval: genv.Key("STORAGE_AUTOSAVE_INTERVAL").Default(10).Int(),
		SavePath:         genv.Key("FILE_STORAGE_PATH").Default(fileStoragePath).String(),
		StorageType:      genv.Key("STORAGE_TYPE").Default("file").String(),
	}

}
