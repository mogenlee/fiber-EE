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
	PackageName string
	StructName  string
	RoutePath   string
}

func toTitle(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func main() {
	moduleName := flag.String("name", "", "模块名称 (必填)")
	routerOut := flag.String("router", "app/router/admin", "router 输出目录")
	serviceOut := flag.String("service", "app/service/admin", "service 输出目录")
	tplDir := flag.String("tpl", "cmd/gen_module", "模板目录")
	flag.Parse()

	if *moduleName == "" {
		fmt.Println("用法: make module name=<模块名>")
		fmt.Println("示例: make module name=article")
		os.Exit(1)
	}

	name := strings.ToLower(*moduleName)
	data := ModuleData{
		PackageName: name,
		StructName:  toTitle(name),
		RoutePath:   name,
	}

	// 生成 router
	routerDir := filepath.Join(*routerOut, name)
	os.MkdirAll(routerDir, 0755)
	generateFile(filepath.Join(*tplDir, "router.go.tpl"), filepath.Join(routerDir, "router.go"), data)

	// 生成 service
	serviceDir := filepath.Join(*serviceOut, name)
	os.MkdirAll(serviceDir, 0755)
	generateFile(filepath.Join(*tplDir, "service.go.tpl"), filepath.Join(serviceDir, "service.go"), data)

	// 自动注入
	injectRouter(name)
	injectService(name)

	fmt.Printf("\n✅ 模块 [%s] 生成完成！\n", name)
}

func generateFile(tplPath, outPath string, data ModuleData) {
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	tpl.Execute(f, data)
	fmt.Printf("  生成: %s\n", outPath)
}

func injectRouter(name string) {
	path := "app/router/build.go"
	content, _ := os.ReadFile(path)
	str := string(content)

	importLine := fmt.Sprintf("\t\"fiber-ee/app/router/admin/%s\"", name)
	routerLine := fmt.Sprintf("\t%s.NewRouter,", name)

	if strings.Contains(str, importLine) {
		return
	}

	// 插入 import
	str = strings.Replace(str,
		"\"fiber-ee/app/router/admin/test\"",
		"\"fiber-ee/app/router/admin/test\"\n"+importLine, 1)

	// 插入 router
	str = strings.Replace(str,
		"var adminRouters = []any{",
		"var adminRouters = []any{\n"+routerLine, 1)

	os.WriteFile(path, []byte(str), 0644)
	fmt.Printf("  注入: %s\n", path)
}

func injectService(name string) {
	path := "app/service/build.go"
	content, _ := os.ReadFile(path)
	str := string(content)

	importLine := fmt.Sprintf("\t\"fiber-ee/app/service/admin/%s\"", name)
	serviceLine := fmt.Sprintf("\t%s.New%sService,", name, toTitle(name))

	if strings.Contains(str, importLine) {
		return
	}

	// 插入 import
	str = strings.Replace(str,
		"\"fiber-ee/app/service/admin/test\"",
		"\"fiber-ee/app/service/admin/test\"\n"+importLine, 1)

	// 插入 service
	str = strings.Replace(str,
		"var adminServices = []any{",
		"var adminServices = []any{\n"+serviceLine, 1)

	os.WriteFile(path, []byte(str), 0644)
	fmt.Printf("  注入: %s\n", path)
}
