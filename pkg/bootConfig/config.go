package bootconfig

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:embed config.txt
var defaultBootConfig []byte

func LoadBootConfig(ctx context.Context, path string, allowDefault bool) ([]byte, error) {
	log := logrus.WithContext(ctx)
	log.Debugf("attempting to read config from %s", path)

	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) && !allowDefault {
		return nil, errors.Wrapf(err, "failed to open %s", path)
	}

	if raw == nil {
		log.Debug("failed to read config from filesystem, using default config.txt")
		raw = defaultBootConfig
	}

	return raw, nil
}

func ReplaceBootConfig(ctx context.Context, newConfig []byte) error {
	log := logrus.WithContext(ctx)
	bootConfigFilePath := fmt.Sprintf("%s/config.txt", bootMountDir)

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		dumpConfig(ctx, bootConfigFilePath, newConfig)
	}

	log.Info("Writing config.txt")

	if err := os.WriteFile(bootConfigFilePath, []byte(newConfig), 0600); err != nil {
		return errors.Wrapf(err, "failed to write boot config")
	}

	return nil
}

func dumpConfig(ctx context.Context, path string, newConfig []byte) {
	log := logrus.WithContext(ctx)
	oldBootConfig, err := os.ReadFile(path)
	if err != nil {
		err = errors.Wrapf(err, "failed to open current boot config.")
		log.WithError(err).Fatal()
	}

	log.Debug("dumping old config, and afterwards the new one")
	fmt.Println("### OLD CONFIG ###")
	fmt.Println(string(oldBootConfig))

	fmt.Println("### NEW CONFIG TO BE WRITTEN ###")
	fmt.Println(string(newConfig))
}
