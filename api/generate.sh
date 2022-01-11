#!/bin/sh

# needs buf: https://docs.buf.build/installation#github-releases

buf mod update
buf generate