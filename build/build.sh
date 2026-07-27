#!/bin/sh

BUILD_DIR=$(cd `dirname $0`; pwd)
SRC_DIR=$BUILD_DIR/../src
OUTPUT_DIR=$BUILD_DIR/../output

mkdir -p $OUTPUT_DIR
rm -rf $OUTPUT_DIR/*
cp $BUILD_DIR/init.sh $OUTPUT_DIR
cp $BUILD_DIR/start.sh $OUTPUT_DIR
cp $BUILD_DIR/sudoer.config $OUTPUT_DIR
cp $BUILD_DIR/Dockerfile $OUTPUT_DIR
cp -r $SRC_DIR/conf $OUTPUT_DIR

dos2unix $OUTPUT_DIR/*
dos2unix $OUTPUT_DIR/conf/*

cd $SRC_DIR
# GOOS=linux GOARCH=arm64 go build --tags noCoreGo -o $OUTPUT_DIR/gids
GOOS=linux GOARCH=arm64 go build -o $OUTPUT_DIR/gids
