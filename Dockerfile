FROM base/archlinux

MAINTAINER kfei <kfei@kfei.net>

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

RUN pacman -Syy && pacman -S --noconfirm gcc git go \
    && go get -u github.com/kfei/sshcam \
    && cd $GOPATH/src/github.com/kfei/sshcam \
    && go build \
    && go install \
    && pacman --noconfirm -R gcc git go \
    && pacman --noconfirm -R $(pacman -Qdtq)

VOLUME ["/.sshcam"]

ENTRYPOINT ["/go/bin/sshcam"]
