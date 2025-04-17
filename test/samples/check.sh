#!/bin/bash
maxPatternSize=4
fn="./trice.bin.sample"
#fn="../../docs/TipUserManual.md"
#fn="../../LICENSE.md"
time ( \
go clean -cache && \
go install ../../... && \
ti_generate -u 7 -n 120 -o ../../src/idTable.c -z ${1:-$maxPatternSize} -i $fn && \
go clean -cache && \
go install ../../... && \
ti_pack -v -i $fn && \
ti_unpack -i $fn.tip &&\
diff -b $fn $fn.tip.untip && \
rm $fn.tip $fn.tip.untip \
)

# -z 8 -u 7 -n 127: 55%
# -z 4 -u 7 -n 127: 55%
# -z 3 -u 7 -n 127: 54%

# -z 4 -u 7 -n 126: 50%
# -z 4 -u 7 -n 120: 35%

# -z 8 -u 6 -n 191: 58%
# -z 4 -u 6 -n 191: 56%
# -z 3 -u 6 -n 191: 56%
