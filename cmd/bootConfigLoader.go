/*
Copyright © 2023 Kevin Hellemun
*/

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	bootconfig "github.com/OGKevin/talos-ext-rpi/pkg/bootConfig"
)

const (
	flagBootConfigPath  = "boot-config-path"
	viperBootConfigPath = "boot.config.path"

	flagAllowDefaultBootConfig  = "boot.allow-default-config"
	viperAllowDefaultBootConfig = "boot.allowDefaultConfig"
)

// init register the command and it's flags.
func init() {
	rootCmd.AddCommand(bootConfigLoaderCmd)

	bootConfigLoaderCmd.Flags().
		String(
			flagBootConfigPath,
			"/var/etc/rpi-boot-config-loader/config.txt",
			"The config.txt file containing the rpi boot config.",
		)
	bootConfigLoaderCmd.Flags().
		Bool(
			flagAllowDefaultBootConfig,
			false,
			"If there is a fauilrue to read the config profided by boot-config, use the default one that this app proveds.",
		)

	_ = viper.BindPFlag(viperBootConfigPath, bootConfigLoaderCmd.Flags().Lookup(flagBootConfigPath))
	_ = viper.BindPFlag(
		viperAllowDefaultBootConfig,
		bootConfigLoaderCmd.Flags().Lookup(flagAllowDefaultBootConfig),
	)
}

// bootConfigLoaderCmd represents the bootConfigLoader command
var bootConfigLoaderCmd = &cobra.Command{
	Use:                        "bootConfigLoader",
	Short:                      "Replaces rpi's boot config.",
	Long:                       `This command can be used to replace https://www.raspberrypi.com/documentation/computers/config_txt.html on Talos Linux for rpi.`,
	RunE:                       runBootConfigLoader,
	SilenceUsage:               true,
	SuggestionsMinimumDistance: 0,
}

// runBootConfigLoader performs the action for this command.
func runBootConfigLoader(c *cobra.Command, args []string) error {
	ctx := c.Context()

	slog.InfoCtx(ctx, "Going to replace boot config...")

	configRaw, err := bootconfig.LoadBootConfig(
		ctx,
		viper.GetString(viperBootConfigPath),
		viper.GetBool(viperAllowDefaultBootConfig),
	)
	if err != nil {
		return errors.Wrapf(err, "failed to load config.txt")
	}

	m, err := bootconfig.MountBootPartition(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to mount boot partition.")
	}

	defer func() {
		if err := m.Unmount(0); err != nil {
			slog.WarnCtx(ctx, "failed to unmount boot partition", slog.Any("error", err))
		}
	}()

	if err := bootconfig.ReplaceBootConfig(ctx, configRaw); err != nil {
		return errors.Wrapf(err, "failed to replace existing boot config")
	}

	slog.InfoCtx(ctx, "config.txt replaced!")
	slog.InfoCtx(ctx, "Don't forget to reboot for changes to take effect.")

	return nil
}
