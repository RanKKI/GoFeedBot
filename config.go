package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
)

type Config struct {
    Token  string `json:"token"`
    Debug  bool   `json:"debug"`
    Proxy  string `json:"proxy"`
    Client *http.Client
}

func loadConfig() *Config {
    log.Println("Loading config file")
    data, err := ioutil.ReadFile("./config.json")
    if err != nil {
        panic(err)
    }
    config := &Config{}
    if err := json.Unmarshal(data, config); err != nil {
        panic(err)
    }

    config.Client = &http.Client{}

    if config.Proxy != "" {
        proxyUrl, err := url.Parse(config.Proxy)
        if err != nil {
            panic(err)
        }
        log.Printf("Using proxy %s", config.Proxy)
        config.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
    }

    return config
}
