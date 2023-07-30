#!/usr/bin/bash

bun run format; rm dist/ -rf; bun run build; rm ../public/* -rf; cp dist/* ../public/ -r
