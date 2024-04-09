package cmd

import (
	"bufio"
	"fmt"
	"github.com/EscanBE/go-ienumerable/goe"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/bcdevtools/devd/cmd/utils"
	"github.com/bcdevtools/devd/constants"
	"github.com/spf13/cobra"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	flagConfirmExecution = "yes"
)

var predefinedAliases map[string]predefinedAlias
var longestUseDesc int

/*
Sample content for alias file .devd_alias:
echo "say-hello	echo \"Hello World\"" >> ~/.devd_alias
devd a say-hello
*/

// aliasCmd represents the 'a' command, it executes commands based on pre-defined input alias
var aliasCmd = &cobra.Command{
	Use:     "a [alias]",
	Aliases: []string{"alias"},
	Short:   "Execute commands based on pre-defined alias",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := types.UnwrapAppContext(cmd.Context())

		userInfo := ctx.GetWorkingUserInfo()

		predefinedAliases = make(map[string]predefinedAlias)

		registerStartupPredefinedAliases(userInfo)
		registerPredefinedAliasesFromFile(userInfo)

		if len(args) < 1 {
			if len(predefinedAliases) < 1 {
				fmt.Println("No registered alias")
				return
			}

			lineFormat := " %-" + fmt.Sprintf("%d", longestUseDesc+1) + "s: %s\n"

			fmt.Println("Registered aliases:")
			for _, alias := range goe.NewIEnumerable[string](libutils.GetKeys(predefinedAliases)...).Order().GetOrderedEnumerable().ToArray() {
				pa := predefinedAliases[alias]
				if pa.overridden {
					fmt.Printf(" *overriden*")
				}
				fmt.Printf(lineFormat, pa.use, strings.Join(pa.command, " "))
			}
			fmt.Printf("Alias can be customized by adding into ~/%s (TSV format with each line content \"<alias><tab><command>\")\n", constants.PREDEFINED_ALIAS_FILE_NAME)
			return
		}

		selectedAlias := args[0]
		pa, found := predefinedAliases[selectedAlias]
		if !found {
			fmt.Printf("Alias '%s' has not been registered before\n", selectedAlias)
			os.Exit(1)
		}

		command := pa.command
		if pa.alwaysInvokeCommandAlter {
			if len(args) < 2 {
				libutils.PrintlnStdErr("ERR: this alias requires argument")
				os.Exit(1)
			}
			command = (pa.alter)(command, args[1:])
		} else if len(args) > 1 && pa.alter != nil {
			command = (pa.alter)(command, args[1:])
		}

		if len(command) < 1 {
			panic("empty command")
		}

		joinedCommand := strings.Join(command, " ")

		confirmExecution, _ := cmd.Flags().GetBool(flagConfirmExecution)

		if confirmExecution {
			const waitingTime = 10
			fmt.Println("Pending execution command:")
			fmt.Printf("> %s\n", joinedCommand)
			fmt.Printf("(actual command: [/bin/bash] [-c] [%s])\n", joinedCommand)
			fmt.Printf("Executing in %d seconds...\n", waitingTime)
			time.Sleep(waitingTime * time.Second)
		} else {
			fmt.Println("Are you sure want to execute the following command?")
			fmt.Printf("> %s\n", joinedCommand)
			fmt.Printf("(actual command: [/bin/bash] [-c] [%s])\n", joinedCommand)
			fmt.Println("Yes/No?")

			reader := bufio.NewReader(os.Stdin)
			yes, err := utils.ReadYesNo(reader)

			if !yes {
				fmt.Println("Aborted!")
				if err != nil {
					libutils.PrintlnStdErr(err)
				}
				os.Exit(1)
			}
		}

		fmt.Println("Executing...")

		ec := utils.LaunchAppWithDirectStd("/bin/bash", []string{"-c", joinedCommand}, nil)
		if ec != 0 {
			os.Exit(ec)
		}
	},
}

