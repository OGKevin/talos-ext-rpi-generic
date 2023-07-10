package bootconfig

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

//go:embed config.txt
var defaultBootConfig []byte

func LoadBootConfig(ctx context.Context, path string, allowDefault bool) ([]byte, error) {
	slog.DebugCtx(ctx, "attempting to read config", slog.String("path", path))

	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) && !allowDefault {
		return nil, errors.Wrapf(err, "failed to open %s", path)
	}

	if raw == nil {
		slog.DebugCtx(ctx, "failed to read config from filesystem, using default config.txt")
		raw = defaultBootConfig
	}

	return raw, nil
}

func ReplaceBootConfig(ctx context.Context, newConfig []byte) error {
	bootConfigFilePath := fmt.Sprintf("%s/config.txt", bootMountDir)

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		dumpConfig(ctx, bootConfigFilePath, newConfig)
	}

	slog.InfoCtx(ctx, "writing config.txt")

	if err := os.WriteFile(bootConfigFilePath, []byte(newConfig), 0600); err != nil {
		return errors.Wrapf(err, "failed to write boot config")
	}

	return nil
}

func dumpConfig(ctx context.Context, path string, newConfig []byte) {
	oldBootConfig, err := os.ReadFile(path)
	if err != nil {
		err = errors.Wrapf(err, "failed to open current boot config.")
		slog.ErrorCtx(ctx, "failed to open current boot config", slog.Any("error", err))

		return
	}

	slog.DebugCtx(ctx, "dumping old config, and afterwards the new one")
	fmt.Println("### OLD CONFIG ###")
	fmt.Println(string(oldBootConfig))

	fmt.Println("### NEW CONFIG TO BE WRITTEN ###")
	fmt.Println(string(newConfig))
}
