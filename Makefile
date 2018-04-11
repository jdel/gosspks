VERSION=`git rev-parse --short HEAD`

.PHONY: default all x-gosspks docker clean

default: gosspks

all: clean gosspks x-gosspks docker

docker:
	@docker build --no-cache -t jdel/gosspks:local --build-arg GOSSPKS_COMMIT="${VERSION}" .
	
gosspks:
	go build -ldflags "-X github.com/jdel/gosspks/cfg.Version=${VERSION}"

x-gosspks:
	@go get github.com/mitchellh/gox
	gox -parallel=1 -osarch="linux/386 linux/amd64 linux/arm darwin/amd64 darwin/386 windows/amd64 windows/386" -output="out/{{.Dir}}-{{.OS}}-{{.Arch}}" -ldflags "-X github.com/jdel/gosspks/cfg.Version=${VERSION}"

clean:
	@rm -rf gosspks debug out