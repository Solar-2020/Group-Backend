#!/bin/sh

branch=$1
# Remove replace ...
sed -i 's/replace/\/\/replace/' go.mod

# Update project mods
#go get -u github.com/Solar-2020/GoUtils
#go get -u "github.com/Solar-2020/Authorization-Backend@$branch"