#!/usr/bin/bash

bun run format; rm dist/ -rf; bun run build
