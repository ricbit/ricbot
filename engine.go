package engine

import "container/list"

const (
  EMPTY = iota
  BLACK
  WHITE
  INVALID
)

type Position int

type Goban interface {
  SizeX() int
  SizeY() int
  GetPosition(y, x int) Position
  GetVisitorMarker() VisitorMarker
}

type VisitorMarker interface {
  ClearMarks()
  SetMark(y, x int)
  IsMarked(y, x int) bool
}

type ArrayGoban struct {
  size_x, size_y int
  board [][]Position
}

func NewArrayGoban(size_y, size_x int, init string) *ArrayGoban {
  goban := new(ArrayGoban)
  goban.size_x = size_x
  goban.size_y = size_y
  goban.board = make([][]Position, size_y)
  conv := map[byte] Position {
    '.': EMPTY,
    'o': WHITE,
    'x': BLACK,
  }
  for j := 0; j < size_y; j++ {
    goban.board[j] = make([]Position, size_x)
    for i := 0; i < size_x; i++ {
      goban.board[j][i] = conv[init[j * size_x + i]]
    }
  }
  return goban
}

func (g *ArrayGoban) SizeX() int {
  return g.size_x
}

func (g *ArrayGoban) GetVisitorMarker() VisitorMarker {
  return g
}

func (g *ArrayGoban) SizeY() int {
  return g.size_y
}

func (g *ArrayGoban) GetPosition(y, x int) Position {
  return g.board[y][x]
}

func (g *ArrayGoban) ClearMarks() {
  for j := 0; j < g.size_y; j++ {
    for i := 0; i < g.size_x; i++ {
      g.board[j][i] &= 0x3
    }
  }
}

func (g *ArrayGoban) SetMark(y, x int) {
  g.board[y][x] += 0x4
}

func (g *ArrayGoban) IsMarked(y, x int) bool {
  return g.board[y][x] & 0x4 > 0
}

func encode(y, x int) int {
  return y * 256 + x
}

func decode(code int) (int, int) {
  return code >> 8, code & 255
}

func valid(y, x, maxY, maxX int) bool {
  return x >= 0 && y >= 0 && x < maxX && y < maxY
}

func CountLiberties(g Goban, y, x int) int {
  liberties := 0
  color := g.GetPosition(y, x)
  dx := []int{1, -1, 0, 0}
  dy := []int{0, 0, 1, -1}
  next := list.New()
  marks := g.GetVisitorMarker()
  marks.ClearMarks()
  marks.SetMark(y, x)
  next.PushBack(encode(y, x))
  for next.Len() != 0 {
    front := next.Front()
    next.Remove(front)
    y, x := decode(front.Value.(int))
    for i := 0; i < 4; i++ {
      nx := x + dx[i]
      ny := y + dy[i]
      if valid(ny, nx, g.SizeY(), g.SizeX()) && !marks.IsMarked(ny, nx) {
        switch g.GetPosition(ny, nx) {
        case EMPTY:
          liberties++
          marks.SetMark(ny, nx)
        case color:
          next.PushBack(encode(ny, nx))
          marks.SetMark(ny, nx)
        }
      }
    }
  }
  return liberties
}

