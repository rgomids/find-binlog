BIN := binlog-finder
MYSQLBINLOG := pkg/bin/mysqlbinlog
MYSQL := pkg/bin/mysql

.PHONY: build test package clean

build:
	go build -o $(BIN) ./cmd

test:
	go test ./...

package: $(MYSQLBINLOG) $(MYSQL)
	test -x $(MYSQLBINLOG)
	test -x $(MYSQL)
	$(MAKE) build
	mkdir -p dist/pkg/bin
	mv $(BIN) dist/$(BIN)
	cp $(MYSQLBINLOG) dist/pkg/bin/
	cp $(MYSQL) dist/pkg/bin/
	cp README.md dist/
	cd ./dist; zip -r ../binlog-finder.zip ./*; cd ../

clean:
	rm -f $(BIN)
	rm -rf dist
