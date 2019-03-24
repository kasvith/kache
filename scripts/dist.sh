#!/bin/bash

echo 'Running cross platform build'
mage kachecrossbuild

mkdir -p dist

FILES=build/*
for f in $FILES
do
  filename=$(basename -- "$f")
  extension="${filename##*.}"
  artifact="${filename%.*}"

  echo "Creating archive for $artifact"
  echo $f
  echo $filename
  mkdir -p dist/$artifact/bin
  cp -v $f dist/$artifact/bin
  cp -r LICENSE config dist/$artifact
  cd dist/$artifact
  zip -r -v ../$artifact.zip *
  cd ../..
  echo
done

echo "Packaging done"
