package config

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "time"
)

type Config struct {
    Token            string `json:"token"`
    Debug            bool   `json:"debug"`
    Interval         string `json:"interval"`
    MaxContentLength int    `json:"max_content_length"`
    Client           *http.Client
}

func (config *Config) LoadFromEnv() {
    config.Client = &http.Client{}
    if os.Getenv("debug") == "1" {
        config.Debug = true
    }
    config.SetToken(os.Getenv("token"))
    config.SetProxy(os.Getenv("proxy"))
    config.SetInterval(os.Getenv("interval"))
}

func (config *Config) SetToken(token string) {
    if token == "EMPTY" || token == "" {
        panic("You must provide a telegram bot token")
    }
    config.Token = token
}

func (config *Config) SetProxy(proxy string) {
    if proxy == "EMPTY" || proxy == "" {
        return
    }
    proxyUrl, err := url.Parse(proxy)
    if err != nil {
        panic(err)
    }
    log.Printf("Using proxy %s", proxy)
    config.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
}

func (config *Config) SetInterval(interval string) {
    _, err := time.ParseDuration(interval)
    if err != nil {
        panic(fmt.Sprintf("invaild interval %s\nsee %s", config.Interval, "https://golang.org/pkg/time/#ParseDuration"))
    }
    config.Interval = interval
}

func LoadConfig() Config {
    config := Config{}
    config.LoadFromEnv()
    return config
}
