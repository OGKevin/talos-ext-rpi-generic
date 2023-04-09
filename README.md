# Talos Generic Raspberrypi Extension
A talos system extension for generic rpi operations.

[TOC]

## Boot config.txt Loader

It is possible overwrite the default config.txt that Talos Linux ships with for rpi.
This is done by mounting the boot partition, replacing the file, and afterwards trigger a reboot.

A word of caution, at the time of writing, there is no build in validation that the provided config.txt file is valid.
This means, that if the file is invalid, the Pi will not boot.
You must re-install the os or manually mount the disk and fix the invalid config.txt file.

With the warning out of the way, here are the instructions on how to use the boot config reloader.

### Install The Config Realoader

#### 1. Update Machine Config With Config.txt File

The default config.txt that Talos Linux ships with can be found [here]( https://github.com/siderolabs/talos/blob/9948a646d20f4ba80916a263ed7bca3e5ca2f0ad/internal/app/machined/pkg/runtime/v1alpha1/board/rpi_generic/config.txt).
Copy it and make the modifications that you need.
Once modified, update the machine config to add a [MachineFile](https://www.talos.dev/v1.3/reference/configuration/#machinefile):
```yaml
  files:
    - permissions: 0o400
      path: /var/etc/rpi-boot-config-loader/config.txt
      op: create
      content: |
        # See https://www.raspberrypi.com/documentation/computers/configuration.html
        # Reduce GPU memory to give more to CPU.
        gpu_mem=32
        # Enable maximum compatibility on both HDMI ports;
        # only the one closest to the power/USB-C port will work in practice.
        hdmi_safe:0=1
        hdmi_safe:1=1
        # Load U-Boot.
        kernel=u-boot.bin
        # Forces the kernel loading system to assume a 64-bit kernel.
        arm_64bit=1
        # Run as fast as firmware / board allows.
        arm_boost=1
        # Enable the primary/console UART.
        enable_uart=1
        # Disable Bluetooth.
        dtoverlay=disable-bt
        
        [pi4]
        # Run as fast as firmware / board allows
        arm_boost=1
        
        [all]
        dtparam=poe_fan_temp0=60000
        dtparam=poe_fan_temp1=70000
        dtparam=poe_fan_temp2=80000
        dtparam=poe_fan_temp3=85000
```

This extension expects the config.txt file to be located at: `/var/etc/rpi-boot-config-loader/config.txt`.
The extension is configured to not be started by Talos Linux until this file exists!

#### 2. Install The Extension

Update the [InstallExtensionConfig](https://www.talos.dev/v1.3/reference/configuration/#installextensionconfig) to include the extension:
```yaml
    extensions:
      - image: ghcr.io/ogkevin/rpi-boot-config-loader:<version>
```

Once this config has been applied with:
```shell
talosctl -n <node-ip> apply-config -f talos/rpi-worker-config.yaml 
```

Issue an upgrade command so that the extension gets installed:
```shell
talosctl upgrade \ 
-n <node-ip> \
--image ghcr.io/siderolabs/installer:<talos-version> \
--preserve \
--wait
```

You can follow the process by looking at dmesg:
```shell
talosctl dmesg -n <node-ip> -f
```

#### 3. Confirm Extension Ran Successfully

After the upgrade process succeeds, the node will reboot.
After this reboot, the extension will run for the first time.
Retrieve it's logs and status to check if everything went as planned:
```shell
❯ talosctl --nodes 192.168.1.150 cat /var/etc/rpi-boot-config-loader/config.txt
# See https://www.raspberrypi.com/documentation/computers/configuration.html
# Reduce GPU memory to give more to CPU.
gpu_mem=32
# Enable maximum compatibility on both HDMI ports;
# only the one closest to the power/USB-C port will work in practice.
hdmi_safe:0=1
hdmi_safe:1=1
# Load U-Boot.
kernel=u-boot.bin
# Forces the kernel loading system to assume a 64-bit kernel.
arm_64bit=1
# Run as fast as firmware / board allows.
arm_boost=1
# Enable the primary/console UART.
enable_uart=1
# Disable Bluetooth.
dtoverlay=disable-bt

[pi4]
# Run as fast as firmware / board allows
arm_boost=1

[all]
dtparam=poe_fan_temp0=60000
dtparam=poe_fan_temp1=70000
dtparam=poe_fan_temp2=80000
dtparam=poe_fan_temp3=85000
```
```shell
❯ talosctl service ext-rpi-boot-config-loader status -n <node-ip>
NODE     <node-ip>
ID       ext-rpi-boot-config-loader
STATE    Finished
HEALTH   ?
EVENTS   [Finished]: Service finished successfully (466959h37m5s ago)
         [Running]: Started task ext-rpi-boot-config-loader (PID 3608) for container ext-rpi-boot-config-loader (466959h37m6s ago)
         [Preparing]: Creating service runner (466959h37m7s ago)
         [Preparing]: Running pre state (466959h37m7s ago)
         [Waiting]: Waiting for file "/var/etc/rpi-boot-config-loader/config.txt" to exist (466959h37m9s ago)
         [Waiting]: Waiting for service "containerd" to be "up", file "/var/etc/rpi-boot-config-loader/config.txt" to exist (466959h37m15s ago)
         [Waiting]: Waiting for service "containerd" to be registered, file "/var/etc/rpi-boot-config-loader/config.txt" to exist (466959h37m21s ago)
         [Waiting]: Waiting for service "containerd" to be "up", file "/var/etc/rpi-boot-config-loader/config.txt" to exist (466959h37m22s ago)
```
```shell
❯ talosctl logs  ext-rpi-boot-config-loader -n <node-ip>         
<node-ip>: time="1970-01-01T00:00:27Z" level=info msg="Going to replace boot config..."
<node-ip>: time="1970-01-01T00:00:27Z" level=info msg="Writing config.txt"
<node-ip>: time="1970-01-01T00:00:27Z" level=info msg="config.txt replaced!"
<node-ip>: time="1970-01-01T00:00:27Z" level=info msg="Don't forget to reboot for changes to take effect."
```

#### 4. Reboot

Once confirmed that everything went as planned, reboot the device so that it boots with the new config.

**Word of caution again**, if the config.txt file is invalid, the node will not boot until it has been fixed!

#### When Upgrading A Node

If you are performing a node upgrade, you must follow steps 3 and 4 again!
This is due to when running the upgrade command, Talos re-creates the boot partition, which means that the changes have been overwritten.

