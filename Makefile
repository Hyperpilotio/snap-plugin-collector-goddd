TIME=$$(date +"%Y%m%d-%H-%M")

default:
	$(MAKE) deps
	$(MAKE) all
test:
	$(MAKE) deps
	bash -c "./scripts/test.sh $(TEST)"
deps:
	bash -c "./scripts/deps.sh"
check:
	$(MAKE) test
all:
	bash -c "./scripts/build.sh $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))"
test-docker:
	docker run --name snap-goddd --rm -v "$$PWD":/go/src/github.com/swhsiang/snap-plugin-collector-goddd -w /go/src/github.com/swhsiang/snap-plugin-collector-goddd golang:1.8 bash -c make
#	docker run --name snap-goddd -d hyperpilot/snap-plugin-collector-goddd:latest
test-docker-clean:
	docker stop snap-goddd && docker rm snap-goddd

