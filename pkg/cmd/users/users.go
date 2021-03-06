package users

import (
	"github.com/spf13/cobra"
	cmdUserChange "github.com/wizedkyle/sumocli/pkg/cmd/users/change"
	cmdUserCreate "github.com/wizedkyle/sumocli/pkg/cmd/users/create"
	cmdUserDelete "github.com/wizedkyle/sumocli/pkg/cmd/users/delete"
	cmdUserDisable "github.com/wizedkyle/sumocli/pkg/cmd/users/disable"
	cmdUserGet "github.com/wizedkyle/sumocli/pkg/cmd/users/get"
	cmdUserList "github.com/wizedkyle/sumocli/pkg/cmd/users/list"
	cmduserReset "github.com/wizedkyle/sumocli/pkg/cmd/users/reset"
	cmdUserUnlock "github.com/wizedkyle/sumocli/pkg/cmd/users/unlock"
	cmdUserUpdate "github.com/wizedkyle/sumocli/pkg/cmd/users/update"
)

func NewCmdUser() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users <command>",
		Short: "Manage users",
		Long:  "Work with Sumo Logic users",
	}

	cmd.AddCommand(cmdUserChange.NewCmdUserChangeEmail())
	cmd.AddCommand(cmdUserCreate.NewCmdUserCreate())
	cmd.AddCommand(cmdUserDelete.NewCmdUserDelete())
	cmd.AddCommand(cmdUserDisable.NewCmdUserDisableMFA())
	cmd.AddCommand(cmdUserGet.NewCmdGetUser())
	cmd.AddCommand(cmdUserList.NewCmdUserList())
	cmd.AddCommand(cmduserReset.NewCmdUserResetPassword())
	cmd.AddCommand(cmdUserUnlock.NewCmdUnlockUser())
	cmd.AddCommand(cmdUserUpdate.NewCmdUserUpdate())
	return cmd
}
