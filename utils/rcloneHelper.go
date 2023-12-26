package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	rcloneVersion = "v1.65.0"
	baseURL       = "https://downloads.rclone.org"
)

func InstallRclone() error {
	//var cmd *exec.Cmd
	installRequiredCommands()
	switch os := getOS(); os {
	case "darwin", "linux":
		// macOS or Linux
		if isRcloneInstalled() {
			fmt.Println("rclone is already installed.")
		} else {
			fmt.Println("Downloading rclone binary...")
			rcloneURL, rcloneDownloadFolder := getRcloneURL(os)
			err := installRcloneOnLinux(rcloneURL, rcloneDownloadFolder)
			//err := installRcloneOnLinuxMacOSUsingScript()
			if err != nil {
				return err
			}
		}
		if !isRcloneRunning() {
			//runRCDcommand("rclone")
			err := createRcloneServiceOnLinux()
			if err != nil {
				return err
			}
			time.Sleep(15 * time.Second)
		}

	case "windows":
		// Windows
		if isRcloneInstalled() {
			fmt.Println("rclone is already installed.")
		} else {
			fmt.Println("Downloading rclone binary...")
			rcloneURL, rcloneDownloadFolder := getRcloneURL(os)
			err := installRcloneInWindows(rcloneURL, rcloneDownloadFolder)
			if err != nil {
				return err
			}
		}
		if !isRcloneRunning() {
			runRCDcommand("rclone.exe")
			time.Sleep(15 * time.Second)
		}

	default:
		return fmt.Errorf("unsupported operating system")
	}

	return nil
}

func getOS() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return os
	}
}
func isRcloneInstalled() bool {
	cmd := exec.Command("rclone", "version")

	// Redirect standard output and standard error
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Attempt to run the command
	err := cmd.Run()

	// Check for errors
	return err == nil
}
func RunPowerShellScript(scriptPath string) error {
	cmd := exec.Command("powershell.exe", "-File", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute PowerShell script: %v", err)
	}
	fmt.Println("rclone installation script executed successfully.")
	return nil
}

func InstallRcloneOnLinuxMacOSUsingScript() error {
	// Step 1: Refresh sudo credentials
	sudoCmd := exec.Command("sudo", "-v")
	sudoCmd.Stdout = os.Stdout
	sudoCmd.Stderr = os.Stderr

	err := sudoCmd.Run()
	if err != nil {
		return fmt.Errorf("error refreshing sudo credentials: %v", err)

	}

	// Step 2: Download and execute rclone install script
	curlCmd := exec.Command("curl", "https://rclone.org/install.sh")
	sudoBashCmd := exec.Command("sudo", "bash")

	// Create a pipe to connect the output of the curl command to the input of the sudo bash command
	curlOutput, err := curlCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating pipe: %v", err)
	}
	sudoBashCmd.Stdin = curlOutput
	sudoBashCmd.Stdout = os.Stdout
	sudoBashCmd.Stderr = os.Stderr

	// Start the curl command and then start the sudo bash command
	err = curlCmd.Start()
	if err != nil {
		return fmt.Errorf("error starting curl command: %v", err)
	}
	err = sudoBashCmd.Run()
	if err != nil {
		return fmt.Errorf("error running sudo bash command: %v", err)
	}
	err = curlCmd.Wait()
	if err != nil {
		return fmt.Errorf("error waiting for curl command to finish: %v", err)
	}

	fmt.Println("rclone installation script executed successfully.")
	return nil
}
func runRCDcommand(rcloneCommand string) {

	cmd := exec.Command(rcloneCommand, "rcd", "--rc-no-auth", "--rc-job-expire-duration=1h", "--log-file=E:\\rclone\\rclone.log", "-vv")

	// Set the working directory if needed
	// cmd.Dir = "/path/to/working/directory"

	// Set environment variables if needed
	// cmd.Env = append(os.Environ(), "KEY=VALUE")

	// Redirect standard output and standard error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func makeRequest(method, url string, payload []byte) ([]byte, error) {
	client := &http.Client{}

	// Create a request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// Make the request
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func QuitRclone() {
	url := "http://127.0.0.1:5572/core/quit"
	method := "POST"

	// Create a JSON payload
	payload := []byte(`{
		"exitCode": 1
	}`)
	// Make a Post request
	response, err := makeRequest(method, url, payload)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	// Print the response body as a string
	fmt.Println("Response:", string(response))
	fmt.Println("rclone closed successfully.")
}

func isRcloneRunning() bool {
	url := "http://127.0.0.1:5572/core/pid"
	method := "POST"
	// Make a Post request
	response, err := makeRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error making request:", err)
		return false
	}
	fmt.Println("rclone process id:", string(response))
	return true
}
func createRcloneServiceOnLinux() error {
	fmt.Println("createRcloneServiceOnLinux")
	//create working directory for service
	rcloneserviceWorkingDIr := "/opt/rclone"
	os.MkdirAll(rcloneserviceWorkingDIr, os.ModePerm)
	// Copy service file

	err := createRcloneService()
	if err != nil {
		return err
	}
	// cpCommand := exec.Command("sudo", "cp", "rclone.service", "/etc/systemd/system/")
	// cpErr := cpCommand.Run()
	// if cpErr != nil {
	// 	fmt.Println("Error while copying rclone.service:", cpErr)
	// 	return cpErr
	// }

	// Reload systemd
	reloadCmd := exec.Command("sudo", "systemctl", "daemon-reload")
	reloadErr := reloadCmd.Run()
	if reloadErr != nil {
		fmt.Println("Error:", reloadErr)
		return reloadErr
	}

	// Enable the service
	enableCmd := exec.Command("sudo", "systemctl", "enable", "rclone")
	enableErr := enableCmd.Run()
	if enableErr != nil {
		fmt.Println("Error:", enableErr)
		return enableErr
	}

	// Start the service
	startCmd := exec.Command("sudo", "systemctl", "start", "rclone")
	startErr := startCmd.Run()
	if startErr != nil {
		fmt.Println("Error:", startErr)
		return startErr
	}

	return nil
}

