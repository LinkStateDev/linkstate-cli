package cmd

import (
	"github.com/LinkStateDev/linkstate-cli/internal/lab"
	"github.com/spf13/cobra"
)

var labCmd = &cobra.Command{
	Use:   "lab",
	Short: "Manage the spine-leaf BGP lab",
}

func labPreRun() {
	if cfg != nil && cfg.Path != "" {
		lab.SetLabDir(cfg.Path)
	}
}

var labUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Create 2 spines + 4 leaves, eBGP, SSH",
	RunE: func(cmd *cobra.Command, args []string) error {
		labPreRun()
		return lab.Up()
	},
}

var labDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Destroy the lab",
	RunE: func(cmd *cobra.Command, args []string) error {
		labPreRun()
		lab.Down()
		return nil
	},
}

var labStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "List devices, loopbacks, ports, peers",
	RunE: func(cmd *cobra.Command, args []string) error {
		lab.Status()
		return nil
	},
}

var labConnectCmd = &cobra.Command{
	Use:   "connect <name>",
	Short: "Attach vtysh console to a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		labPreRun()
		return lab.Connect(args[0])
	},
}

var labSSHCmd = &cobra.Command{
	Use:   "ssh <name>",
	Short: "SSH to a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		labPreRun()
		return lab.SSH(args[0])
	},
}

var labVethLinkCmd = &cobra.Command{
	Use:    "veth-link <idx> <ctr1> <ctr2> <if1> <if2>",
	Short:  "Create veth pair (internal)",
	Hidden: true,
	Args:   cobra.ExactArgs(5),
	RunE: func(cmd *cobra.Command, args []string) error {
		return lab.VethLink(args)
	},
}

func init() {
	labCmd.AddCommand(labUpCmd)
	labCmd.AddCommand(labDownCmd)
	labCmd.AddCommand(labStatusCmd)
	labCmd.AddCommand(labConnectCmd)
	labCmd.AddCommand(labSSHCmd)
	labCmd.AddCommand(labVethLinkCmd)
	rootCmd.AddCommand(labCmd)
}
