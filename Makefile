DIST=$(shell pwd)/dist
OUTPUT=$(DIST)/out
BUILD=$(shell pwd)/build
BUILD_DEPS=$(BUILD)/deps

BUSYBOX_VER=1.21.1

QEMU_BIN=qemu-system-x86_64
QEMU_IMG=qemu-img
BOOT_ARGS="console=ttyS0 rdinit=/init"
DEBUG_BOOT_ARGS="debug console=ttyS0 rdinit=/init"
KERNEL_PATH=$(OUTPUT)/boot/vmlinuz
INITRAMFS_PATH=$(OUTPUT)/initramfs/initramfs.cpio

DEV_FLAGS=-kernel $(KERNEL_PATH) -initrd $(INITRAMFS_PATH)

DOCKER_BIN=docker
DOCKER_FILE=build/Dockerfile
DOCKER_LOCAL_IMAGE="container-os:$(shell date +%s)"
DOCKER_CONTAINER_NAME="container-os"

.PHONY: clean
.SILENT:
clean:
	rm -rf $(DIST)

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

.PHONY: disk
disk:
	@mkdir -p $(OUTPUT)/disk
	@$(QEMU_IMG) create -f qcow2 $(OUTPUT)/disk/boot.qcow2 10G

.PHONY: run
run:
	$(QEMU_BIN) $(DEV_FLAGS) -append $(BOOT_ARGS) -nographic

.PHONY: run-debug
run-debug:
	$(QEMU_BIN) $(DEV_FLAGS) -drive file=$(OUTPUT)/disk/boot.qcow2,if=ide,format=qcow2 -append $(DEBUG_BOOT_ARGS) -nographic
