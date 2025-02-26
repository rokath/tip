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
* 1. [TiP - Why and How?](#tip---why-and-how?)
  * 1.1. [Initial Situation](#initial-situation)
    * 1.1.1. [Framing](#framing)
    * 1.1.2. [ Very Small Buffer Data Compession](#-very-small-buffer-data-compession)
  * 1.2. [Bytes and Numbers](#bytes-and-numbers)
  * 1.3. [The TiP Idea](#the-tip-idea)
    * 1.3.1. [Packing](#packing)
    * 1.3.2. [Unpacking](#unpacking)
* 2. [ID Table Generation](#id-table-generation)
  * 2.1. [ID Table Generation Questions](#id-table-generation-questions)
* 3. [Improvement Ideas](#improvement-ideas)
  * 3.1. [Reserve some IDs for Run-Length Encoding](#reserve-some-ids-for-run-length-encoding)
  * 3.2. [Minimize Worst-Case Size](#minimize-worst-case-size)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

</div></ol></details><!-- TABLE OF CONTENTS END -->

---

![./images/logo.png](../images/logo.png)

---

## 1. <a id='tip---why-and-how?'></a>TiP - Why and How?

### 1.1. <a id='initial-situation'></a>Initial Situation

#### 1.1.1. <a id='framing'></a>Framing

For low level buffer storage or MCU transfers some kind of framing is needed for resynchronization after failure. An old variant is to declare a special character as escape sign and to start each package with it. And if the escape sign is part of the buffer data, add an escape sign there too. Even the as escape sign selected character occurs seldom in the buffer data, a careful design should consider the possibility of a buffer containing only such characters.

[COBS](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) is a newer and much better approach, to achieve framing. It transformes the buffer data containing 256 different characters into a sequence of 255 only characters. That allows to use the spare character as frame delimiter. Usually `0` is used for that.

#### 1.1.2. <a id='-very-small-buffer-data-compession'></a> Very Small Buffer Data Compession

A compression followed by COBS framing would do perfectly. But when it comes to very short buffers, like 4 or 20 bytes, **normal zip code fails** to reduce the buffer size.

To combine the COBS technique with compression especially for very short buffers, some additional spare characters are needed. That's done with [TCOBS](https://github.com/rokath/tcobs) in a manual coded way, meaning, expected special data properties are reflected in the TCOBS code. See the [TCOBS User Manual](https://github.com/rokath/tcobs/blob/master/docs/TCOBSv2Specification.md) for more details.

There is also [SMAZ](https://github.com/antirez/smaz), but suitable only for text buffers mainly in English.

[RZCOBS](https://github.com/Dirbaio/rzcobs), assumes many zeroes and tries some compression this way.

An adaptive solution would be nice, meaning, not depending on a specific data structure like English text or many integers.

### 1.2. <a id='bytes-and-numbers'></a>Bytes and Numbers

COBS and TCOBS are starting or ending with some control characters and these are linked togeter to distinguish them from data bytes. But there is also an other option.

If there is a buffer of, let's say 20 bytes, we can consider it as a 20-digit number with 256 ciphers. To free like 8 characters for special usage, we could transform the 20 times 256 cipher number into a 21 or 22 times 248 ciphers number. This transformation is possible, but very computing intensive because of many divisions by 248, or a different base number. So this is no solution for small MCUs. But a division by 128 is cheap! If we transform the 256 base into a 128 base, we only need to perform a shift operation for the conversion. This way we get 128 special characters usable for compressing and framing.

### 1.3. <a id='the-tip-idea'></a>The TiP Idea

#### 1.3.1. <a id='packing'></a>Packing

Find the 127 most common pattern in sample data, similar to the real data, and assign the IDs 1-127 to them. This is done once offline and the generated ID table gets part of the tiny packer code as well as for the tiny unpacker code. For that task a generator tool was build.

At runtime the actual buffer is searched for matching patterns from the ID table beginning with the longest ones. All these found patterns get replaced by the IDs later. All unreplacable bytes are collected into one separate buffer. N unreplacable bytes occupy N\*8 bits. These are distributed onto N\*8/7 7-bit bytes, all getting the MSBit set to avoid zeroes and to distinguish them later from the ID bytes. In fact we do not change these N\*8 bits, we simply reorder them slightly. This bit reordering is de-facto the number transformation to the base 128, mentioned above.

After replacing all found patterns with their IDs, which all have MSBit=0, the unreplacable bytes are replaced with the bit-reordered unreplacable bytes, having MSBit=1.

Now the packing is done and no zeroes are existing anymore.

#### 1.3.2. <a id='unpacking'></a>Unpacking

On the receiver side all bytes with MSBit=0 are identified as IDs and are replaced with the patterns they stay for. All bytes with MSBit=1 are carying the unreplacable bytes bits. These are ordered back to restore the unreplacable bytes which fill the wholes between the patterens then.

## 2. <a id='id-table-generation'></a>ID Table Generation

* We create a bunch of test files with data similar to those we want to pack in the future.
  * For now `ti_generate` takes a longer text file, splits it into sentences and treats each as a separate sample buffer.
  * Also a folder name is accepted and all files are processed then.
* We assume a longest pattern, like 8 for example.
  * `ti_generate` accepts it as parameter.
  * The longest possible pattern is 255 bytes long.
* We take the first 8 bytes of some sample data and move that window in 1-byte steps over the sample data and build a histogram over all found pattern and their occurances count.
* The same is done with all smaller pattern sizes, ergo 7, 6, 5, 4, 3, 2. Not interesting are 1-byte patterns, because their replacement by an ID gives no compression effect.
* The 127 most often occuring pattern are sorted by descending size and are used to create the file `idTable.c`.

### 2.1. <a id='id-table-generation-questions'></a>ID Table Generation Questions

* It is not clear, if the this way created ID table is optimal. Especially, when pattern are sub-pattern of other patterns. That is easily the case with sample data containing the same bytes in longer rows.
* Also it could make sense to use the length of a pattern as weigth. If, for example a 5-bytes long pattern occurs 100 times and a 2-bytes long pattern exists 200 times in the sample data - which should get preceedence to get into the ID table? My guess is, to multiply the pattern length with its occureances count gives a good approximation.
* We could also just determine all pattern from 2 to 8 bytes length and then go byte by byte through the sample data and increment for each byte the pattern counter for the pattern containing this byte on the right place.
* It could make sense, to build several ID tables and then measure how good the packing is with the different tables.

<! --

10 bytes: 123456789a 


+|-|length|pattern|no pattern|weight
-|-|-|-|-|-
10 | 0 | 1er | 1 ... a ||1
9 | 1 |2er| 12 23 ... 9a| a1| 1/2
8 | 2 |3er |123 234 ... 89a| 9a1 a12 | 1/3
4 |6| 7er|1234567 2345678 3456789 456789a||1/7
1 | 9 | 8er | 123456789a | 23456789a1...a123456789|1/8

-->

## 3. <a id='improvement-ideas'></a>Improvement Ideas

### 3.1. <a id='reserve-some-ids-for-run-length-encoding'></a>Reserve some IDs for Run-Length Encoding

* Example:

| ID sequence                               | Meaning                                                       |
| ----------------------------------------- | ------------------------------------------------------------- |
| ID `7F` + count `1...15`                  | 3 to 17 zeroes                                                |
| ID `7F` + count `16...24`                 | 3 to 11 FFs                                                   |
| ID `7F` + count `25...63` + byte `XX`!=0  | 4 to 42 `XX`s, `XX` is any non-zero byte, all `XX` are equal  |
| ID `7F` + `64...255` + `?`                | reserved                                                      |


* The tiny unpack routine first regards all bytes with MSBit=0 as IDs.
* The ID `7F` is followed by a count byte and optional other bytes. These are regarded as part of this ID too during tip package interpretation.
  * The count is guarantied not to be zero and also some optional additional bytes.



<!--
### 3.2. <a id='minimize-worst-case-size'></a>Minimize Worst-Case Size

* If data are containing no ID table pattern at all, they are getting bigger by the factor 8/7. Thats a result of treating the data in 8 bit units (bytes).
* If we change that to 16-bit units, by accepting an optional padding byte, we can reduce this increase factor to 16/15.
* We still have IDs 1-127
* An existing ID 127 just tells if there is a padding byte in the unreplacable data.
* When unpacking, the first set MSBit tells that this byte and the next are unreplaceable. So we get N 16-bit groups of unreplacable data. BUT the 2nd byte could be zero!
* 



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
