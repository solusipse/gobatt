gobatt
============

Lightweight battery tray icon for Linux.

![gobatt](http://solusipse.net/misc/gobatt.gif)

-------------------------------------------------------------

## Installation ##

### Requirements ###

- Go compiler (https://code.google.com/p/go/downloads/list)
- go-gtk (https://github.com/mattn/go-gtk)

### Installation ###

```
git clone https://github.com/solusipse/gobatt.git && cd gobatt
make go-gtk
make
sudo make install
```

-------------------------------------------------------------

## Usage ##

Add this line to your custom startup script or to `.xinitrc`:

```
gobatt
```

-------------------------------------------------------------

## License ##

See `LICENSE`.
