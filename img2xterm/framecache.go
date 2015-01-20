package img2xterm

const (
	cacheXSize int = 360
	cacheYSize int = 120
)

type FrameCache [cacheXSize][cacheYSize][2]int

var fCache FrameCache

func ClearCache() {
	// TODO: Can this be more efficient?
	for i := 0; i < cacheXSize; i++ {
		for j := 0; j < cacheYSize; j++ {
			fCache[i][j][0] = 0
			fCache[i][j][1] = 0
		}
	}
}
