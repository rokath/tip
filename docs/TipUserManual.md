# Tip User Manual

(to do)


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

-->

