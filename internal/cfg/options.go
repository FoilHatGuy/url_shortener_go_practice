package cfg

import "flag"

// WithStorage
//
//	@Description: replaces the StorageT config with supplied one
//	@param data
func WithStorage(data StorageT) ConfigOption {
	if !flag.Parsed() {
		flag.Parse()
	}

	return func(c *ConfigT) *ConfigT {
		c.Storage = data
		return c
	}
}

// WithShortener
//
//	@Description: replaces the ShortenerT config with supplied one
//	@param data
func WithShortener(data ShortenerT) ConfigOption {
	if !flag.Parsed() {
		flag.Parse()
	}

	return func(c *ConfigT) *ConfigT {
		c.Shortener = data
		return c
	}
}

// WithServer
//
//	@Description: replaces the ServerT config with supplied one
//	@param data
func WithServer(data ServerT) ConfigOption {
	if !flag.Parsed() {
		flag.Parse()
	}

	return func(c *ConfigT) *ConfigT {
		c.Server = data
		return c
	}
}
