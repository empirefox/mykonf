package mykonf

import (
	"testing"
)

const envPrefix = "SFTPGO_HOOK_"

// Config all config in this
type Config struct {
	Listen string `yaml:"listen" default:":8080"`

	BaseHomeDir string `yaml:"base_home_dir" default:"/srv/sftpgo/data"`
	Uid         int    `yaml:"uid" default:"1000"`
	Gid         int    `yaml:"gid" default:"1000"`
	// 1G = 1073741824
	Quota int64 `yaml:"quota"`

	Gitea struct {
		Url    string `yaml:"url"`
		ApiKey string `yaml:"api_key"`
	} `yaml:"gitea"`
}

var conf *Config

func Get() *Config {
	if conf == nil {
		conf = new(Config)
		Load(envPrefix, conf)
	}
	return conf
}

func TestLoadConfig(t *testing.T) {
	t.Setenv("SFTPGO_HOOK_SERVER_CONFIG", "../config-test.yaml")
	t.Setenv("SFTPGO_HOOK_UID", "2000")
	t.Setenv("SFTPGO_HOOK_GITEA_URL", "url_here")
	t.Setenv("SFTPGO_HOOK_GITEA_API_KEY", "api_key_here")

	conf := Get()

	if conf.Listen != ":8080" {
		t.Fatalf("conf.Listen != defaultListen")
	}
	if conf.Uid != 2000 {
		t.Fatalf("conf.Gitea.Url should be '2000', but got: '%d'",
			conf.Uid)
	}
	if conf.Gitea.Url != "url_here" {
		t.Fatalf("conf.Gitea.Url should be 'url_here', but got: '%s'",
			conf.Gitea.Url)
	}
	if conf.Gitea.ApiKey != "api_key_here" {
		t.Fatalf("conf.Gitea.ApiKey should be 'api_key_here', but got: '%s'",
			conf.Gitea.ApiKey)
	}
}
