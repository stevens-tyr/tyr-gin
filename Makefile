GO = go
GO_FILES = $(wildcard *.go)
FMT = $(GO)imports -w 
LINT = $(GO)lint
TEST = $(GO) test
VET = $(GO) vet -composites=false
GET = $(GO) get

.PHONY: get fmt lint test all

get:
	$(GET)
fmt:
	$(FMT) $(GO_FILES)
lint:
	$(VET) $(GO_FILES)
	$(LINT) $(GO_FILES)
test:
	$(TEST)
clean:
	rm -f log.json
	rm -f *~
	rm -f \#*\#
all: fmt lint test clean
