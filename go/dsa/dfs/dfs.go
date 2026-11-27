package dfs

/**
 * 733. Flood Fill
 */
func floodFill(image [][]int, sr int, sc int, newColor int) [][]int {
	var dfs func(int, int, int, int)
	dfs = func(x int, y int, sourceColor int, newColor int) {
		if x < 0 || x >= len(image) || y < 0 || y >= len(image[0]) ||
			image[x][y] != sourceColor {
			return
		}
		image[x][y] = newColor
		dfs(x-1, y, sourceColor, newColor)
		dfs(x+1, y, sourceColor, newColor)
		dfs(x, y-1, sourceColor, newColor)
		dfs(x, y+1, sourceColor, newColor)
	}
	sourceColor := image[sr][sc]
	if sourceColor != newColor {
		dfs(sr, sc, sourceColor, newColor)
	}

	return image
}
