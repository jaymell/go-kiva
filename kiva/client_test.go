package kiva

import "testing"

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config.BaseURL != "http://api.kivaws.org" {
		t.Error("wrong default base url")
	}
	if config.AppID != "" {
		t.Error("app id shoud be nil by default")
	}
}

func TestNewClientDefaultConfig(t *testing.T) {
	cli := New(nil)
	if cli.baseURL.String() != "http://api.kivaws.org" {
		t.Error("wrong default base url")
	}
	if cli.appID != "" {
		t.Error("app id shoud be nil by default")
	}
}

func TestNewClientFullConfig(t *testing.T) {
	clientConfig := Config{
		BaseURL: "http://api.fakesite.com",
		AppID: "com.fakesite.app",
	}
	cli := New(&clientConfig)
	if cli.baseURL.String() != "http://api.fakesite.com" {
		t.Error("wrong base url")
	}
	if cli.appID != "com.fakesite.app" {
		t.Error("wrong app id")
	}
}

func TestNewClientBaseURLOnly(t *testing.T) {
	clientConfig := Config{
		BaseURL: "http://api.fakesite.com",
	}
	cli := New(&clientConfig)
	if cli.baseURL.String() != "http://api.fakesite.com" {
		t.Error("wrong base url")
	}
	if cli.appID != "" {
		t.Error("wrong app id")
	}
}

func TestNewClientAppIDOnly(t *testing.T) {
	clientConfig := Config{
		AppID: "com.fakesite.app",
	}
	cli := New(&clientConfig)
	if cli.baseURL.String() != "http://api.kivaws.org" {
		t.Error("wrong base url")
	}
	if cli.appID != "com.fakesite.app" {
		t.Error("wrong app id")
	}
}

