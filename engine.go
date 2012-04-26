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
  SetPosition(y, x int, color Position)
  GetVisitorMarker() VisitorMarker
}

type VisitorMarker interface {
  ClearMarks()
  SetMark(y, x int)
  IsMarked(y, x int) bool
}

type Stack interface {
  Push(int)
  Pop() int
  Empty() bool
}

type ListStack struct {
  stack *list.List
}

func (s *ListStack) Push(value int) {
  s.stack.PushBack(value)
}

func (s *ListStack) Pop() int {
  value := s.stack.Front()
  s.stack.Remove(value)
  return value.Value.(int)
}

func (s *ListStack) Empty() bool {
  return s.stack.Len() == 0
}

func NewListStack() *ListStack {
  stack := new(ListStack)
  stack.stack = list.New()
  return stack
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
  return g.board[y][x] & 0x3
}

func (g *ArrayGoban) SetPosition(y, x int, color Position) {
  g.board[y][x] = g.board[y][x] & (^0x3) | color
}

func (g *ArrayGoban) ClearMarks() {
  for j := 0; j < g.size_y; j++ {
    for i := 0; i < g.size_x; i++ {
      g.board[j][i] &= 0x3
    }
  }
}

func (g *ArrayGoban) SetMark(y, x int) {
  g.board[y][x] |= 0x4
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

var dx = []int{1, -1, 0, 0}
var dy = []int{0, 0, 1, -1}

func iterateNeighbours(y, x, maxY, maxX int, callback func(y, x int)) {
  for i := 0; i < 4; i++ {
    nx, ny := x + dx[i], y + dy[i]
    if valid(ny, nx, maxY, maxX) {
      callback(ny, nx)
    }
  }
}

func CountLiberties(g Goban, y, x int) int {
  liberties := 0
  color := g.GetPosition(y, x)
  next := NewListStack()
  marks := g.GetVisitorMarker()
  marks.ClearMarks()
  marks.SetMark(y, x)
  next.Push(encode(y, x))
  for !next.Empty() {
    y, x := decode(next.Pop())
    iterateNeighbours(y, x, g.SizeY(), g.SizeX(), func (ny, nx int) {
      if !marks.IsMarked(ny, nx) {
        switch g.GetPosition(ny, nx) {
        case EMPTY:
          liberties++
          marks.SetMark(ny, nx)
        case color:
          next.Push(encode(ny, nx))
          marks.SetMark(ny, nx)
        }
      }
    })
  }
  return liberties
}

func iterateAll(g Goban, color Position, callback func(y, x int)) {
  for j := 0; j < g.SizeY(); j++ {
    for i := 0; i < g.SizeX(); i++ {
      if g.GetPosition(j, i) == color {
        callback(j, i)
      }
    }
  }
}

func Opposite(color Position) Position {
  if color == BLACK {
    return WHITE
  }
  return BLACK
}

func Suicide(g Goban, y, x int, color Position) bool {
  liberties := 0
  iterateNeighbours(y, x, g.SizeY(), g.SizeX(), func (ny, nx int) {
    if g.GetPosition(ny, nx) == EMPTY {
      liberties++
    }
  })
  if liberties > 0 {
    return false
  }
  g.SetPosition(y, x, color)
  defer func() { g.SetPosition(y, x, EMPTY) }()
  suicide := true
  opponent := Opposite(color)
  iterateNeighbours(y, x, g.SizeY(), g.SizeX(), func (ny, nx int) {
    if suicide && g.GetPosition(ny, nx) == opponent {
      liberties := CountLiberties(g, ny, nx)
      if liberties == 0 {
        suicide = false
      }
    }
  })
  if !suicide {
    return false
  }
  return CountLiberties(g, y, x) == 0
}

func ValidMoves(g Goban, color Position, callback func (y, x int)) {
  iterateAll(g, EMPTY, func (y, x int) {
    if !Suicide(g, y, x, color) {
      callback(y, x)
    }
  })
}
