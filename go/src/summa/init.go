package summa

import (
	"log"
	"path/filepath"
)

// SetAuthProvider sets the function to call to authenticate users
func SetAuthProvider(ap AuthProvider) {
	config.SetAuthProvider(ap)
}

// Init loads the Summa configuration file and performs some base
// initialization tasks on the config settings
func Init(configFile string) error {
	configFilePath, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}

	if err = configLoad(configFilePath); err != nil {
		return err
	}

	configFileDir := filepath.Dir(configFilePath)

	// Resolve all directory path config settings
	// into absolute paths, making sure that the
	// directory exists
	for k, v := range config.DirPaths {
		config.DirPaths[k], err = ResolvePath(v, configFileDir)
		if err != nil {
			log.Fatalf("Could not resolve %s: %s", k, err)
		}

		if !IsDir(config.DirPaths[k]) {
			log.Fatalf("%s is not a directory", k)
		}
	}

	// Resolve all file path config settings
	// into absolute paths, making sure their
	// parent directory exists
	for k, v := range config.FilePaths {
		config.FilePaths[k], err = ResolvePath(v, configFileDir)
		if err != nil {
			log.Fatalf("Could not resolve %s: %s", k, err)
		}

		if !IsDir(filepath.Dir(config.FilePaths[k])) {
			log.Fatalf("Directory for %s does not exist", k)
		}
	}

	err = startLogging(config.LogFile())
	if err != nil {
		log.Fatalf("Could not setup log file: %s", err)
	}

	infoLog.Printf("summa.Init()")
	infoLog.Printf("Loaded configuration from %s", configFilePath)

	return nil
}
