TIME=$$(date +"%Y%m%d-%H-%M")

default:
	$(MAKE) dep-glide
	$(MAKE) all
test:
	bash -c "./scripts/test.sh $(TEST)"
deps:
	bash -c "./scripts/deps.sh"
dep-glide:
	glide install
check:
	$(MAKE) dep-glide
	$(MAKE) test
all:
	bash -c "./scripts/build.sh $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))"
