package cfg

import "flag"

func WithStorage(data StorageT) ConfigOption {
	if !flag.Parsed() {
		flag.Parse()
	}

	return func(c *ConfigT) *ConfigT {
		c.Storage = data
		return c
	}
}

func WithShortener(data ShortenerT) ConfigOption {
	if !flag.Parsed() {
		flag.Parse()
	}

	return func(c *ConfigT) *ConfigT {
		c.Shortener = data
		return c
	}
}

func WithServer(data ServerT) ConfigOption {
	if !flag.Parsed() {
		flag.Parse()
	}

	return func(c *ConfigT) *ConfigT {
		c.Server = data
		return c
	}
}
