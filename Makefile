RELEASE ?= "release"

DIST=$(shell pwd)/out/apathy-os

.PHONY: clean
.SILENT:
clean:
	rm -rf $(DIST)/image

.PHONY: snapshot
snapshot:
	goreleaser release --clean --snapshot

.PHONY: clean-uefi
.SILENT:
	rm -rf $(DIST)/uefi

.PHONY: image
image:
	docker buildx bake image

.PHONY: initd
initd:
	docker buildx bake go-initd

.PHONY: initramfs
initramfs:
	docker buildx bake initramfs

.PHONY: kernel
kernel:
	docker buildx bake kernel

.PHONY: rootfs
rootfs:
	docker buildx bake rootfs

.PHONY: serviced
serviced:
	docker buildx bake go-serviced

.PHONY: toolchain
toolchain:
	docker buildx bake toolchain

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
	qemu-system-x86_64 \
				-display gtk,show-tabs=on \
				-machine type=q35,accel=kvm \
				-drive if=pflash,format=raw,readonly=on,file=$(DIST)/uefi/OVMF_CODE.fd \
				-drive if=pflash,format=raw,readonly=on,file=$(DIST)/uefi/OVMF_VARS.fd \
				-drive file=$(DIST)/image/disk.img,format=raw,if=virtio \
				-m 2G \
				-smp 1
