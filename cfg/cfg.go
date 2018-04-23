package cfg

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jdel/go-syno"
	"github.com/jdel/gosspks/cache"
	"github.com/jdel/gosspks/stc"
	"github.com/jdel/gosspks/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Version is the current application version.
// This variable is populated when building the binary with:
// -ldflags "-X github.com/jdel/gosspks/cfg.Version=${VERSION}"
var Version string

// SynoOptions provides options for go-syno library
var SynoOptions *syno.Options
var goSSPKSHome string

var loggerCfg = log.WithFields(log.Fields{
	"module": util.WhereAmI(),
})

// InitConfig loads the config file according
// to cfgFile and homeDir flags from cmd/root.go
func InitConfig(cfgFile string, homeDir string) {
	// Instantiate goSSPKSHome ASAP
	goSSPKSHome = getOrCreateHome(homeDir)

	if cfgFile != "" {
		// Use config file from the flag if present
		viper.SetConfigFile(cfgFile)
	} else {
		// Otherwise use goSSPKSHome from flag if present
		if homeDir != "" {
			viper.AddConfigPath(homeDir)
			viper.SetConfigName("gosspks")
		} else {
			// Search config in home directory with name gosspks
			viper.AddConfigPath(goSSPKSHome)
			viper.AddConfigPath(".")
			viper.SetConfigName("gosspks")
		}
	}

	// Read the config
	if err := viper.ReadInConfig(); err != nil {
		e, ok := err.(viper.ConfigParseError)
		if ok {
			loggerCfg.Error(e)
		}

		configFileToWrite := cfgFile
		if cfgFile == "" {
			if homeDir != "" {
				configFileToWrite = filepath.Join(homeDir, "gosspks.yml")
			} else {
				configFileToWrite = filepath.Join(goSSPKSHome, "gosspks.yml")
			}
		}

		loggerCfg.Warn("No config file used, writing ", configFileToWrite, " with default values")
		settings, _ := yaml.Marshal(viper.AllSettings())
		if err := ioutil.WriteFile(configFileToWrite, settings, 0644); err != nil {
			loggerCfg.Error(err)
		}
	}

	logLevel := parseLogLevel(GetLogLevel())
	log.SetLevel(logLevel)
	loggerCfg.Info("Using config file: ", viper.ConfigFileUsed())

	// Create directories
	cacheDir := GetCacheDir()
	if err := util.CreateDir(cacheDir); err != nil {
		loggerCfg.Error("Couldn't create cache directory ", cacheDir)
	}
	packagesDir := GetPackagesDir()
	if err := util.CreateDir(packagesDir); err != nil {
		loggerCfg.Error("Couldn't create packages directory ", packagesDir)
	}

	loggerCfg.Info("Home is: ", goSSPKSHome)
	loggerCfg.Info("Cache directory is: ", cacheDir)
	loggerCfg.Info("Packages directory is: ", packagesDir)
	// Write the debug.png and png file, always overwrite
	data, _ := base64.StdEncoding.DecodeString(stc.DebugPackagePng)
	ioutil.WriteFile(filepath.Join(GetCacheDir(), "debug.png"), data, 0664)
	data, _ = base64.StdEncoding.DecodeString(stc.DefaultPackagePng)
	ioutil.WriteFile(filepath.Join(GetCacheDir(), "default.png"), data, 0664)

	// Make sure it has slashes
	viper.Set("gosspks.router.static", addSlashPrefix(GetStaticPrefix()))
	viper.Set("gosspks.router.download", addSlashPrefix(GetDownloadPrefix()))

	// Set options for go-syno library
	SynoOptions = syno.GetOptions()
	SynoOptions.PackagesDir = packagesDir
	SynoOptions.CacheDir = cacheDir
	SynoOptions.MD5 = GetMD5()

	loggerCfg.Info("Found ", len(GetGPGKeys()), " GPG key(s)")

	// Watch for changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		loggerCfg.Info("Config file changed: ", e.Name)
		logLevel := parseLogLevel(GetLogLevel())
		log.SetLevel(logLevel)
	})
}

func parseLogLevel(level string) log.Level {
	var logLevel log.Level
	var err error
	loggerCfg.WithField("log-level", level).Info("Parsing log level")
	if logLevel, err = log.ParseLevel(level); err != nil {
		logLevel = log.ErrorLevel
		loggerCfg.WithField("log-level", level).Error("Cannot parse log level, setting to Error")
	}
	return logLevel
}

// GetDebugPackage returns true if the
// --debug-package flag is present
func GetDebugPackage() bool {
	return viper.GetBool("gosspks.debug-package")
}

// GetMD5 returns true if the
// --md5 flag is present
func GetMD5() bool {
	return viper.GetBool("gosspks.md5")
}

// GetDownloadPrefix returns the prefix
// for download URLs (packages)
func GetDownloadPrefix() string {
	return viper.GetString("gosspks.router.download")
}

