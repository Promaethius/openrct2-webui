#!/bin/bash
x11vnc -forever -create &
openrct2 "$@"