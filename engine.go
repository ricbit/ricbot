package engine

import "container/list"

const (
  EMPTY = iota
  BLACK
  WHITE
  INVALID
)

type Color int

type Goban interface {
  SizeX() int
  SizeY() int
  GetColor(y, x int) Color
  SetColor(y, x int, color Color)
  GetVisitorMarker() VisitorMarker
}

type VisitorMarker interface {
  ClearMarks()
  SetMark(y, x int)
  IsMarked(y, x int) bool
}

type GameState struct {
  goban Goban
  komi int
  captured_white, captured_black int
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
  board [][]Color
}

func NewArrayGoban(size_y, size_x int, init string) *ArrayGoban {
  goban := new(ArrayGoban)
  goban.size_x = size_x
  goban.size_y = size_y
  goban.board = make([][]Color, size_y)
  conv := map[byte] Color {
    '.': EMPTY,
    'o': WHITE,
    'x': BLACK,
  }
  for j := 0; j < size_y; j++ {
    goban.board[j] = make([]Color, size_x)
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

func (g *ArrayGoban) GetColor(y, x int) Color {
  return g.board[y][x] & 0x3
}

func (g *ArrayGoban) SetColor(y, x int, color Color) {
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

func iterateGroup(g Goban, y, x int,
                  group_callback func(ny, nx int),
                  border_callback func(ny, nx int)) {
  color := g.GetColor(y, x)
  next := NewListStack()
  marks := g.GetVisitorMarker()
  marks.ClearMarks()
  marks.SetMark(y, x)
  group_callback(y, x)
  next.Push(encode(y, x))
  for !next.Empty() {
    y, x := decode(next.Pop())
    iterateNeighbours(y, x, g.SizeY(), g.SizeX(), func (ny, nx int) {
      if !marks.IsMarked(ny, nx) {
        if g.GetColor(ny, nx) == color {
          group_callback(ny, nx)
          next.Push(encode(ny, nx))
        } else {
          border_callback(ny, nx)
        }
        marks.SetMark(ny, nx)
      }
    })
  }
}

func CountLiberties(g Goban, y, x int) int {
  liberties := 0
  iterateGroup(g, y, x, func (ny, nx int) {}, func (ny, nx int) {
    if g.GetColor(ny, nx) == EMPTY {
      liberties++
    }
  })
  return liberties
}

func iterateAll(g Goban, color Color, callback func(y, x int)) {
  for j := 0; j < g.SizeY(); j++ {
    for i := 0; i < g.SizeX(); i++ {
      if g.GetColor(j, i) == color {
        callback(j, i)
      }
    }
  }
}

func Opposite(color Color) Color {
  if color == BLACK {
    return WHITE
  }
  return BLACK
}

func Suicide(g Goban, y, x int, color Color) bool {
  // It's not suicide if you have an empty cell next to you.
  liberties := 0
  iterateNeighbours(y, x, g.SizeY(), g.SizeX(), func (ny, nx int) {
    if g.GetColor(ny, nx) == EMPTY {
      liberties++
    }
  })
  if liberties > 0 {
    return false
  }
  g.SetColor(y, x, color)
  defer func() { g.SetColor(y, x, EMPTY) }()
  // It's not suicide if you are capturing something.
  suicide := true
  opponent := Opposite(color)
  iterateNeighbours(y, x, g.SizeY(), g.SizeX(), func (ny, nx int) {
    if suicide && g.GetColor(ny, nx) == opponent {
      liberties := CountLiberties(g, ny, nx)
      if liberties == 0 {
        suicide = false
      }
    }
  })
  if !suicide {
    return false
  }
  // It's not suicide if you connect to a group with liberties.
  return CountLiberties(g, y, x) == 0
}

func ValidMoves(g Goban, color Color, callback func (y, x int)) {
  iterateAll(g, EMPTY, func (y, x int) {
    if !Suicide(g, y, x, color) {
      callback(y, x)
    }
  })
}

func RemoveGroup(g Goban, y, x int) int {
  captured := 0
  iterateGroup(g, y, x, func (ny, nx int) {
    captured++
    g.SetColor(ny, nx, EMPTY)
  }, func (ny, nx int) {})
  return captured
}

func Play(state *GameState, y, x int, color Color) {
  state.goban.SetColor(y, x, color)
  my, mx := state.goban.SizeY(), state.goban.SizeX()
  iterateNeighbours(y, x, my, mx, func (ny, nx int) {
    if state.goban.GetColor(ny, nx) == Opposite(color) {
      if CountLiberties(state.goban, ny, nx) == 0 {
        RemoveGroup(state.goban, ny, nx)
      }
    }
  })
}
