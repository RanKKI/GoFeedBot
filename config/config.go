package config

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
)

type Config struct {
    Token    string `json:"token"`
    Debug    bool   `json:"debug"`
    Proxy    string `json:"proxy"`
    Interval int    `json:"interval"`
    Client   *http.Client
}

func loadFile(filename string) []byte {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    return data
}

func setupClient(config *Config) {
    config.Client = &http.Client{}
    if config.Proxy != "" {
        proxyUrl, err := url.Parse(config.Proxy)
        if err != nil {
            panic(err)
        }
        log.Printf("Using proxy %s", config.Proxy)
        config.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
    }
}

func parseConfig(data []byte) Config {
    config := Config{}
    if err := json.Unmarshal(data, &config); err != nil {
        panic(err)
    }
    setupClient(&config)
    return config
}

func LoadConfig(filename string) Config {
    log.Println("Loading config file", filename)
    data := loadFile(filename)
    return parseConfig(data)
}
