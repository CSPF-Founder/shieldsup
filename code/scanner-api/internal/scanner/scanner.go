package scanner

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/CSPF-Founder/shieldsup/scanner-api/config"
	"github.com/CSPF-Founder/shieldsup/scanner-api/enums"
	"github.com/CSPF-Founder/shieldsup/scanner-api/logger"
	"github.com/CSPF-Founder/shieldsup/scanner-api/utils/iputils"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type ScannerModule struct {
	Logger         *logger.Logger
	Config         *config.Config
	TargetAddress  string
	TemplateFolder string
	DockerName     string
	FileNameHash   string
	OutputDir      string
	TempStateFile  string
	OutputFileName string
	OutputFile     string
}

func NewScannerModule(conf *config.Config, logger *logger.Logger, target string) *ScannerModule {
	fileNameHash := fmt.Sprintf("%x", sha256.Sum256([]byte(target)))
	logger.Info(fmt.Sprintf("Target: %s, Hash: %s", target, fileNameHash))
	outputDir := conf.LocalTmpDir

	tempStateFile := filepath.Join(outputDir, fileNameHash+"_temp"+".json")

	outputFileName := fileNameHash + ".json"
	outputFile := filepath.Join(outputDir, outputFileName)

	fileCreateErr := os.MkdirAll(filepath.Dir(outputFile), 0755)
	if fileCreateErr != nil {
		logger.Error("Error creating directory:", fileCreateErr)
	}

	return &ScannerModule{
		Logger:         logger,
		TargetAddress:  target,
		TemplateFolder: conf.TemplateFolder,
		DockerName:     conf.DockerName,
		FileNameHash:   fileNameHash,
		OutputDir:      outputDir,
		TempStateFile:  tempStateFile,
		OutputFileName: outputFileName,
		OutputFile:     outputFile,
		Config:         conf,
	}
}

func (s *ScannerModule) UpdateScanStatus(status enums.ScanStatus) {

	data := map[string]enums.ScanStatus{"status": status}

	jsonData, err := json.Marshal(data)
	if err != nil {
		s.Logger.Error("Error Marshal Status", err)
	}

	fileCreateErr := os.MkdirAll(filepath.Dir(s.TempStateFile), 0666)
	if fileCreateErr != nil {
		s.Logger.Error("Error creating directory:", err)
	}

	tempStateFile, err := os.OpenFile(s.TempStateFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		s.Logger.Error("Error Opening TempStateFile", err)
	}

	defer tempStateFile.Close()

	_, err = tempStateFile.Write(jsonData)
	if err != nil {
		s.Logger.Error("Error Writing in TempStateFile", err)
	}
}

func (s *ScannerModule) GetScanStatus() enums.ScanStatus {

	_, err := os.Stat(s.TempStateFile)
	if os.IsNotExist(err) {
		s.Logger.Error("TempStateFile Does Not Exists", err)
		return enums.ScanStatusDoesNotExist
	}

	outfile, err := os.OpenFile(s.TempStateFile, os.O_RDONLY, 0)
	if err != nil {
		s.Logger.Error("Error Opening TempStateFile", err)
		return enums.ScanStatusDoesNotExist
	}

	defer outfile.Close()

	// Decoding JSON data from the file
	var jsonData map[string]int
	err = json.NewDecoder(outfile).Decode(&jsonData)

	if err != nil {
		s.Logger.Error("Error decoding JSON:", err)
		return enums.ScanStatusDoesNotExist
	}

	// Check if JSON data is empty
	if len(jsonData) == 0 {
		s.Logger.Info("Scan status does not exist")
		return enums.ScanStatusDoesNotExist
	}

	status, ok := jsonData["status"]
	if !ok {
		s.Logger.Info("Scan status does not exist")
		return enums.ScanStatusDoesNotExist
	}

	if status == int(enums.ScanStatusStarted) {
		return enums.ScanStatusNotCompleted
	}
	if status == int(enums.ScanStatusCompleted) {
		return enums.ScanStatusCompleted
	}

	return enums.ScanStatusDoesNotExist
}

func (s *ScannerModule) CalculateScanTimeout() int {
	scanTimeout := 60 * 30 // 30 minutes
	ipCount, _ := iputils.GetIPCountIfRange(s.TargetAddress)

	if ipCount > 1 {
		scanTimeout += 60 * 15 * ipCount
	}
	return scanTimeout
}

// getRAMSize returns the total RAM size
func getRAMSize() (uint64, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 0, err
	}
	totalRAM := info.Totalram * uint64(info.Unit)

	return totalRAM, nil
}

