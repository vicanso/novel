package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// ENV 配置的环境变量
var env = os.Getenv("GO_ENV")
var configPath = os.Getenv("CONFIG")
var viperInitTest = os.Getenv("VIPER_INIT_TEST")

// 初始化配置
func viperInit(path string) error {
	v := viper.New()
	v.SetConfigName("default")
	v.AddConfigPath(".")
	v.AddConfigPath(path)
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	configs := v.AllSettings()
	for k, v := range configs {
		viper.SetDefault(k, v)
	}
	if env != "" {
		viper.SetConfigName(env)
		viper.AddConfigPath(".")
		viper.AddConfigPath(path)
		viper.SetConfigType("yml")
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
	}
	return nil
}

func setDefaultForTest() {
	viper.Set("locationByIP", "http://ip.taobao.com/service/getIpInfo.php")
	viper.Set("redis", "127.0.0.1:6379")
}

func init() {
	if viperInitTest != "" {
		setDefaultForTest()
		return
	}
	if configPath == "" {
		runPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		configPath = runPath + "/asset"
	}
	err := viperInit(configPath)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

// GetENV get the go env
func GetENV() string {
	return env
}

// GetInt viper get int
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetIntDefault get int with default value
func GetIntDefault(key string, defaultValue int) int {
	v := viper.GetInt(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetString viper get string
func GetString(key string) string {
	return viper.GetString(key)
}

// GetStringDefault get string with default value
func GetStringDefault(key, defaultValue string) string {
	v := viper.GetString(key)
	if v != "" {
		return v
	}
	return defaultValue
}

// GetDuration viper get duration
func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

// GetDurationDefault get duration with default value
func GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	v := viper.GetDuration(key)
	if v.Nanoseconds() != 0 {
		return v
	}
	return defaultValue
}

// GetStringSlice viper get string slice
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// GetTrackKey get track key
func GetTrackKey() string {
	v := viper.GetString("track")
	if v == "" {
		return "jt"
	}
	return v
}

// GetSessionKeys get the encrypt keys for session
func GetSessionKeys() []string {
	v := viper.GetStringSlice("session.keys")
	if len(v) == 0 {
		return []string{
			"cuttlefish",
		}
	}
	return v
}

// GetSessionCookie get the session cookie's name
func GetSessionCookie() string {
	v := viper.GetString("session.cookie.name")
	if v == "" {
		return "sess"
	}
	return v
}

// GetCookiePath get the cookie's path
func GetCookiePath() string {
	v := viper.GetString("session.cookie.path")
	if v == "" {
		return "/"
	}
	return v
}
