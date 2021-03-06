package collectors

import (
	"github.com/spf13/cobra"
	cmdCollectorCreate "github.com/wizedkyle/sumocli/pkg/cmd/collectors/create"
	cmdCollectorDelete "github.com/wizedkyle/sumocli/pkg/cmd/collectors/delete"
	cmdCollectorGet "github.com/wizedkyle/sumocli/pkg/cmd/collectors/get"
	cmdCollectorList "github.com/wizedkyle/sumocli/pkg/cmd/collectors/list"
	cmdCollectorUpdate "github.com/wizedkyle/sumocli/pkg/cmd/collectors/update"
	cmdCollectorUpgrade "github.com/wizedkyle/sumocli/pkg/cmd/collectors/upgrade"
)

func NewCmdCollectors() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collectors <command>",
		Short: "Manages collectors",
	}

	cmd.AddCommand(cmdCollectorCreate.NewCmdCollectorCreate())
	cmd.AddCommand(cmdCollectorDelete.NewCmdCollectorDelete())
	cmd.AddCommand(cmdCollectorGet.NewCmdCollectorGet())
	cmd.AddCommand(cmdCollectorList.NewCmdCollectorList())
	cmd.AddCommand(cmdCollectorUpdate.NewCmdCollectorUpdate())
	cmd.AddCommand(cmdCollectorUpgrade.NewCmdUpgradeCollectors())
	return cmd
}
