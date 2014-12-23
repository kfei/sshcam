# sshcam

Webcam live streaming in SSH terminal.

[DEMO](https://asciinema.org/a/14968)

## Install

**Install binary from GitHub** if you don't have a Go development environment:

```bash
curl -sL https://github.com/kfei/sshcam/releases/download/v1.0.0/sshcam-v1.0.0-x64.tar.bz | tar xj
```

**Build from source** (also requires GCC):

```bash
go get -u github.com/kfei/sshcam
cd $GOPATH/src/github.com/kfei/sshcam
go build
go install
```

## Quick Start

Preview the live stream locally via `sshcam` (without any argument). Also note
that you should run `sshcam` as an user with privilege to open `/dev/videoX`.

Start the SSH server with all default settings:

```bash
sshcam --server
```

Then on the client-side, run:

```bash
ssh sshcam@your.server.ip -p 5566 # Default login: sshcam/p@ssw0rd
```

Note that there are more configurable settings for server, have a look at
`sshcam -h` for more details. For example,

```bash
sshcam --server --listen=127.0.0.1 --port=22222 \
       --user=john --pass=nhoj \
       --device=/dev/video0 --size=1280x720
```

## Requirements

  - On the client-side, a standard `ssh` utility is enough.
  - Video device is supported by 
    [V4L2](https://www.kernel.org/doc/Documentation/video4linux/v4l2-framework.txt),
    which means the server is currently Linux only.

## TODO

There are still lots of interesting works to be done. Discussions and pull
requests are both welcome. :)

  - *P2P functions*: So that it can be a terminal video messenger.
  - *Higher resolution*: Can this be done by using the Unicode block
    characters?
  - *Colors*: In a *xterm-256color* terminal, display images with colors.
  - *Port to other platforms*, maybe QT?
  - *Improve algorithm*: Currently the image scaling is inefficient and lossy.

## Credits

This project is inspired from [txtcam](https://github.com/dhotson/txtcam) and
[Star Wars Movie in Telnet](telnet://towel.blinkenlights.nl).

## License

The MIT License (MIT)
