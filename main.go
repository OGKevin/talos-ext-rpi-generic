package main

import (
	"context"
	"fmt"
	"os"

	_ "embed"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/sirupsen/logrus"
	"github.com/u-root/u-root/pkg/mount/block"
)

//go:embed config.txt
var defaultBootConfig string

func main() {
	ctx := context.Background()

	log := logrus.StandardLogger()
	log.SetLevel(logrus.DebugLevel)

	log.Print("running")

	partitions, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		err = errors.Wrapf(err, "could not list partitions.")
		log.WithError(err).Fatal(err.Error())
	}

	for _, parition := range partitions {
		log.Debugf("found parition %s mounted at %s", parition.Device, parition.Mountpoint)
	}

	if err := listDir("/dev"); err != nil {
		log.WithError(err).Fatal()
	}

	const (
		bootMountDir  = "/mnt/boot"
		bootPartition = "/dev/mmcblk0p1"
	)

	d, err := block.Device(bootPartition)
	if err != nil {
		err = errors.Wrapf(err, "could not init block device")
		log.WithError(err).Fatal(err.Error())
	}

	log.Debugf("mounting %s to %s", bootPartition, bootMountDir)
	if err := ensureDirectory(bootMountDir); err != nil {
		log.WithError(err).Fatal(err.Error())
	}

	m, err := d.Mount(bootMountDir, 0)
	if err != nil {
		err = errors.Wrapf(err, "could mount %s", d.DevName())
		log.WithError(err).Fatal(err.Error())
	}

	defer func() {
		if err := m.Unmount(0); err != nil {
			log.WithError(err).Warn("failed to umnount disk")
		}
	}()

	partitions, err = disk.PartitionsWithContext(ctx, true)
	if err != nil {
		err = errors.Wrapf(err, "could not list partitions.")
		log.WithError(err).Fatal(err.Error())
	}

	for _, parition := range partitions {
		log.Debugf("found parition %s mounted at %s", parition.Device, parition.Mountpoint)
	}

	if err := listDir(bootMountDir); err != nil {
		log.WithError(err).Fatal()
	}

	bootConfigFilePath := fmt.Sprintf("%s/config.txt", bootMountDir)

	oldBootConfig, err := os.ReadFile(bootConfigFilePath)
	if err != nil {
		err = errors.Wrapf(err, "failed to open current boot config.")
		log.WithError(err).Fatal()
	}

	log.Debug("dumping old config, and afterwards the new one")
	if log.IsLevelEnabled(logrus.DebugLevel) {
		fmt.Println("### OLD CONFIG ###")
		fmt.Println(string(oldBootConfig))

		fmt.Println("### NEW CONFIG TO BE WRITTEN ###")
		fmt.Println(string(defaultBootConfig))
	}

	if err := os.WriteFile(bootConfigFilePath, []byte(defaultBootConfig), 0600); err != nil {
		err = errors.Wrapf(err, "failed to write boot config")
		log.WithError(err).Fatal()
	}

	if log.IsLevelEnabled(logrus.DebugLevel) {
		oldBootConfig, err := os.ReadFile(bootConfigFilePath)
		if err != nil {
			err = errors.Wrapf(err, "failed to open current boot config.")
			log.WithError(err).Fatal()
		}

		fmt.Println("### NEWLY WRITTEN CONFIG ###")
		fmt.Println(string(oldBootConfig))
	}

	log.Info("boot config written")
}

func listDir(path string) error {
	log := logrus.StandardLogger()

	efiDirs, err := os.ReadDir(path)
	if err != nil {
		err = errors.Wrapf(err, "failed to list dirs in %s", path)

		return err
	}

	log.Debugf("listing dirs for %s", path)
	for _, dir := range efiDirs {
		log.Debugf("listing dir %s/%s", path, dir.Name())
	}

	return nil
}

func ensureDirectory(target string) (err error) {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		logrus.Debugf("ensuring dir %s exists", target)
		if err = os.MkdirAll(target, 0o755); err != nil {
			return errors.Wrapf(err, "error creating mount point dir for %s", target)
		}
	}

	return nil
}
