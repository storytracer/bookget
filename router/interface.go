// router/interface.go
package router

import (
	"bookget/app"
	"bookget/config"
	"bookget/pkg/util"
	"errors"
	"strings"
	"sync"
)

type RouterInit interface {
	GetRouterInit(sUrl string) (map[string]interface{}, error)
}

var (
	Router = make(map[string]RouterInit)
	doInit sync.Once
)

// FactoryRouter factory function for creating routers
func FactoryRouter(siteID string, sUrl string) (map[string]interface{}, error) {
	// Auto-detection logic
	if config.Conf.DownloaderMode == 1 {
		siteID = "bookget"
	} else if config.Conf.DownloaderMode == 2 || strings.Contains(sUrl, ".json") {
		siteID = "iiif.io"
	}

	if strings.Contains(sUrl, "tiles/infos.json") {
		siteID = "dzicnlib"
	}

	// Initialize routers (thread-safe)
	doInit.Do(func() {
		//[China] National Library of China
		Router["read.nlc.cn"] = app.NewChinaNlc()
		Router["mylib.nlc.cn"] = app.NewChinaNlc()
		Router["guji.nlc.cn"] = app.NewNlcGuji()

		//[China] Taiwan Chinese E-book Repository
		Router["taiwanebook.ncl.edu.tw"] = app.NewHuawen()

		//[China] Chinese University of Hong Kong Library
		Router["repository.lib.cuhk.edu.hk"] = app.NewCuhk()

		//[China] Hong Kong University of Science and Technology Library
		Router["lbezone.hkust.edu.hk"] = app.NewUsthk()

		//[China] Luoyang City Library
		Router["111.7.82.29:8090"] = app.NewLuoyang()

		//[China] Wenzhou City Library
		Router["oyjy.wzlib.cn"] = app.NewWzlib()
		Router["arcgxhpv7cw0.db.wzlib.cn"] = app.NewWzlib()

		//[China] Shenzhen Library - Ancient Books
		Router["yun.szlib.org.cn"] = app.NewSzLib()

		//[China] Guangzhou Dadian
		Router["gzdd.gzlib.gov.cn"] = app.NewGzlib()
		Router["gzdd.gzlib.org.cn"] = app.NewGzlib()

		//[China] Tianyi Pavilion Museum Ancient Books Digitization Platform
		Router["gj.tianyige.com.cn"] = app.NewTianyige()

		//[China] Jiangsu Colleges Precious Ancient Books Digital Library
		Router["jsgxgj.nju.edu.cn"] = app.NewNjuedu()

		//[China] China Roots Network - National Library
		Router["ouroots.nlc.cn"] = app.NewOuroots()

		//[China] National Center for Philosophy and Social Sciences Documentation
		Router["www.ncpssd.org"] = app.NewNcpssd()
		Router["www.ncpssd.cn"] = app.NewNcpssd()

		//[China] Shandong University of Traditional Chinese Medicine Digital Ancient Books Library
		Router["gjsztsg.sdutcm.edu.cn"] = app.NewSdutcm()
		//[China] Shandong Province Ancient Books Digital Resource Platform
		Router["guji.sdlib.com"] = app.NewSdlib()

		//[China] Tianjin Library Historical Literature Digital Resource Database
		Router["lswx.tjl.tj.cn:8001"] = app.NewTjlswx()

		//[China] Yunnan Digital Local Gazetteer
		Router["dfz.yn.gov.cn"] = app.NewYndfz()

		//[China] University of Hong Kong Digital Library
		Router["digitalrepository.lib.hku.hk"] = app.NewHkulib()

		//[China] Zhucheng City Library, Shandong Province
		Router["124.134.220.209:8100"] = app.NewZhuCheng()
		//[China] Central Academy of Fine Arts
		Router["dlibgate.cafa.edu.cn"] = app.NewCafaEdu()
		Router["dlib.cafa.edu.cn"] = app.NewCafaEdu()

		//[China] Anti-Japanese War and Sino-Japanese Relations Literature Database Platform
		Router["www.modernhistory.org.cn"] = app.NewWar1931()
		//}}} -----------------------------------------------------------------

		//---------------Japan--------------------------------------------------
		//[Japan] National Diet Library
		Router["dl.ndl.go.jp"] = app.NewNdlJP()

		//[Japan] e-Museum National Treasures
		Router["emuseum.nich.go.jp"] = app.NewEmuseum()

		//[Japan] Imperial Household Agency Archives and Mausolea Department (Chinese Books Collection)
		Router["db2.sido.keio.ac.jp"] = app.NewKeio()

		//[Japan] University of Tokyo Institute for Oriental Culture (Chinese Rare Books Database)
		Router["shanben.ioc.u-tokyo.ac.jp"] = app.NewUtokyo()

		//[Japan] National Archives of Japan (Cabinet Library)
		Router["www.digital.archives.go.jp"] = app.NewNationaljp()

		//[Japan] Toyo Bunko (Oriental Library)
		Router["dsr.nii.ac.jp"] = app.NewNiiac()

		//[Japan] Waseda University Library
		Router["archive.wul.waseda.ac.jp"] = app.NewWaseda()

		//[Japan] Kokusho Database (Classical Books)
		Router["kokusho.nijl.ac.jp"] = app.NewKokusho()

		//[Japan] Kyoto University Institute for Research in Humanities - Digital Library Museum of Oriental Studies
		Router["kanji.zinbun.kyoto-u.ac.jp"] = app.NewKyotou()

		//[Japan] Komazawa University Electronic Rare Books Collection
		Router["repo.komazawa-u.ac.jp"] = app.NewIiifRouter()

		//[Japan] Kansai University Library
		Router["www.iiif.ku-orcas.kansai-u.ac.jp"] = app.NewIiifRouter()

		//[Japan] Keio University Library
		Router["dcollections.lib.keio.ac.jp"] = app.NewIiifRouter()

		//[Japan] National Museum of Japanese History
		Router["khirin-a.rekihaku.ac.jp"] = app.NewKhirin()

		//[Japan] Yonezawa City Library
		Router["www.library.yonezawa.yamagata.jp"] = app.NewYonezawa()
		Router["webarchives.tnm.jp"] = app.NewTnm()

		//[Japan] Ryukoku University
		Router["da.library.ryukoku.ac.jp"] = app.NewRyukoku()
		//}}} -----------------------------------------------------------------

		//{{{---------------United States, Europe--------------------------------------------------
		//[United States] Harvard University Library
		Router["iiif.lib.harvard.edu"] = app.NewHarvard()
		Router["listview.lib.harvard.edu"] = app.NewHarvard()
		Router["curiosity.lib.harvard.edu"] = app.NewHarvard()

		//[United States] HathiTrust Digital Library
		Router["babel.hathitrust.org"] = app.NewHathitrust()

		//[United States] Princeton University Library
		Router["catalog.princeton.edu"] = app.NewPrinceton()
		Router["dpul.princeton.edu"] = app.NewPrinceton()

		//[United States] Library of Congress
		Router["www.loc.gov"] = app.NewLoc()

		//[United States] Stanford University Library

		//[United States] Utah Genealogy (FamilySearch)
		Router["www.familysearch.org"] = app.NewFamilysearch()

		//[Germany] Berlin State Library
		Router["digital.staatsbibliothek-berlin.de"] = app.NewBerlin()

		//[Germany] Bavarian State Library East Asian Digital Collections
		Router["ostasien.digitale-sammlungen.de"] = app.NewSammlungen()
		Router["www.digitale-sammlungen.de"] = app.NewSammlungen()

		//[United Kingdom] Oxford University Bodleian Library
		Router["digital.bodleian.ox.ac.uk"] = app.NewOxacuk()

		//[United Kingdom] British Library Manuscripts
		Router["www.bl.uk"] = app.NewBluk()

		//Smithsonian Institution
		Router["ids.si.edu"] = app.NewSiEdu()
		Router["www.si.edu"] = app.NewSiEdu()
		Router["iiif.si.edu"] = app.NewSiEdu()
		Router["asia.si.edu"] = app.NewSiEdu()

		//[United States] UC Berkeley East Asian Library
		Router["digicoll.lib.berkeley.edu"] = app.NewBerkeley()

		//[Austria] Austrian National Library
		Router["digital.onb.ac.at"] = app.NewOnbDigital()
		//}}} -----------------------------------------------------------------

		//{{{---------------Others--------------------------------------------------
		//International Dunhuang Project
		Router["idp.nlc.cn"] = app.NewIdp()
		Router["idp.bl.uk"] = app.NewIdp()
		Router["idp.orientalstudies.ru"] = app.NewIdp()
		Router["idp.afc.ryukoku.ac.jp"] = app.NewIdp()
		Router["idp.bbaw.de"] = app.NewIdp()
		Router["idp.bnf.fr"] = app.NewIdp()
		Router["idp.korea.ac.kr"] = app.NewIdp()

		//[Korea]
		Router["kyudb.snu.ac.kr"] = app.NewKyudbSnu()
		Router["lod.nl.go.kr"] = app.NewLodNLGoKr()

		//[Korea] Korea University
		Router["kostma.korea.ac.kr"] = app.NewKorea()

		//[Russia] Russian State Library
		Router["viewer.rsl.ru"] = app.NewRslRu()

		//[Vietnam] Vietnamese Han-Nom Ancient Books Digital Preservation Project
		Router["lib.nomfoundation.org"] = app.NewNomfoundation()

		//[Vietnam] National Library of Vietnam Han-Nom Library
		Router["hannom.nlv.gov.vn"] = app.NewHannomNlv()
		//}}} -----------------------------------------------------------------

		Router["bookget"] = app.NewImageDownloader()
		Router["dzicnlib"] = app.NewDziCnLib()
		Router["iiif.io"] = app.NewIiifRouter()
	})

	// Check if router exists
	if _, ok := Router[siteID]; !ok {
		urlType := util.GetHeaderContentType(sUrl)
		if urlType == "json" {
			siteID = "iiif.io"
		} else if urlType == "bookget" {
			siteID = "bookget"
		}

		if _, ok := Router[siteID]; !ok {
			return nil, errors.New("unsupported URL: " + sUrl)
		}
	}

	return Router[siteID].GetRouterInit(sUrl)

}