// getRAMForContainer returns the RAM size for the container
func getRAMForContainer() (int64, error) {
	totalRAM, err := getRAMSize()
	if err != nil {
		return 0, err
	}

	// 1/2th of total RAM
	containerRAM := int64(totalRAM / 2)

	if containerRAM < 1024 {
		containerRAM = 1024
	}

	return containerRAM, nil
}

// getNanoCPUForContainer returns the CPU size for the container
// Divides the number of CPUs by 2
func getNanoCPUForContainer() int64 {
	halfCPU := runtime.NumCPU() / 2
	if halfCPU <= 1 {
		return 1e9
	}
	// numCPU := float64(halfCPU) * 0.8
	// return int64(numCPU) * 1e9
	return int64(halfCPU) * 1e9
}

func (s *ScannerModule) StartScan() {
	scanCompleted := false
	defer func() {
		// If scan is not completed, update the status
		if !scanCompleted {
			s.UpdateScanStatus(enums.ScanStatusCouldNotCompleted)
		}
	}()
	os.Setenv("n_scaninput", s.TargetAddress)
	os.Setenv("n_scanstatus", "RUNNING")
	s.UpdateScanStatus(enums.ScanStatusStarted)

	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		s.Logger.Error("Error Initiating Docker", err)
		return
	}

	containerRAM, err := getRAMForContainer()
	if err != nil {
		s.Logger.Error("Error getting RAM for container", err)
		return
	}
	containerCPU := getNanoCPUForContainer()

	config := &container.Config{
		Image: s.Config.DockerName,
		Cmd:   []string{"-u=" + s.TargetAddress, "-o=/nucout/" + s.OutputFileName, "-j"},
	}

	hostConfig := &container.HostConfig{
		NetworkMode: "host",
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: s.TemplateFolder,
				Target: "/root/nuclei-templates",
			},
			{
				Type:   mount.TypeBind,
				Source: s.OutputDir,
				Target: "/nucout/",
			},
		},
		Resources: container.Resources{
			Memory:   containerRAM,
			NanoCPUs: containerCPU,
		},
	}

	containerName := "scan-" + s.FileNameHash

	resp, err := dockerClient.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		s.Logger.Error("Error Creating Docker Container", err)
		return
	}

	containerID := resp.ID

	err = dockerClient.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		s.Logger.Error("Error Starting Docker Container", err)
		return
	}

	_, err = dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		s.Logger.Error("Error During Docker Container Inspects", err)
		return
	}

	scanWaitTime := s.CalculateScanTimeout()

	tEnd := time.Now().Add(time.Duration(scanWaitTime) * time.Second)

	cleanExit := false

	for time.Now().Before(tEnd) {
		// Get the target container
		targetContainer, err := dockerClient.ContainerInspect(context.Background(), containerID)
		if err != nil {
			s.Logger.Error("Error getting container", err)
			cleanExit = true
			break
		}

		if targetContainer.State.Status == "exited" {
			cleanExit = true
			break
		}

		time.Sleep(10 * time.Second)
	}

	os.Setenv("n_scanstatus", "NONE")
	// Remove the container if not clean exit
	if !cleanExit {
		if err := dockerClient.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true}); err != nil {
			s.Logger.Error("Error removing container", err)
		}
		return
	}

	// Update scan status
	scanCompleted = true
	s.UpdateScanStatus(enums.ScanStatusCompleted)

}

// RetrieveResults Function
func (s ScannerModule) RetrieveResults(force bool) map[string]any {
	scanStatus := s.GetScanStatus()
	if force {
		return s.ReturnScanResults()
	}
	switch scanStatus {
	case enums.ScanStatusDoesNotExist:
		return map[string]any{
			"success":     false,
			"scan_status": enums.ScanStatusDoesNotExist,
			"messages":    []string{"Scan does not exist"},
		}
	case enums.ScanStatusNotCompleted:
		return map[string]any{
			"success":     false,
			"scan_status": enums.ScanStatusNotCompleted,
			"messages":    []string{"Scan not completed"},
		}
	case enums.ScanStatusCompleted:
		return s.ReturnScanResults()
	default:
		return map[string]any{
			"success":     false,
			"scan_status": enums.ScanStatusDoesNotExist,
			"messages":    []string{"Scan does not exist"},
		}
	}
}

type ScanResult struct {
	Success    bool
	ScanStatus enums.ScanStatus
	Data       []map[string]any
}

func (s ScannerModule) ReturnScanResults() map[string]any {

	parser := NewAlertsParser()
	alerts := parser.Parse(s.OutputFile, s.Logger)

	results := make(map[string]any)
	results["success"] = true
	results["scan_status"] = enums.ScanStatusCompleted
	results["data"] = alerts

	return results

}
