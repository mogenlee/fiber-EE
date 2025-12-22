package bootstrap

import (
	"fmt"
	"os"
	"strings"

	"fiber-ee/config"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func NewConfig() *config.Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config") // 在 config 目录中查找配置文件

	// 检查环境变量以覆盖配置路径
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		v.AddConfigPath(configPath)
	}

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("致命错误：无法读取配置文件: %w", err))
	}

	// 绑定环境变量（支持 FIBER_ 前缀覆盖配置）
	v.SetEnvPrefix("FIBER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	cfg := &config.Config{}
	if err := v.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("无法解析配置到结构体: %w", err))
	}

	// 开启热更新
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("[CONFIG] 配置文件已热更新: %s\n", e.Name)
		if err := v.Unmarshal(cfg); err != nil {
			fmt.Printf("[CONFIG] 热更新解析配置失败: %v\n", err)
		}
	})

	fmt.Printf("[CONFIG] 配置加载成功 env=%s\n", cfg.App.Env)
	return cfg
}
