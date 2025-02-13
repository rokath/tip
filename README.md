<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
![GitHub issues](https://img.shields.io/github/issues/rokath/tip)
![GitHub downloads](https://img.shields.io/github/downloads/rokath/tip/total)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/rokath/tip)
![GitHub watchers](https://img.shields.io/github/watchers/rokath/tip?label=watch)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
[![Go Report Card](https://goreportcard.com/badge/github.com/rokath/tip)](https://goreportcard.com/report/github.com/rokath/tip)

<!--
![GitHub release (latest by date)](https://img.shields.io/github/v/release/rokath/tip)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/rokath/tip/latest)
[![Coverage Status](https://coveralls.io/repos/github/rokath/tip/badge.svg?branch=master)](https://coveralls.io/github/rokath/tip?branch=master)
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]
https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/rokath/tip">
    <img src="images/logo.png" alt="Logo"> <!--width="80" height="80"-->
  </a>

</div>


# TiP - Tiny Packer for small buffers

```diff

+ Compress small buffers fast and efficient with Zeroes Elemination for delimitting. +
-->  Works with big buffers too but will not compress like establisched zip tools ❗

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


## Project Status

```diff
--> Early experimental state!
```


<!-- ABOUT THE PROJECT -->

## About The Project

* Usual compressors cannot succeed on small data, because they add a translation table into the data.
* **TiP** is an adaptable small buffer packer suitable for embedded devices. Like [COBS](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) it removes all zeroes from the data, but additionally tries data compression. 
  * [COBS](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) adds on byte overhead per each starting 255 bytes.
  * The TiP worst-case overhead is 1 byte per each starting 127 bytes for uncompressable data.
* Like [TCOBS](github.com/rokath/tcobs), TiP can already compress 3 bytes into 2 bytes but is expected to do better in a general way.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### TiP Components

* Tiny C-Code usable on embedded devices inside `src` folder containing
  * `idTable.c` - generated data specific translation table
  * `pack.c` and `unpack.c` - separately or together compilable
* `tipTable` - ID Table Generator to create a suitable `idTable.c` file
* `tip` - tiny pack app using the **pack** C code
* `tiu` - tiny unpack app using the **unpack** C code

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

* For now install [Go](https://golang.org/) to easily build the executables.
* You need some files containing typical data you want to tip.
  * Just to try TiP out, you can use a folder containing ASCII texts.

### Built TipTable Generator `tipTable`

* `cd tipTable; go build -ldflags "-w" ./...`
* Run `tipTable` on the data files to get an `idTable.c` file.

### Build `tip` and `tiu`

* Copy the generated `idTable.c` file into the `src` folder.
* Run `go build ./...` or `go install ./...`.

### Try `tip` and `tiu`

* Run `tip -i myFile -v` to get `myFile.tip`.
* Run `tiu -i myFile.tip -v` to get `myFile.tip.untip`.
* `myFile` and `myFile.tip.untip` are equal.

### Installation

* Add `src` folder to your project and compile.
* `pack.h` and `unpack.h` is the user interface.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- USAGE EXAMPLES -->
## Usage

Please refer to the [Tip User Manual](./docs/TipUserManual.md) (not written yet)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [x] Create `tipTable.h` Generator `tipTable`.
- [x] Create `pack.c` and `unpack.c` and test.
- [ ] Improve `tipTable.h` Generator `tipTable`.
- [ ] Write [Tip User Manual](./docs/TipUserManual.md).
- [ ] Write fuzzy tests.
- [ ] Remove 255 bytes limitation.
- [ ] Improve implementation for less RAM usage.


See the [open issues](https://github.com/rokath/tip/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Top contributors:

<a href="https://github.com/rokath/tip/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=rokath/tip" alt="contrib.rocks image" />
</a>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Thomas Höhenleitner - th@seerose.net

Project Link: [https://github.com/rokath/tip](https://github.com/rokath/tip)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [COBS](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing)
* []() to do
* []() todo
* to do

<!--

Use this space to list resources you find helpful and would like to give credit to. I've included a few of my favorites to kick things off!

* [Choose an Open Source License](https://choosealicense.com)
* [GitHub Emoji Cheat Sheet](https://www.webpagefx.com/tools/emoji-cheat-sheet)
* [Malven's Flexbox Cheatsheet](https://flexbox.malven.co/)
* [Malven's Grid Cheatsheet](https://grid.malven.co/)
* [Img Shields](https://shields.io)
* [GitHub Pages](https://pages.github.com)
* [Font Awesome](https://fontawesome.com)
* [React Icons](https://react-icons.github.io/react-icons/search)

-->

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

<!--

[contributors-shield]: https://img.shields.io/github/contributors/rokath/tip.svg?style=for-the-badge
[contributors-url]: https://github.com/rokath/tip/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/rokath/tip.svg?style=for-the-badge
[forks-url]: https://github.com/rokath/tip/network/members
[stars-shield]: https://img.shields.io/github/stars/rokath/tip.svg?style=for-the-badge
[stars-url]: https://github.com/rokath/tip/stargazers
[issues-shield]: https://img.shields.io/github/issues/rokath/tip.svg?style=for-the-badge
[issues-url]: https://github.com/rokath/tip/issues
[license-shield]: https://img.shields.io/github/license/rokath/tip.svg?style=for-the-badge
[license-url]: https://github.com/rokath/tip/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/linkedin_username
[product-screenshot]: images/screenshot.png


[contributors-shield]: https://img.shields.io/github/contributors/othneildrew/Best-README-Template.svg?style=for-the-badge
[contributors-url]: https://github.com/othneildrew/Best-README-Template/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/othneildrew/Best-README-Template.svg?style=for-the-badge
[forks-url]: https://github.com/othneildrew/Best-README-Template/network/members
[stars-shield]: https://img.shields.io/github/stars/othneildrew/Best-README-Template.svg?style=for-the-badge
[stars-url]: https://github.com/othneildrew/Best-README-Template/stargazers
[issues-shield]: https://img.shields.io/github/issues/othneildrew/Best-README-Template.svg?style=for-the-badge
[issues-url]: https://github.com/othneildrew/Best-README-Template/issues
[license-shield]: https://img.shields.io/github/license/othneildrew/Best-README-Template.svg?style=for-the-badge
[license-url]: https://github.com/othneildrew/Best-README-Template/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/othneildrew
[product-screenshot]: images/screenshot.png

-->


<!--
<h3 align="center">Tiny Packer</h3>
  <p align="center">
    for small buffers
    <br />
    <a href="https://github.com/rokath/tip"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/rokath/tip">View Demo</a>
    ·
    <a href="https://github.com/rokath/tip/issues">Report Bug</a>
    ·
    <a href="https://github.com/rokath/tip/issues">Request Feature</a>
  </p>
-->