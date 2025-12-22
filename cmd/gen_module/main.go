package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type ModuleData struct {
	PackageName string // 包名，如 article
	StructName  string // 结构体名，如 Article（首字母大写）
	RoutePath   string // 路由路径，如 article
}

// 首字母大写
func toTitle(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func main() {
	// 命令行参数
	moduleName := flag.String("name", "", "模块名称 (必填)")
	handlerOut := flag.String("handler", "internal/handler/admin", "handler 输出目录")
	serviceOut := flag.String("service", "internal/service/admin", "service 输出目录")
	tplDir := flag.String("tpl", "cmd/gen_module", "模板目录")
	flag.Parse()

	if *moduleName == "" {
		fmt.Println("用法: go run cmd/gen_module/main.go -name=<模块名> [选项]")
		fmt.Println("\n选项:")
		flag.PrintDefaults()
		fmt.Println("\n示例:")
		fmt.Println("  go run cmd/gen_module/main.go -name=article")
		fmt.Println("  go run cmd/gen_module/main.go -name=article -handler=internal/handler/app")
		os.Exit(1)
	}

	name := strings.ToLower(*moduleName)
	data := ModuleData{
		PackageName: name,
		StructName:  toTitle(name),
		RoutePath:   name,
	}

	// 生成 handler
	handlerDir := filepath.Join(*handlerOut, name)
	if err := os.MkdirAll(handlerDir, 0755); err != nil {
		panic(err)
	}
	generateFile(filepath.Join(*tplDir, "router.go.tpl"), filepath.Join(handlerDir, "handler.go"), data)

	// 生成 service
	serviceDir := filepath.Join(*serviceOut, name)
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		panic(err)
	}
	generateFile(filepath.Join(*tplDir, "service.go.tpl"), filepath.Join(serviceDir, "service.go"), data)

	fmt.Printf("\n模块 [%s] 生成完成！\n", name)
	fmt.Println("\n请手动添加到注册列表：")
	fmt.Printf("  internal/bootstrap/buildRouter.go:\n")
	fmt.Printf("    import %s \"fiber-ee/internal/handler/admin/%s\"\n", name, name)
	fmt.Printf("    var adminRouters = []any{ ..., %s.NewRouter }\n\n", name)
	fmt.Printf("  internal/bootstrap/buildService.go:\n")
	fmt.Printf("    import %s \"fiber-ee/internal/service/admin/%s\"\n", name, name)
	fmt.Printf("    var adminServices = []any{ ..., %s.NewService }\n", name)
}

func generateFile(tplPath, outPath string, data ModuleData) {
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		panic(fmt.Errorf("解析模板失败 %s: %w", tplPath, err))
	}

	f, err := os.Create(outPath)
	if err != nil {
		panic(fmt.Errorf("创建文件失败 %s: %w", outPath, err))
	}
	defer f.Close()

	if err := tpl.Execute(f, data); err != nil {
		panic(fmt.Errorf("生成文件失败 %s: %w", outPath, err))
	}

	fmt.Printf("  生成: %s\n", outPath)
}
