package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/plugin/rpcplugin"
)

type HelloWorldPlugin struct{
	api plugin.API
	configuration atomic.Value
	hogeru string
}

func (p *HelloWorldPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("Hello %s!", p.hogeru))
}

func (p *HelloWorldPlugin) OnActivate(api plugin.API) error {
	p.api = api
	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	config := p.configuration.Load().(*Configuration)
	if err := config.IsValid(); err != nil {
		return err
	}

	p.hogeru = config.Hogeru
	return nil
}

func (p *HelloWorldPlugin)OnConfigurationChange() error {
	var configuration Configuration
	err := p.api.LoadPluginConfiguration(&configuration)
	p.configuration.Store(&configuration)
	return err
}

func main() {
	rpcplugin.Main(&HelloWorldPlugin{})
}