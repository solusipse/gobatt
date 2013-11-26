GOCOMMAND=go
GOFLAGS=build
prefix=/usr/local

all:
		$(GOCOMMAND) $(GOFLAGS) wm-batt-tray.go

install:
		install -m 0755 wm-batt-tray $(prefix)/bin

clean:
		rm -rf wm-batt-tray