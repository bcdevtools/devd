package debug

import (
	"fmt"
	"github.com/bcdevtools/devd/cmd/types"
	"github.com/spf13/cobra"
)

func GetUserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Aliases: []string{"u"},
		Short:   "Get current user info",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := types.UnwrapAppContext(cmd.Context())

			operationUserInfo := ctx.GetOperationUserInfo()

			if operationUserInfo.IsSameUser {
				fmt.Println("User info:")
				fmt.Println("- Username:", operationUserInfo.EffectiveUserInfo.Username)
				fmt.Println("- User ID:", operationUserInfo.EffectiveUserInfo.UserId)
				fmt.Println("- Home Dir:", operationUserInfo.EffectiveUserInfo.HomeDir)
			} else {
				fmt.Println("Effective user info:")
				fmt.Println("- Username:", operationUserInfo.EffectiveUserInfo.Username)
				fmt.Println("- User ID:", operationUserInfo.EffectiveUserInfo.UserId)
				fmt.Println("- Home Dir:", operationUserInfo.EffectiveUserInfo.HomeDir)
				fmt.Println("- Super user:", operationUserInfo.EffectiveUserInfo.IsSuperUser)
				fmt.Println("Real user info:")
				fmt.Println("- Username:", operationUserInfo.RealUserInfo.Username)
				fmt.Println("- User ID:", operationUserInfo.RealUserInfo.UserId)
				fmt.Println("- Home Dir:", operationUserInfo.RealUserInfo.HomeDir)
				fmt.Println("- Super user:", operationUserInfo.RealUserInfo.IsSuperUser)
			}

			workingUserInfo := ctx.GetWorkingUserInfo()
			fmt.Println("Working user info:")
			fmt.Println("- Username:", workingUserInfo.Username)
			fmt.Println("- User ID:", workingUserInfo.UserId)
			fmt.Println("- Home Dir:", workingUserInfo.HomeDir)
			fmt.Println("- Super user:", workingUserInfo.IsSuperUser)

			fmt.Println("Super user:", operationUserInfo.IsSuperUser)
			fmt.Println("Operating as super user:", operationUserInfo.OperatingAsSuperUser)
		},
	}

	return cmd
}
