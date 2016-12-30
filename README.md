# Installation

## Install latest Golang

## go-bindata

    go get -u github.com/jteeuwen/go-bindata/...

## Compile

    cd $GOPATH/src/vserv
    go-bindata tmpl/ && go install

# Usage

    vserv video.mp4
