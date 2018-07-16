clean:
	rm -rf build
	rm -rf dist

halt_vms:
	vagrant halt centos ubuntu windows

unit_test: clean
	go test ./...

test_%: clean
	vagrant rsync $*; PLATFORM=$* ./scripts/test_vagrant.sh
	
test_no_halt: unit_test  test_ubuntu  test_centos  test_windows
test: unit_test halt_vms test_ubuntu halt_vms test_centos halt_vms test_windows