func getRcloneURL(osType string) (string, string) {

	switch osType {
	case "windows":
		rcloneDownloadedFolder := fmt.Sprintf("rclone-%s-windows-amd64", rcloneVersion)
		return fmt.Sprintf("%s/%s/%s.zip", baseURL, rcloneVersion, rcloneDownloadedFolder), rcloneDownloadedFolder
	case "linux":
		rcloneDownloadedFolder := fmt.Sprintf("rclone-%s-linux-amd64", rcloneVersion)
		return fmt.Sprintf("%s/%s/%s.zip", baseURL, rcloneVersion, rcloneDownloadedFolder), rcloneDownloadedFolder
	case "darwin":
		rcloneDownloadedFolder := fmt.Sprintf("rclone-%s-osx-amd64", rcloneVersion)
		return fmt.Sprintf("%s/%s/%s.zip", baseURL, rcloneVersion, rcloneDownloadedFolder), rcloneDownloadedFolder
	default:
		return "", ""
	}
}
func downloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func unzipFile(zipFile, destFolder string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destFolder, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}

	return nil
}
func installRcloneInWindows(rcloneURL string, rcloneDownloadFolder string) error {
	rcloneDownloadedZipFolder := fmt.Sprintf("%s.zip", rcloneDownloadFolder)
	err := downloadFile(rcloneURL, rcloneDownloadedZipFolder)
	if err != nil {
		fmt.Println("Error downloading rclone:", err)
		return err
	}
	fmt.Println("Unzipping rclone binary...")
	err = unzipFile(rcloneDownloadedZipFolder, ".")
	if err != nil {
		fmt.Println("Error unzipping rclone:", err)
		return err
	}
	// Copy binary file
	fmt.Println("Copying rclone binary to C:\\ProgramFiles\\rclone\\...")
	rcloneExeDestPath := filepath.Join(os.Getenv("ProgramFiles"), "rclone")
	// Create the destination directory if it doesn't exist
	os.MkdirAll(rcloneExeDestPath, os.ModePerm)

	cmd := exec.Command("cmd", "/c", "copy", filepath.Join(rcloneDownloadFolder, "rclone.exe"), rcloneExeDestPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error copying rclone binary:", err)
		return err
	}
	// Add rclone directory to the system PATH
	os.Setenv("Path", os.Getenv("Path")+string(filepath.ListSeparator)+rcloneExeDestPath)
	fmt.Println("rclone installation completed successfully!")

	// Now you can run rclone.exe directly
	rclonecmd := exec.Command("rclone", "version")
	rclonecmd.Stdout = os.Stdout
	rclonecmd.Stderr = os.Stderr
	err = rclonecmd.Run()
	if err != nil {
		fmt.Println("Error running rclone:", err)
		return err
	}

	return nil
}

