package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/0xM3H51N/MetaGo/internal"
	"github.com/0xM3H51N/MetaGo/internal/core"
)

const Version = "MetaGo v0.9.0"

func fileFilter(f any) bool {
	switch v := f.(type) {
	case os.DirEntry:
		info, err := v.Info()
		if err != nil {
			log.Printf("Failed to get info for %s: %v", v.Name(), err)
			return false
		}
		return info.Mode().IsRegular()
	case os.FileInfo:
		return v.Mode().IsRegular()
	default:
		return false
	}

}

func getFileMeta(file string, hashtype string) (core.FileMeta, error) {
	info, err := os.Stat(file)
	if err != nil {
		return core.FileMeta{}, fmt.Errorf("could not get file info %w", err)
	}

	if !fileFilter(info) {
		return core.FileMeta{}, fmt.Errorf("not a regular file")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return core.FileMeta{}, fmt.Errorf("could not read file %w", err)
	}

	hash, err := internal.GetFileHash(data, hashtype)
	if err != nil {
		hash = "Could not hash the file for some reason!"
	}

	fileMeta := core.FileMeta{Name: info.Name(), Size: info.Size(), Hash: hash, ModTime: info.ModTime().String()}

	return fileMeta, nil
}

func collectDirFiles(dir string) ([]string, []string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read directory %q: %w", dir, err)
	}

	var filesList []string
	var dirList []string

	for _, f := range files {
		if f.IsDir() {
			dirList = append(dirList, filepath.Join(dir, f.Name()))
			continue
		}
		if fileFilter(f) {
			filesList = append(filesList, filepath.Join(dir, f.Name()))
		}

	}

	return filesList, dirList, nil
}

func jsonOutput(meta core.FileMeta) {
	if err := json.NewEncoder(os.Stdout).Encode(meta); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func output(meta core.FileMeta) {
	fmt.Printf("Name: %s\nSize: %d\nHash: %s\nModTime: %s\n===============\n", meta.Name, meta.Size, meta.Hash, meta.ModTime)
}

func Run(cfg core.Config) error {
	if cfg.FilePath != "" && cfg.DirPath != "" {
		return fmt.Errorf("please provide either -f or -d, not both")
	}
	if cfg.FilePath != "" {
		meta, err := getFileMeta(cfg.FilePath, cfg.HashType)
		if err != nil {
			return err
		}
		if cfg.Json {
			jsonOutput(meta)
			return nil
		}
		output(meta)
		return nil
	} else if cfg.DirPath != "" {
		filesList, dirList, err := collectDirFiles(cfg.DirPath)
		if err != nil {
			return err
		}
		var wg sync.WaitGroup
		if cfg.Recursive {
			for _, dir := range dirList {
				childCfg := cfg
				childCfg.DirPath = dir
				childCfg.FilePath = ""
				wg.Add(1)
				go func(c core.Config) {
					defer wg.Done()
					if err := Run(c); err != nil {
						log.Printf("error scanning %s: %v", c.DirPath, err)
					}
				}(childCfg)
			}
		}

		for _, file := range filesList {
			wg.Add(1)
			go func(f string, h string) {
				defer wg.Done()
				meta, err := getFileMeta(f, h)
				if err != nil {
					log.Printf("could not get file meta: %v", err)
				}
				if cfg.Json {
					jsonOutput(meta)
				} else {
					output(meta)
				}
			}(file, cfg.HashType)
		}
		wg.Wait()
	} else {
		return fmt.Errorf("no input provided : use --help for help")
	}
	return nil
}

func parseFlags() core.Config {
	filePath := flag.String("f", "", "Path to the file to analyze")
	dirPath := flag.String("d", "", "Path to a directory containing files to analyze")
	jsonOutput := flag.Bool("json", false, "Output as JSON")
	hashType := flag.String("h", "SHA256", "Choose hash type (MD5, SHA256)")
	recursive := flag.Bool("r", false, "scanning entire folders")
	version := flag.Bool("v", false, "Print version and exit")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	config := core.Config{FilePath: *filePath, DirPath: *dirPath, Json: *jsonOutput, HashType: *hashType, Recursive: *recursive}

	return config
}

func Execute() {

	cfg := parseFlags()
	err := Run(cfg)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
