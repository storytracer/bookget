package config

import (
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type Input struct {
	DownloaderMode int  // Auto-detect download URL. Values [0|1|2]: 0=default; 1=generic batch download (like IDM/Thunder); 2=IIIF manifest.json auto-detect image download
	UseDzi         bool // Enable Dezoomify for IIIF downloads

	DUrl       string
	UrlsFile   string // Deprecated
	CookieFile string // Input chttp.txt
	HeaderFile string // Input header.txt

	Seq      string // Page range 4:434
	SeqStart int
	SeqEnd   int
	Volume   string // Volume range 4:434
	VolStart int
	VolEnd   int

	Sleep     int    // Rate limiting
	Directory string // Download directory, defaults to Downloads folder in current directory
	Format    string // For high-res image downloads, specify width pixels (16K paper 185mm*260mm, pixels 2185*3071)
	UserAgent string // Custom UserAgent

	Threads       int
	MaxConcurrent int
	PageRate      int           // Page concurrency for IIIF mode
	Timeout       time.Duration // Timeout seconds
	Retries       int           // Retry count

	FileExt string // Specify download file extension
	Quality int    // JPG quality

	Help    bool
	Version bool
}

func Init(ctx context.Context) bool {

	dir, _ := os.Getwd()

	// Directory path validation for Windows - avoid special characters that cause issues
	if os.PathSeparator == '\\' {
		matched, _ := regexp.MatchString(`([^A-z0-9_\\/\-:.]+)`, dir)
		if matched {
			fmt.Println("Software directory path cannot contain spaces, Chinese characters, or other special symbols. Recommended: D:\\bookget")
			fmt.Println("Press Enter to exit...")
			endKey := make([]byte, 1)
			os.Stdin.Read(endKey)
			os.Exit(0)
		}
	}

	pflag.StringVarP(&Conf.DUrl, "input", "i", "", "Download URL")
	pflag.StringVarP(&Conf.UrlsFile, "input-file", "I", "", "Download URLs from file")
	pflag.StringVarP(&Conf.Directory, "dir", "O", path.Join(dir, "downloads"), "Save files to directory")

	pflag.StringVarP(&Conf.Seq, "sequence", "p", "", "Page range, e.g. 4:434")
	pflag.StringVarP(&Conf.Volume, "volume", "v", "", "Multi-volume books, e.g. 10:20 volumes, download only volumes 10 to 20")

	pflag.StringVar(&Conf.Format, "format", "full/full/0/default.jpg", "IIIF image request URI")

	pflag.StringVarP(&Conf.UserAgent, "user-agent", "U", defaultUserAgent, "HTTP header user-agent")

	pflag.BoolVarP(&Conf.UseDzi, "dzi", "d", true, "Use IIIF/DeepZoom tile download")

	pflag.StringVarP(&Conf.CookieFile, "cookies", "C", path.Join(dir, "cookie.txt"), "Cookie file")
	pflag.StringVarP(&Conf.HeaderFile, "headers", "H", path.Join(dir, "header.txt"), "Header file")

	pflag.IntVarP(&Conf.Threads, "threads", "n", 1, "Maximum threads per task")
	pflag.IntVarP(&Conf.MaxConcurrent, "concurrent", "c", 16, "Maximum concurrent tasks")
	pflag.IntVar(&Conf.PageRate, "page-rate", 1, "Page concurrency for IIIF mode, default 1 (sequential download)")

	pflag.IntVar(&Conf.Quality, "quality", 80, "JPG quality, default 80")
	pflag.StringVar(&Conf.FileExt, "ext", ".jpg", "Specify file extension [.jpg|.tif|.png] etc.")

	pflag.IntVar(&Conf.Retries, "retries", 3, "Download retry count")

	pflag.DurationVarP(&Conf.Timeout, "timeout", "T", 300, "Network timeout (seconds)")
	pflag.IntVar(&Conf.Sleep, "sleep", 3, "Interval sleep seconds, typical range 3-20")

	pflag.IntVarP(&Conf.DownloaderMode, "downloader_mode", "m", 0, "Download mode. Values [0|1|2]: 0=default;\n1=generic batch download (like IDM/Thunder);\n2=IIIF manifest.json auto-detect image download")

	pflag.BoolVarP(&Conf.Help, "help", "h", false, "Show help")
	pflag.BoolVarP(&Conf.Version, "version", "V", false, "Show version")
	pflag.Parse()

	k := len(os.Args)
	if k == 2 {
		if Conf.Version {
			printVersion()
			return false
		}
		if Conf.Help {
			printHelp()
			return false
		}
	}
	v := pflag.Arg(0)
	if strings.HasPrefix(v, "http") {
		Conf.DUrl = v
	}
	if Conf.UrlsFile != "" && !strings.Contains(Conf.UrlsFile, string(os.PathSeparator)) {
		Conf.UrlsFile = path.Join(dir, Conf.UrlsFile)
	}
	initSeqRange()
	initVolumeRange()
	// Create download directory
	_ = os.Mkdir(Conf.Directory, os.ModePerm)
	//_ = os.Mkdir(CacheDir(), os.ModePerm)
	return true
}

func printHelp() {
	printVersion()
	fmt.Println(`Usage: bookget [OPTION]... [URL]...`)
	pflag.PrintDefaults()
	fmt.Println()
	fmt.Println("Originally written by zhudw <zhudwi@outlook.com>.")
	fmt.Println("https://github.com/deweizhu/bookget/")
}

func printVersion() {
	fmt.Printf("bookget v%s\n", Version)
}
