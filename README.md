Checking wifi AP for connected stations, controlling monitor on/off state.

# Build
Requires go >= 1.11

```
go build
```

# Usage
```
$ ./wifi-screen-control
NAME:
   wifi-screen-control - Checking wifi AP for connected stations, controlling monitor on/off state.

USAGE:
   wifi-screen-control [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     watch, w  watch wifi AP for stations to connect
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
```
$ ./wifi-screen-control help watch
NAME:
   wifi-screen-control watch - watch wifi AP for stations to connect

USAGE:
   wifi-screen-control watch [command options] DEVICE

OPTIONS:
   --interval value, -n value  Polling interval for the wifi status, in seconds (default: 10)
```

Example
```
./wifi-screen-control watch wlp0s29f7u2 --interval 1
```

# Systemd Setup
```
sudo cp wifi-screen-control /usr/local/bin/
sudo cp wifi-screen-control@.service /etc/systemd/user/

systemctl --user enable wifi-screen-control@wlp0s29f7u2
systemctl --user start wifi-screen-control@wlp0s29f7u2
```
