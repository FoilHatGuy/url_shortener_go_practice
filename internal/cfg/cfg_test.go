package cfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/mcuadros/go-defaults"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestNew() {
	config1 := New(FromDefaults())
	config2 := new(ConfigT)
	defaults.SetDefaults(config2)
	s.Assert().Equal(config1, config2)
}

func (s *ConfigTestSuite) TestWithServer() {
	source := ServerT{
		Address:        "SERVER_ADDRESS_VALUE",
		Port:           "SERVER_PORT_VALUE",
		BaseURL:        "BASE_URL_VALUE",
		CookieLifetime: 20,
	}
	config1 := New(
		FromDefaults(),
		WithServer(source))

	s.Assert().Equal(config1.Server, source)
}

func (s *ConfigTestSuite) TestWithShortener() {
	source := ShortenerT{
		Secret:    "SECRET_VALUE",
		URLLength: 10,
	}
	config1 := New(
		FromDefaults(),
		WithShortener(source))

	s.Assert().Equal(config1.Shortener, source)
}

func (s *ConfigTestSuite) TestWithStorage() {
	source := StorageT{
		AutosaveInterval: 10,
		SavePath:         "FILE_STORAGE_PATH_VALUE",
		DatabaseDSN:      "DATABASE_DSN_VALUE",
	}
	config1 := New(
		FromDefaults(),
		WithStorage(source))

	s.Assert().Equal(config1.Storage, source)
}

func (s *ConfigTestSuite) TestFromEnv() {
	t := s.T()

	const (
		secret                  = "SECRET_VALUE"
		shortURLLength          = 10
		serverAddress           = "SERVER_ADDRESS_VALUE"
		serverPort              = "SERVER_PORT_VALUE"
		baseURL                 = "BASE_URL_VALUE"
		serverCookieLifetime    = 20
		storageAutosaveInterval = 30
		fileStoragePath         = "FILE_STORAGE_PATH_VALUE"
		databaseDsn             = "DATABASE_DSN_VALUE"
	)
	t.Setenv("SECRET", secret)
	t.Setenv("SHORT_URL_LENGTH", strconv.Itoa(shortURLLength))
	t.Setenv("SERVER_ADDRESS", serverAddress)
	t.Setenv("SERVER_PORT", serverPort)
	t.Setenv("BASE_URL", baseURL)
	t.Setenv("SERVER_COOKIE_LIFETIME", strconv.Itoa(serverCookieLifetime))
	t.Setenv("STORAGE_AUTOSAVE_INTERVAL", strconv.Itoa(storageAutosaveInterval))
	t.Setenv("FILE_STORAGE_PATH", fileStoragePath)
	t.Setenv("DATABASE_DSN", databaseDsn)

	config1 := &ConfigT{
		Shortener: ShortenerT{
			Secret:    secret,
			URLLength: shortURLLength,
		},

		Server: ServerT{
			Address:        serverAddress,
			Port:           serverPort,
			BaseURL:        baseURL,
			CookieLifetime: serverCookieLifetime,
		},

		Storage: StorageT{
			AutosaveInterval: storageAutosaveInterval,
			SavePath:         fileStoragePath,
			DatabaseDSN:      databaseDsn,
		},
	}

	config2 := New(FromDefaults(),
		FromEnv(),
	)
	s.Assert().Equal(*config1, *config2)
}

func (s *ConfigTestSuite) TestFromFlags() {
	const (
		address     = "Address"
		baseURL     = "BaseURL"
		databaseDSN = "DatabaseDSN"
		savePath    = "SavePath"
		isHTTPS     = true
	)

	err := flag.Set("a", address)
	s.Assert().NoError(err)
	err = flag.Set("b", baseURL)
	s.Assert().NoError(err)
	err = flag.Set("d", databaseDSN)
	s.Assert().NoError(err)
	err = flag.Set("f", savePath)
	s.Assert().NoError(err)
	err = flag.Set("s", fmt.Sprintf("%t", isHTTPS))
	s.Assert().NoError(err)

	config1 := New(FromDefaults(),
		FromFlags(),
	)

	s.Assert().Equal(config1.Server.Address, address)
	s.Assert().Equal(config1.Server.BaseURL, baseURL)
	s.Assert().Equal(config1.Server.IsHTTPS, isHTTPS)
	s.Assert().Equal(config1.Storage.DatabaseDSN, databaseDSN)
	s.Assert().Equal(config1.Storage.SavePath, savePath)
}

func (s *ConfigTestSuite) TestFromJSONFile() {
	const filePath = "./test.json"
	origin := fileJSONT{
		ServerAddress:      "1",
		ServerBaseURL:      "2",
		ServerEnableHTTPS:  true,
		StorageSavePath:    "3",
		StorageDatabaseDSN: "4",
	}

	New(
		FromJSON(),
	) // cause an error

	file, _ := json.MarshalIndent(origin, "", "\t")

	_ = os.WriteFile(filePath, file, 0o600)
	defer func() {
		err := os.Remove(filePath)
		s.Assert().NoError(err)
	}()

	err := flag.Set("c", filePath)
	s.Assert().NoError(err)

	data, err := os.ReadFile(configPath)
	s.Assert().NoError(err)
	s.Assert().NotEmpty(data)

	config1 := New(
		FromJSON(),
	)

	s.Assert().Equal(config1.Server.Address, origin.ServerAddress)
	s.Assert().Equal(config1.Server.BaseURL, origin.ServerBaseURL)
	s.Assert().Equal(config1.Server.IsHTTPS, origin.ServerEnableHTTPS)
	s.Assert().Equal(config1.Storage.DatabaseDSN, origin.StorageDatabaseDSN)
	s.Assert().Equal(config1.Storage.SavePath, origin.StorageSavePath)

	// broken file
	file, _ = json.MarshalIndent(origin, "\"", "\"")
	//nolint:gosec
	_ = os.WriteFile(filePath, file, 0o300)
	New(
		FromJSON(),
	)
}

func (s *ConfigTestSuite) TestParseFlagsTwice() {
	const (
		address     = "Address"
		baseURL     = "BaseURL"
		databaseDSN = "DatabaseDSN"
		savePath    = "SavePath"
		isHTTPS     = true
	)

	err := flag.Set("a", address)
	s.Assert().NoError(err)
	err = flag.Set("b", baseURL)
	s.Assert().NoError(err)
	err = flag.Set("d", databaseDSN)
	s.Assert().NoError(err)
	err = flag.Set("f", savePath)
	s.Assert().NoError(err)
	err = flag.Set("s", fmt.Sprintf("%t", isHTTPS))
	s.Assert().NoError(err)

	config1 := New(
		FromFlags(),
	)

	config2 := New(
		FromFlags(),
	)

	s.Assert().Equal(config1.Server.Address, config2.Server.Address)
	s.Assert().Equal(config1.Server.BaseURL, config2.Server.BaseURL)
	s.Assert().Equal(config1.Server.IsHTTPS, config2.Server.IsHTTPS)
	s.Assert().Equal(config1.Storage.DatabaseDSN, config2.Storage.DatabaseDSN)
	s.Assert().Equal(config1.Storage.SavePath, config2.Storage.SavePath)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
