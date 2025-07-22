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

func fileFilter(f core.FileLike) bool {
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

func collectDirFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read directory")
	}

	var filesList []string
	for _, file := range files {
		if fileFilter(file) {
			filesList = append(filesList, filepath.Join(dir, file.Name()))
		}
	}

	return filesList, nil
}

func jsonOutput(meta core.FileMeta) {
	json.NewEncoder(os.Stdout).Encode(meta)
}

func output(meta core.FileMeta) {
	fmt.Printf("Name: %s\nSize: %d\nHash: %s\nModTime: %s\n===============\n", meta.Name, meta.Size, meta.Hash, meta.ModTime)
}

func Run(cfg core.Config) {
	if cfg.FilePath != "" {
		meta, err := getFileMeta(cfg.FilePath, cfg.HashType)
		if err != nil {
			log.Fatal(err)
		}
		if cfg.Json {
			jsonOutput(meta)
			return
		}
		output(meta)
		return
	} else if cfg.DirPath != "" {
		filesList, err := collectDirFiles(cfg.DirPath)
		if err != nil {
			log.Fatal(err)
		}
		var wg sync.WaitGroup
		for _, file := range filesList {
			wg.Add(1)
			go func(f string, h string) {
				defer wg.Done()
				meta, err := getFileMeta(f, h)
				if err != nil {
					log.Print("could not get file meta %w", err)
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
		log.Fatal("Please provide path to file check --help")
	}
}

func parseFlags() core.Config {
	filePath := flag.String("file", "", "Path to the file to analyze")
	dirPath := flag.String("dir", "", "Path to a directory containing files to analyze")
	jsonOutput := flag.Bool("json", false, "Output as JSON")
	hashType := flag.String("hash", "SHA256", "Choose hash type (MD5, SHA256)")

	flag.Parse()

	config := core.Config{FilePath: *filePath, DirPath: *dirPath, Json: *jsonOutput, HashType: *hashType}

	return config
}

func Execute() {

	cfg := parseFlags()
	Run(cfg)
}
