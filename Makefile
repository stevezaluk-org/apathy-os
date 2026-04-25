RELEASE ?= "release"
LIBC_FLAVOR ?= "static"

DIST=$(shell pwd)/dist
OUTPUT=$(DIST)/out
BUILD=$(shell pwd)/build
BUILD_DEPS=$(BUILD)/deps

BUSYBOX_VER=1.21.1

BOOT_ARGS_release="root=87965c58-36e0-11f1-8de9-0242ac120002 console=ttyS0 rdinit=/init early"
BOOT_ARGS_development="debug root=87965c58-36e0-11f1-8de9-0242ac120002 console=ttyS0 rdinit=/init early"

BOOT_ARGS=$(BOOT_ARGS_$(RELEASE))

.PHONY: clean
.SILENT:
clean:
	rm -rf $(OUTPUT)

.PHONY: clean-deps
.SILENT:
clean-deps:
	rm -rf $(BUILD_DEPS)

.PHONY: build
build:
	@mkdir -p $(OUTPUT)
	docker buildx bake \
		--file build/docker-bake.hcl \
		--set "_common.args.RELEASE=$(RELEASE)" \
		--set "_common.args.LIBC_FLAVOR=$(LIBC_FLAVOR)" \
		_common

.PHONY: deps
deps:
	mkdir -p $(BUILD_DEPS)
	echo "[deps] Downloading busybox ($(BUSYBOX_VER)). Warning: This may fail due to issues with BusyBox content delivery. A mirror should really be added here"
	@wget https://busybox.net/downloads/binaries/$(BUSYBOX_VER)/busybox-x86_64 -O $(BUILD_DEPS)/busybox-x86_64
	@chmod +x $(BUILD_DEPS)/busybox-x86_64

.PHONY: image
image:
	mkdir -p /tmp/container_os/esp/EFI/BOOT
	mkdir -p /tmp/container_os/root
	mkdir -p $(DIST)/out/img

	ukify build \
		--linux=$(DIST)/out/$(RELEASE)-$(LIBC_FLAVOR)/boot/vmlinuz \
		--initrd=$(DIST)/out/$(RELEASE)-$(LIBC_FLAVOR)/boot/initrd.img \
		--cmdline=$(BOOT_ARGS) \
		--output=/tmp/container_os/esp/EFI/BOOT/BOOTX64.EFI

	cp -r $(DIST)/out/$(RELEASE)-$(LIBC_FLAVOR)/* /tmp/container_os/root

	genimage \
		--config $(BUILD)/genimage/release.cfg \
		--rootpath=/tmp/container_os/root \
		--tmppath=/tmp/genimage-tmp \
		--inputpath=/tmp/container_os/esp \
		--outputpath=$(DIST)/out/img

	rm -rf /tmp/container_os

.PHONY: uefi
uefi:
	bash -c ' \
		if [ ! -f "/usr/share/OVMF/OVMF_CODE.fd" ] || [ ! -f "/usr/share/OVMF/OVMF_VARS.fd" ]; then \
			echo "[error] OVMF firmware files do not exist at: /usr/share/OVMF"; \
			echo "Install them through your package manager"; \
		else \
			mkdir -p $(DIST)/uefi; \
			cp /usr/share/OVMF/OVMF_CODE.fd $(DIST)/uefi; \
			cp /usr/share/OVMF/OVMF_VARS.fd $(DIST)/uefi; \
		fi'

.PHONY: run
run:
	qemu-system-x86_64 -machine type=q35,accel=kvm \
				-drive if=pflash,format=raw,readonly=on,file=$(DIST)/uefi/OVMF_CODE.fd \
				-drive if=pflash,format=raw,readonly=on,file=$(DIST)/uefi/OVMF_VARS.fd \
				-drive file=$(OUTPUT)/img/disk.img,format=raw,if=virtio \
				-m 2G \
				-smp 1
