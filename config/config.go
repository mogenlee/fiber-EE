package config

type Config struct {
	App      App      `mapstructure:"app"`
	Log      Log      `mapstructure:"log"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
	JWT      JWT      `mapstructure:"jwt"`
	Casbin   Casbin   `mapstructure:"casbin"`
}

type App struct {
	Name  string `mapstructure:"name"`
	Env   string `mapstructure:"env"`
	Port  string `mapstructure:"port"`
	Debug bool   `mapstructure:"debug"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type Database struct {
	Driver          string `mapstructure:"driver"`
	Source          string `mapstructure:"source"`
	TablePrefix     string `mapstructure:"tablePrefix"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	// Gen 代码生成配置
	GenOutPath      string `mapstructure:"gen_out_path"`
	GenModelPkgPath string `mapstructure:"gen_model_pkg_path"`
	GenYamlPath     string `mapstructure:"gen_yaml_path"`
}

type Redis struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWT struct {
	Secret        string `mapstructure:"secret"`
	Expire        int64  `mapstructure:"expire"`
	RefreshExpire int64  `mapstructure:"refresh_expire"`
}

type Casbin struct {
	ModelPath string `mapstructure:"model_path"`
}
