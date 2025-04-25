#!/bin/bash
maxPatternSize=8
fn=Deutsch.txt
time (
go clean -cache
go install ../../...
echo tokenize $fn...
ti_generate -t $fn
echo ti_generate...
ti_generate -u 7 -n 127 -o ../../src/idTable.c -z ${1:-$maxPatternSize} -i $fn.SAMPLES
go clean -cache
go install ../../...
echo "ti_pack...ti_unpack..."
for filename in $fn.SAMPLES/*.txt; do
    ti_pack   -i $filename
    ti_unpack -i $filename.tip
    cmp --silent $filename $filename.tip.untip || echo $filename.tip $filename.tip.untip files are different
    ls -l $filename $filename.tip
done
filesize="$(du -bc $fn.SAMPLES/*.txt | grep total | cut -f 1)"
tipsize="$(du -bc $fn.SAMPLES/*.txt.tip | grep total | cut -f 1)"
echo filesize sum $filesize, tip size sum $tipsize, rate = $(( 100 * tipsize / filesize ))
rm $filename.tip.untip
)


# -i ./trice.bin.sample3000 -z 8 -u 7 -n 127: 55%
# -i ./trice.bin.sample3000 -z 4 -u 7 -n 127: 55%
# -i ./trice.bin.sample3000 -z 3 -u 7 -n 127: 54%
# -i ./trice.bin.sample3000 -z 4 -u 7 -n 126: 50%
# -i ./trice.bin.sample3000 -z 4 -u 7 -n 120: 35%
# -i ./trice.bin.sample3000 -z 8 -u 6 -n 191: 58%
# -i ./trice.bin.sample3000 -z 4 -u 6 -n 191: 56%
# -i ./trice.bin.sample3000 -z 3 -u 6 -n 191: 56%
# -t Deutsch.txt            -z 8 -u 6 -n 190: 63%
# -t Deutsch.txt            -z 8 -u 7 -n 127: 65%
