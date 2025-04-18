#!/bin/bash
maxPatternSize=4
#fn="./trice.bin.sample3000"
fn="../../docs/TipUserManual.md"
#fn="../../LICENSE.md"
time ( \
go clean -cache
go install ../../...
echo tokenize $fn...
ti_generate -t $fn
echo ti_generate...
ti_generate -u 7 -n 120 -o ../../src/idTable.c -z ${1:-$maxPatternSize} -i $fn.SAMPLES
go clean -cache
go install ../../...
echo "ti_pack...ti_unpack..."
for filename in $fn.SAMPLES/*.txt; do
    ti_pack -v -i $filename
    ti_unpack -i $filename.tip
    #diff -b $filename.tip $filename.tip.untip
    ls -l $filename $filename.tip
    #rm $filename.tip $filename.tip.untip 
done
)


# -i ./trice.bin.sample3000 -z 8 -u 7 -n 127: 55%
# -i ./trice.bin.sample3000 -z 4 -u 7 -n 127: 55%
# -i ./trice.bin.sample3000 -z 3 -u 7 -n 127: 54%
# -i ./trice.bin.sample3000 -z 4 -u 7 -n 126: 50%
# -i ./trice.bin.sample3000 -z 4 -u 7 -n 120: 35%
# -i ./trice.bin.sample3000 -z 8 -u 6 -n 191: 58%
# -i ./trice.bin.sample3000 -z 4 -u 6 -n 191: 56%
# -i ./trice.bin.sample3000 -z 3 -u 6 -n 191: 56%
