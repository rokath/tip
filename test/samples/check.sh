#!/bin/bash
maxPatternSize=12
fn="./trice.bin.sample"
#fn="../../docs/TipUserManual.md"
fn="../../LICENSE.md"
time ( \
go clean -cache && \
go install ../../... && \
ti_generate -u 6 -n 186 -o ../../src/idTable.c -z ${1:-$maxPatternSize} -i $fn && \
go clean -cache && \
go install ../../... && \
ti_pack -v -i $fn && \
ti_unpack -i $fn.tip &&\
diff -b $fn $fn.tip.untip && \
rm $fn.tip $fn.tip.untip \
)
