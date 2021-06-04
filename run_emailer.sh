#!/bin/sh

set -e

# Get latest release
curl -s https://api.github.com/repos/jason-adam/go-rss-feed/releases/latest |
    grep "browser_download_url" |
    cut -d '"' -f 4 |
    wget -qi -

# Unzip and move to repo
tar -xvf go-rss-feed-linux-arm32.tar.gz &&
    mv go-rss-feed-linux-arm32 "$HOME"/code/go-rss-feed &&
    cd "$HOME"/code/go-rss-feed &&
    ./go-rss-feed-linux-arm32

# Cleanup
rm "$HOME"/go-rss-feed-linux-arm32.tar.gz
rm "$HOME"/code/go-rss-feed/go-rss-feed-linux-arm32
