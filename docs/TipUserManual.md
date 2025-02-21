# TiP - Tiny Packer - User Manual

(work in progress)

```diff

+ Compress small buffers fast and efficient with Zeroes Elemination +
--> Works with big buffers too but will not compress like establisched zip tools ❗

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
<!-- vscode-markdown-toc-config
    numbering=true
    autoSave=true
    /vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

</div></ol></details><!-- TABLE OF CONTENTS END -->

---

![./images/logo.png](../images/logo.png)

---

<!--
12312

     123.   12 23. 31
--------------
12   22.    2 
23.   21.      1
31.  1 2.           1
123  111
231. 111
312. 111

--> 

## TiP - Why and How?

### Current Situation

For low level buffer storage or MCU transfers some kind of framing is needed for resynchronization after failure. An old variant is to declare a special character as escape sign and to start each package with it. And if the escape sign is part of the buffer data, add an escape sign there too. Even the as escape sign selected character occurs seldom in the buffer data, a careful design should consider the possibility of a buffer containing only such characters.

[COBS](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) is a newer and much better approach, to achieve framing. It transformes the buffer data containing 256 different characters into a sequence of 255 only characters. That allows to use the spare character as frame delimiter. Usually `0` is used for that.

To combine the COBS technique with compression especially for very short buffers, some additional spare characters are needed. That's done with [TCOBS](https://github.com/rokath/tcobs) more or less in a "manual" way, meaning, expected special data properties are reflected in the TCOBS code. See the [TCOBS User Manual](https://github.com/rokath/tcobs/blob/master/docs/TCOBSv2Specification.md) for more details.

There is also [SMAZ](https://github.com/antirez/smaz), but suitable only for text buffers mainly in English.

The TiP approach is more generic, meaning, not depending on a specific data structure but expecting any data structure.

### Bytes and Numbers

COBS and TCOBS are starting or ending with some control characters and these are linked togeter to distinguish them form data bytes. But there is also an other otion.

If there is a buffer of, let's say 20 bytes, we can consider it as a 20-digit number with 256 ciphers. To free like 8 characters for special usage, we could transform the 20 times 256 cipher number into a 21 or 22 times 248 ciphers number. This transformation is possible, but very computing intensive because of many divisions by 248, or a different base number. So this is no solution for small MCUs. But a division by 128 is cheap! If we transform the 256 base into a 128 base, we only need to perform a shift operation for the conversion. This way we get 128 special characters usable for compressing and framing.

### The TiP Idea

#### Packing

Find the 127 most common pattern in similar sample data and assign the IDs 1-127 to them. This is done once offline and the generated ID table gets part of the tiny packer code as well as for the tiny unpacker code. For that task a generator tool was build.

At runtime the actual buffer is searched for matching patterns from the ID table beginning with the longest ones. All these found patterns get replaced by the IDs later. All unreplacable bytes are collected into one separate buffer. N unreplacable bytes occupy N\*8 bits. These are distributed onto N\*8/7 7-bit bytes, all having the MSBit set. In fact we do not change these N\*8 bits, we simply reorder them slightly. This way we have no zeros in the result and we can distinguish bytes carrying unreplacable bits from ID bytes, which replaced patterns. This bit reordering, to free the MSBit, is de-facto the number transformation to the base 128, mentioned above.

After replacing all found patterns with their IDs, which all have MSBit=0, the unreplacable bytes are replaced with the bit-reordered unreplacable bytes, having MSBit=1.

Now the packing is done and no zeroes are existing anymore.

#### Unpacking

On the receiver side all bytes with MSBit=0 are identified as IDs and are replaced with the patterns they stay for. All bytes with MSBit=1 are carying the unreplacable bytes bits. These are reordered to restore the unreplacable bytes. These are then used to fill the wholes between the patterens.

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

With 2 reserved bytes, zA and fA is this possible:
* 1: 00                            Z1
* 2: 00 00                         Z2
* 3: 00 00 00                      Z3
* 4: 00 00 00 00                   Z1 zA
* 5: 00 00 00 00 00                zA Z1
* 6: 00 00 00 00 00 00             Z2 zA
* 7: 00 00 00 00 00 00 00          za Z2
* 8: 00 00 00 00 00 00 00 00       Z3 zA
* 9: 00 00 00 00 00 00 00 00 00    zA Z3
* 
* 1: FF                            FF
* 2: FF FF                         F2
* 3: FF FF FF                      F3
* 4: FF FF FF FF                   F4
* 5: FF FF FF FF FF                F2 fA
* 6: FF FF FF FF FF FF             fA F2
* 7: FF FF FF FF FF FF FF          F3 fA
* 8: FF FF FF FF FF FF FF FF       fA F3
* 9: FF FF FF FF FF FF FF FF FF    F4 fA
* A: FF FF FF FF FF FF FF FF FF FF fA F4

### How to reduce short buffers

* Lets imagine to have some reserved bytes like 00, Z1, Z2, Z3, Z4, F1==FF, F2, F3, F4
* 00 we want eleminate
* We replace 00...00 00 00 00 with Z1...Z4
* We replace 5...21 00 with Z1 Z1...Z4 Z4
* We replace 5...21 FF with F1 F1...F4 F4
* What if we have more than 21 00 or FF in a row? Probabli that is ok.
* We extract the remaining bytes. Example: x4 x3 00 00 x2 FF FF FF x1 x0, so we have x4 x3 x2 x1 x0
* x4...x0 is a 5 digit number N using 256 ciphers. We need to translate N into yn...y0 with 128 ciphers.
* This costs computing effort: x4*256^4 + ... x0*256^0
* N0/128 = N>>7 = yn
* N0-yn = N1 ... N1/128 y(n-1) ...
* We put yn...y0 into the place of x4...x0 and append the ciphers up to n.
* In general we translate 40 bit (x0...x4) into 42 bit (yn...y0), so yn is y5
* If we say all shortcut bytes have a MSB 0 and all y have a MSB 1 we can
* Use 127 schortcut bytes and replace common pattern with shortcut bytes.
* Then we take the x4...x0 and translate to y5...0 by just bit shifting
* No we have a sewuence with mixed MSB 0 or 1.
* To decompress we change y5...y0 (the bytes with MSB1) into x4...x0.
* We replce all shortcuts (the bytes with MSB0) and we are done.
* 00 is not used at all.
* 1...127 are shortcut bits.
* We take binary data and automatically determine a good shortcut set.
* The shortcut set is de-facto a pattern list.


 tiPack converts in to out and returns final lenth.

 Algorithm:
 * Start with tip list longest pattern and try to find a match inside in.
 * If a longest possible pattern match was found we have afterwards:
   - preBytes match postBytes
   - start over with preBytes and postBytes and so on until we cannot replace any pattern anymore
   - Then we have: xx xx p7 x p0 p0 xx xx xx for example, where pp are any pattern replace bytes,
     which all != 0 and all have MSB==0. The xx are the remaining bytes, which can have any values.
     Of course we need the position information like:

 (A) in:  xx xx xx xx xx xx xx xx xx xx xx xx xx xx xx xx
 (B) in:  xx xx P7 P7 P7 P7 xx P0 P0 P0 P0 P0 P0 xx xx xx
 (C) ref:  0  0  1  1  1  1  0  1  1  1  1  1  1  0  0  0
 (D) (in) xx xx      p7     xx    p0    p0       xx xx xx
 * (A) is in and (C) is the result of the first
 Using (C) we collect the remaing bytes: xx xx xx xx xx xx in this example
 We convert them to yy yy yy yy yy yy yy

Worst case length, when no compression is possible:

in | bits |     7-bits | out | 7*out | 7*o/8 | out/7 | out%7 | msbits | in%7 | delta to previous | out delta to in
--:|-----:|-----------:|----:|:-----:|:-----:|:-----:|:-----:|:------:|:----:|:-----------------:|----------------
 0 |    0 |  0 * 7 + 0 |   0 |   0   |   0   |   0   |   0   |   0    |  0   |                   |
 1 |    8 |  1 * 7 + 1 |   2 |  14   |   1   |   0   |   2   |   1    |  1   |        +2         | 1
 2 |   16 |  2 * 7 + 2 |   3 |  21   |   2   |   0   |   3   |   2    |  2   |        +1         | 1
 3 |   24 |  3 * 7 + 3 |   4 |  28   |   3   |   0   |   4   |   3    |  3   |        +1         | 1
 4 |   32 |  4 * 7 + 4 |   5 |  35   |   4   |   0   |   5   |   4    |  4   |        +1         | 1
 5 |   40 |  5 * 7 + 5 |   6 |  42   |   5   |   0   |   6   |   5    |  5   |        +1         | 1
 6 |   48 |  6 * 7 + 6 |   7 |  49   |   6   |   1   |   0   |   6    |  6   |        +1         | 1
 7 |   56 |  7 * 7 + 7 |   8 |  56   |   7   |   1   |   1   |   0    |  0   |        +1         | 1
 8 |   64 |  9 * 7 + 1 |  10 |  70   |   8   |   1   |   3   |   1    |  1   |        +2         | 2
 9 |   72 | 10 * 7 + 2 |  11 |  77   |   9   |   1   |   4   |   2    |  2   |        +1         | 2
10 |   80 | 11 * 7 + 3 |  12 |       |       |   1   |   5   |   3    |  3   |        +1         | 2
11 |   88 | 12 * 7 + 4 |  13 |       |       |   1   |   6   |   4    |  4   |        +1         | 2
12 |   96 | 13 * 7 + 5 |  14 |       |       |   2   |   0   |   5    |  5   |        +1         | 2
13 |  104 | 14 * 7 + 6 |  15 |       |       |   2   |   1   |   6    |  6   |        +1         | 2
14 |  112 | 15 * 7 + 7 |  16 |       |       |   2   |   2   |   0    |  0   |        +1         | 2
15 |  120 | 17 * 7 + 1 |  18 |       |       |   2   |   4   |   1    |  1   |        +2         | 3
16 |  128 | 18 * 7 + 2 |  19 |       |       |   2   |   5   |   2    |  2   |        +1         | 3
17 |  136 | 19 * 7 + 3 |  20 |       |       |   2   |   6   |   3    |  3   |        +1         | 3
18 |  144 | 20 * 7 + 4 |  21 |       |       |   3   |   0   |   4    |  4   |        +1         | 3
19 |  152 | 21 * 7 + 5 |  22 |       |       |   3   |   1   |   5    |  5   |        +1         | 3
20 |  160 | 22 * 7 + 6 |  23 |       |       |   3   |   2   |   6    |  6   |        +1         | 3
21 |  168 | 23 * 7 + 7 |  24 |       |       |   3   |   3   |   0    |  0   |        +1         | 3
22 |  176 | 25 * 7 + 1 |  26 |       |       |   3   |   5   |   1    |  1   |        +2         | 4
23 |  184 | 26 * 7 + 2 |  27 |       |       |   3   |   6   |   2    |  2   |        +1         | 4

Compute in from out: in = (7*out/8)
Compute out from in: out = (8*in)/7 + (8*in)%7

msbits = in%7 = (7*out/8)%7
msbits = msbits ? msbits : 7 
-->

