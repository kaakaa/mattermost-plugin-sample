package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/plugin/rpcplugin"
)

type HelloWorldPlugin struct{
	api plugin.API
	configuration atomic.Value
	hogeru string
}

func (p *HelloWorldPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/set":
		p.handleSet(w,r)
	case "/get":
		p.handleGet(w,r)
	case "/delete":
		p.handleDelete(w,r)
	case "/hello":
		fmt.Fprintf(w, fmt.Sprintf("Hello %s!", p.hogeru))
	}
}

type KeyValue struct {
	key string
	value string
}

func (p HelloWorldPlugin) handleSet(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var data KeyValue
	if err := decoder.Decode(&data); err != nil {
		fmt.Fprintf(w, "Request err: " + err.Error())
	}

	if err := p.api.KeyValueStore().Set(data.key, []byte(data.value)); err != nil {
		fmt.Fprintf(w, "KeyValue set err: " + err.Error())		
	}
}

func (p HelloWorldPlugin) handleGet(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var data KeyValue
	if err := decoder.Decode(&data); err != nil {
		fmt.Fprintf(w, "Request err: " + err.Error())
	}

	b, err := p.api.KeyValueStore().Get(data.key)
	if err != nil {
		fmt.Fprintf(w, "KeyValue get err: " + err.Error())		
	} else {
		fmt.Fprintf(w, "Read: " + string(b))
	}
}

func (p HelloWorldPlugin) handleDelete(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var data KeyValue
	if err := decoder.Decode(&data); err != nil {
		fmt.Fprintf(w, "Request err: " + err.Error())
	}

	if err := p.api.KeyValueStore().Delete(data.key,); err != nil {
		fmt.Fprintf(w, "KeyValue delete err: " + err.Error())		
	}
}

func (p HelloWorldPlugin) ExecuteCommand(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	log.Printf("%v", args)
	return &model.CommandResponse{
		ResponseType: "in_channel",
		Text: "test",
	}, nil
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

	t, err := p.api.GetTeamByName("sandbox")
	if err != nil {
		return err
	}
	log.Printf("$$$$$$$$$$jfldjsa;fjd;$$$$$$$")
	command := &model.Command{
		Id: "hoge_plugin_command",
		TeamId: t.Id,
		Trigger: "hoge",
		Method: "POST",
		URL: "/plugins/hoge1/hello",
		AutoComplete: true,
		AutoCompleteDesc: "hohohohogehgoehogheoghoehoge",
		AutoCompleteHint: "ohogehogehogehoge",
	}
	if err := p.api.RegisterCommand(command); err != nil {
		return err
	}
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