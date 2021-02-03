#!/usr/bin/env sh

set -e;

curl -s https://api.github.com/repos/jason-adam/go-rss-feed/releases/latest \
    | grep "browser_download_url" \
    | cut -d '"' -f 4 \
    | wget -qi -;

tar -xvf go-rss-feed-linux-arm32.tar.gz && \
    set -o allexport && \
    source ~/envs/go-rss-feed.env && \
    set +o allexport && \
    ./go-rss-feed-linux-arm32;
