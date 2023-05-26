package cfg

// ConfigT
// Parent structure for all configuration structs. provides config separation into
// ShortenerT, ServerT and StorageT for the ease of use
type ConfigT struct {
	Shortener ShortenerT
	Server    ServerT
	Storage   StorageT
}

// ShortenerT
// Contains the required URLLength to return to user and Secret for cookie encryption.
// Can be accessed via a structure of type ConfigT
type ShortenerT struct {
	Secret    string
	URLLength int
}

// ServerT
// Contains server launch configuration, namely server Address, Port and BaseURL used for URL construction.
// Additionally, stores CookieLifetime used in session IDs.
// Can be accessed via a structure of type ConfigT
type ServerT struct {
	Address        string
	Port           string
	BaseURL        string
	CookieLifetime int
}

// StorageT
// Contains database configuration. DatabaseDSN contains string used for connection to Postgres DB.
// If said variable is not provided StorageType defaults to 'file',
// that performs auto saves to SavePath every AutosaveInterval seconds
// Can be accessed via a structure of type ConfigT
type StorageT struct {
	AutosaveInterval int
	SavePath         string
	StorageType      string
	DatabaseDSN      string
}
