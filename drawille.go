package drawille

//import "code.google.com/p/goncurses"
import (
	"math"
	"strings"
)

var pixel_map = [4][2]int{
	{0x01, 0x08},
	{0x02, 0x10},
	{0x04, 0x20},
	{0x40, 0x80}}

// Braille chars start at 0x2800
var braille_char_offset = 0x2800

func getPixel(y, x int) int {
	var cy, cx int
	if y >= 0 {
		cy = y % 4
	} else {
		cy = 3 + ((y + 1) % 4)
	}
	if x >= 0 {
		cx = x % 2
	} else {
		cx = 1 + ((x + 1) % 2)
	}
	return pixel_map[cy][cx]
}

type Canvas struct {
	cfg Config

	LineEnding string
	chars      map[int]map[int]int
	colors     map[int]map[int]int
}

// Make a new canvas
func NewCanvas(os ...Option) Canvas {
	c := Canvas{LineEnding: "\n"}
	c.Clear(os...)
	return c
}

func (c Canvas) MaxY() int {
	max := 0
	for k, _ := range c.chars {
		if k > max {
			max = k
		}
	}
	return max * 4
}

func (c Canvas) MinY() int {
	min := 0
	for k, _ := range c.chars {
		if k < min {
			min = k
		}
	}
	return min * 4
}

func (c Canvas) MaxX() int {
	max := 0
	for _, v := range c.chars {
		for k, _ := range v {
			if k > max {
				max = k
			}
		}
	}
	return max * 2
}

func (c Canvas) MinX() int {
	min := 0
	for _, v := range c.chars {
		for k, _ := range v {
			if k < min {
				min = k
			}
		}
	}
	return min * 2
}

// Clear all pixels
func (c *Canvas) Clear(os ...Option) {
	c.chars = make(map[int]map[int]int)
	c.colors = make(map[int]map[int]int)
	for _, o := range os {
		o(&c.cfg)
	}
}

// Convert x,y to cols, rows
func (c Canvas) get_pos(x, y int) (int, int) {
	return (x / 2), (y / 4)
}

// Set a pixel of c
func (c *Canvas) Set(x, y int, color ...int) {
	px, py := c.get_pos(x, y)
	if m := c.chars[py]; m == nil {
		c.chars[py] = make(map[int]int)
	}
	val := c.chars[py][px]
	mapv := getPixel(y, x)
	c.chars[py][px] = val | mapv
	if len(color) > 0 {
		if m := c.colors[py]; m == nil {
			c.colors[py] = make(map[int]int)
		}
		c.colors[py][px] = color[0]
	}
}

// Unset a pixel of c
func (c *Canvas) UnSet(x, y int, color ...int) {
	px, py := c.get_pos(x, y)
	x, y = int(math.Abs(float64(x))), int(math.Abs(float64(y)))
	if m := c.chars[py]; m == nil {
		c.chars[py] = make(map[int]int)
	}
	c.chars[py][px] = c.chars[py][px] &^ getPixel(y, x)
	if len(color) > 0 {
		c.colors[py][px] = color[0]
	}
}

// Toggle a point
func (c *Canvas) Toggle(x, y int, color ...int) {
	px, py := c.get_pos(x, y)
	if (c.chars[py][px] & getPixel(y, x)) != 0 {
		c.UnSet(x, y)
	} else {
		c.Set(x, y)
	}
}

// Set text to the given coordinates
func (c *Canvas) SetText(x, y int, text string) {
	x, y = x/2, y/4
	if m := c.chars[y]; m == nil {
		c.chars[y] = make(map[int]int)
	}
	for i, char := range text {
		c.chars[y][x+i] = int(char) - braille_char_offset
	}
}

// Get pixel at the given coordinates
func (c Canvas) Get(x, y int) bool {
	dot_index := pixel_map[y%4][x%2]
	x, y = x/2, y/4
	char := c.chars[y][x]
	return (char & dot_index) != 0
}

// Get character at the given screen coordinates
func (c Canvas) GetScreenCharacter(x, y int) rune {
	return rune(c.chars[y][x] + braille_char_offset)
}

// Get character for the given pixel
func (c Canvas) GetCharacter(x, y int) rune {
	return c.GetScreenCharacter(x/4,y/4)
}

func (c Canvas) colorize(rownum, x int,  s string) string {
	if c.cfg.pallette == nil {
		return s
	}

	if fnId, ok := c.colors[rownum][x]; ok {
		fn := c.cfg.pallette[fnId]
		return fn(s)
	}
	return s
}

// Retrieve the rows from a given view
func (c Canvas) Rows(minX, minY, maxX, maxY int) []string {
	minrow, maxrow := minY/4, (maxY)/4
	mincol, maxcol := minX/2, (maxX)/2
	var sb strings.Builder

	ret := make([]string, 0)
	for rownum := minrow; rownum < (maxrow + 1); rownum = rownum + 1 {
		for x := mincol; x < (maxcol + 1); x = x + 1 {
			char := c.chars[rownum][x]
			// TODO(cwolff): make this more compact
			s := string(rune(char + braille_char_offset))

			s = c.colorize(rownum, x, s)
			sb.WriteString(s)
		}
		ret = append(ret, sb.String())
		sb.Reset()
	}
	return ret
}

// Retrieve a string representation of the frame at the given parameters
func (c Canvas) Frame(minX, minY, maxX, maxY int) string {
	var ret string
	for _, row := range c.Rows(minX, minY, maxX, maxY) {
		ret += row
		ret += c.LineEnding
	}
	return ret
}

func (c Canvas) String() string {
	return c.Frame(c.MinX(), c.MinY(), c.MaxX(), c.MaxY())
}

func (c *Canvas) DrawLine(x1, y1, x2, y2 float64, color ...int) {
	xdiff := math.Abs(x1 - x2)
	ydiff := math.Abs(y2 - y1)

	var xdir, ydir float64
	if x1 <= x2 {
		xdir = 1
	} else {
		xdir = -1
	}
	if y1 <= y2 {
		ydir = 1
	} else {
		ydir = -1
	}

	r := math.Max(xdiff, ydiff)

	for i := 0; i < round(r)+1; i = i + 1 {
		x, y := x1, y1
		if ydiff != 0 {
			y += (float64(i) * ydiff) / (r * ydir)
		}
		if xdiff != 0 {
			x += (float64(i) * xdiff) / (r * xdir)
		}
		if c.cfg.overwrite == overwriteToggle {
			c.Toggle(round(x), round(y), color...)
		} else {
			c.Set(round(x), round(y), color...)
		}
	}
}

func (c *Canvas) DrawPolygon(center_x, center_y, sides, radius float64) {
	degree := 360 / sides
	for n := 0; n < int(sides); n = n + 1 {
		a := float64(n) * degree
		b := float64(n+1) * degree

		x1 := (center_x + (math.Cos(radians(a)) * (radius/2 + 1)))
		y1 := (center_y + (math.Sin(radians(a)) * (radius/2 + 1)))
		x2 := (center_x + (math.Cos(radians(b)) * (radius/2 + 1)))
		y2 := (center_y + (math.Sin(radians(b)) * (radius/2 + 1)))

		c.DrawLine(x1, y1, x2, y2)
	}
}

func radians(d float64) float64 {
	return d * (math.Pi / 180)
}

func round(x float64) int {
	return int(x + 0.5)
}
