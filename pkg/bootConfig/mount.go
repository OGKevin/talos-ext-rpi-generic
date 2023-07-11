package bootconfig

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/siderolabs/go-blockdevice/blockdevice/probe"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"golang.org/x/exp/slog"
)

const (
	bootMountDir = "/mnt/boot"
)

func MountBootPartition(ctx context.Context) (*mount.MountPoint, error) {
	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		if err := logCurrentlyMountedPartitions(ctx); err != nil {
			return nil, errors.Wrapf(err, "failed to list currenlty mounted partitions.")
		}
	}

	bootPartitionDevice, err := probe.GetDevWithFileSystemLabel("BOOT")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get device with BOOT partition")
	}

	d, err := block.Device(bootPartitionDevice.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not init block device")
	}

	slog.DebugCtx(
		ctx, "mounting partition", slog.String("disk", bootPartitionDevice.Path), slog.String("mount_dir", bootMountDir),
	)
	if err := ensureDirectory(ctx, bootMountDir); err != nil {
		return nil, errors.Wrapf(err, "failed to ensure mount dir exists.")
	}

	m, err := d.Mount(bootMountDir, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "could mount %s", d.DevName())
	}

	return m, nil
}

func logCurrentlyMountedPartitions(ctx context.Context) error {
	partitions, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		return errors.Wrapf(err, "could not list partitions.")
	}

	for _, parition := range partitions {
		slog.DebugCtx(
			ctx, "found mounted partition", slog.String("device", parition.Device),
			slog.String("mount_point", parition.Mountpoint),
		)
	}

	return nil
}

func ensureDirectory(ctx context.Context, target string) (err error) {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		slog.DebugCtx(ctx, "ensuring dir exists", slog.String("dir", target))
		if err = os.MkdirAll(target, 0o755); err != nil {
			return errors.Wrapf(err, "error creating mount point dir for %s", target)
		}
	}

	return nil
}
