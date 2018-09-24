GO = go
GO_FILES = $(wildcard *.go)
FMT = $(GO)fmt -w
LINT = $(GO)lint

fmt:
	$(FMT) $(GO_FILES)
lint:
	$(LINT) $(GO_FILES)
all: fmt lint
