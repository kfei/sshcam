# sshcam

Webcam live streaming in SSH terminal.

[See it in action](https://asciinema.org/a/14968)

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
more information. As an example:

```bash
sshcam --server --listen=127.0.0.1 --port=22222 \
       --user=john --pass=nhoj \
       --device=/dev/video0 --size=1280x720 --color
```

## Requirements

  - On the client-side, a `ssh` utility with 256-colors support is enough.
  - Video device is supported by 
    [V4L2](https://www.kernel.org/doc/Documentation/video4linux/v4l2-framework.txt),
    which means the server is currently Linux only.

## Installation

There are several ways to install `sshcam`.

**Install binary from GitHub**:

```bash
curl -sL https://github.com/kfei/sshcam/releases/download/0.9.8/sshcam-x64.tar.bz | tar xj
mv sshcam /usr/local/bin/
```

**Build from source** if you have a Go development environment:

```bash
# Build passed on Go version 1.4 and GCC version 4.9.2
go get -u github.com/kfei/sshcam
cd $GOPATH/src/github.com/kfei/sshcam
go build
go install
```

**Build and run in Docker container**:

```bash
git clone https://github.com/kfei/sshcam
cd sshcam
docker build -t sshcam .
# After built, you can run sshcam via the Docker container
# The privileged flag is for /dev/videoX access (FIXME)
alias sshcam='docker run -it -p 5566:5566 --priviliged sshcam'
```

## TODO

There are still lots of interesting works to be done. Discussions and pull
requests are both welcome. :)

  - **Better performance**: Currently the pixel rendering is inefficient.
  - **Even higher resolution**: Use Unicode quadrant block characters [2596 to
    259F](http://www.alanwood.net/unicode/block_elements.html).
  - **Port to other platforms**: Maybe by using QT?
  - **P2P functions**: So that it can be a terminal video messenger.

## License

The MIT License (MIT)
