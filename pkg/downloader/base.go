package downloader

const (
	maxConcurrent = 16 // Maximum concurrent downloads
	userAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0"
	minFileSize   = 1024 // Minimum file size (1KB)

	maxRetries = 3
	JPGQuality = 90
)

// Add in downloader.go or related files
type Vec2d struct {
	x int // 或 float64 根据需求
	y int // 或 float64
}

// Optional: add constructor function
func NewVec2d(x, y int) Vec2d {
	return Vec2d{x: x, y: y}
}

// Optional: add common methods
func (v Vec2d) Width() int  { return v.x }
func (v Vec2d) Height() int { return v.y }

type TileSizeFormat int

const (
	WidthHeight TileSizeFormat = iota // "width,height"
	Width                             // "width,"
)

// Quality preference order (least to most preferred)
var qualityOrder = []string{"default", "native"}

// Format preference order (least to most preferred)
var formatOrder = []string{"jpg", "png"}
