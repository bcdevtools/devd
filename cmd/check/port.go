package check

import (
	"fmt"
	"github.com/bcdevtools/devd/v2/cmd/utils"
	psnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"os"
	"strconv"
)

func GetCheckPortCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "port [?port]",
		Short: `List ports
or check specific port currently open and holding by a process.`,
		Args: cobra.RangeArgs(0, 1),
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				listPorts()
			} else {
				checkPort(args[0])
			}
		},
	}

	return cmd
}

func listPorts() {
	connections, err := psnet.Connections("all")
	utils.ExitOnErr(err, "failed to get connections")

	slices.SortFunc(connections, func(l, r psnet.ConnectionStat) bool {
		return l.Laddr.Port < r.Laddr.Port
	})

	for _, conn := range connections {
		fmt.Printf("%5d", conn.Laddr.Port)
		fmt.Printf(" | %-11s", conn.Status)
		fmt.Printf(" | PID %-5d", conn.Pid)

		processName := "(ERR)"
		proc, err := process.NewProcess(conn.Pid)
		if err == nil && proc != nil {
			name, err := proc.Name()
			if err == nil {
				processName = name
			}
		}
		fmt.Printf(" | PN %-30s", processName)

		remoteIP := "-"
		remotePort := "-"
		if conn.Raddr.Port != 0 || conn.Raddr.IP != "" {
			remoteIP = conn.Raddr.IP
			remotePort = fmt.Sprintf("%d", conn.Raddr.Port)
		}
		fmt.Printf(" | REMOTE %30s:%-5s", remoteIP, remotePort)

		fmt.Println()
	}
}

func checkPort(portStr string) {
	port64, err := strconv.ParseInt(portStr, 10, 64)
	utils.ExitOnErr(err, "failed to read, port is not a number")

	if port64 < 0 || port64 > 65535 {
		utils.PrintlnStdErr("ERR: port must be in range [1, 65535]")
		os.Exit(1)
	}

	var isOpen, anyErr bool

	port32 := uint32(port64)

	connections, err := psnet.Connections("all")
	utils.ExitOnErr(err, "failed to get connections")

	for _, conn := range connections {
		if conn.Laddr.Port != port32 {
			continue
		}

		isOpen = true

		fmt.Println(conn.Status)
		fmt.Println("PID:", conn.Pid)

		if conn.Pid > 0 {
			proc, err := process.NewProcess(conn.Pid)
			if err != nil || proc == nil {
				utils.PrintlnStdErr("ERR: Failed to get process or process not found")
				anyErr = true
			} else {
				name, err := proc.Name()
				if err == nil {
					fmt.Println("PROC:", name)
				} else {
					utils.PrintlnStdErr("ERR: Failed to get process name")
					anyErr = true
				}
				cmdLine, err := proc.Cmdline()
				if err == nil {
					fmt.Println("PROC CLI:", cmdLine)
				} else {
					utils.PrintlnStdErr("ERR: Failed to get process command line")
					anyErr = true
				}
			}
		}

		fmt.Println("LOCAL IP:", conn.Laddr.IP)
		fmt.Print("REMOTE: ")
		if conn.Raddr.Port == 0 && conn.Raddr.IP == "" {
			fmt.Println("NONE")
		} else {
			fmt.Println(conn.Raddr.IP, " ", conn.Raddr.Port)
		}
		fmt.Println("FD:", conn.Fd)
		fmt.Println("FAMILY:", conn.Family)
		fmt.Println("TYPE:", conn.Type)
		fmt.Println("UIDs:", conn.Uids)

		break
	}

	if !isOpen {
		fmt.Println("Not open")
	}

	if anyErr {
		os.Exit(1)
	}
}
