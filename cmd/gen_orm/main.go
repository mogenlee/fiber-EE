package main

import (
	"fmt"
	"os"

	"fiber-ee/config"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	yamlgen "github.com/we7coreteam/gorm-gen-yaml"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	// 读取配置
	cfg := loadConfig()

	// 设置默认值
	outPath := cfg.Database.GenOutPath
	if outPath == "" {
		outPath = "./internal/model/query"
	}
	modelPkgPath := cfg.Database.GenModelPkgPath
	if modelPkgPath == "" {
		modelPkgPath = "./internal/model/entity"
	}
	genYamlPath := cfg.Database.GenYamlPath
	if genYamlPath == "" {
		genYamlPath = "./cmd/gen_orm/gen.yaml"
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           outPath,
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		ModelPkgPath:      modelPkgPath,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	// 连接数据库
	db := connectDB(cfg)
	g.UseDB(db)

	// 添加自定义类型的 import
	g.WithImportPkgPath("fiber-ee/internal/pkg/types")

	// 自定义字段的数据类型
	dataMap := map[string]func(detailType gorm.ColumnType) (dataType string){
		"timestamp": func(detailType gorm.ColumnType) (dataType string) { return "LocalTime" },
	}
	g.WithDataTypeMap(dataMap)

	// GORM 标签配置
	createTimeGormTag := gen.FieldGORMTag("created_at", func(tag field.GormTag) field.GormTag {
		tag.Append("autoCreateTime")
		return tag
	})
	updateTimeGormTag := gen.FieldGORMTag("updated_at", func(tag field.GormTag) field.GormTag {
		tag.Append("autoUpdateTime")
		return tag
	})

	// 时间字段使用自定义 Timestamp 类型
	createdAtType := gen.FieldType("created_at", "types.Timestamp")
	updatedAtType := gen.FieldType("updated_at", "types.Timestamp")
	deletedAtType := gen.FieldType("deleted_at", "soft_delete.DeletedAt")

	modelOpts := []gen.ModelOpt{
		createdAtType, updatedAtType, deletedAtType,
		createTimeGormTag, updateTimeGormTag,
	}

	// 生成代码
	yamlgen.NewYamlGenerator(genYamlPath).UseGormGenerator(g).Generate(modelOpts...)
	g.Execute()

	fmt.Printf("代码生成完成\n  Query: %s\n  Entity: %s\n", outPath, modelPkgPath)
}

// loadConfig 加载配置文件
func loadConfig() *config.Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config")

	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		v.AddConfigPath(configPath)
	}

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("无法读取配置文件: %w", err))
	}

	cfg := &config.Config{}
	if err := v.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("无法解析配置: %w", err))
	}

	return cfg
}

// connectDB 根据配置连接数据库
func connectDB(cfg *config.Config) *gorm.DB {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
			cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.DBName,
			cfg.Database.Port,
			cfg.Database.SSLMode,
		)
		dialector = postgres.Open(dsn)
	default:
		dialector = sqlite.Open(cfg.Database.Source)
	}

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 表前缀
	if cfg.Database.TablePrefix != "" {
		gormConfig.NamingStrategy = schema.NamingStrategy{
			TablePrefix: cfg.Database.TablePrefix,
		}
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		panic(fmt.Errorf("数据库连接失败: %w", err))
	}

	fmt.Printf("数据库连接成功 [%s] 表前缀: %s\n", cfg.Database.Driver, cfg.Database.TablePrefix)
	return db
}
