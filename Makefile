DIST=$(shell pwd)/dist
OUTPUT=$(DIST)/out
BUILD=$(shell pwd)/build
BUILD_DEPS=$(BUILD)/deps

BUSYBOX_VER=1.21.1

QEMU_BIN=qemu-system-x86_64
QEMU_IMG=qemu-img
BOOT_ARGS="root=UUID=87965c58-36e0-11f1-8de9-0242ac120002 console=ttyS0 rdinit=/init early"
DEBUG_BOOT_ARGS="debug root=UUID=87965c58-36e0-11f1-8de9-0242ac120002 console=ttyS0 rdinit=/init early"
KERNEL_PATH=$(OUTPUT)/boot/vmlinuz
INITRAMFS_PATH=$(OUTPUT)/initramfs/initramfs.cpio

DEV_FLAGS=-kernel $(KERNEL_PATH) -initrd $(INITRAMFS_PATH)

DOCKER_BIN=docker
DOCKER_FILE=build/Dockerfile
DOCKER_LOCAL_IMAGE="container-os:$(shell date +%s)"
DOCKER_CONTAINER_NAME="container-os"

UKIFY_BIN=ukify
GENIMG_BIN=genimage

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
	$(DOCKER_BIN) build . --platform linux/amd64 -f $(DOCKER_FILE) -t $(DOCKER_LOCAL_IMAGE)
	@$(DOCKER_BIN) create --name $(DOCKER_CONTAINER_NAME) $(DOCKER_LOCAL_IMAGE) sh
	@$(DOCKER_BIN) export $(DOCKER_CONTAINER_NAME) > $(DIST)/dist.tar 
	@$(DOCKER_BIN) rm $(DOCKER_CONTAINER_NAME) 
	
	# Export and Remove additional directories that 'docker export' adds
	@tar -xf $(DIST)/dist.tar -C $(OUTPUT)
	@rm -rf $(OUTPUT)/dev $(OUTPUT)/sys $(OUTPUT)/proc $(OUTPUT)/etc $(OUTPUT)/.dockerenv

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

	$(UKIFY_BIN) build \
		--linux=$(DIST)/out/boot/vmlinuz \
		--initrd=$(DIST)/out/initramfs/initramfs.cpio \
		--cmdline=$(BOOT_ARGS) \
		--output=/tmp/container_os/esp/EFI/BOOT/BOOTX64.EFI

	cp $(DIST)/out/rootfs/rootfs.tar.gz /tmp/container_os/root
	tar -xvf /tmp/container_os/root/rootfs.tar.gz -C /tmp/container_os/root --strip-components 2

	$(GENIMG_BIN) \
		--config $(BUILD)/genimage.cfg \
		--rootpath=/tmp/container_os/root \
		--tmppath=/tmp/genimage-tmp \
		--inputpath=/tmp/container_os/esp \
		--outputpath=$(DIST)/out/img

	rm -rf /tmp/container_os

# This really needs to be condensed as we are repeating way too much from .PHONY: image
# This will do for the time being, but is subjected to change in a later commit
.PHONY: debug-image
debug-image:
	mkdir -p /tmp/container_os/esp/EFI/BOOT
	mkdir -p /tmp/container_os/root
	mkdir -p $(DIST)/out/img

	$(UKIFY_BIN) build \
		--linux=$(DIST)/out/boot/vmlinuz \
		--initrd=$(DIST)/out/initramfs/initramfs.cpio \
		--cmdline=$(DEBUG_BOOT_ARGS) \
		--output=/tmp/container_os/esp/EFI/BOOT/BOOTX64.EFI

	cp $(DIST)/out/rootfs/rootfs.tar.gz /tmp/container_os/root
	tar -xvf /tmp/container_os/root/rootfs.tar.gz -C /tmp/container_os/root --strip-components 2

	$(GENIMG_BIN) \
		--config $(BUILD)/genimage.debug.cfg \
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
	$(QEMU_BIN) -machine type=q35,accel=kvm \
				-drive if=pflash,format=raw,readonly=on,file=$(DIST)/uefi/OVMF_CODE.fd \
				-drive if=pflash,format=raw,readonly=on,file=$(DIST)/uefi/OVMF_VARS.fd \
				-drive file=$(OUTPUT)/img/disk.img,format=raw,if=virtio \
				-m 2G \
				-smp 1

.PHONY: run-debug
run-debug:
	$(QEMU_BIN) -machine type=q35,accel=kvm \
				-drive if=pflash,format=raw,readonly=on,file=$(DIST)/uefi/OVMF_CODE.fd \
				-drive if=pflash,format=raw,readonly=on,file=$(DIST)/uefi/OVMF_VARS.fd \
				-drive file=$(OUTPUT)/img/debug-disk.img,format=raw,if=virtio \
				-m 2G \
				-smp 1

.PHONY: run
run-simple:
	$(QEMU_BIN) $(DEV_FLAGS) -append $(DEBUG_BOOT_ARGS) -nographic
