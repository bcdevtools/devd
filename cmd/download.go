package cmd

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	flagConcurrentDownload = "concurrent"
	flagWorkingDir         = "working-directory"
	flagOutputFileName     = "output"
)

// downloadCmd represents the version command, it downloads file from URL
var downloadCmd = &cobra.Command{
	Use:     "download",
	Aliases: []string{"dl"},
	Args:    cobra.ExactArgs(1),
	Short:   "Download a file from URL",
	Run: func(cmd *cobra.Command, args []string) {
		concurrent, _ := cmd.Flags().GetInt(flagConcurrentDownload)
		concurrent = libutils.MaxInt(1, concurrent)

		workingDir, _ := cmd.Flags().GetString(flagWorkingDir)
		if len(workingDir) > 0 {
			workingDirInfo, err := os.Stat(workingDir)
			if os.IsNotExist(err) {
				libutils.PrintlnStdErr("ERR: specified working directory does not exists:", workingDir)
				os.Exit(1)
			}
			if !workingDirInfo.IsDir() {
				libutils.PrintlnStdErr("ERR: specified working directory is not a directory:", workingDir)
				os.Exit(1)
			}
		}

		outputFileName, _ := cmd.Flags().GetString(flagOutputFileName)
		if len(outputFileName) > 0 {
			outputInfo, err := os.Stat(outputFileName)
			if err == nil && outputInfo.IsDir() {
				libutils.PrintlnStdErr("ERR: specified output is a directory:", outputFileName)
				os.Exit(1)
			}
			actualFileName := filepath.Base(outputFileName)
			if actualFileName != outputFileName {
				libutils.PrintfStdErr("ERR: specified output file name is not just a file name, should be: %s. You can specify the directory to download into using flag --%s\n", actualFileName, flagWorkingDir)
				os.Exit(1)
			}
		}

		_, contentLengthDesc, err := utils.TryGetDownloadFileSizeInBytes(args[0])
		if err == nil {
			fmt.Println("File size:", contentLengthDesc)
		} else {
			libutils.PrintlnStdErr("WARN: failed to get file size:", err)
		}

		if utils.HasBinaryName("aria2c") {
			launchDownloadAria2c(args[0], concurrent, workingDir, outputFileName)
			return
		} else {
			fmt.Println("WARN: it is recommended to install aria2c for better performance: sudo apt install -y aria2")
			time.Sleep(10 * time.Second)
		}

		if utils.HasBinaryName("wget") {
			fmt.Println("WARN: aria2c is not installed, fallback to wget")
			launchDownloadWget(args[0], workingDir, outputFileName)
			return
		}

		if utils.HasBinaryName("curl") {
			fmt.Println("WARN: aria2c & wget are not installed, fallback to curl")
			if len(outputFileName) < 1 {
				r, _ := http.NewRequest("GET", args[0], nil)
				outputFileName = path.Base(r.URL.Path)
				fmt.Println("File will be downloaded as", outputFileName)
			}
			launchDownloadCurl(args[0], workingDir, outputFileName)
			return
		}

		libutils.PrintlnStdErr("Neither tool installed: aria2c & wget & curl! Please install at least one of them first!")
		os.Exit(1)
	},
}

func launchDownloadAria2c(url string, concurrent int, workingDir, outputFileName string) {
	split := libutils.MinInt(5, concurrent)
	maxConnectionPerServer := concurrent

	args := []string{
		url,
		"-x", fmt.Sprintf("%d", maxConnectionPerServer),
		"-s", fmt.Sprintf("%d", split),
	}

	if len(workingDir) > 0 {
		args = append(args, "-d", workingDir)
	}
	if len(outputFileName) > 0 {
		args = append(args, "-o", outputFileName)
	}

	launchDownloader("aria2c", args, workingDir)
}

func launchDownloadWget(url string, workingDir, outputFileName string) {
	args := []string{url}

	if len(outputFileName) > 0 {
		if len(workingDir) > 0 {
			args = append(args, "-O", path.Join(workingDir, outputFileName))
		} else {
			args = append(args, "-O", outputFileName)
		}
	}

	launchDownloader("wget", args, workingDir)
}

func launchDownloadCurl(url string, workingDir, outputFileName string) {
	args := []string{url}

	if len(outputFileName) < 1 {
		panic("output file name must be specified")
	}

	if len(workingDir) > 0 {
		args = append(args, "-o", path.Join(workingDir, outputFileName))
	} else {
		args = append(args, "-o", outputFileName)
	}

	launchDownloader("curl", args, workingDir)
}

func launchDownloader(appName string, args []string, workingDir string) {
	if len(workingDir) > 0 {
		fmt.Println("Launching downloader at", workingDir)
	} else {
		fmt.Println("Launching downloader")
	}

	fmt.Printf("> %s %s\n", appName, strings.Join(args, " "))

	ec := utils.LaunchAppWithSetup(appName, args, func(launchCmd *exec.Cmd) {
		if len(workingDir) > 0 {
			launchCmd.Dir = workingDir
		}
		launchCmd.Stdin = os.Stdin
		launchCmd.Stdout = os.Stdout
		launchCmd.Stderr = os.Stderr
	})

	if ec != 0 {
		os.Exit(ec)
	}
}

func init() {
	downloadCmd.Flags().IntP(
		flagConcurrentDownload, "c",
		4,
		"number of concurrent downloads (only supported for aria2c)",
	)

	downloadCmd.Flags().StringP(
		flagWorkingDir, "D",
		"",
		"working directory",
	)

	downloadCmd.Flags().StringP(
		flagOutputFileName, "o",
		"",
		"output document name",
	)

	rootCmd.AddCommand(downloadCmd)
}
