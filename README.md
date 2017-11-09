du-flamegraph
===========================================

Visualize disk usage as flamegraph

Example Flame Graph
------------------------------------------

[![Inception](http://yuroyoro.net/du-flamegraph.svg)](http://yuroyoro.net/du-flamegraph.svg)

Usage:
------------------------------------------

```
NAME:
   du-flamegraph - visualize disk usage as flamegraph

USAGE:
   du-flamegraph [global options] [FILE]

VERSION:
   0.0.0

GLOBAL OPTIONS:
   --width value, -w value    width of image (default 1200) (default: 1200)
   --height value, -h value   height of each frame (default 16) (default: 16)
   --flamegraph-script value  path of flamegraph.pl. if not given, find the script from $PATH
   --out value                distination path of grenerated flamegraph. default is ./du-flamegraph.svg (default: "./du-flamegraph.svg")
   --verbose                  show verbose log
   --version, -v              print the version
```


Installation
------------------------------------------

```
$ go get -u github.com/yuroyoro/du-flamegraph
```

du-flamegraph requires `flamegraph.pl` to render flamegraph svg.
Download the script from `https://github.com/brendangregg/FlameGraph` by using git clone, and set the path of `flamegraph.pl` script to `--flamegraph-script` option, or put it to your $PATH.

