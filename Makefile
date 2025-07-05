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
	mv $(BIN) dist/$(BIN)
	cp $(MYSQLBINLOG) dist/pkg/bin/
	cp README.md dist/
	cd ./dist; zip -r ../binlog-finder.zip ./*; cd ../

clean:
	rm -f $(BIN)
	rm -rf dist
