package commands

import "github.com/abiosoft/ishell"
import gmwi "github.com/ski7777/gomultiwa/internal/gomultiwa/interface"

// GetCmdExit returns the command to exit gomultiwa
func GetCmdExit(gmw gmwi.GoMultiWAInterface) *ishell.Cmd {
	return &ishell.Cmd{
		Name: "exit",
		Help: "Exit gomultiwa",
		Func: func(c *ishell.Context) {
			gmw.Stop()
		},
	}
}
