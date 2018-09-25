GO = go
GO_FILES = $(wildcard *.go)
FMT = $(GO)imports -w 
LINT = $(GO)lint
TEST = $(GO) test
VET = $(GO) vet

fmt:
	$(FMT) $(GO_FILES)
lint:
	$(VET) $(GO_FILES)
	$(LINT) $(GO_FILES)
test:
	$(TEST)
all: fmt lint test
