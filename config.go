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

func loadFile(filename string) []byte {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    return data
}

func setupClient(config *Config) *http.Client {
    client := &http.Client{}
    if config.Proxy != "" {
        proxyUrl, err := url.Parse(config.Proxy)
        if err != nil {
            panic(err)
        }
        log.Printf("Using proxy %s", config.Proxy)
        client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
    }
    return client
}

func parseConfig(data []byte) *Config {
    config := &Config{}
    if err := json.Unmarshal(data, config); err != nil {
        panic(err)
    }
    config.Client = setupClient(config)
    return config
}

func loadConfig() *Config {
    log.Println("Loading config file")
    data := loadFile("./config.json")
    config := parseConfig(data)
    return config
}
