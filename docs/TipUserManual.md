*# TiP - Tiny Packer - User Manual

```diff

+ Compress very small buffers fast and efficient including Zeroes Elemination +
--> Works with several-KByte buffers too but will not compress like established zip tools ❗

```

---
<h2>Table of Contents</h2>
<details><summary>(click to expand)</summary><ol><!-- TABLE OF CONTENTS START -->

<!--
Table of Contents Generation:
* Install vsCode extension "Markdown TOC" from dumeng
* Use Shift-Command-P "markdownTOC:generate" to get the automatic numbering.
* replace "<a name" with "<a id"
* replace "##" followed by 2 spaces with "## "‚
-->

<!-- vscode-markdown-toc -->
* 1. [TiP - Why and How? Initial Situation](#tip---why-and-how?-initial-situation)
  * 1.1. [Framing](#framing)
  * 1.2. [ Very Small Buffer Data Compession](#-very-small-buffer-data-compession)
* 2. [Bytes, Numbers and the TiP Idea](#bytes,-numbers-and-the-tip-idea)
* 3. [ID Table Generation](#id-table-generation)
  * 3.1. [ID Table Generation Algorithm](#id-table-generation-algorithm)
  * 3.2. [ID Table Generation Questions](#id-table-generation-questions)
* 4. [The TiP Algorithm](#the-tip-algorithm)
  * 4.1. [ID Position Table Generation](#id-position-table-generation)
  * 4.2. [ID Position Table Processing](#id-position-table-processing)
  * 4.3. [Packing - Unreplacable Bytes Handling](#packing---unreplacable-bytes-handling)
  * 4.4. [Unpacking](#unpacking)
* 5. [Getting Started](#getting-started)
  * 5.1. [Prerequisites](#prerequisites)
  * 5.2. [Built TipTable Generator `ti_generate`](#built-tiptable-generator-`ti_generate`)
  * 5.3. [Build `ti_pack` and `ti_unpack`](#build-`ti_pack`-and-`ti_unpack`)
  * 5.4. [Try `ti_pack` and `ti_unpack`](#try-`ti_pack`-and-`ti_unpack`)
  * 5.5. [Installation](#installation)
* 6. [ TiP in Action](#-tip-in-action)
    * 6.1. [Training](#training)
    * 6.2. [Test Preparation](#test-preparation)
    * 6.3. [Test Execution](#test-execution)
    * 6.4. [Test Results Interpretation](#test-results-interpretation)
* 7. [Improvement Thoughts](#improvement-thoughts)
  * 7.1. [Additional Indirect Dictionary](#additional-indirect-dictionary)
  * 7.2. [Reserve ID `7f` for Run-Length Encoding](#reserve-id-`7f`-for-run-length-encoding)
  * 7.3. [Minimize Worst-Case Size by using 16-bit transfer units with 2 zeroes as delimiter.](#minimize-worst-case-size-by-using-16-bit-transfer-units-with-2-zeroes-as-delimiter.)
  * 7.4. [Do not remove zeroes in favour of better compression as an option or a separate project.](#do-not-remove-zeroes-in-favour-of-better-compression-as-an-option-or-a-separate-project.)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

</div></ol></details><!-- TABLE OF CONTENTS END -->

---

![./images/logo.png](../images/logo.png)

---

##  1. <a id='tip---why-and-how?-initial-situation'></a>TiP - Why and How? Initial Situation

###  1.1. <a id='framing'></a>Framing

For low level buffer storage or MCU transfers some kind of framing is needed for resynchronization after failure. An old variant is to declare a special character as escape sign and to start each package with it. And if the escape sign is part of the buffer data, add an escape sign there too. Even the as escape sign selected character occurs seldom in the buffer data, a careful design should consider the possibility of a buffer containing only such characters.

[COBS](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) is a newer and much better approach, to achieve framing. It transformes the buffer data containing 256 different characters into a sequence of 255 only characters. That allows to use the spare character as frame delimiter. Usually `0` is used for that.

###  1.2. <a id='-very-small-buffer-data-compession'></a> Very Small Buffer Data Compession

A compression and then [COBS](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) framing would do perfectly. But when it comes to very short buffers, like 4 or 20 bytes, **normal zip code fails** to reduce the buffer size.

To combine the [COBS](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) technique with compression especially for very short buffers, some additional spare characters are needed. That's done with [TCOBS](https://github.com/rokath/tcobs) in a manual coded way, meaning, expected special data properties are reflected in the [TCOBS](https://github.com/rokath/tcobs) code. See the [TCOBS User Manual](https://github.com/rokath/tcobs/blob/master/docs/TCOBSv2Specification.md) for more details.

There is also [smaz](https://github.com/antirez/smaz), but suitable only for text buffers mainly in English.

[RZCOBS](https://github.com/Dirbaio/rzcobs), assumes many zeroes and tries some compression this way.

An adaptive solution would be nice, meaning, not depending on a specific data structure like English text or many integers. [shoco](https://ed-von-schleck.github.io/shoco/) is a way to go but focusses more on strings.

##  2. <a id='bytes,-numbers-and-the-tip-idea'></a>Bytes, Numbers and the TiP Idea

[COBS](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) and [TCOBS](https://github.com/rokath/tcobs) are starting or ending with some control characters and these are linked togeter to distinguish them from data bytes. But there is also an other option.

If there is a buffer of, let's say 20 bytes, we can consider it as a 20-digit number with 256 ciphers. To free some 8 characters for special usage, we could transform the 20 times 256 cipher number into a 21 or 22 times 248 ciphers number. This transformation is possible, but very computing intensive because of many divisions by 248, or a different base number. So this is no solution for small MCUs. But a division by 128 is cheap! If we transform the 256 base into a 128 base, we only need to perform a shift operation for the conversion. This way we get 128 special characters usable for compressing and framing:

* Byte `00` is not used at all. One aim of TiP is, to get rid of all zeroes in the TiP packets to be able to use `00` as a package delimiter.
* Bytes `01` to `7f` are used as pattern IDs. These IDs are used as pattern replacements.
* Before we pack the buffer data, we try to find pattern from the ID table, we can then replace with IDs. See [The TiP Algorithm](./docsTipUserManual.md#the-tip-algorithm) for the how-to-do.
* _Unreplacable_ bytes need a transformation in a way, that no bytes in the range 0-127 remain. That is our tranformation to the 128 base. We simply collect them and do a bit shifting in a way, that no most significant bit is used anymore. The MSBits of the reordered unreplacable bytes are all set to 1 and so we have only bytes `80` to `ff` left.

The `ti_unpack` then sees bytes `01` to `7f` and knows, that these are IDs, intermixed with bytes `80` to `ff` and knows, that the 7 least significant bits are the unreplacable bytes. The byte places are containing the position informtion for the unreplacable bytes.

##  3. <a id='id-table-generation'></a>ID Table Generation

###  3.1. <a id='id-table-generation-algorithm'></a>ID Table Generation Algorithm

* We create a bunch of test files with data similar to those we want to pack in the future.
  * `ti_generate` takes a single file and treats it as a separate sample buffer.
  * Also a folder name is accepted and all files are processed then.
* We assume a longest pattern, like N=8 for example.
  * `ti_generate` accepts it as parameter.
  * The longest possible pattern is 255 bytes long.
  * For very short buffers, 2 to 8 bytes as maximum is recommended as max size N.
* We take the first N bytes of some sample data and move that window in 1-byte steps over the sample data and build a histogram over all found pattern and their occurances count.
* The same is done with all smaller pattern sizes, ergo N, ..., 3, 2. Not interesting are 1-byte patterns, because their replacement by an ID gives no compression effect.
* The 127 most often occuring pattern are sorted by descending size and are used to create the file `idTable.c`.

###  3.2. <a id='id-table-generation-questions'></a>ID Table Generation Questions

* It is not clear, if the this way created ID table is optimal. Especially, when pattern are sub-pattern of other patterns. That is easily the case with sample data containing the same bytes in longer rows.
* Also it could make sense to use the length of a pattern as weigth. If, for example a 5-bytes long pattern occurs 100 times and a 2-bytes long pattern exists 200 times in the sample data - which should get preceedence to get into the ID table? My guess is, to multiply the pattern length with its occureances count gives a good approximation.
* We could also just determine all pattern from 2 to N bytes length and then go byte by byte through the sample data and increment for each byte the pattern counter for the pattern containing this byte on the right place.
* It could make sense, to build several ID tables and then measure how good the packing is with the different tables.

<!--
``` c
// 1 2 3 4 -> 12:1 23:1 34:1 123:1 234:1 1234:1 -> weighted: 12:2 23:2 34:2 123:3 234:3 1234:4
//         -> 12:0 23:- 34:0 123:0 234:0 1234:1 -> weighted: 12:0 23:- 34:0 123:0 234:0 1234:4
// 1 1 1 1 -> 11:3           111:2       1111:1 -> weighted: 11:6           111:6       1111:4
//         -> 11:2           111:1       1111:1 -> weighted: 11:4           111:3       1111:4
```

-####  3.2.1. <a id='10-bytes:-123456789a'></a>10 bytes: 123456789a 

p   | m   | length | pattern                         | no pattern    | byte usage count | equ. factor
----|-----|--------|---------------------------------|---------------|------------------|------------
10  | 0   | 1er    | 1 ... a                         |               | 1                | 10/1
9   | 1   | 2er    | 12 23 ... 9a                    | a1            | 2                | 10/2
8   | 2   | 3er    | 123 234 ... 89a                 | 9a1 a12       | 3                | 10/3
... | ... | ...    | ...                             | ...           | ...              | ...
4   | 6   | 7er    | 1234567 2345678 3456789 456789a | 56789a1...    | 7                | 10/7
3   | 7   | 8er    | 12345678 23456789 3456789a      | 456789a1...   | 8                | 10/8
2   | 8   | 9er    | 123456789 23456789a             | 3456789a1...  | 9                | 10/9
1   | 9   | 10er   | 123456789a                      | 23456789a1... | 10               | 10/10

1234

count           | balance factor | hist                | reduced             | \*length
----------------|----------------|---------------------|---------------------|---------------
1:1,2:1,3:1,4:1 | \*4/1          | all:4               | 1:3,2:2,3:2,4:3     | =
12:1,23:1,34:1  | \*4/2          | all:2               | 12:1,23:0,34:1      | 12:2,23:0,34:2
123:1, 234:1    | \*4/3          | 123:1.333,234:1.333 | 123:0.333,234:0.333 | 123:1,234:1
1234:1          | \*4/4          | 1234:1              | 1234:1              | 1234:4

table: 1234, 12, 34, 123, 234, 23

1111

count       | balance factor | hist      | reduced   | \*length
------------|----------------|-----------|-----------|---------
1:4         | \*4/1          | 1:16      | 1:4       | 1:4
11:3        | \*4/2          | 11:6      | 0.666     | 1.333
111:2       | \*4/3          | 111:2.666 | 111:0.666 | 111:2
1111:1\*4/4 | 1111:1         | 1111:1    | 1111:4    | 1111:4

table: 1111 111 11

aa0000bb0000cc maxSize 4
  ----  ----

pattern | count | balance factor | balanced      | remark
--------|-------|----------------|---------------|---------------
0000    | 1     | 4/2            | 4000/2 = 2000 | gets negative!
aa0000  | 1     | 4/3            | 4000/3 = 1333 | contains 0000
0000bb  | 1     | 4/3            | 4000/3 = 1333 | contains 0000
bb0000  | 1     | 4/3            | 4000/3 = 1333 | contains 0000
0000cc  | 1     | 4/3            | 4000/3 = 1333 | contains 0000
-->

##  4. <a id='the-tip-algorithm'></a>The TiP Algorithm

###  4.1. <a id='id-position-table-generation'></a>ID Position Table Generation

* Step byte by byte thru the `slen` `src` buffer and check if a pattern from the (into `ti_pack` and `ti_unpack`) compiled [./src/idTable.c](../src/idTable.c) matches and build a sorted ID position table. Its max length is slen-1. Example for file 43.bin (see below):

IDPositionTable:
idx | ID  | pos | ASCII
----|-----|-----|-------
0   | 52  | 4   | '    '
1   | 95  | 4   | '   '
2   | 127 | 4   | '  '
3   | 51  | 5   | '   ■'
4   | 95  | 5   | '   '
5   | 127 | 5   | '  '
6   | 43  | 6   | '  ■ '
7   | 94  | 6   | '  ■'
8   | 127 | 6   | '  '
9   | 35  | 7   | ' ■  '
... | ... | ... | ...

* The pattern in the IDPositionTable:
  * are a subset of the [./src/idTable.c](../src/idTable.c) pattern 
  * can occur several times at different positions, example: ID 127 at pos 4, 5 and 6
  * can overlap, example: IDs 52, 95, 127, 51, 55 all cover position 5

###  4.2. <a id='id-position-table-processing'></a>ID Position Table Processing

* To build a TiP packet, many different ID position sequences are possible, maybe interrupted by some _unreplacable_ bytes. The TiP packer starts creating a full `srcMap` containing all possible paths. For that it traverses the (by incrementing position sorted) IDPositionTable and checks, if the current ID position is appenable to any paths. If so, these paths are forked and the ID position is appended to the fork. That fork is needed, because the same path is extendable with different ID positions. If the current ID position did not fit to any path, a new path is created. After processing an ID position, a new path may exist or some paths have been foked and the forked paths are extended with this ID position. Before going to the next ID position from the IDPositionTable, obsolete `srcMap` paths are deleted. Obsolete are paths, if their limit plus the maximum pattern size is smaller than biggest existing path limit. Obsolete paths are too those path, which have an equal limit but wuld result in a bigger (partial) TiP packet. Even if they would result in an equal TiP packet size, it is only one of them needed for futher ID position provessing.
* When the PositionTable was processed completely, a few paths are remaining. A path, which would result in the smallest TiP packet is selected to create the TiP packet.

###  4.3. <a id='packing---unreplacable-bytes-handling'></a>Packing - Unreplacable Bytes Handling

The selected path covers no, some or all bytes with ID pattern. Bytes not covered, are unreplacable bytes.
All unreplacable bytes are collected into one separate buffer. N unreplacable bytes occupy N\*8 bits. These bits are distributed onto N\*8/7 7-bit bytes, all getting the MSBit set to avoid zeroes and to distinguish them later from the ID bytes. In fact we do not change these N\*8 bits, we simply reorder them slightly. This bit reordering is de-facto the number transformation to the base 128, mentioned above. By setting the most significant bits, also is guarantied, that no `00` bytes exist anymore.

Next all found patterns are replaced with their IDs, which all have MSBit=0. The unreplacable bytes are replaced with the bit-reordered unreplacable bytes, having MSBit=1. The bit-reordered unreplacable bytes fill the wholes between the IDs.

###  4.4. <a id='unpacking'></a>Unpacking

On the receiver side all bytes with MSBit=0 are identified as IDs and are replaced with the patterns they stay for. All bytes with MSBit=1 are carying the unreplacable bytes bits. These are ordered back to restore the unreplacable bytes which fill the wholes between the patterens then.

##  5. <a id='getting-started'></a>Getting Started

###  5.1. <a id='prerequisites'></a>Prerequisites

* For now install [Go](https://golang.org/) to easily build the executables.
* You need some files containing typical data you want to pack and unpack.
  * Just to try out TiP, you can use a folder containing ASCII texts.

###  5.2. <a id='built-tiptable-generator-`ti_generate`'></a>Built TipTable Generator `ti_generate`

* `cd ti_generate && go build -ldflags "-w" ./...`
* Run `ti_generate` on the data files to get an `idTable.c` file.

###  5.3. <a id='build-`ti_pack`-and-`ti_unpack`'></a>Build `ti_pack` and `ti_unpack`

* Copy the generated `idTable.c` file into the `src` folder.
* Run `go clean -cache`.
* Run `go build ./...` or `go install ./...`.

###  5.4. <a id='try-`ti_pack`-and-`ti_unpack`'></a>Try `ti_pack` and `ti_unpack`

* Run `ti_pack -i myFile -v` to get `myFile.tip`.
* Run `ti_unpack -i myFile.tip -v` to get `myFile.tip.untip`.
* `myFile` and `myFile.tip.untip` are expected to be equal.

###  5.5. <a id='installation'></a>Installation

* Add `src` folder to your project and compile.
* `pack.h` and `unpack.h` is the user interface.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

##  6. <a id='-tip-in-action'></a> TiP in Action

> **Follow these steps with your own data, to see quickly if it makes sense for your project.**

####  6.1. <a id='training'></a>Training 

* Find the 127 most common pattern in some sample data, similar to the real data expected later, and assign the IDs 1-127 to them. This is done once offline and the generated ID table gets part of the tiny packer code as well as for the tiny unpacker code. For that task a generator tool `ti_generate` was build.
* Sample data specific result: [./src/idTable.c](../src/idTable.c)

> 🛑 The current ID table generation might not give an optimal result and is matter of further investigation❗

* Training data example (binary [Trice](https://github.com/rokath/trice) output file)

```bash
$ xxd -g 1 trice.bin.sample
00000000: 3d 73 2a 00 3e 73 2b 04 ff ff ff ff 3f 73 2c 08  =s*.>s+.....?s,.
00000010: ff ff ff ff fe ff ff ff 40 73 2d 0c ff ff ff ff  ........@s-.....
00000020: fe ff ff ff fd ff ff ff 41 73 2e 10 ff ff ff ff  ........As......
00000030: fe ff ff ff fd ff ff ff fc ff ff ff 42 73 2f 14  ............Bs/.
00000040: ff ff ff ff fe ff ff ff fd ff ff ff fc ff ff ff  ................
00000050: fb ff ff ff 43 73 30 18 ff ff ff ff fe ff ff ff  ....Cs0.........
00000060: fd ff ff ff fc ff ff ff fb ff ff ff fa ff ff ff  ................
...

$ ti_generate.exe -i trice.bin.sample -z 4 -v -o ../../src/idTable.c
go clean -cache && go install ../../...
```

* The maximum allowed pattern size `-z 4` has influence on the TiP pack results and the best value depends on the data. 

####  6.2. <a id='test-preparation'></a>Test Preparation

* Create some sample files: In this example, the messages are starting with `3d`, `3e`, `3f`, `40`, `41`, `42`, `43`, ... (see [TriceUserManual # Package Format](https://github.com/rokath/trice/blob/master/docs/TriceUserManual.md#package-format)). So we cut out a few single binary Trice messages. 


```bash
cat trice.bin.sample | dd bs=1 skip=0 count=4 > 3d.bin
$ xxd -g1 3d.bin
00000000: 3d 73 2a 00                                      =s*.
# ID -----^^-^^
# cycle --------^^
# payloadsize -----^^
$ cat trice.bin.sample | dd bs=1 skip=4 count=8 > 3e.bin
$ xxd -g1 3e.bin
00000000: 3e 73 2b 04 ff ff ff ff                          >s+.....
# ID -----^^-^^
# cycle --------^^
# payloadsize -----^^
# payload ------------^^-^^-^^-^^
$ cat trice.bin.sample | dd bs=1 skip=12 count=12 > 3f.bin
$ xxd -g1 3f.bin
00000000: 3f 73 2c 08 ff ff ff ff fe ff ff ff              ?s,.........
# ID -----^^-^^
# cycle --------^^
# payloadsize -----^^
# payload ------------^^-^^-^^-^^-^^-^^-^^-^^
$ cat trice.bin.sample | dd bs=1 skip=24 count=16 > 40.bin
$ xxd -g1 40.bin
00000000: 40 73 2d 0c ff ff ff ff fe ff ff ff fd ff ff ff  @s-.............
# ID -----^^-^^
# cycle --------^^
# payloadsize -----^^
# payload ------------^^-^^-^^-^^-^^-^^-^^-^^-^^-^^-^^-^^
$ cat trice.bin.sample | dd bs=1 skip=40 count=20 > 41.bin
$ xxd -g1 41.bin
00000000: 41 73 2e 10 ff ff ff ff fe ff ff ff fd ff ff ff  As..............
00000010: fc ff ff ff                                      ....

$ cat trice.bin.sample | dd bs=1 skip=60 count=24 > 42.bin
$ xxd -g1 42.bin
00000000: 42 73 2f 14 ff ff ff ff fe ff ff ff fd ff ff ff  Bs/.............
00000010: fc ff ff ff fb ff ff ff                          ........

$ cat trice.bin.sample | dd bs=1 skip=84 count=28 > 43.bin
$ xxd -g1 43.bin
00000000: 43 73 30 18 ff ff ff ff fe ff ff ff fd ff ff ff  Cs0.............
00000010: fc ff ff ff fb ff ff ff fa ff ff ff              ............
```

####  6.3. <a id='test-execution'></a>Test Execution

```bash
$ ti_pack.exe -v -i 3d.bin
file size 4 changed to 1 (rate 25 percent)

$ ti_pack.exe -v -i 3e.bin
file size 8 changed to 6 (rate 75 percent)

$ ti_pack.exe -v -i 3f.bin
file size 12 changed to 8 (rate 66 percent)

$ ti_pack.exe -v -i 40.bin
file size 16 changed to 9 (rate 56 percent)

$ ti_pack.exe -v -i 41.bin
file size 20 changed to 10 (rate 50 percent)

$ ti_pack.exe -v -i 42.bin
file size 24 changed to 11 (rate 45 percent)

$ ti_pack.exe -v -i 43.bin
file size 28 changed to 12 (rate 42 percent)
```

####  6.4. <a id='test-results-interpretation'></a>Test Results Interpretation

If the real data are similar to the training data, an average packed size of about 50\% is expected.

##  7. <a id='improvement-thoughts'></a>Improvement Thoughts

###  7.1. <a id='additional-indirect-dictionary'></a>Additional Indirect Dictionary

For example we can limit direct pattern count to 120 (instead of 127) and use their order in such a way:

* ID 1...120                    -> at least 2-bytes pattern <= 50% compressed
* ID 121 followed by id 1...255 -> at least 3-bytes pattern <= 67% compressed
* ID 122 followed by id 1...255 -> at least 3-bytes pattern <= 67% compressed
* ID 123 followed by id 1...255 -> at least 3-bytes pattern <= 67% compressed
* ID 124 followed by id 1...255 -> at least 3-bytes pattern <= 67% compressed
* ID 125 followed by id 1...255 -> at least 3-bytes pattern <= 67% compressed
* ID 126 followed by id 1...255 -> at least 3-bytes pattern <= 67% compressed
* ID 127 reserved

This allows 120 at least 2-bytes pattern and 1525 longer pattern.

<!--
* the MSBit = 0|1 after a first ID 121-126 are the indiret table indices
* the MSBit = 1   after ID 1...120 are the unreplacable (bit-shfted) bytes
* the MSBit = 0 not after ID 121-126 are the direct table indices
-->

 On unpacking:

* START 
  * Next byte with MSBit=1 is unreplaceable, goto START
  * Next byte=1...120 is direct 2-byte pattern ID, goto START
  * Next byte=121...126 is followed by indirect pattern ID, goto START
  * Next byte=127 is followed by runlength code, goto START

To implement add to [tipConfig.h](../src.config/tipConfig.h):

```C
//! INDIRECT_DICTIONARY_COUNT adds a number of indirect dictionaries.
//! An indirect dictionary needs a 2-byte reference and therefore only pattern with at least 3 bytes make sense there.
//! Each indirect dictionary adds 255 >= 2-bytes reference pattern and reduces the direct pattern space by one.
//! The max possible value is 127, but that would not allow any direct references at all.
// Values making sense are probably in the range 0...10. The optimum depend on the kind of data.
#define INDIRECT_DICTIONARY_COUNT 0 
```

> **Consideration:** Promizing

###  7.2. <a id='reserve-id-`7f`-for-run-length-encoding'></a>Reserve ID `7f` for Run-Length Encoding

* Example:

| ID sequence                              | Meaning                                                      |
|------------------------------------------|--------------------------------------------------------------|
| ID `7F` + count `1...15`                 | 3 to 17 zeroes                                               |
| ID `7F` + count `16...24`                | 3 to 11 FFs                                                  |
| ID `7F` + count `25...63` + byte `XX`!=0 | 4 to 42 `XX`s, `XX` is any non-zero byte, all `XX` are equal |
| ID `7F` + `64...255` + `?`               | reserved                                                     |

* The tiny unpack routine first regards all bytes with MSBit=0 as IDs.
* The ID `7F` is followed by a count byte and optional other bytes. These are regarded as part of this ID too during TiP package interpretation.
  * The count is guarantied not to be zero and also some optional additional bytes.

To implement add to [tipConfig.h](../src.config/tipConfig.h):

```C
#define RUN_LENGTH_ID 127
//! TODO: define ranges here
```

> **Consideration:** Possible, but currenly no aim.

###  7.3. <a id='minimize-worst-case-size-by-using-16-bit-transfer-units-with-2-zeroes-as-delimiter.'></a>Minimize Worst-Case Size by using 16-bit transfer units with 2 zeroes as delimiter

* If data are containing no ID table pattern at all, they are getting bigger by the factor 8/7 (+14\%). That is a result of treating the data in 8 bit units (bytes).
* If we change that to 16-bit units, by accepting an optional padding byte, we can reduce this increasing factor to 16/15 (+7\%).
* We still have IDs 1-127
* An existing ID 127 just tells if there is a padding byte in the unreplacable data.
* When unpacking, the first set MSBit tells that this byte and the next are unreplaceable. So we get N 16-bit groups of unreplacable data.
* BUT we need 2 frame delimiter bytes then!

> **Consideration:** Not a good idea, because we get other overhead.

###  7.4. <a id='do-not-remove-zeroes-in-favour-of-better-compression-as-an-option-or-a-separate-project.'></a>Do not remove zeroes in favour of better compression as an option or a separate project

[smaz](https://github.com/antirez/smaz):

* IDs 0...253 are coding 254 >= 2-bytes patteren
* ID 254 -> next byte is unreplacable
* ID 255 -> next byte is a count of following 2...257 unreplacable bytes

Modificate [smaz](https://github.com/antirez/smaz) and add indirect indices:

* IDs 0...239 are coding 240 >= 2-bytes patteren
* ID 240 -> next byte is one of 256 indicies in indirect table 0
* ...
* ID 249 -> next byte is one of 256 indicies in indirect table 9
* ID 250 -> reserved for run-length code
* ID 251 -> next byte is unreplacable
* ID 252 -> next 2 bytes are unreplacable
* ID 253 -> next 3 bytes are unreplacable
* ID 254 -> next 4 bytes are unreplacable
* ID 255 -> next byte is count of 5...231 unreplacable bytes

This allows 2560 additional pattern for the price 14 less 2-bytes pattern and the need for 2 bytes for the 2560 additional patterns. The details could be configurable.

> **Consideration:** Interesting extension but we want elemminate zeroes in one shot to keep the overall overhead small. This could make sense to improve SMAZ in an universal way, by providing a pattern table generator, which could be practically the same. The pattern table generator could get an option to use some internet data for the table generation.

<!--

https://jwakely.github.io/pkg-gcc-latest/

```bash
wget --content-disposition https://kayari.org/gcc-latest/gcc-latest.deb
cd ~/Downloads
sudo dpkg -i gcc-latest_15.0.0-20250112gitf4fa0b7d493a.deb
cd /opt
ls -l # gcc-latest
cd /etc/profile.d # ls -l
sudo echo export PATH=/opt/gcc-latest/bin/:$PATH > # /etc/profile.d/gccpath.go 
```
-->

<p align="right">(<a href="#readme-top">back to top</a>)</p>
