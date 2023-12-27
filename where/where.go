package where

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/key"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

const EnvConfigPath = "MANGAL_CONFIG_PATH"

// mkdir creates a directory and all parent directories if they don't exist
// will return the path of the directory
func mkdir(path string) string {
	if filesystem.Api().MkdirAll(path, os.ModePerm) != nil {
		log.Fatalf("Error: could not create directory %s", path)
	}
	return path
}

// Config path
// Will create the directory if it doesn't exist
func Config() string {
	var path string

	if customDir, present := os.LookupEnv(EnvConfigPath); present {
		path = customDir
	} else {
		path = filepath.Join(lo.Must(os.UserConfigDir()), constant.Mangal)
	}

	return mkdir(path)
}

// Sources path
// Will create the directory if it doesn't exist
func Sources() string {
	return mkdir(filepath.Join(Config(), "sources"))
}

func AnilistBinds() string {
	return filepath.Join(Config(), "anilist.json")
}

// Logs path
// Will create the directory if it doesn't exist
func Logs() string {
	return mkdir(filepath.Join(Config(), "logs"))
}

// Queries path
// Will create the directory if it doesn't exist
func Queries() string {
	return filepath.Join(Cache(), "queries.json")
}

// History path to the file
// Will create the directory if it doesn't exist
func History() string {
	return filepath.Join(Config(), "history.json")
}

// Downloads path
// Will create the directory if it doesn't exist
func Downloads() string {
	path, err := filepath.Abs(viper.GetString(key.DownloaderPath))

	if err != nil {
		path, err = os.Getwd()
		if err != nil {
			path = "."
		}
	}

	return mkdir(path)
}

// Cache path
// Will create the directory if it doesn't exist
func Cache() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = filepath.Join(".", "cache")
	}

	cacheDir = filepath.Join(cacheDir, constant.Mangal)
	return mkdir(cacheDir)
}

// Temp path
// Will create the directory if it doesn't exist
func Temp() string {
	tempDir := filepath.Join(os.TempDir(), constant.Mangal)
	return mkdir(tempDir)
}
