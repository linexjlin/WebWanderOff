package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name           string   `yaml:"name"`
	ListenAddr     string   `yaml:"listen_addr"`
	DefaultServer  string   `yaml:"default_server"`
	DefaultScheme  string   `yaml:"default_scheme"`
	CacheRoot      string   `yaml:"cache_root"`
	OfflineDomains []string `yaml:"offline_domains"`
}

func main() {
	//set log format show line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 获取当前目录下所有的 inpaint.yaml 配置文件名
	files, err := filepath.Glob("*.yaml")
	if err != nil {
		fmt.Printf("failed to read config files: %v", err)
		return
	}

	// 遍历每个文件，并读取配置内容
	for _, file := range files {
		// 读取 YAML 文件内容
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("failed to read file %s: %v\n", file, err)
			continue
		}

		// 解析 YAML 文件数据到 Config 结构体实例
		var conf Config
		if err := yaml.Unmarshal(content, &conf); err != nil {
			fmt.Printf("failed to unmarshal file %s: %v\n", file, err)
			continue
		}

		// 打印该配置文件的内容
		fmt.Printf("Configuration of %s:\n", file)
		fmt.Printf("  Name: %s\n", conf.Name)
		fmt.Printf("  ListenAddr: %s\n", conf.ListenAddr)
		fmt.Printf("  DefaultServer: %s\n", conf.DefaultServer)
		fmt.Printf("  DefaultScheme: %s\n", conf.DefaultScheme)
		fmt.Printf("  CacheRoot: %s\n", conf.CacheRoot)
		fmt.Printf("  OfflineDomains: %v\n", conf.OfflineDomains)

		cs := CacheSystem{ListenAddr: conf.ListenAddr, DefaultServer: conf.DefaultServer, DefaultScheme: conf.DefaultScheme, CacheRoot: conf.CacheRoot, OfflineDomains: conf.OfflineDomains}
		go cs.Listen()
	}

	select {}
}
