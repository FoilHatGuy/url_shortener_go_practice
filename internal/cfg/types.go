package cfg

type fileJSONT struct {
	ServerAddressHTTP  string `json:"server_address_http"`
	ServerAddressGRPC  string `json:"server_address_grpc"`
	ServerBaseURL      string `json:"base_url"`
	ServerEnableHTTPS  bool   `json:"enable_https"`
	StorageSavePath    string `json:"file_storage_path"`
	StorageDatabaseDSN string `json:"database_dsn"`
	TrustedSubnet      string `json:"trusted_subnet"`
}

// ConfigT
// Parent structure for all configuration structs. provides config separation into
// ShortenerT, ServerT and StorageT for the ease of use
type ConfigT struct {
	Shortener *ShortenerT
	Server    *ServerT
	Storage   *StorageT
}

// ShortenerT
// Contains the required URLLength to return to user and Secret for cookie encryption.
// Can be accessed via a structure of type ConfigT
type ShortenerT struct {
	Secret    string `default:"12345qwerty"`
	URLLength int    `default:"10"`
}

// ServerT
// Contains server launch configuration, namely server Address, Port and BaseURL used for URL construction.
// Additionally, stores CookieLifetime used in session IDs.
// Can be accessed via a structure of type ConfigT
type ServerT struct {
	AddressHTTP string `default:"localhost:8080"`
	AddressGRPC string `default:"localhost:8081"`

	Port           string `default:"8080"`
	BaseURL        string `default:"http://localhost:8080"`
	CookieLifetime int    `default:"30 * 24 * 60 * 60"`
	IsHTTPS        bool   `default:"false"`
	TrustedSubnet  string `default:"127.0.0.1/32"`
}

// StorageT
// Contains database configuration. DatabaseDSN contains string used for connection to Postgres DB.
// Performs auto saves to SavePath every AutosaveInterval seconds
// Can be accessed via a structure of type ConfigT
type StorageT struct {
	AutosaveInterval int    `default:"-1"`
	SavePath         string `default:"./data/data"`
	DatabaseDSN      string `default:""`
}
