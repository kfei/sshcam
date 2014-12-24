# sshcam

Webcam live streaming in SSH terminal.

[Watch a Demo](https://asciinema.org/a/14968)

## Install

There are several ways to install `sshcam`.

**Install binary from GitHub** if you don't have a Go development environment:

```bash
curl -sL https://github.com/kfei/sshcam/releases/download/0.9.1/sshcam-x64.tar.bz | tar xj
mv sshcam /usr/local/bin/
```

**Build from source**:

```bash
# Build passed on Go version 1.4 and GCC version 4.9.2
go get -u github.com/kfei/sshcam
cd $GOPATH/src/github.com/kfei/sshcam
go build
go install
```

**Build in Docker container**:

```bash
git clone https://github.com/kfei/sshcam
cd sshcam
docker build -t sshcam .
# After built, you can run sshcam via the Docker container.
# FIXME: The privileged flag is for access /dev/videoX
alias sshcam='docker run -it -p 5566:5566 --priviliged sshcam'
```

## Quick Start

Run `sshcam` directly (without any argument) to preview the live stream. Also
note that you should run `sshcam` as an user with privilege to open
`/dev/videoX`.

Start the SSH server with all default settings:

```bash
sshcam --server
```

Then on the client-side, run:

```bash
ssh sshcam@your.server.ip -p 5566 # Default login: sshcam/p@ssw0rd
```

There are some configurable settings for server, have a look at `sshcam -h` for
more details. For example,

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

  - **P2P functions**: So that it can be a terminal video messenger.
  - **Higher resolution**: Can this be done by using the Unicode block
    characters?
  - **Colors**: In a *xterm-256color* terminal, display images with colors.
  - **Port to other platforms**: Maybe by using QT?
  - **Better algorithm**: Currently the rendering is inefficient and lossy.

## Credits

This project is inspired from [txtcam](https://github.com/dhotson/txtcam) and
**Star Wars Movie in Telnet** (telnet://towel.blinkenlights.nl).

## License

The MIT License (MIT)
