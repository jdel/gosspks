package cmd

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/jdel/gosspks/cache"
	"github.com/jdel/gosspks/cfg"
	"github.com/jdel/gosspks/handler/v1"
	"github.com/jdel/gosspks/rtr"
	"github.com/jdel/gosspks/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

var cfgFile, goSSPKSHome string
var cfgDebugPackage bool

var loggerCmd = log.WithFields(log.Fields{
	"module": util.WhereAmI(),
})

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gosspks",
	Short: "Serving your Synology Packages.",
	Long:  `Serving your Synology Packages.`,
	Run: func(cmd *cobra.Command, args []string) {
		router := rtr.NewRouter()
		allowedOrigins := handlers.AllowedOrigins([]string{"*"})
		allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "OPTIONS"})
		allowedHeaders := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "X-Requested-With"})
		corsHandler := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)

		if t := cfg.GetPackagesCacheRefreshRate(); t != 0 {
			ticker := time.NewTicker(t)
			go func() {
				for _ = range ticker.C {
					c := cache.Cache
					for k, v := range c.Items() {
						if strings.Contains(k, "-request") {
							bits := strings.Split(k, "-")
							_, e, _ := c.GetWithExpiration(bits[0] + "-" + bits[1])
							if time.Now().Add(t).UnixNano() > e.UnixNano() {
								loggerCmd.Debugf("Packages cache %s-%s expires before next refresh in %s, refreshing", bits[0], bits[1], time.Until(e))
								v1.GetPackagesFromCacheOrFileSystem(v.Object.(*http.Request), bits[1], true)
							}
						}
					}
				}
			}()
		}

		if t := cfg.GetModelsCacheRefreshRate(); t != 0 {
			ticker := time.NewTicker(t)
			go func() {
				for _ = range ticker.C {
					i, e, _ := cache.Cache.GetWithExpiration("synoModels")
					if i != nil && time.Now().Add(t).UnixNano() > e.UnixNano() {
						loggerCmd.Debugf("Models cache expires before next refresh in %s, refreshing", time.Until(e))
						v1.GetModelsFromCacheOrWeb(true)
					}
				}
			}()
		}

		loggerCmd.Info("Using config:")
		settings, _ := yaml.Marshal(viper.AllSettings())
		for _, setting := range strings.Split(string(settings), "\n") {
			loggerCmd.Info(setting)
		}
		loggerCmd.Info("Listening on: ", cfg.GetPort())
		loggerCmd.Fatal(http.ListenAndServe(":"+cfg.GetPort(), corsHandler(router)))
	},
}

// Execute runs the main command that
// serves synology packages
func Execute() {
	RootCmd.Execute()
}

func init() {
	// CMD line args > ENV VARS > Config file
	cobra.OnInitialize(func() { cfg.InitConfig(cfgFile, goSSPKSHome) })
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", "", "config file (default is $HOME/.gosspks/config.yml)")
	RootCmd.PersistentFlags().StringVarP(&goSSPKSHome, "home", "H", "", "gosspks home (default is $HOME/.gosspks/")
	// Optional flags
	RootCmd.PersistentFlags().IntP("port", "p", 8080, "port to listen to")
	RootCmd.PersistentFlags().BoolP("debug-package", "d", false, "generates a debug package visible in Synology Package Center")
	RootCmd.PersistentFlags().String("packages", "packages", "packages directory")
	RootCmd.PersistentFlags().String("cache", "cache", "cache directory (gosspks extracts INFO and images here)")
	RootCmd.PersistentFlags().Int("packages-cache-count", 15, "im-memory cache size (0 to read packages from disk every time)")
	RootCmd.PersistentFlags().String("packages-cache-duration", "5m", "packages in-memory cache TTL")
	RootCmd.PersistentFlags().String("models-cache-duration", "7d", "models in-memory cache TTL")
	RootCmd.PersistentFlags().String("packages-cache-refresh", "1m", "packages in-memory cache automatic refresh rate")
	RootCmd.PersistentFlags().String("models-cache-refresh", "", "models in-memory cache automatic refresh rate")
	RootCmd.PersistentFlags().String("models", "models.yml", "models file")
	RootCmd.PersistentFlags().String("hostname", "", "hostname to use when generating urls")
	RootCmd.PersistentFlags().String("scheme", "http", "scheme to use when generating urls")
	RootCmd.PersistentFlags().StringP("log-level", "l", "Error", "log level [Error,Warn,Info,Debug]")
	RootCmd.PersistentFlags().String("static", "static", "prefix to serve static images")
	RootCmd.PersistentFlags().String("download", "download", "prefix to serve packages")
	RootCmd.PersistentFlags().Bool("md5", false, "enable md5 calculation")
	// Bind flags to config
	viper.BindPFlag("gosspks.port", RootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("gosspks.debug-package", RootCmd.PersistentFlags().Lookup("debug-package"))
	viper.BindPFlag("gosspks.filesystem.packages", RootCmd.PersistentFlags().Lookup("packages"))
	viper.BindPFlag("gosspks.filesystem.cache", RootCmd.PersistentFlags().Lookup("cache"))
	viper.BindPFlag("gosspks.filesystem.models", RootCmd.PersistentFlags().Lookup("models"))
	viper.BindPFlag("gosspks.cache.packages.count", RootCmd.PersistentFlags().Lookup("packages-cache-count"))
	viper.BindPFlag("gosspks.cache.packages.duration", RootCmd.PersistentFlags().Lookup("packages-cache-duration"))
	viper.BindPFlag("gosspks.cache.models.duration", RootCmd.PersistentFlags().Lookup("models-cache-duration"))
	viper.BindPFlag("gosspks.cache.packages.refresh", RootCmd.PersistentFlags().Lookup("packages-cache-refresh"))
	viper.BindPFlag("gosspks.cache.models.refresh", RootCmd.PersistentFlags().Lookup("models-cache-refresh"))
	viper.BindPFlag("gosspks.hostname", RootCmd.PersistentFlags().Lookup("hostname"))
	viper.BindPFlag("gosspks.scheme", RootCmd.PersistentFlags().Lookup("scheme"))
	viper.BindPFlag("gosspks.log-level", RootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("gosspks.router.static", RootCmd.PersistentFlags().Lookup("static"))
	viper.BindPFlag("gosspks.router.download", RootCmd.PersistentFlags().Lookup("download"))
	viper.BindPFlag("gosspks.md5", RootCmd.PersistentFlags().Lookup("md5"))
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
}