func registerStartupPredefinedAliases(_ *types.UserInfo) {
	registerPredefinedAlias(
		"delete-git-tag [tag]",
		[]string{"git", "tag", "-d", "$TAG", "&&", "git", "push", "-d", "origin", "$TAG"},
		func(_, args []string) []string {
			if len(args) != 1 {
				libutils.PrintlnStdErr("ERR: this alias requires exactly one argument is the git tag you want to delete")
				os.Exit(1)
			}
			tag := args[0]
			if !regexp.MustCompile("^[a-z\\d]*[a-z\\d.-_]+").MatchString(tag) {
				libutils.PrintlnStdErr("ERR: invalid git tag format:", tag)
				os.Exit(1)
			}
			return []string{
				"git", "tag", "-d", tag,
				"&&", "git", "push", "-d", "origin", tag,
			}
		},
		true,
	)

	if utils.IsDarwin() {
		registerPredefinedAlias(
			"awake [time]",
			[]string{"caffeinate", "-d", "-t", "$TIME"},
			func(_, args []string) []string {
				if len(args) != 1 {
					libutils.PrintlnStdErr("ERR: this alias requires exactly one argument is the amount of time you want your Mac to stay awake")
					os.Exit(1)
				}
				duration, err := time.ParseDuration(args[0])
				if err != nil {
					libutils.PrintlnStdErr("ERR: invalid duration format:", args[0])
					os.Exit(1)
				}

				seconds := int(duration.Seconds())
				return []string{
					"caffeinate", "-d", "-t", fmt.Sprintf("%d", seconds),
				}
			},
			true,
		)
	}
}

func registerPredefinedAliasesFromFile(userInfo *types.UserInfo) {
	aliasFile := path.Join(userInfo.HomeDir, constants.PREDEFINED_ALIAS_FILE_NAME)

	file, errFile := os.Stat(aliasFile)
	if errFile != nil {
		if os.IsNotExist(errFile) {
			return
		}
		fmt.Printf("ERR: unable to check alias file %s: %s\n", aliasFile, errFile.Error())
		return
	}

	if file.IsDir() {
		return
	}

	bz, err := os.ReadFile(aliasFile)
	if err != nil {
		fmt.Printf("ERR: failed to read alias file %s: %s\n", aliasFile, err.Error())
		return
	}

	fmt.Println("Loading aliases from", aliasFile, "...")

	tsvLines := goe.NewIEnumerable(strings.Split(string(bz), "\n")...).Select(func(line string) any {
		return strings.TrimSpace(line)
	}).CastString()

	regexReplaceContinousSpace := regexp.MustCompile("[\\s\\t]+")

	for _, line := range tsvLines.ToArray() {
		if strings.HasPrefix(line, "#") {
			continue
		}
		if libutils.IsBlank(line) {
			continue
		}

		spl := strings.Split(
			regexReplaceContinousSpace.ReplaceAllString(strings.Replace(line, "\t", " ", -1), " "),
			" ",
		)

		if len(spl) < 2 {
			panic(fmt.Errorf("malformed %s", constants.PREDEFINED_ALIAS_FILE_NAME))
		}

		alias := spl[0]
		command := spl[1:]
		if pa, found := predefinedAliases[alias]; found {
			pa.command = command
			pa.use = alias
			pa.overridden = true
			predefinedAliases[alias] = pa
		} else {
			registerPredefinedAlias(alias, command, nil, false)
		}
	}
}

func registerPredefinedAlias(use string, defaultCommand []string, alter commandAlter, alwaysInvokeCommandAlter bool) {
	spl := strings.Split(use, " ")
	alias := spl[0]
	predefinedAliases[alias] = predefinedAlias{
		alias:                    alias,
		use:                      use,
		command:                  defaultCommand,
		alter:                    alter,
		alwaysInvokeCommandAlter: alter != nil && alwaysInvokeCommandAlter,
	}
	longestUseDesc = libutils.MaxInt(longestUseDesc, len(use))
}

func init() {
	aliasCmd.Flags().BoolP(
		flagConfirmExecution,
		"y",
		false,
		"skip confirmation before executing the command, but wait few seconds before executing",
	)

	rootCmd.AddCommand(aliasCmd)
}

type predefinedAlias struct {
	alias                    string
	use                      string
	command                  []string
	alter                    commandAlter
	alwaysInvokeCommandAlter bool
	overridden               bool
}

type commandAlter func(command, args []string) []string
