GGircd
=======

GGircd is a minimalistic IRC server with a primary goal of being easy to
configure and deploy.

Implementation Status
---------------------

GGircd currently implements a subset of IRC protocol as defined in RFC 1459, RFC
2812, and various mailing lists. Notable missing features are any form of
server-to-server connections, banning users and user modes

Configuration
-------------

Configuration options are specified in a JSON encoded file. The default location
is `/etc/ggircd/ggircd.conf`. An example configuration file is given below:

    {
      "Name": "ggirc",
      "Network": "Greatest Guys",
      "Port": 6667,
      "MOTD": "/etc/ggircd/motd",

      "DefaultChannelMode": "t"
    }
