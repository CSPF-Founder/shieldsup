package scanner

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	TargetHash     string
	OutputDir      string
	targetType     enums.TargetType
}

func NewScannerModule(conf *config.Config, logger *logger.Logger, target string) *ScannerModule {
	targetHash := stringToHash(target)
	logger.Info(fmt.Sprintf("Target: %s, Hash: %s", target, targetHash))

	// TODO: return error for invalid target
	targetType := enums.ParseTargetType(target)

	return &ScannerModule{
		Logger:         logger,
		TargetAddress:  target,
		TemplateFolder: conf.TemplateFolder,
		DockerName:     conf.DockerName,
		TargetHash:     targetHash,
		OutputDir:      conf.LocalTmpDir,
		Config:         conf,
		targetType:     targetType,
	}
}

func stringToHash(target string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(target)))
}

func (s *ScannerModule) getOutputFileName() string {
	return s.TargetHash + ".json"
}

func (s *ScannerModule) getOutputFilePath() string {
	return filepath.Join(s.OutputDir, s.getOutputFileName())
}

func (s *ScannerModule) getDastFileName() string {
	return s.TargetHash + "_dast.json"
}

func (s *ScannerModule) getDastFilePath() string {
	return filepath.Join(s.OutputDir, s.getDastFileName())
}

func (s *ScannerModule) getWebScanFileName() string {
	return s.TargetHash + "_web.json"
}

func (s *ScannerModule) getWebScanFilePath() string {
	return filepath.Join(s.OutputDir, s.getWebScanFileName())
}

func (s *ScannerModule) getTargetsFileName() string {
	return s.TargetHash + "_targets.text"
}

func (s *ScannerModule) getTargetsFilePath() string {
	return filepath.Join(s.OutputDir, s.getTargetsFileName())
}

func (s *ScannerModule) getTempStateFileName() string {
	return s.TargetHash + "_temp.json"
}

func (s *ScannerModule) getTempStateFilePath() string {
	return filepath.Join(s.OutputDir, s.getTempStateFileName())
}

func (s *ScannerModule) UpdateScanStatus(status enums.ScanStatus) {

	data := map[string]enums.ScanStatus{"status": status}

	jsonData, err := json.Marshal(data)
	if err != nil {
		s.Logger.Error("Error Marshal Status", err)
	}

	tempStateFile, err := os.OpenFile(s.getTempStateFilePath(), os.O_WRONLY|os.O_CREATE, 0666)
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

	_, err := os.Stat(s.getTempStateFilePath())
	if os.IsNotExist(err) {
		s.Logger.Error("TempStateFile Does Not Exists", err)
		return enums.ScanStatusDoesNotExist
	}

	outfile, err := os.OpenFile(s.getTempStateFilePath(), os.O_RDONLY, 0)
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

// runPortScanner runs the port scanner before the web scan
func (s *ScannerModule) runPortScanner(ctx context.Context) error {
	ipCount, _ := iputils.GetIPCountIfRange(s.TargetAddress)
	portScanTimeout := 60 * time.Second
	if ipCount > 1 {
		portScanTimeout = time.Duration(ipCount) * portScanTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, portScanTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "naabu", "-host", s.TargetAddress, "-rate", "100", "-silent", "-o", s.getTargetsFilePath())

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error in Naabu: %w", err)
	}

	if stdOut.String() == "" {
		return fmt.Errorf("No ports found by Naabu")
	}

	return nil

}

// runCrawler runs the crawler before the web scan
func (s *ScannerModule) runCrawler(ctx context.Context) error {
	katanatime := 60 * time.Second
	ctx, cancel := context.WithTimeout(ctx, katanatime)
	defer cancel()

	cmd := exec.CommandContext(ctx, "katana", "-silent", "-kf", "all", "-jsl", "-aff", "-ct", "119", "-c", "10", "-d", "5", "-jsonl", "-u", s.TargetAddress, "-o", s.getTargetsFilePath())

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error in Katana: %w", err)
	}

	if stdOut.String() == "" {
		return fmt.Errorf("No URL found by Katana")
	}

	return nil

}

// getContainerResources returns the resources for the container
func getContainerResources() (*container.Resources, error) {
	containerRAM, err := getRAMForContainer()
	if err != nil {
		return nil, fmt.Errorf("Error getting RAM for container: %w", err)
	}
	containerCPU := getNanoCPUForContainer()

	resources := &container.Resources{
		Memory:   containerRAM,
		NanoCPUs: containerCPU,
	}

	return resources, nil
}

// getContainerHostConfig returns the host configuration for the container
func (s *ScannerModule) getContainerHostConfig() (*container.HostConfig, error) {
	ctrResource, err := getContainerResources()
	if err != nil {
		s.Logger.Error("Error getting container resources", err)
		return nil, err
	}

	ctrVolumes := []mount.Mount{
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
	}

	hostConfig := &container.HostConfig{
		NetworkMode: "host",
		Mounts:      ctrVolumes,
		Resources:   *ctrResource,
		AutoRemove:  true, // Remove the container after it exits
	}
	return hostConfig, nil
}

// StartDASTScan runs the DAST scan
func (s *ScannerModule) StartDASTScan(ctx context.Context, dockerClient *client.Client) error {
	config := &container.Config{
		Image: s.Config.DockerName,
		Cmd: []string{
			"-im", "jsonl",
			"-l", "/nucout/" + s.getTargetsFileName(),
			"-o", "/nucout/" + s.getDastFileName(),
			"-j",
			"-dast",
		},
	}

	containerName := "scan-dast-" + s.TargetHash
	err := s.runScanCtr(ctx, dockerClient, config, containerName)
	if err != nil {
		return fmt.Errorf("Error running container: %w", err)
	}

	return nil
}

