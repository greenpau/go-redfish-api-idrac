.PHONY: test ctest covdir coverage docs linter qtest clean dep release logo
APP_VERSION:=$(shell cat VERSION | head -1)
GIT_COMMIT:=$(shell git describe --dirty --always)
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD -- | head -1)
BUILD_USER:=$(shell whoami)
BUILD_DATE:=$(shell date +"%Y-%m-%d")
APP:="go-redfish-api-idrac-client"
BINARY:="go-redfish-api-idrac-client"
VERBOSE:=-v
ifdef TEST
	TEST:="-run ${TEST}"
endif

all:
	@echo "Version: $(APP_VERSION), Branch: $(GIT_BRANCH), Revision: $(GIT_COMMIT)"
	@echo "Build on $(BUILD_DATE) by $(BUILD_USER)"
	@mkdir -p bin/
	@CGO_ENABLED=0 go build -o bin/$(BINARY) $(VERBOSE) \
		-ldflags="-w -s \
		-X main.appName=$(BINARY) \
		-X main.appVersion=$(APP_VERSION) \
		-X main.gitBranch=$(GIT_BRANCH) \
		-X main.gitCommit=$(GIT_COMMIT) \
		-X main.buildUser=$(BUILD_USER) \
		-X main.buildDate=$(BUILD_DATE)" \
		-gcflags="all=-trimpath=$(GOPATH)/src" \
		-asmflags="all=-trimpath $(GOPATH)/src" cmd/$(APP)/*
	@echo "Done!"

linter:
	@echo "Running lint checks"
	@golint -set_exit_status cmd/$(APP)/*.go
	@golint -set_exit_status pkg/client/*.go
	@golint -set_exit_status internal/client/*.go
	@for f in `find ./pkg -type f -name '*.go'`; do echo $$f; go fmt $$f; golint -set_exit_status $$f; done
	@echo "PASS: golint"

test: covdir linter
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out ./pkg/client/*.go

ctest: covdir linter
	@#richgo version || go get -u github.com/kyoh86/richgo
	@time richgo test $(VERBOSE) $(TEST) -coverprofile=.coverage/coverage.out ./pkg/client/*.go

covdir:
	@echo "Creating .coverage/ directory"
	@mkdir -p .coverage

coverage:
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html
	@go test -covermode=count -coverprofile=.coverage/coverage.out ./pkg/client/*.go
	@go tool cover -func=.coverage/coverage.out | grep -v "100.0"

docs:
	@mkdir -p .doc
	@go doc -all pkg/client > .doc/index.txt
	@cat .doc/index.txt

clean:
	@rm -rf .doc
	@rm -rf .coverage
	@rm -rf bin/

qtest:
	@echo "Perform quick tests ..."
	@#go test -v -run TestClient ./pkg/client/*.go
	@#go test -v -run TestParseInfoJsonOutput ./pkg/client/*.go
	@#go test -v -run TestParseResourceJsonOutput ./pkg/client/*.go
	@#go test -v -run TestParseComputerSystemJsonOutput ./pkg/client/*.go
	@#time richgo test -v -run TestComputerSystemStruct ./pkg/client/*.go
	@go test -v -run TestComputerSystemStruct ./pkg/client/*.go

dep:
	@echo "Making dependencies check ..."
	@go get -u golang.org/x/lint/golint
	@go get -u golang.org/x/tools/cmd/godoc
	@go get -u github.com/kyoh86/richgo

release:
	@echo "Making release"
	@if [ $(GIT_BRANCH) != "main" ]; then echo "cannot release to non-main branch $(GIT_BRANCH)" && false; fi
	@git diff-index --quiet HEAD -- || ( echo "git directory is dirty, commit changes first" && false )
	@versioned -patch
	@git add VERSION
	@git commit -m 'updated VERSION file'
	@versioned -sync cmd/${APP}/main.go
	@echo "Patched version"
	@git add cmd/${APP}/main.go
	@git commit -m "released v`cat VERSION | head -1`"
	@git tag -a v`cat VERSION | head -1` -m "v`cat VERSION | head -1`"
	@git push
	@git push --tags
	@echo "If necessary, run the following commands:"
	@echo "  git push --delete origin v$(APP_VERSION)"
	@echo "  git tag --delete v$(APP_VERSION)"

logo:
	@mkdir -p assets/docs/images
	@gm convert -background black -font Bookman-Demi \
		-size 640x320 "xc:black" \
		-pointsize 72 \
		-draw "fill white gravity center text 0,0 'iDRAC API'" \
		assets/docs/images/logo.png
