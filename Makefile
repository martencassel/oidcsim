BINARY=oidcsim
SRC=main.go

all: build

build:
	/usr/local/go/bin/go build -o $(BINARY) $(SRC)

run: build
	./$(BINARY) --listen $(ADDR) --port $(PORT) --tls-cert $(CERT) --tls-key $(KEY)

clean:
	rm -f $(BINARY)

.PHONY: all build clean run

gen-cert:
	@./gen-cert.sh /tmp/ca.pem /tmp/ca.key server.pem server.key  DNS:idp.local
	@openssl x509 -in server.pem -text -noout

allocate-ip:
	@sudo ./add-ip.sh 172.21.22.200/20
