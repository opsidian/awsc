package command

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/opsidian/awsc/awsc/autoscaling"
	"github.com/spf13/cobra"
)

var (
	ecsCluster  string
	maxInFlight int
)

var autoScalingCmd = &cobra.Command{
	Use:   "autoscaling command <params>",
	Short: "AWS auto scaling commands",
}

var migrateCmd = &cobra.Command{
	Use:        "migrate <auto scaling group name>",
	Short:      "Migrate",
	ArgAliases: []string{"name"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("auto scaling group name is missing")
		}
		if len(args) > 1 {
			return errors.New("too many arguments")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		config := &aws.Config{}
		if Region != "" {
			config.Region = aws.String(Region)
		}
		migrateService := autoscaling.NewMigrateService(config, cmd.OutOrStdout())
		return migrateService.MigrateInstances(args[0], ecsCluster, maxInFlight)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	migrateCmd.PersistentFlags().StringVarP(&ecsCluster, "ecs-cluster", "", "", "If any instance is part of an ECS cluster it will be drained first")
	migrateCmd.PersistentFlags().IntVarP(&maxInFlight, "min-healthy-percent", "m", 50, "Minimum percent of instances to keep healthy during the migration")
	autoScalingCmd.AddCommand(migrateCmd)
	RootCmd.AddCommand(autoScalingCmd)
}
