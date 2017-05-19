#!/bin/bash

if [ ! -d "env" ]; then
  echo "Create virtual environment..."
  mkdir -p env/bin
  cp ./.bootstrap/.activate env/bin/activate
fi
source ./env/bin/activate

echo "Install dependencies..."
gopm get -v -u -g

echo "Install assert"
gopm get -u -v -g github.com/stretchr/testify/assert

echo "Finished. Now go hack!"
