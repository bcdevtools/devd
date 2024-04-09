package files

import (
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	flagDirectory = "directory"
)

// DecompressLz4Command registers a sub-tree of commands
func DecompressLz4Command() *cobra.Command {
	//goland:noinspection SpellCheckingInspection
	cmd := &cobra.Command{
		Use:     "decompress-lz4 [input lz4 file]",
		Short:   "Decompress .tar.lz4 file",
		Aliases: []string{"decompress"},
		Args:    cobra.ExactArgs(1),
		PreRun:  preRunCompressDecompressLz4,
		Run:     decompressLz4,
	}

	cmd.Flags().StringP(flagDirectory, "C", "", "output directory, default is the current one")

	return cmd
}

func decompressLz4(cmd *cobra.Command, args []string) {
	input := args[0]

	// ensure input is lz4 file
	if !strings.HasSuffix(strings.ToLower(input), fileExtensionTarLz4) {
		libutils.PrintlnStdErr("ERR: input file must have", fileExtensionTarLz4, "extension:", input)
		os.Exit(1)
	}

	// ensure input is exists
	fi, err := os.Stat(input)
	if err != nil {
		if os.IsNotExist(err) {
			libutils.PrintlnStdErr("ERR: input file does not exists:", input)
		} else {
			libutils.PrintlnStdErr("ERR: failed to stat", input)
			libutils.PrintlnStdErr("ERR:", err)
		}
		os.Exit(1)
	} else if fi.IsDir() {
		libutils.PrintlnStdErr("ERR: input file is a directory:", input)
		os.Exit(1)
	}

	var outputDirPath string
	outputDirPath, err = cmd.Flags().GetString(flagDirectory)
	libutils.ExitIfErr(err, "failed to read flag --"+flagDirectory)
	outputDirPath = strings.TrimSpace(outputDirPath)

	if len(outputDirPath) > 0 {
		exists, _, err := utils.IsDirAndExists(outputDirPath)
		libutils.ExitIfErr(err, "failed to check existence of the specified output directory "+outputDirPath)
		if !exists {
			libutils.PrintlnStdErr("ERR: the specified output directory", outputDirPath, "is not exists")
			os.Exit(1)
		}
	}

	execArgs := []string{"lz4", "-c", "-d", input, "|", "tar", "-x"}
	if len(outputDirPath) > 0 {
		execArgs = append(execArgs, "-C", outputDirPath)
	}

	executeCompressionDecompressionCommand(execArgs)
}
