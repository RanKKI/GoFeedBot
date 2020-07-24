package main

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestLoadConfig(t *testing.T) {
    ass := assert.New(t)

    // File should not exists
    ass.Panics(func() {
        loadFile("some_file_lalalal")
    })

    // config.example.json must exists,
    ass.NotEmpty(loadFile("./config.example.json"))

    // Test given json data
    config := parseConfig([]byte("{\n  \"token\": \"test_token\",\n  \"debug\": false,\n  \"proxy\": \"\"\n}"))

    // Invalid json format
    ass.Panics(func() {
        parseConfig([]byte("{'token':'test_token'}}"))
    })

    ass.Equal("test_token", config.Token)
    ass.Equal(false, config.Debug)

    // Invalid proxy
    ass.Panics(func() {
        config.Proxy = "kkk://123:-1"
        setupClient(config)
    })

    ass.NotPanics(func() {
        config.Proxy = "http://127.0.0.1:1234"
        setupClient(config)
    })
}
