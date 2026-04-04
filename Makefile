DIST=$(shell pwd)/dist
OUTPUT=$(DIST)/out

QEMU_BIN=/usr/bin/qemu-system-x86_64
BOOT_ARGS="console=ttyS0 rdinit=/init"
DEBUG_BOOT_ARGS="debug console=ttyS0 rdinit=/init"
KERNEL_PATH=$(OUTPUT)/boot/vmlinuz
INITRAMFS_PATH=$(OUTPUT)/initramfs/initramfs.cpio

DEV_FLAGS=-kernel $(KERNEL_PATH) -initrd $(INITRAMFS_PATH)

DOCKER_BIN=/usr/bin/docker
DOCKER_FILE=build/Dockerfile
DOCKER_LOCAL_IMAGE="container-os:$(shell date +%s)"
DOCKER_CONTAINER_NAME="container-os"

.PHONY: clean
.SILENT:
clean:
	rm -rf $(DIST)

.PHONY: build
build:
	@mkdir -p $(OUTPUT)
	$(DOCKER_BIN) build . -f $(DOCKER_FILE) -t $(DOCKER_LOCAL_IMAGE)
	@$(DOCKER_BIN) create --name $(DOCKER_CONTAINER_NAME) $(DOCKER_LOCAL_IMAGE) sh
	@$(DOCKER_BIN) export $(DOCKER_CONTAINER_NAME) > $(DIST)/dist.tar 
	@$(DOCKER_BIN) rm $(DOCKER_CONTAINER_NAME) 
	
	# Export and Remove additional directories that 'docker export' adds
	@tar -xf $(DIST)/dist.tar -C $(OUTPUT)
	@rm -rf $(OUTPUT)/dev $(OUTPUT)/sys $(OUTPUT)/proc $(OUTPUT)/etc $(OUTPUT)/.dockerenv

.PHONY: run
run:
	$(QEMU_BIN) $(DEV_FLAGS) -append $(BOOT_ARGS) -nographic

.PHONY: run-debug
run-debug:
	$(QEMU_BIN) $(DEV_FLAGS) -append $(DEBUG_BOOT_ARGS) -nographic
