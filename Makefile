.PHONY: all clean install lint protos protos_clean test

cmd_bin = $(shell go list ./src/cmd/... | xargs -I path sh -c 'echo "bin/$$(basename path)"')
src = $(shell find ./src/ -iname "*.go" -not -wholename "*/src/cmd/*")

protos = $(shell find src -iname "*.proto")
protos_go = ${protos:.proto=.pb.go}
service_name=tempulate
now=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# For now we run generate in all our dirs (could filter this in future)
generate_dirs=$(shell find src -maxdepth 1 -type d -not -iname "cmd" -not -iname "." -not -iname "src")
generate_dirs_flag_file=$(addsuffix /.generate,$(generate_dirs))

# Get the branch from git, or env var
ifndef BRANCH
	branch=$(shell git symbolic-ref --short HEAD)
else
	branch=$(BRANCH)
endif

# Get git hash
ifndef GIT_HASH
	git_hash=$(shell git rev-parse HEAD)
else
	git_hash=$(GIT_HASH)
endif

# Set semantic version, if not set from build job (for building on dev machines)
$(shell git show origin/version:version > version || true) # keep calm and carry on - this will fail in CI, ok with this
ifndef SEMVER
	semver=$(shell cat version)
else
	semver=$(SEMVER)
endif

# Build tag
ifndef GIT_TAG
	git_tag=$(semver)-$(branch)
else
	git_tag=$(GIT_TAG)
endif

all_build=$(protos_go) $(generate_dirs_flag_file) $(cmd_bin)
kitchen_sink=$(all_build)
all: $(kitchen_sink)

protos: $(protos_go)

# Command binaries depend on it's own package and other non-command package go files
bin/%:src/cmd/%/*.go $(src)
	go build -ldflags "-X main.version=$(semver)" -o $(subst /cmd,,$@) ./$(subst bin/,src/cmd/,$@)

# Generated protobufs go files depends on protos files
%.pb.go:%.proto
	protoc $< --go_out="plugins=grpc:."

# Handle go generate ?
genfiles=$(shell find src -type f -not -name ".generate" -not -name "*.pb.go")
src/%/.generate: $(genfiles)
	if [ $(shell find ./$(dir $@) -maxdepth 1 -name "*.go" | wc -l) -gt 0 ]; then go generate ./$(basename $@);	fi;
	touch ./$@

clean: protos_clean
	rm -rf $(cmd_bin)
	rm -rf bin
	rm $(generate_dirs_flag_file) || true

protos_clean:
	rm -rf $(protos_go)

lint: $(GOPATH)/bin/golint
	golint ./src/...

test: $(protos_go)
	go test -v ./...


$(GOPATH)/bin/golint:
	go get -u golang.org/x/lint/golint
