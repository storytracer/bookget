# bookget

> **⚠️ Important Notice:** This is an English-translated fork of the original bookget project. This fork may contain substantive alterations, feature removals, and compatibility changes compared to the original version. Future development may diverge significantly from the upstream project.

## Original Project Attribution

This project is a fork of the original [bookget](https://github.com/deweizhu/bookget) by [deweizhu](https://github.com/deweizhu). All original code, features, and documentation are credited to the original author. This fork exists solely to provide an English version for international users and may be modified for specific use cases.

**Original Repository:** https://github.com/deweizhu/bookget  
**Original Author:** [deweizhu](https://github.com/deweizhu)

---

# Introduction

**bookget** is a powerful Go-based digital ancient book download tool that supports 50+ digital libraries and cultural institutions worldwide. It downloads high-resolution images from academic and cultural digital collections, with specialized support for various formats including IIIF, DZI (Deep Zoom Images), and custom institutional APIs.

## Key Features

- **Multi-institutional Support**: Compatible with major libraries, museums, and archives across China, Japan, United States, Europe, Korea, Russia, and Vietnam
- **IIIF Protocol Support**: Full International Image Interoperability Framework (IIIF) compatibility, similar to dezoomify functionality
- **Concurrent Downloads**: Multi-threaded downloading with configurable thread pools and resume capability
- **Format Detection**: Automatic site identification from URL patterns and content types
- **Authentication Support**: Cookie and header management for authenticated sessions
- **Progress Tracking**: Real-time download progress visualization
- **Proxy Support**: Respects HTTP_PROXY/HTTPS_PROXY environment variables

## IIIF Compatibility

bookget provides comprehensive IIIF (International Image Interoperability Framework) support, functioning similarly to dezoomify for downloading high-resolution images from IIIF-compliant repositories. The tool automatically detects IIIF manifests and can process:

- **IIIF Image API** endpoints for tiled high-resolution images
- **IIIF Presentation API** manifests for complete documents
- **Auto-detection** of IIIF resources from URLs containing `.json` files
- **Multi-format support** including JPEG, PNG, and WebP image formats

This makes bookget compatible with any institution using standard IIIF protocols, extending support beyond the explicitly listed sites.

#### Source Code Compilation
Building from source code is for computer programmers reference only. Regular users can skip this section.
Read the [golang official documentation](https://golang.google.cn/doc/install) to install the golang development environment on your computer.
```shell
git clone https://github.com/storytracer/bookget.git
cd bookget

# Can be run directly during local development
make linux-amd64    # Compile Linux version
make windows-amd64  # Compile Windows version
make release        # Compile all platforms
```

## Explictly supported domains

- 111.7.82.29:8090
- 124.134.220.209:8100
- archive.wul.waseda.ac.jp
- arcgxhpv7cw0.db.wzlib.cn
- asia.si.edu
- babel.hathitrust.org
- catalog.princeton.edu
- curiosity.lib.harvard.edu
- da.library.ryukoku.ac.jp
- db2.sido.keio.ac.jp
- dcollections.lib.keio.ac.jp
- dfz.yn.gov.cn
- digicoll.lib.berkeley.edu
- digital.bodleian.ox.ac.uk
- digital.onb.ac.at
- digital.staatsbibliothek-berlin.de
- digitalrepository.lib.hku.hk
- dl.ndl.go.jp
- dlib.cafa.edu.cn
- dlibgate.cafa.edu.cn
- dpul.princeton.edu
- dsr.nii.ac.jp
- emuseum.nich.go.jp
- gjsztsg.sdutcm.edu.cn
- gj.tianyige.com.cn
- guji.nlc.cn
- guji.sdlib.com
- gzdd.gzlib.gov.cn
- gzdd.gzlib.org.cn
- hannom.nlv.gov.vn
- ids.si.edu
- idp.afc.ryukoku.ac.jp
- idp.bbaw.de
- idp.bl.uk
- idp.bnf.fr
- idp.korea.ac.kr
- idp.nlc.cn
- idp.orientalstudies.ru
- iiif.lib.harvard.edu
- iiif.si.edu
- jsgxgj.nju.edu.cn
- kanji.zinbun.kyoto-u.ac.jp
- khirin-a.rekihaku.ac.jp
- kokusho.nijl.ac.jp
- kostma.korea.ac.kr
- kyudb.snu.ac.kr
- lbezone.hkust.edu.hk
- lib.nomfoundation.org
- listview.lib.harvard.edu
- lod.nl.go.kr
- lswx.tjl.tj.cn:8001
- mylib.nlc.cn
- ostasien.digitale-sammlungen.de
- ouroots.nlc.cn
- oyjy.wzlib.cn
- read.nlc.cn
- repo.komazawa-u.ac.jp
- repository.lib.cuhk.edu.hk
- shanben.ioc.u-tokyo.ac.jp
- taiwanebook.ncl.edu.tw
- viewer.rsl.ru
- webarchives.tnm.jp
- www.bl.uk
- www.digital.archives.go.jp
- www.digitale-sammlungen.de
- www.familysearch.org
- www.iiif.ku-orcas.kansai-u.ac.jp
- www.library.yonezawa.yamagata.jp
- www.loc.gov
- www.modernhistory.org.cn
- www.ncpssd.cn
- www.ncpssd.org
- www.si.edu
- yun.szlib.org.cn