GO="$(shell which go)"
ROOT="$(shell pwd)"
CMD="$(ROOT)/cmd"
BIN="$(ROOT)/bin"
TARGETS = apiserver kubelet kubectl 

.PHONY: all

build:
	@$(foreach target,$(TARGETS), \
	$(GO) build -o $(BIN)/$(target) $(CMD)/$(target)/$(target).go;)

all: build

clean: 
	@$(foreach target,$(TARGETS), \
	rm $(BIN)/$(target);)
