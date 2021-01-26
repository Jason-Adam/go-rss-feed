#!/usr/bin/env bash

zip -r go-rss-feed.zip . \
    -x "dev.env" \
    -x ".git/**" \
    -x ".github/**" \
    -x ".gitignore" \
    -x ".DS_Store"
