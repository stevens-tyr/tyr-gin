GO = go
GO_FILES = $(wildcard *.go)
# go get golang.org/x/tools/cmd/goimports
FMT = $(GO)imports -w 
LINT = $(GO)lint
TEST = $(GO) test
VET = $(GO) vet

fmt:
	$(FMT) $(GO_FILES)
lint:
	$(VET)
	$(LINT) $(GO_FILES)
test:
	$(TEST)
all: fmt lint test
