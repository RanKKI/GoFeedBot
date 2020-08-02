package config

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
    ass.NotEmpty(loadFile("../config.example.json"))
}

func TestData(t *testing.T) {
    ass := assert.New(t)
    // JSON Data
    data := "{\n  \"token\": \"test_token_abcde\",\n  \"debug\": false,\n  \"proxy\": \"\",\n  \"interval\": \"10m\",\n  \"max_content_length\": 100\n}"

    // Test given json data
    config := parseConfig([]byte(data))

    // Invalid json format
    ass.Panics(func() {
        parseConfig([]byte("{'token':'test_token'}}"))
    })

    ass.Equal("test_token_abcde", config.Token)
    ass.Equal(false, config.Debug)

    data = "{\n  \"token\": \"\",\n  \"debug\": false,\n  \"proxy\": \"\",\n  \"interval\": \"10h12m\",\n  \"max_content_length\": -12\n}"
    ass.Panics(func() {
        parseConfig([]byte(data))
    })
}

func TestInterval(t *testing.T) {
    ass := assert.New(t)
    data := "{\n  \"token\": \"\",\n  \"debug\": false,\n  \"proxy\": \"\",\n  \"interval\": \"10\",\n  \"max_content_length\": 100\n}"
    ass.Panics(func() {
        parseConfig([]byte(data))
    })

    data = "{\n  \"token\": \"\",\n  \"debug\": false,\n  \"proxy\": \"\",\n  \"interval\": \"10m\",\n  \"max_content_length\": 100\n}"
    ass.NotPanics(func() {
        parseConfig([]byte(data))
    })

    data = "{\n  \"token\": \"\",\n  \"debug\": false,\n  \"proxy\": \"\",\n  \"interval\": \"10h12m\",\n  \"max_content_length\": 100\n}"
    ass.NotPanics(func() {
        parseConfig([]byte(data))
    })
}

func TestProxy(t *testing.T) {
    ass := assert.New(t)
    config := Config{}
    // Invalid proxy
    ass.Panics(func() {
        config.Proxy = "kkk://123:-1"
        setupClient(&config)
    })

    ass.NotPanics(func() {
        config.Proxy = "http://127.0.0.1:1234"
        setupClient(&config)
    })

}
