package cmd

import (
	"github.com/agent-kit/q/cmd/apikeys"
	"github.com/agent-kit/q/cmd/compile"
	"github.com/agent-kit/q/cmd/demo"
	"github.com/agent-kit/q/cmd/doctor"
	"github.com/agent-kit/q/cmd/inspect"
	"github.com/agent-kit/q/cmd/logging"
	"github.com/agent-kit/q/cmd/monitoring"
	"github.com/agent-kit/q/cmd/servicecontrol"
	"github.com/agent-kit/q/cmd/servicemanagement"
	"github.com/agent-kit/q/cmd/serviceusage"
	"github.com/agent-kit/q/cmd/translate"
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
