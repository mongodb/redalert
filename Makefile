clean:
	find . -name '*.pyc' -delete
	find . -name '__pycache__' -delete
	rm -rf build
	rm -rf dist
	rm -rf redalert.egg-info
	rm -rf *.log


halt_vms:
	vagrant halt centos ubuntu windows

unit_test: clean
	pytest

test_%: clean
	PLATFORM=$* ./scripts/test_vagrant.sh
	
test_no_halt: unit_test  test_ubuntu  test_centos  test_windows
test: unit_test halt_vms test_ubuntu halt_vms test_centos halt_vms test_windows