// GetStaticPrefix returns the prefix
// for static URLs (images)
func GetStaticPrefix() string {
	return viper.GetString("gosspks.router.static")
}

// GetHostName returns the hostname
func GetHostName() string {
	return viper.GetString("gosspks.hostname")
}

// GetScheme returns the scheme
func GetScheme() string {
	return viper.GetString("gosspks.scheme")
}

// GetCacheDir returns the cache directory
func GetCacheDir() string {
	cacheDir := viper.GetString("gosspks.filesystem.cache")
	if match, _ := regexp.MatchString("^/", cacheDir); !match {
		return filepath.Join(goSSPKSHome, cacheDir)
	}
	return cacheDir
}

// GetPackagesDir returns the packages directory
func GetPackagesDir() string {
	packagesDir := viper.GetString("gosspks.filesystem.packages")
	if match, _ := regexp.MatchString("^/", packagesDir); !match {
		return filepath.Join(goSSPKSHome, packagesDir)
	}
	return packagesDir
}

// GetModelsFile returns the prefix
// for download URLs (packages)
func GetModelsFile() string {
	modelsFile := viper.GetString("gosspks.filesystem.models")
	if match, _ := regexp.MatchString("^/", modelsFile); !match {
		return filepath.Join(goSSPKSHome, modelsFile)
	}
	return modelsFile
}

// GetPackagesCacheCount returns the number of
// packages collections to cache (one for each language)
func GetPackagesCacheCount() int {
	return viper.GetInt("gosspks.cache.packages.count")
}

// GetPackagesCacheDuration returns the duration
// of the packages cache
func GetPackagesCacheDuration() time.Duration {
	var err error
	packagesDuration := viper.GetString("gosspks.cache.packages.duration")
	cacheDuration := cache.DefaultCacheDuration
	if cacheDuration, err = time.ParseDuration(packagesDuration); err != nil {
		loggerCfg.Warnf("Cannot parse string %s to duration. Defaulting to %s.", packagesDuration, cache.DefaultCacheDuration)
		return 0
	}
	return cacheDuration
}

// GetPackagesCacheRefreshRate returns the refresh rate
// of the packages cache
func GetPackagesCacheRefreshRate() time.Duration {
	packagesRefreshRate := viper.GetString("gosspks.cache.packages.refresh")
	cacheDuration, err := time.ParseDuration(packagesRefreshRate)
	if err != nil {
		loggerCfg.Warnf("Cannot parse string %s to duration.", packagesRefreshRate)
		return 0
	}
	return cacheDuration
}

// GetModelsCacheDuration returns the duration
// of the models cache
func GetModelsCacheDuration() time.Duration {
	var err error
	modelsDuration := viper.GetString("gosspks.cache.models.duration")
	cacheDuration := cache.DefaultCacheDuration
	if cacheDuration, err = time.ParseDuration(modelsDuration); err != nil {
		loggerCfg.Warnf("Cannot parse string %s to duration. Defaulting to %s.", modelsDuration, cache.DefaultCacheDuration)
		return 0
	}
	return cacheDuration
}

// GetModelsCacheRefreshRate returns the refresh rate
// of the packages cache
func GetModelsCacheRefreshRate() time.Duration {
	modelsRefreshRate := viper.GetString("gosspks.cache.models.refresh")
	cacheDuration, err := time.ParseDuration(modelsRefreshRate)
	if err != nil {
		loggerCfg.Warnf("Cannot parse string %s to duration.", modelsRefreshRate)
		return 0
	}
	return cacheDuration
}

// GetPort returns the port
// default value is "8080"
func GetPort() string {
	return viper.GetString("gosspks.port")
}

// GetLogLevel returns the log level.
// default value is "Error"
func GetLogLevel() string {
	return viper.GetString("gosspks.log-level")
}

// GetGPGKeys returns a []string of GPG keys
// Keys have to be specified as an array in the config file
func GetGPGKeys() []string {
	return viper.GetStringSlice("gosspks.gpg")
}

// Prefixes a string with / if the string
// doesn't start with /
func addSlashPrefix(value string) string {
	if match, _ := regexp.MatchString("^/", value); !match {
		value = fmt.Sprintf("/%s", value)
	}
	return value
}

// getOrCreateHome returns .gosspks subdir from
// user's home directory and creates it if required
func getOrCreateHome(goSSPKSHome string) string {
	var home string
	if goSSPKSHome != "" {
		home = goSSPKSHome
	} else {
		usr, err := user.Current()
		if err != nil {
			loggerCfg.Fatal(err)
		}
		loggerCfg.Info("Current user: ", usr.Username)
		home = filepath.Join(usr.HomeDir, "/gosspks/")
	}

	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := util.CreateDir(home); err != nil {
			loggerCfg.Error("Couldn't create home directory ", home)
		}
	}
	return home
}

func init() {
	// Sets logrus options
	formatter := &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "06/01/02 15:04:05.000",
	}
	log.SetFormatter(formatter)
	log.SetOutput(os.Stderr)
}