func installRcloneOnLinux(rcloneURL string, rcloneDownloadFolder string) error {

	rcloneDownloadedZipFolder := fmt.Sprintf("%s.zip", rcloneDownloadFolder)
	err := downloadFile(rcloneURL, rcloneDownloadedZipFolder)
	if err != nil {
		fmt.Println("Error downloading rclone:", err)
		return err
	}
	fmt.Println("Unzipping rclone binary...")
	err = unzipFile(rcloneDownloadedZipFolder, ".")
	if err != nil {
		fmt.Println("Error unzipping rclone:", err)
		return err
	}
	// Copy binary file
	fmt.Println("Copying rclone binary to /usr/bin/...")
	cmd := exec.Command("sudo", "scp", rcloneDownloadFolder+"/rclone", "/usr/bin/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error copying rclone binary:", err)
		return err
	}

	// Set permissions
	fmt.Println("Setting permissions for rclone binary...")
	cmd = exec.Command("sudo", "chown", "root:root", "/usr/bin/rclone")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error setting permissions for rclone binary:", err)
		return err
	}

	cmd = exec.Command("sudo", "chmod", "755", "/usr/bin/rclone")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error setting permissions for rclone binary:", err)
		return err
	}

	// Install manpage
	fmt.Println("Installing rclone manpage...")
	cmd = exec.Command("sudo", "mkdir", "-p", "/usr/local/share/man/man1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error creating manpage directory:", err)
		return err
	}

	cmd = exec.Command("sudo", "scp", rcloneDownloadFolder+"/rclone.1", "/usr/local/share/man/man1/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error copying rclone manpage:", err)
		return err
	}

	// Update manpage database
	fmt.Println("Updating manpage database...")
	cmd = exec.Command("sudo", "mandb")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error updating manpage database:", err)
		return err
	}

	fmt.Println("rclone installation completed successfully!")
	// Now you can run rclone.exe directly
	rclonecmd := exec.Command("rclone", "version")
	rclonecmd.Stdout = os.Stdout
	rclonecmd.Stderr = os.Stderr
	err = rclonecmd.Run()
	if err != nil {
		fmt.Println("Error running rclone:", err)
		return err
	}

	return nil
}
func createRcloneService() error {
	// Define the content of the rclone.service file as a string
	content := `[Unit]
Description= rclone test server

Wants=network.target
After=syslog.target network-online.target

[Service]
Type=simple

WorkingDirectory=/opt/rclone
ExecStart=/usr/bin/rclone rcd --rc-no-auth --rc-job-expire-duration=365d --log-file=/var/log/rclone.log -vv
Restart=on-failure
RestartSec=10
KillMode=process
TimeoutStopSec=5

[Install]
WantedBy=multi-user.target`

	// Create the rclone.service file at /etc/systemd/system/ with 0644 permissions
	fmt.Println("creating rclone.service in /etc/systemd/system/")
	file, err := os.Create("/etc/systemd/system/rclone.service")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	// Sync the file to disk
	err = file.Sync()
	if err != nil {
		fmt.Println("Error syncing file:", err)
		return err
	}
	return nil
}

func installRequiredCommands() {
	requiredCommands := []string{"systemctl", "curl", "scp"}

	for _, command := range requiredCommands {
		if !isCommandAvailable(command) {
			fmt.Printf("%s is not installed. Installing...\n", command)
			err := installPackage(command)
			if err != nil {
				fmt.Printf("Failed to install %s: %v\n", command, err)
			} else {
				fmt.Printf("%s installed successfully.\n", command)
			}
		} else {
			fmt.Printf("%s is already installed.\n", command)
		}
	}
}
func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func installPackage(packageName string) error {
	cmd := exec.Command("sudo", "apt", "install", "-y", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
