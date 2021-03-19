.POSIX:
fast:
	@./do.sh fast
dev:
	@./do.sh dev
install:
	@./do.sh fast
	@sudo cp ./dxhd /usr/bin/dxhd
	@sudo mkdir -p /usr/share/licenses/dxhd/
	@sudo cp LICENSE /usr/share/licenses/dxhd/
	@test -d /usr/lib/systemd/user \
		&& sudo cp systemd/dxhd.service /usr/lib/systemd/user/ \
		|| true
	@echo installed
check:
	@./do.sh check
