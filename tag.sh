#!/bin/bash

tag=$1
message=$1

if [ ! -z "$2" ]
then
    message=$2
fi

git push --delete origin $tag
git tag -d $tag
git tag -a $tag -m "$message"
git push origin $tag