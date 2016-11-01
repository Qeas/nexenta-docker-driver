NEDGE_DEST = /usr/bin/
NEDGE_ETC = /etc/nvd/
NDVOL_EXE = ndvol

build: 
	go get -v github.com/docker/go-plugins-helpers/...
	cd $(GOPATH)/src/github.com/docker/go-plugins-helpers/volume; git checkout 60d242c
	go get -v github.com/Nexenta/nedge-docker-volume/...

lint:
	GOPATH=$(shell pwd) GOROOT=$(GO_INSTALL) $(GO) get -v github.com/golang/lint/golint
	for file in $$(find . -name '*.go' | grep -v vendor | grep -v '\.pb\.go' | grep -v '\.pb\.gw\.go'); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

install: .build
	cp -f $(GOPATH)/bin/$(NDVOL_EXE) $(NEDGE_DEST)

uninstall:
	rm -f $(NEDGE_ETC)/ndvol.json
	rm -f $(NEDGE_DEST)/ndvol

clean:
	go clean github.com/Nexenta/nedge-docker-volume
