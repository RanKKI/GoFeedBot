package config

import (
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestData(t *testing.T) {
    ass := assert.New(t)
    // JSON Data
    _ = os.Setenv("token", "123")
    _ = os.Setenv("interval", "11m1s")
    _ = os.Setenv("proxy", "http://127.0.0.1:1234")

    config := Config{}
    config.LoadFromEnv()

    ass.Equal(config.Token, "123")
    ass.Equal(config.Interval, "11m1s")
    ass.Equal(config.Debug, false)

    _ = os.Setenv("debug", "1")
    config.LoadFromEnv()
    ass.Equal(config.Debug, true)
}

func TestInterval(t *testing.T) {
    ass := assert.New(t)
    config := Config{}
    config.SetInterval("")
    ass.Equal(config.Interval, "30m")
    ass.Panics(func() {
        config.SetInterval("10,")
    })
    ass.NotPanics(func() {
        config.SetInterval("10m")
    })
    ass.NotPanics(func() {
        config.SetInterval("10m30s")
    })
}

func TestProxy(t *testing.T) {
    ass := assert.New(t)
    config := Config{}
    config.SetProxy("")

    // Invalid proxy
    ass.Panics(func() {
        config.SetProxy("kkk://123:-1")
    })

    ass.NotPanics(func() {
        config.SetProxy("http://127.0.0.1:1234")
    })

}
