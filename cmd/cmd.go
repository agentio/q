package cmd

import (
	"github.com/agentio/q/cmd/apikeys"
	"github.com/agentio/q/cmd/compile"
	"github.com/agentio/q/cmd/demo"
	"github.com/agentio/q/cmd/doctor"
	"github.com/agentio/q/cmd/inspect"
	"github.com/agentio/q/cmd/logging"
	"github.com/agentio/q/cmd/monitoring"
	"github.com/agentio/q/cmd/servicecontrol"
	"github.com/agentio/q/cmd/servicemanagement"
	"github.com/agentio/q/cmd/serviceusage"
	"github.com/agentio/q/cmd/translate"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "q",
		Short: "Manage APIs with Service Infrastructure",
	}
	cmd.AddCommand(apikeys.Cmd())
	cmd.AddCommand(inspect.Cmd())
	cmd.AddCommand(compile.Cmd())
	cmd.AddCommand(demo.Cmd())
	cmd.AddCommand(doctor.Cmd())
	cmd.AddCommand(logging.Cmd())
	cmd.AddCommand(monitoring.Cmd())
	cmd.AddCommand(servicecontrol.Cmd())
	cmd.AddCommand(servicemanagement.Cmd())
	cmd.AddCommand(serviceusage.Cmd())
	cmd.AddCommand(translate.Cmd())
	return cmd
}
