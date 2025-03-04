package commands

import (
	"fmt"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/malice-plugins/pkgs/database/elasticsearch"
	"github.com/malice-plugins/pkgs/utils"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/database"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	"github.com/maliceio/malice/malice/persist"
	"github.com/maliceio/malice/plugins"
	"github.com/pkg/errors"
)

// cmdScan scans a sample with all appropriate malice plugins
func cmdScan(path string, logs bool) error {

	if len(path) > 0 {
		// Check that file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Fatal(path + ": no such file or directory")
		}

		docker := client.NewDockerClient()

		// clean stale containers from previous runs
		containers, err := container.List(docker, true)
		if err != nil {
			return errors.Wrap(err, "failed to list containers")
		}

		for _, contr := range containers {
			if utils.StringInSlice("malice", contr.Names) {
				err = container.Remove(docker, contr.ID, true, true, true)
				if err != nil {
					return errors.Wrapf(err, "failed to remove container: %s", contr.Names[0])
				}
			}
		}
		elasticsearchInDocker := false
		es := elasticsearch.Database{
			Index:    utils.Getopt("MALICE_ELASTICSEARCH_INDEX", "malice"),
			Type:     utils.Getopt("MALICE_ELASTICSEARCH_TYPE", "samples"),
			URL:      utils.Getopt("MALICE_ELASTICSEARCH_URL", config.Conf.DB.URL),
			Username: utils.Getopt("MALICE_ELASTICSEARCH_USERNAME", config.Conf.DB.Username),
			Password: utils.Getopt("MALICE_ELASTICSEARCH_PASSWORD", config.Conf.DB.Password),
		}

		// This assumes you haven't set up an elasticsearch instance and that malice should create one
		if strings.EqualFold(es.URL, "http://localhost:9200") {
			elasticsearchInDocker = true
			// Check that database is running
			if _, running, _ := container.Running(docker, config.Conf.DB.Name); !running {
				log.Error("database is NOT running, starting now...")
				err := database.Start(docker, es, logs)
				if err != nil {
					return errors.Wrap(err, "failed to start to database")
				}
			}
		}

		// Initialize the malice database
		es.Init()

		// Check Plugin Status
		if plugins.InstalledPluginsCheck(docker) {
			log.Debug("All enabled plugins are installed.")
		} else {
			// Prompt user to install all plugins?
			fmt.Println("All enabled plugins not installed would you like to install them now? (yes/no)")
			fmt.Println("[Warning] This can take a while if it is the first time you have ran Malice.")
			if utils.AskForConfirmation() {
				plugins.UpdateEnabledPlugins(docker)
			}
		}

		es.Plugins = database.GetPluginsByCategory()

		file := persist.File{Path: path}
		file.Init()

		// Output File Hashes
		file.ToMarkdownTable()
		// fmt.Println(string(file.ToJSON()))

		//////////////////////////////////////
		// Copy file to malice volume
		container.CopyToVolume(docker, file)

		//////////////////////////////////////
		// Write all file data to the Database
		resp, err := es.StoreFileInfo(structs.Map(file))
		if err != nil {
			return errors.Wrap(err, "scan cmd failed to store file info")
		}

		scanID := resp.Id

		/////////////////////////////////////////////////////////////////
		// Run all Intel Plugins on the md5 hash associated with the file
		plugins.RunIntelPlugins(docker, file.SHA1, scanID, true, elasticsearchInDocker)

		// Get file's mime type
		mimeType, err := persist.GetMimeType(docker, file.SHA256)
		if err != nil {
			return errors.Wrap(err, "failed to get file's mime type")
		}

		log.Debug("looking for plugins that will run on: ", mimeType)
		// Iterate over all applicable installed plugins
		pluginsForMime := plugins.GetPluginsForMime(mimeType, true)
		log.Debug("found these plugins: ")
		for _, plugin := range pluginsForMime {
			log.Debugf(" - %v", plugin.Name)
		}

		var wg sync.WaitGroup
		wg.Add(len(pluginsForMime))

		for _, plugin := range pluginsForMime {
			log.Debugf(">>>>> RUNNING Plugin: %s >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", plugin.Name)
			// Start Plugin Container
			// TODO: don't use the default of true for --logs
			go plugin.StartPlugin(docker, file.SHA256, scanID, true, elasticsearchInDocker, &wg)
		}

		wg.Wait() // this waits for the counter to be 0
		log.Debug("Done with plugins.")
	} else {
		log.Error("please supply a valid file to scan")
	}

	return nil
}

// APIScan is an API wrapper for cmdScan
func APIScan(file string) error {
	return cmdScan(file, false)
}
