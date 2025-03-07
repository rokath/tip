#!/bin/bash
fn="./trice.bin.sample"
maxPatternSize=4
#fn="../../docs/TipUserManual.md"
time ( \
go clean -cache && \
go install ../../... && \
ti_generate -i $fn -z ${1:-$maxPatternSize} && \
cp idTable.c ../../src && \
go clean -cache && \
go install ../../... && \
ti_pack -v -i $fn && \
ti_unpack -i $fn.tip &&\
diff -b $fn $fn.tip.untip && \
rm $fn.tip $fn.tip.untip \
)
