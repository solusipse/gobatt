GOCOMMAND=go
GOFLAGS=build
prefix=/usr/local

all:
		$(GOCOMMAND) $(GOFLAGS) wb-batt-tray.go

install:
		install -m 0755 wb-batt-tray $(prefix)/bin

clean:
		rm -rf wb-batt-tray