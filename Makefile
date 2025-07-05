BIN := binlog-finder
MYSQLBINLOG := pkg/bin/mysqlbinlog

.PHONY: build test package clean

build:
	go build -o $(BIN) ./cmd

test:
	go test ./...

package: $(MYSQLBINLOG)
	test -x $(MYSQLBINLOG)
	$(MAKE) build
	mkdir -p dist/pkg/bin
	cp $(BIN) dist/$(BIN)
	cp $(MYSQLBINLOG) dist/pkg/bin/
	cp README.md dist/

clean:
	rm -f $(BIN)
	rm -rf dist
