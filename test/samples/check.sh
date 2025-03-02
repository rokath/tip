#!/bin/bash
time ( \
go clean -cache && \
go install ../../... && \
ti_generate -v -i ../../docs/TipUserManual.md -z $1 && \
cp idTable.c ../../src && \
go clean -cache && \
go install ../../... && \
ti_pack -v -i ../../docs/TipUserManual.md \
)
