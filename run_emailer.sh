#!/bin/sh

set -e

curl -s https://api.github.com/repos/jason-adam/go-rss-feed/releases/latest |
    grep "browser_download_url" |
    cut -d '"' -f 4 |
    wget -qi -

tar -xvf go-rss-feed-linux-arm32.tar.gz &&
    mv go-rss-feed-linux-arm32 "$HOME"/code/go-rss-feed &&
    cd "$HOME"/code/go-rss-feed &&
    ./go-rss-feed-linux-arm32
