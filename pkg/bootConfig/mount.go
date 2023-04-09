package bootconfig

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/sirupsen/logrus"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
)

const (
	bootMountDir  = "/mnt/boot"
	bootPartition = "/dev/mmcblk0p1"
)

func MountBootPartition(ctx context.Context) (*mount.MountPoint, error) {
	log := logrus.WithContext(ctx)

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		if err := logCurrentlyMountedPartitions(ctx); err != nil {
			return nil, errors.Wrapf(err, "failed to list currenlty mounted partitions.")
		}
	}

	d, err := block.Device(bootPartition)
	if err != nil {
		return nil, errors.Wrapf(err, "could not init block device")
	}

	log.Debugf("mounting %s to %s", bootPartition, bootMountDir)
	if err := ensureDirectory(bootMountDir); err != nil {
		return nil, errors.Wrapf(err, "failed to ensure mount dir exists.")
	}

	m, err := d.Mount(bootMountDir, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "could mount %s", d.DevName())
	}

	return m, nil
}

func logCurrentlyMountedPartitions(ctx context.Context) error {
	log := logrus.WithContext(ctx)

	partitions, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		return errors.Wrapf(err, "could not list partitions.")
	}

	for _, parition := range partitions {
		log.Debugf("found parition %s mounted at %s", parition.Device, parition.Mountpoint)
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
