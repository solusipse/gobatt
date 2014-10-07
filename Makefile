GOCOMMAND=go
GOFLAGS=build
prefix=/usr/local

all:
	$(GOCOMMAND) $(GOFLAGS) gobatt.go

go-gtk:
	go get github.com/mattn/go-gtk/gtk

install:
	install -m 0755 gobatt $(prefix)/bin

clean:
	rm -rf gobatt
