#!/usr/bin/env bash

cd ~/code/go-rss-feed && \
	set -o allexport && \
	source dev.env && \
	set +o allexport && \
	./go-rss-feed;
