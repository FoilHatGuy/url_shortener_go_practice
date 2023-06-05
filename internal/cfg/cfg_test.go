package cfg

import (
	"flag"
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
		StorageType:      "STORAGE_TYPE_VALUE",
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
		storageType             = "STORAGE_TYPE_VALUE"
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
	t.Setenv("STORAGE_TYPE", storageType)
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
			StorageType:      storageType,
			DatabaseDSN:      databaseDsn,
		},
	}

	config2 := New(FromDefaults(),
		FromEnv(),
	)
	s.Assert().Equal(*config1, *config2)
}

func (s *ConfigTestSuite) TestFromFlags() {
	address := "Address"
	baseURL := "BaseURL"
	databaseDSN := "DatabaseDSN"
	savePath := "SavePath"

	err := flag.Set("a", address)
	s.Assert().NoError(err)
	err = flag.Set("b", baseURL)
	s.Assert().NoError(err)
	err = flag.Set("d", databaseDSN)
	s.Assert().NoError(err)
	err = flag.Set("f", savePath)
	s.Assert().NoError(err)

	config1 := New(FromDefaults(),
		FromFlags(),
	)

	s.Assert().Equal(config1.Server.Address, address)
	s.Assert().Equal(config1.Server.BaseURL, baseURL)
	s.Assert().Equal(config1.Storage.DatabaseDSN, databaseDSN)
	s.Assert().Equal(config1.Storage.SavePath, savePath)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
