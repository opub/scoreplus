#!/bin/bash

MAJOR=0
MINOR=$(date +%y%m)
PATCH=$(date +%d%H%M)

govvv build -pkg $(go list ./util) -version $MAJOR.$MINOR.$PATCH -o ./bin/scoreplus
