package cfg

import (
	"flag"
	"fmt"

	"github.com/sakirsensoy/genv"
	_ "github.com/sakirsensoy/genv/dotenv/autoload" // import for automatic loading of .env config
)

var (
	serverAddress   string
	databaseDSN     string
	baseURL         string
	fileStoragePath string
	storageType     string

	config *ConfigT
)

// Initialize
// Performs initial initialisation of configurations based on flags and/or env variables.
// If the initialisation was already performed somewhere else, returns the config object of type ConfigT
func Initialize() *ConfigT {
	if config != nil {
		return config
	}
	fmt.Println("cfg initialized")
	flag.StringVar(&databaseDSN, "d", "", "help message for flagname")
	flag.StringVar(&serverAddress, "a", "localhost:8080", "help message for flagname")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "help message for flagname")
	flag.StringVar(&fileStoragePath, "f", "./data/data", "help message for flagname")
	flag.Parse()

	databaseDSN = genv.Key("DATABASE_DSN").Default(databaseDSN).String()
	fmt.Println("DSN:\t", databaseDSN)
	if databaseDSN == "" {
		fmt.Println("FILE SELECTED AS STORAGE TYPE DUE TO NO DSN")
		storageType = "file"
	} else {
		storageType = genv.Key("STORAGE_TYPE").Default("database").String()
		fmt.Println(storageType, "SELECTED AS STORAGE TYPE")
	}

	fmt.Println(serverAddress, baseURL, fileStoragePath)
	fmt.Println(
		genv.Key("SERVER_ADDRESS").Default("NO SUCH FIELD").String(),
		genv.Key("BASE_URL").Default("NO SUCH FIELD").String(),
		genv.Key("FILE_STORAGE_PATH").Default("NO SUCH FIELD").String())
	config = &ConfigT{
		Shortener: ShortenerT{
			Secret:    genv.Key("SECRET").Default("12345qwerty").String(),
			URLLength: genv.Key("SHORT_URL_LENGTH").Default(10).Int(),
		},

		Server: ServerT{
			Address:        genv.Key("SERVER_ADDRESS").Default(serverAddress).String(),
			Port:           genv.Key("SERVER_PORT").Default("8080").String(),
			BaseURL:        genv.Key("BASE_URL").Default(baseURL).String(),
			CookieLifetime: 30 * 24 * 60 * 60,
		},

		Storage: StorageT{
			AutosaveInterval: genv.Key("STORAGE_AUTOSAVE_INTERVAL").Default(-1).Int(),
			SavePath:         genv.Key("FILE_STORAGE_PATH").Default(fileStoragePath).String(),
			StorageType:      genv.Key("STORAGE_TYPE").Default(storageType).String(),
			DatabaseDSN:      databaseDSN,
		},
	}
	return config
}
