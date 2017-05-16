package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/spf13/cobra"
)

var pidFlag int

var targetsCmd = &cobra.Command{
	Use:   "targets",
	Short: "Get the current targets registered with hoverctl",
	Long: `
Get the current targets registered with hoverctl
`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(config.Targets) == 0 {
			handleIfError(errors.New("No targets registered"))
		}

		data := [][]string{
			{"Target name", "Pid", "Host", "Admin port", "Proxy port", "Default"},
		}

		for key, target := range config.Targets {
			defaultMarker := ""
			if target.Name == config.DefaultTarget {
				defaultMarker = "X"
			}

			data = append(data, []string{key, strconv.Itoa(target.Pid), target.Host, strconv.Itoa(target.AdminPort), strconv.Itoa(target.ProxyPort), defaultMarker})
		}

		drawTable(data, true)
	},
}

var targetsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete target",
	Long: `
Delete target"
`,

	Run: func(cmd *cobra.Command, args []string) {
		checkArgAndExit(args, "Cannot delete a target without a name", "targets delete")

		if !askForConfirmation("Are you sure you want to delete the target " + args[0] + "?") {
			return
		}
		config.DeleteTarget(configuration.Target{
			Name: args[0],
		})

		handleIfError(config.WriteToFile(hoverflyDirectory))

		targetsCmd.Run(cmd, args)
	},
}

var targetsNewCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new target",
	Long: `
Create target"
`,

	Run: func(cmd *cobra.Command, args []string) {
		checkArgAndExit(args, "Cannot create a target without a name", "targets new")

		if config.GetTarget(args[0]) != nil {
			handleIfError(fmt.Errorf("Target %s already exists\n\nUse a different target name or run `hoverctl targets update %[1]s`", args[0]))
		}

		newTarget := configuration.NewTarget(args[0], hostFlag, adminPortFlag, proxyPortFlag)
		newTarget.Pid = pidFlag

		config.NewTarget(*newTarget)

		handleIfError(config.WriteToFile(hoverflyDirectory))

		targetsCmd.Run(cmd, args)
	},
}

var targetsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update target",
	Long: `
Update target
`,

	Run: func(cmd *cobra.Command, args []string) {
		checkArgAndExit(args, "Cannot update a target without a name", "targets update")

		if config.GetTarget(args[0]) == nil {
			handleIfError(fmt.Errorf("Target %s does not exist\n\nUse a different target name or run `hoverctl targets create %[1]s`", args[0]))
		}

		newTarget := configuration.NewTarget(args[0], hostFlag, adminPortFlag, proxyPortFlag)
		newTarget.Pid = pidFlag

		config.NewTarget(*newTarget)

		handleIfError(config.WriteToFile(hoverflyDirectory))

		targetsCmd.Run(cmd, args)
	},
}

var targetsDefaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Get and set the default target",
	Long: `
Without a target name, default will print the configuration
of the current default target."
`,

	Run: func(cmd *cobra.Command, args []string) {
		if config.GetTarget("") == nil {
			handleIfError(errors.New("No targets registered"))
		}

		if len(args) > 0 {
			checkTarget := config.GetTarget(args[0])
			if checkTarget == nil {
				handleIfError(fmt.Errorf("%[1]s is not a target\n\nRun `hoverctl targets new %[1]s`", args[0]))
			}
			config.DefaultTarget = args[0]
		}

		data := [][]string{
			[]string{"Target name", "Pid", "Host", "Admin port", "Proxy port"},
		}

		defaultTarget := config.GetTarget("")
		data = append(data, []string{defaultTarget.Name, strconv.Itoa(defaultTarget.Pid), defaultTarget.Host, strconv.Itoa(defaultTarget.AdminPort), strconv.Itoa(defaultTarget.ProxyPort)})

		drawTable(data, true)
	},
}

func init() {
	RootCmd.AddCommand(targetsCmd)

	targetsCmd.AddCommand(targetsDeleteCmd)
	targetsCmd.AddCommand(targetsNewCmd)
	targetsCmd.AddCommand(targetsUpdateCmd)
	targetsCmd.AddCommand(targetsDefaultCmd)

	targetsNewCmd.Flags().IntVar(&pidFlag, "pid", 0, "Process id for a running instance of Hoverfly")
	targetsUpdateCmd.Flags().IntVar(&pidFlag, "pid", 0, "Process id for a running instance of Hoverfly")
}
