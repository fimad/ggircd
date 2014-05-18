# GGircd #

GGircd is a minimalistic IRC server written in Go with the primary goal of being
feature compatible with the [zhenya_bot](https://github.com/numberten/zhenya_bot)
IRC bot and a secondary goal of being easy to configure and deploy.

## Implementation Status ##

GGircd currently implements a subset of IRC protocol as defined in RFC 1459, RFC
2812, and various mailing lists. Notable missing features are any form of
server-to-server connections, banning users and various status messages.

## Installation ##

Eventually there will be two ways of installing GGircd: from source and from a
pre-built binary. Pre-built binaries don't currently exist, so you are stuck
with from source.

### Installing From Source ###

Installing from source will give you the most up to date version of GGircd
available. Every effort is made to keep this bug free, but sometimes mistakes
are made :(

These instructions assumes you already have installed Go. If you haven't
installed Go, you can find easy to follow instructions
[here](http://golang.org/doc/install).

If you already have Go installed, installing GGircd is as easy as:

    go get -u github.com/fimad/ggircd/ggircd
    go install github.com/fimad/ggircd/ggircd

## Running ##

After installing the `ggircd` binary will be placed in `$GOPATH/bin`. This
binary takes the following command line arguments:

  * --config takes a path to the ggircd.conf config file. The default value is
    `/etc/ggircd/ggircd.conf`.

  * --log takes an integer and controls the verbosity of GGircd's logging. The
    higher the number the more verbose the logging. The default value is 3.

## Configuration ##

Before you can run GGircd, you'll need to configure you're server.
Configuration is done in via a JSON encoded file, the default location of which
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
the moment because multi-server IRC networks are not currently supported.

Example: "My IRC Network"

#### Port ####

The port that GGircd should listen on.

Example: 6667

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

If given all connected users will be listed as having this host name.

Example: "i.love.plan9.bell-labs.com"