// StartWebScan runs the web scan
// 1. Run the crawler
// 2. If the crawler fails, skips the DAST scan
// 3. If the crawler is successful, run the DAST scan
// 4. Run the normal scan
// 5. Merge the results of the DAST and normal scan
func (s *ScannerModule) StartWebScan(ctx context.Context, dockerClient *client.Client) error {
	// if the target is a URL, run the crawler
	err := s.runCrawler(ctx)
	if err != nil {
		// if the crawler fails, run the nuclei scan with the target address
		s.Logger.Error("Error running crawler", err)
	} else {
		// if the crawler is successful, run the nuclei scan with the crawled URLs
		err = s.StartDASTScan(ctx, dockerClient)
		if err != nil {
			return err
		}
	}

	config := &container.Config{
		Image: s.Config.DockerName,
		Cmd:   []string{"-l", "/nucout/" + s.getTargetsFileName(), "-o", "/nucout/" + s.getOutputFileName(), "-j"},
	}
	// Run normal scan
	config.Cmd = []string{
		"-u", s.TargetAddress,
		"-o", "/nucout/" + s.getWebScanFileName(),
		"-j",
	}

	containerName := "scan-" + s.TargetHash
	err = s.runScanCtr(ctx, dockerClient, config, containerName)
	if err != nil {
		return fmt.Errorf("Error running container: %w", err)
	}

	if err = s.mergeWebResults(); err != nil {
		return err
	}

	return nil
}

func (s *ScannerModule) mergeWebResults() error {

	dst, err := os.OpenFile(s.getOutputFilePath(), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Error opening final output file: %w", err)
	}
	defer dst.Close()

	// Open the first file for reading if it exists
	dastSrc, err := os.Open(s.getDastFilePath())
	if os.IsNotExist(err) {
		s.Logger.Info("Dast Scan output does not exist, skipping")
	} else if err != nil {
		return fmt.Errorf("Error opening first file: %w", err)
	} else {
		defer dastSrc.Close()
		// Copy the content from the first file to the merged file
		if _, err := io.Copy(dst, dastSrc); err != nil {
			return fmt.Errorf("Error copying content from first file: %w", err)
		}
	}

	// Open the second file for reading if it exists
	normalSrc, err := os.Open(s.getWebScanFilePath())
	if os.IsNotExist(err) {
		s.Logger.Info("Web Scan output does not exist, skipping")
	} else if err != nil {
		return fmt.Errorf("Error opening second file: %w", err)
	} else {
		defer normalSrc.Close()
		// Copy the content from the second file to the merged file
		if _, err := io.Copy(dst, normalSrc); err != nil {
			return fmt.Errorf("Error copying content from second file: %w", err)
		}
	}

	return nil
}

// StartNetworkScan runs the network scan
func (s *ScannerModule) StartNetworkScan(ctx context.Context, dockerClient *client.Client) error {
	ctrCmd := []string{"-l", "/nucout/" + s.getTargetsFileName(), "-o", "/nucout/" + s.getOutputFileName(), "-j"}
	err := s.runPortScanner(ctx)
	if err != nil {
		// If port scan file does not exist, run the nuclei scan with the target address
		s.Logger.Error("Error running port scanner", err)
		ctrCmd = []string{"-u", s.TargetAddress, "-o", "/nucout/" + s.getOutputFileName(), "-j"}
	}

	// Default input is the port scan file
	config := &container.Config{
		Image: s.Config.DockerName,
		Cmd:   ctrCmd,
	}

	containerName := "scan-" + s.TargetHash
	err = s.runScanCtr(ctx, dockerClient, config, containerName)
	if err != nil {
		return fmt.Errorf("Error running container: %w", err)
	}

	return nil

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

	switch s.targetType {
	case enums.TargetTypeIPRange, enums.TargetTypeIP:
		err = s.StartNetworkScan(ctx, dockerClient)
		if err != nil {
			s.Logger.Error("Error Starting Network Scan", err)
			return
		}
	case enums.TargetTypeURL:
		err = s.StartWebScan(ctx, dockerClient)
		if err != nil {
			s.Logger.Error("Error Starting Web Scan", err)
			return
		}
	default:
		s.Logger.Error("Invalid target type", nil)
		return
	}

	os.Setenv("n_scanstatus", "NONE")

	// Update scan status
	scanCompleted = true
	s.UpdateScanStatus(enums.ScanStatusCompleted)

}

func (s ScannerModule) runScanCtr(ctx context.Context, dockerClient *client.Client, config *container.Config, containerName string) error {

	hostConfig, err := s.getContainerHostConfig()
	if err != nil {
		return fmt.Errorf("Error getting container host config: %w", err)
	}

	resp, err := dockerClient.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return fmt.Errorf("Error Creating Docker Container: %w", err)
	}

	containerID := resp.ID

	err = dockerClient.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("Error Starting Docker Container: %w", err)
	}

	_, err = dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		return fmt.Errorf("Error During Docker Container Inspects: %w", err)
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

	// Remove the container if not clean exit
	if !cleanExit {
		if err := dockerClient.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true}); err != nil {
			s.Logger.Error("Error removing container", err)
		}
		return fmt.Errorf("The container did not exit cleanly")
	}

	return nil
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
	alerts := parser.Parse(s.getOutputFilePath(), s.Logger)

	results := make(map[string]any)
	results["success"] = true
	results["scan_status"] = enums.ScanStatusCompleted
	results["data"] = alerts

	return results

}
