package cfg

import "github.com/sakirsensoy/genv"

type shortCfg struct {
	UrlLength int
}

type serverCfg struct {
	Host string
	Port string
}
type routerCfg struct {
	BaseURL string
}

var Shortener = &shortCfg{
	UrlLength: genv.Key("SHORT_URL_LENGTH").Default(10).Int(),
}

var Server = &serverCfg{
	Host: genv.Key("SERVER_HOST").Default("localhost").String(),
	Port: genv.Key("SERVER_PORT").Default("8080").String(),
}
var Router = &routerCfg{
	BaseURL: genv.Key("BASE_URL").Default("/").String(),
}
