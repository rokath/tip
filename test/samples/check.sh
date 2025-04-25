#!/bin/bash
maxPatternSize=6
fn="./trice.bin.sample3000"
#fn="../../docs/TipUserManual.md"
#fn="../../LICENSE.md"
time ( \
go clean -cache && \
go install ../../... && \
echo ti_generate...
ti_generate -u 6 -n 180 -o ../../src/idTable.c -z ${1:-$maxPatternSize} -i $fn && \
go clean -cache && \
go install ../../... && \
echo "ti_pack...(can take a while for bigger files but that is TiP not made for)"
ti_pack -v -i $fn && \
echo ti_unpack...
ti_unpack -i $fn.tip &&\
diff -b $fn $fn.tip.untip && \
ls -l $fn.tip $fn.tip.untip \
rm $fn.tip $fn.tip.untip \
)


# -i ./trice.bin.sample3000 -z 8 -u 7 -n 127: 55%
# -i ./trice.bin.sample3000 -z 4 -u 7 -n 127: 55%
# -i ./trice.bin.sample3000 -z 3 -u 7 -n 127: 54%
# -i ./trice.bin.sample3000 -z 4 -u 7 -n 126: 50%
# -i ./trice.bin.sample3000 -z 4 -u 7 -n 120: 35%
# -i ./trice.bin.sample3000 -z 8 -u 6 -n 191: 58%
# -i ./trice.bin.sample3000 -z 4 -u 6 -n 191: 56%
# -i ./trice.bin.sample3000 -z 3 -u 6 -n 191: 56%
# -i ./trice.bin.sample3000 -z 6 -u 6 -n 180: 25%
