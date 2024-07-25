package cleanup

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/CSPF-Founder/shieldsup/scanner-api/config"
	"github.com/CSPF-Founder/shieldsup/scanner-api/logger"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func CleanupAll(log *logger.Logger, config *config.Config) (bool, error) {

	os.Setenv("n_scaninput", "NONE")
	os.Setenv("n_scanstatus", "NONE")
	os.Setenv("docker_running_id", "NONE")

	dockerName := config.DockerName

	cli, err := client.NewClientWithOpts()
	if err != nil {
		log.Error("Error initializing docker client API", err)
		return false, err
	}

	stopList, err := cli.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("ancestor", dockerName)),
	})
	if err != nil {
		log.Error("Error getting docker containers", err)
		return false, err
	}

	for _, dockContainer := range stopList {
		if dockContainer.ID != "" {
			err := cli.ContainerStop(context.Background(), dockContainer.ID, container.StopOptions{})
			if err != nil {
				log.Error("Error stopping docker container", err)
				return false, err
			}
			err = cli.ContainerRemove(context.Background(), dockContainer.ID, container.RemoveOptions{Force: true})
			if err != nil {
				log.Error("Error removing docker container", err)
				return false, err
			}
		}
	}

	return resetTheTempDir(log, config)
}

// resetTheTempDir clears the temporary directory used by the scanner
func resetTheTempDir(log *logger.Logger, config *config.Config) (bool, error) {
	if _, err := os.Stat(config.LocalTmpDir); err == nil {
		if err := os.RemoveAll(config.LocalTmpDir); err != nil {
			log.Error("Error removing logs directory", err)
			return false, err
		}
	}

	if _, err := os.Stat(config.LocalTmpDir); os.IsNotExist(err) {
		if err := os.MkdirAll(config.LocalTmpDir, 0755); err != nil {
			log.Error("Error creating logs directory", err)
			return false, err
		}
	}

	return true, nil
}

func CleanupLogs(log *logger.Logger, config *config.Config) error {
	list_of_files, dirErr := os.ReadDir(config.LocalTmpDir)
	if dirErr != nil {
		log.Error("Error reading logs directory", dirErr)
		return dirErr
	}

	currentTime := time.Now()
	day := 24 * time.Hour // Duration representing one day

	for _, fileInfo := range list_of_files {
		fileLocation := filepath.Join(config.LocalTmpDir, fileInfo.Name())
		fileInfo, err := os.Stat(fileLocation)
		if err != nil {
			log.Error("Error reading file info", err)
			continue
		}

		fileTime := fileInfo.ModTime()
		if fileTime.Before(currentTime.Add(-day)) {
			log.Info("Deleting file: " + fileInfo.Name())
			// fmt.Printf("Delete: %s\n", fileInfo.Name())
			err := os.Remove(fileLocation)
			if err != nil {
				log.Error("Error deleting file", err)
				continue
			}
		}
	}

	return nil
}
