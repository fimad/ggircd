# GGircd [![Build Status](https://travis-ci.org/fimad/ggircd.svg?branch=master)](https://travis-ci.org/fimad/ggircd) #

GGircd is a minimalistic IRC server written in Go with the primary goal of being
feature compatible with the [zhenya_bot](https://github.com/numberten/zhenya_bot)
IRC bot and a secondary goal of being easy to configure and deploy.

## Implementation Status ##

GGircd currently implements a subset of the IRC protocol as defined in
[RFC 1459](https://www.ietf.org/rfc/rfc1459),
[RFC 2812](https://tools.ietf.org/html/rfc2812), and various mailing lists.
Notable missing features are any form of
server-to-server connections, banning users, and various status messages.

## Installation ##

Eventually, there will be two ways of installing GGircd: from source and from a
pre-built binary. Pre-built binaries don't currently exist, so you are stuck
with building from source.

### Installing From Source ###

Installing from source will give you the most up-to-date version of GGircd
available. Every effort is made to keep this bug free, but sometimes mistakes
are made. :(

These instructions assume you have already installed Go. If you haven't
installed Go, you can find easy-to-follow instructions
[here](http://golang.org/doc/install).

If you already have Go installed, installing GGircd is as easy as:

    go get -u github.com/fimad/ggircd/ggircd
    go install github.com/fimad/ggircd/ggircd

## Running ##

After installing, the `ggircd` binary will be placed in `$GOPATH/bin`. This
binary takes the following command line arguments:

  * --config takes a path to the ggircd.conf config file. The default value is
    `/etc/ggircd/ggircd.conf`.

  * --log takes an integer and controls the verbosity of GGircd's logging. The
    higher the number, the more verbose the logging. The default value is 3.

## Configuration ##

Before you can run GGircd, you'll need to configure your server.
Configuration is done via a JSON encoded file, the default location of which
is `/etc/ggircd/ggircd.conf`.

The `ggircd.conf` file at the root of this repository provides the minimal
configuration necessary to get a server up and running.

### Configuration Options ###

The following are the possible fields that can appear in the ggircd.conf file
along with descriptions and example values.

#### Name ####

This is the name of your IRC server.

Example: "ggircd"

#### Network ####

This is the name of your IRC network. This field is basically meaningless at
the moment because multi-server IRC networks are not supported.

Example: "My IRC Network"

#### Port ####

The port that GGircd should listen on.

Example: 6667

#### SSLPort ####

The port that GGircd should listen on for SSL connections.

Example: 6697

#### SSLCertificate ####

An object containing the paths for the private key and certificate files (PEM
encoded) that are to be used for SSL connections. If this option is empty, then
no SSL port will be opened.

Example: {"KeyFile": "./key.pem", "CertFile": "./cert.pem"}

#### MOTD ####

The path to a file containing the message of the day. This file should be line
wrapped at 80 columns.

Example: "/etc/ggircd/motd"

#### DefaultChannelMode ####

This is the default mode all channels will have when they are first created.

Example: "t"

#### DefaultUserMode ####

This is the default mode all users will have when they first connect.

Example: "i"

#### SpoofHostName (Optional) ####

If given, all connected users will be listed as having this host name.

Example: "i.love.plan9.bell-labs.com"

#### PingFrequency (Optional) ####

The frequency (in seconds) to send out PING messages to all connected clients.
If no value is given, this defaults to 30 seconds.

Example: 30

#### PongMaxLatency (Optional) ####

The maximum latency a client can respond to a PING message with without being
forcefully disconnected. If no value is given, this defaults to 5 seconds.

Example: 5

#### Prometheus.Port (Optional) ####

If supplied the port that Prometheus metrics will be exposed on. If zero or
unspecified Prometheus metrics will not be collected.

Example: 4000

#### Logs.Path (Optional) ####

A path to a directory that the IRC server will write logs files to.

Example: /var/log/irc/

#### Logs.LogChannelMessages (Optional) ####

If true, a log file will be created in the log directory for each channel. Each
file will contain JSON serialized IRC message.

Example: true
