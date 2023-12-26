package main

import (
	"fmt"

	"github.com/PCCloudnAndRMM/RcloneSetup/utils"
)

func main() {
	fmt.Println("Set up rclone")
	utils.InstallRclone()
}
