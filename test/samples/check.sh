#!/bin/bash
time ( \
go clean -cache && \
go install ../../... && \
ti_generate -v -i ./trice.bin.sample -z $1 && \
cp idTable.c ../../src && \
go clean -cache && \
go install ../../... && \
ti_pack -v -i ./trice.bin.sample \
)
