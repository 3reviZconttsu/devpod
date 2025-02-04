package cmd

import (
	"context"
	"fmt"

	"github.com/loft-sh/devpod/cmd/flags"
	client2 "github.com/loft-sh/devpod/pkg/client"
	"github.com/loft-sh/devpod/pkg/config"
	workspace2 "github.com/loft-sh/devpod/pkg/workspace"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

// StopCmd holds the destroy cmd flags
type StopCmd struct {
	*flags.GlobalFlags
}

// NewStopCmd creates a new destroy command
func NewStopCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &StopCmd{
		GlobalFlags: flags,
	}
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops an existing workspace",
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()
			devPodConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
			if err != nil {
				return err
			}

			client, err := workspace2.GetWorkspace(devPodConfig, args, false, log.Default)
			if err != nil {
				return err
			}

			return cmd.Run(ctx, client)
		},
	}

	return stopCmd
}

// Run runs the command logic
func (cmd *StopCmd) Run(ctx context.Context, client client2.BaseWorkspaceClient) error {
	// lock workspace
	client.Lock()
	defer client.Unlock()

	// get instance status
	instanceStatus, err := client.Status(ctx, client2.StatusOptions{})
	if err != nil {
		return err
	} else if instanceStatus != client2.StatusRunning {
		return fmt.Errorf("cannot stop workspace because it is '%s'", instanceStatus)
	}

	// stop environment
	err = client.Stop(ctx, client2.StopOptions{})
	if err != nil {
		return err
	}

	return nil
}
