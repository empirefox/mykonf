package mykonf

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

// Load does:
// 1. load yaml
// 2. set with env
// 3. load defaults
func Load(envPrefix string, conf any) error {
	return LoadPath(envPrefix, ConfigPath(envPrefix), conf)
}

// LoadPath does:
// 1. load yaml
// 2. set with env
// 3. load defaults
func LoadPath(envPrefix, path string, conf any) error {
	k := koanf.New(".")

	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		err = k.Load(Provider(path), yaml.Parser())
		if err != nil {
			return err
		}
	}

	envToKey := EnvToKey(conf, "yaml")
	err = k.Load(env.Provider(".", env.Opt{
		Prefix: envPrefix,
		TransformFunc: func(k, v string) (string, any) {
			e := strings.TrimPrefix(k, envPrefix)
			key, ok := envToKey[e]
			if !ok {
				return strings.ToLower(e), v
			}
			return key, v
		},
	}), nil)
	if err != nil {
		return err
	}

	err = k.UnmarshalWithConf("", conf, koanf.UnmarshalConf{Tag: "yaml",
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				StringToJsonHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.TextUnmarshallerHookFunc(),
			),
			Metadata:         nil,
			WeaklyTypedInput: true,
		}})
	if err != nil {
		return err
	}

	return defaults.Set(conf)
}

const defaultConfigEnv = "SERVER_CONFIG"
const defaultConfigPath = "config.yaml"

func ConfigPath(envPrefix string) string {
	p := os.Getenv(envPrefix + defaultConfigEnv)
	if p != "" {
		return p
	}
	return defaultConfigPath
}

func StringToJsonHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		var r reflect.Value
		switch t.Kind() {
		case reflect.Map, reflect.Struct:
			r = reflect.New(t)
		default:
			return data, nil
		}
		v := data.(string)
		if v != "" {
			err := json.Unmarshal([]byte(v), r.Interface())
			if err != nil {
				return nil, err
			}
		}

		return r.Elem(), nil
	}
}
