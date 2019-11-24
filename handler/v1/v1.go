package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
	"jdel.org/go-syno"
	"jdel.org/gosspks/cache"
	"jdel.org/gosspks/cfg"
	"jdel.org/gosspks/util"
	log "github.com/sirupsen/logrus"
)

var filename string

var loggerV1 = log.WithFields(log.Fields{
	"module": util.WhereAmI(),
})

func getScheme(r *http.Request) string {
	scheme := "http"
	if r.URL.IsAbs() {
		scheme = r.URL.Scheme
	}
	return scheme
}

func buildFromFiles(r *http.Request, files []os.FileInfo) (syno.Packages, error) {
	var synoPkgs = make(syno.Packages, 0)
	cfg.SynoOptions.Language = r.URL.Query().Get("language")

	loggerV1.Debug("Processing ", len(files), " files")

	// Deal with biggest files first
	sort.Sort(util.FileInfoBySizeDesc(files))

	// Waitgroup to sync routines
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(files))

	for _, file := range files {
		loggerV1.WithFields(log.Fields{
			"filename": file.Name(),
			"size":     humanize.Bytes(uint64(file.Size())),
		}).Debug("Started processing")
		go func(filename string) {
			defer func(start time.Time) {
				loggerV1.WithFields(log.Fields{
					"filename": filename,
					"duration": util.Duration(time.Since(start)).Round(time.Millisecond),
				}).Debug("Finished processing")
			}(time.Now())

			defer waitGroup.Done()

			if filepath.Ext(filename) == ".spk" {
				newSynoPkg, err := syno.NewPackage(filename)
				if err != nil {
					loggerV1.WithFields(log.Fields{
						"filename": filename,
					}).Debug(err)
					return
				}
				hostname := cfg.GetHostName()
				if hostname == "" {
					hostname = r.Host
				}

				scheme := cfg.GetScheme()
				if scheme == "" {
					scheme = getScheme(r)
				}

				newSynoPkg.Link = fmt.Sprintf(
					"%s://%s/download/%s",
					scheme, hostname, filename)

				if len(newSynoPkg.Thumbnail) == 0 {
					newSynoPkg.Thumbnail = append(newSynoPkg.Thumbnail, "default.png")
				}
				if len(newSynoPkg.ThumbnailRetina) == 0 {
					newSynoPkg.ThumbnailRetina = append(newSynoPkg.ThumbnailRetina, "default.png")
				}
				// add url prefix to images
				for thumbIndex, thumb := range newSynoPkg.Thumbnail {
					if thumb == "default.png" {
						newSynoPkg.Thumbnail[thumbIndex] = fmt.Sprintf(
							"%s://%s%s/%s",
							scheme, hostname, cfg.GetStaticPrefix(), thumb)
					} else {
						newSynoPkg.Thumbnail[thumbIndex] = fmt.Sprintf(
							"%s://%s%s/%s/%s",
							scheme, hostname, cfg.GetStaticPrefix(), filename, thumb)
					}
				}
				for thumbRetinaIndex, thumbRetina := range newSynoPkg.ThumbnailRetina {
					if thumbRetina == "default.png" {
						newSynoPkg.ThumbnailRetina[thumbRetinaIndex] = fmt.Sprintf(
							"%s://%s%s/%s",
							scheme, hostname, cfg.GetStaticPrefix(), thumbRetina)
					} else {
						newSynoPkg.ThumbnailRetina[thumbRetinaIndex] = fmt.Sprintf(
							"%s://%s%s/%s/%s",
							scheme, hostname, cfg.GetStaticPrefix(), filename, thumbRetina)
					}
				}
				for screenIndex, screen := range newSynoPkg.Snapshot {
					newSynoPkg.Snapshot[screenIndex] = fmt.Sprintf(
						"%s://%s%s/%s/%s",
						scheme, hostname, cfg.GetStaticPrefix(), filename, screen)
				}
				synoPkgs = append(synoPkgs, newSynoPkg)
			} else {
				loggerV1.Debug("Ignoring non-spk file ", filename)
			}
		}(file.Name())
	}
	waitGroup.Wait()
	return synoPkgs, nil
}

// GetPackagesFromCacheOrFileSystem gets packages from cache
// Or gets them from filesystem and populate cache
func GetPackagesFromCacheOrFileSystem(r *http.Request, paramLanguage string, forceRefresh bool) (syno.Packages, error) {
	var synoPkgs syno.Packages
	var files []os.FileInfo
	var err error
	if cachedSynoPkgs, cacheExpiration, found := cache.Cache.GetWithExpiration(fmt.Sprintf("synoPkgs-%s", paramLanguage)); found && !forceRefresh {
		synoPkgs = cachedSynoPkgs.(syno.Packages)
		loggerV1.Debugf("Pulled syno packages (%s) from cache, valid for %s", paramLanguage, time.Until(cacheExpiration))
	} else {
		// Read all files
		if files, err = ioutil.ReadDir(cfg.GetPackagesDir()); err != nil {
			loggerV1.Error("Error walking packages directory. ", err)
			return synoPkgs, err
		}
		// Pass them to the builder

		synoPkgs, err = buildFromFiles(r, files)
		if err != nil {
			loggerV1.Error(err)
		}

		// Purge old cache entries
		if cache.Cache.ItemCount() > cfg.GetPackagesCacheCount() {
			items := cache.Cache.Items()
			oldestItem := ""
			for name, item := range items {
				if oldestItem == "" {
					oldestItem = name
				}
				if item.Expiration < items[oldestItem].Expiration {
					oldestItem = name
				}
			}
			cache.Cache.Delete(oldestItem)
			loggerV1.Infof("Purged oldest item %s from the cache", oldestItem)
		}
		duration := cfg.GetPackagesCacheDuration()
		cache.Cache.Set(fmt.Sprintf("synoPkgs-%s", paramLanguage), synoPkgs, duration)
		cache.Cache.Set(fmt.Sprintf("synoPkgs-%s-request", paramLanguage), r, 0)
		loggerV1.Debugf("Stored syno packages (%s) in cache for %s", paramLanguage, duration)
	}
	return synoPkgs, nil
}

// GetModelsFromCacheOrWeb gets packages from cache
// Or gets them from filesystem and populate cache
func GetModelsFromCacheOrWeb(forceRefresh bool) syno.Models {
	var synoModels syno.Models
	var err error
	if cachedSynoModels, cacheExpiration, found := cache.Cache.GetWithExpiration("synoModels"); found && !forceRefresh {
		synoModels = cachedSynoModels.(syno.Models)
		loggerV1.Debugf("Pulled syno models from cache, valid for %s", time.Until(cacheExpiration))
	} else {
		synoModels, err = syno.GetModels(forceRefresh)
		if err != nil {
			loggerV1.Warn("Could not get models ", err)
		}

		if synoModels != nil {
			duration := cfg.GetModelsCacheDuration()
			cache.Cache.Set("synoModels", synoModels, duration)
			loggerV1.Debugf("Stored syno models in cache for %s", duration)
		}
	}
	return synoModels
}
