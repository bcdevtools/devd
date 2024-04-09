package files

import (
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

const fileExtensionTarLz4 = ".tar.lz4"

const osPathSeparator = string(os.PathSeparator)

// CompressLz4Command registers a sub-tree of commands
func CompressLz4Command() *cobra.Command {
	//goland:noinspection SpellCheckingInspection
	cmd := &cobra.Command{
		Use:     "compress-lz4 [input] [?output]",
		Short:   "Compress directory or file into .tar.lz4 file",
		Aliases: []string{"compress"},
		Args:    cobra.RangeArgs(1, 2),
		PreRun:  preRunCompressDecompressLz4,
		Run:     compressLz4,
	}

	return cmd
}

func preRunCompressDecompressLz4(_ *cobra.Command, _ []string) {
	if !utils.HasBinaryName("lz4") {
		libutils.PrintlnStdErr("ERR: lz4 is not installed, please install it first")
		os.Exit(1)
	}
	if !utils.HasBinaryName("tar") {
		libutils.PrintlnStdErr("ERR: tar is not installed, please install it first")
		os.Exit(1)
	}
}

func compressLz4(_ *cobra.Command, args []string) {
	var input, output string
	input = args[0]

	for {
		if !strings.HasSuffix(input, osPathSeparator) {
			break
		}
		input = strings.TrimSuffix(input, osPathSeparator)
	}

	// ensure input path can not have .tar.lz4 extension
	if strings.HasSuffix(input, fileExtensionTarLz4) {
		libutils.PrintlnStdErr("ERR: input path can not have", fileExtensionTarLz4, "extension:", input)
		os.Exit(1)
	}

	// ensure input exists
	fi, err := os.Stat(input)
	if err != nil {
		if os.IsNotExist(err) {
			libutils.PrintlnStdErr("ERR: input path does not exists:", input)
		} else {
			libutils.PrintlnStdErr("ERR: failed to stat input", input)
			libutils.PrintlnStdErr("ERR:", err)
		}
		os.Exit(1)
	}

	// build output file name
	if len(args) > 1 {
		output = args[1]

		if !strings.HasSuffix(output, fileExtensionTarLz4) {
			libutils.PrintlnStdErr("ERR: output file, if specified, must ends with", fileExtensionTarLz4, "extension:", output)
			os.Exit(1)
		}
	} else { // not supplied
		if fi.IsDir() {
			spl := strings.Split(input, osPathSeparator)
			output = spl[len(spl)-1] + fileExtensionTarLz4
		} else {
			_, fileName := filepath.Split(input)
			output = fileName + fileExtensionTarLz4
		}
	}

	// ensure output path does not exists
	_, err = os.Stat(output)
	if err != nil {
		if os.IsNotExist(err) {
			// good
		} else {
			libutils.PrintlnStdErr("ERR: failed to stat", output)
			libutils.PrintlnStdErr("ERR:", err)
			os.Exit(1)
		}
	} else {
		libutils.PrintlnStdErr("ERR: expected output lz4 file is already exists:", output)
		os.Exit(1)
	}

	execArgs := []string{
		"tar", "cvf", "-", input, "|", "lz4", "-", output,
	}
	executeCompressionDecompressionCommand(execArgs)
}

func executeCompressionDecompressionCommand(execArgs []string) {
	ec := utils.LaunchAppWithDirectStd(
		"/bin/bash",
		[]string{
			"-c",
			strings.Join(execArgs, " "),
		},
		nil,
	)
	if ec != 0 {
		libutils.PrintlnStdErr("ERR: failed to compress/decompress")
		os.Exit(ec)
	}
}
