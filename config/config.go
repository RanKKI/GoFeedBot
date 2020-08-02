package config

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "time"
)

type Config struct {
    Token            string `json:"token"`
    Debug            bool   `json:"debug"`
    Proxy            string `json:"proxy"`
    Interval         string `json:"interval"`
    MaxContentLength int    `json:"max_content_length"`
    Client           *http.Client
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
    _, err := time.ParseDuration(config.Interval)
    if err != nil {
        panic(fmt.Sprintf("invaild interval %s\nsee %s", config.Interval, "https://golang.org/pkg/time/#ParseDuration"))
    }
    if config.MaxContentLength < 0 {
        panic(fmt.Sprintf("Max Content Length must >= 0 not %d", config.MaxContentLength))
    }
    setupClient(&config)
    return config
}

func LoadConfig(filename string) Config {
    log.Println("Loading config file", filename)
    data := loadFile(filename)
    return parseConfig(data)
}
