package cfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mcuadros/go-defaults"
	"github.com/sakirsensoy/genv"
	_ "github.com/sakirsensoy/genv/dotenv/autoload" // import for automatic loading of .env config
)

var (
	serverAddress   string
	databaseDSN     string
	baseURL         string
	fileStoragePath string
	isHTTPS         bool
	configPath      string
)

func init() {
	flag.StringVar(&serverAddress, "a", "", "Server running address")
	flag.StringVar(&databaseDSN, "d", "", "BaseURL for shortened links")
	flag.StringVar(&baseURL, "b", "", "DSN for database")
	flag.StringVar(&fileStoragePath, "f", "", "File storage path")
	flag.BoolVar(&isHTTPS, "s", false, "run server as HTTPS")
	flag.StringVar(&configPath, "c", "", "path to JSON config")
}

// ConfigOption
// Various options that can be used in New() to set up configs
type ConfigOption func(*ConfigT) *ConfigT

// New
// Accepts config creation options from package.
// Returns the basic config with default values of ConfigT.
func New(opts ...ConfigOption) *ConfigT {
	cfg := &ConfigT{
		Server:    ServerT{},
		Shortener: ShortenerT{},
		Storage:   StorageT{},
	}

	if !flag.Parsed() {
		flag.Parse()
	}

	for _, o := range opts {
		cfg = o(cfg)
	}

	return cfg
}

// FromDefaults
//
//	@Description: Initializes default values of type ConfigT
func FromDefaults() ConfigOption {
	return func(c *ConfigT) *ConfigT {
		defaults.SetDefaults(c)
		return c
	}
}

// FromJSON
//
//	@Description: Overwrites existing values with values from environment (if present)
func FromJSON() ConfigOption {
	return func(c *ConfigT) *ConfigT {
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Printf("opening JSON failed. Details: %v", err)
			return nil
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		data := make([]byte, 0)
		_, err = io.ReadFull(file, data)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		tempConfig := fileJSONT{}
		err = json.Unmarshal(data, &tempConfig)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if tempConfig.ServerAddress != "" {
			c.Server.Address = tempConfig.ServerAddress
		}
		if tempConfig.ServerBaseURL != "" {
			c.Server.BaseURL = tempConfig.ServerBaseURL
		}
		if tempConfig.ServerEnableHTTPS {
			c.Server.IsHTTPS = true
		}
		if tempConfig.StorageSavePath != "" {
			c.Storage.SavePath = tempConfig.StorageSavePath
		}
		if tempConfig.StorageDatabaseDsn != "" {
			c.Storage.DatabaseDSN = tempConfig.StorageDatabaseDsn
		}

		return c
	}
}

// FromEnv
//
//	@Description: Overwrites existing values with values from environment (if present)
func FromEnv() ConfigOption {
	return func(c *ConfigT) *ConfigT {
		c = &ConfigT{
			Shortener: ShortenerT{
				Secret:    genv.Key("SECRET").Default(c.Shortener.Secret).String(),
				URLLength: genv.Key("SHORT_URL_LENGTH").Default(c.Shortener.URLLength).Int(),
			},

			Server: ServerT{
				Address:        genv.Key("SERVER_ADDRESS").Default(c.Server.Address).String(),
				Port:           genv.Key("SERVER_PORT").Default(c.Server.Port).String(),
				BaseURL:        genv.Key("BASE_URL").Default(c.Server.BaseURL).String(),
				CookieLifetime: genv.Key("SERVER_COOKIE_LIFETIME").Default(c.Server.CookieLifetime).Int(),
				IsHTTPS:        genv.Key("ENABLE_HTTPS").Default(c.Server.IsHTTPS).Bool(),
			},

			Storage: StorageT{
				AutosaveInterval: genv.Key("STORAGE_AUTOSAVE_INTERVAL").Default(c.Storage.AutosaveInterval).Int(),
				SavePath:         genv.Key("FILE_STORAGE_PATH").Default(c.Storage.SavePath).String(),
				DatabaseDSN:      genv.Key("DATABASE_DSN").Default(c.Storage.DatabaseDSN).String(),
			},
		}
		return c
	}
}

// FromFlags
//
//	@Description: Overwrites existing values with values from flags (if present)
func FromFlags() ConfigOption {
	return func(c *ConfigT) *ConfigT {
		if serverAddress != "" {
			c.Server.Address = serverAddress
		}
		if baseURL != "" {
			c.Server.BaseURL = baseURL
		}
		if isHTTPS {
			c.Server.IsHTTPS = true
		}
		if databaseDSN != "" {
			c.Storage.DatabaseDSN = databaseDSN
		}
		if fileStoragePath != "" {
			c.Storage.SavePath = fileStoragePath
		}
		return c
	}
}
