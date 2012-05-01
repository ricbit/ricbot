// Copyright (C) 2012 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Ricardo Bittencourt (bluepenguin@gmail.com)

package engine

import "math/rand"
import "strings"
import "fmt"
import "time"
import "runtime"

type GameState struct {
  goban Goban
  komi float32
  captured_white, captured_black int
}

type ArrayGoban struct {
  size_x, size_y int
  board [][]Color
}

type SliceGoban struct {
  size_x, size_y int
  board []Color
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

func NewSliceGoban(size_y, size_x int, init string) *SliceGoban {
  goban := new(SliceGoban)
  goban.size_x = size_x
  goban.size_y = size_y
  goban.board = make([]Color, size_y * size_x)
  conv := map[byte] Color {
    '.': EMPTY,
    'o': WHITE,
    'x': BLACK,
  }
  for j := 0; j < len(init); j++ {
    goban.board[j] = conv[init[j]]
  }
  return goban
}

func (g *ArrayGoban) SetAllocator(alloc StackAllocator) {
}

func (g *SliceGoban) SetAllocator(alloc StackAllocator) {
}

func (g *ArrayGoban) Copy() Goban {
  new_goban := new(ArrayGoban)
  new_goban.size_x = g.size_x
  new_goban.size_y = g.size_y
  new_goban.board = make([][]Color, new_goban.size_y)
  for i := 0; i < new_goban.size_y; i++ {
    new_goban.board[i] = make([]Color, new_goban.size_x)
    copy(new_goban.board[i], g.board[i])
  }
  return new_goban
}

func (g *SliceGoban) Copy() Goban {
  new_goban := new(SliceGoban)
  new_goban.size_x = g.size_x
  new_goban.size_y = g.size_y
  new_goban.board = make([]Color, new_goban.size_y * new_goban.size_x)
  copy(new_goban.board, g.board)
  return new_goban
}

func (g *ArrayGoban) SizeX() int {
  return g.size_x
}

func (g *SliceGoban) SizeX() int {
  return g.size_x
}

func (g *ArrayGoban) GetVisitorMarker() VisitorMarker {
  return g
}

func (g *SliceGoban) GetVisitorMarker() VisitorMarker {
  return g
}

func (g *ArrayGoban) SizeY() int {
  return g.size_y
}

func (g *SliceGoban) SizeY() int {
  return g.size_y
}

func (g *ArrayGoban) GetColor(y, x int) Color {
  return g.board[y][x] & 0x3
}

func (g *SliceGoban) GetColor(y, x int) Color {
  return g.board[y * g.size_x + x] & 0x3
}

func (g *ArrayGoban) SetColor(y, x int, color Color) {
  g.board[y][x] = g.board[y][x] & (^0x3) | color
}

func (g *SliceGoban) SetColor(y, x int, color Color) {
  g.board[y * g.size_x + x] = g.board[y * g.size_x + x] & (^0x3) | color
}

func (g *ArrayGoban) ClearMarks() {
  for j := 0; j < g.size_y; j++ {
    for i := 0; i < g.size_x; i++ {
      g.board[j][i] &= 0x3
    }
  }
}

func (g *SliceGoban) ClearMarks() {
  for j := 0; j < len(g.board); j++ {
    g.board[j] &= 0x3
  }
}

func (g *ArrayGoban) SetMark(y, x int) {
  g.board[y][x] |= 0x4
}

func (g *SliceGoban) SetMark(y, x int) {
  g.board[y * g.size_x + x] |= 0x4
}

func (g *ArrayGoban) IsMarked(y, x int) bool {
  return g.board[y][x] & 0x4 > 0
}

func (g *SliceGoban) IsMarked(y, x int) bool {
  return g.board[y * g.size_x + x] & 0x4 > 0
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

func iterateNeighbours(g Goban, y, x int, callback func(y, x int)) {
  for i := 0; i < 4; i++ {
    nx, ny := x + dx[i], y + dy[i]
    if valid(ny, nx, g.SizeY(), g.SizeX()) {
      callback(ny, nx)
    }
  }
}

var diagx = []int{1, -1, 1, -1}
var diagy = []int{1, 1, -1, -1}

func iterateDiagonals(g Goban, y, x int, callback func(y, x int)) {
  for i := 0; i < 4; i++ {
    nx, ny := x + diagx[i], y + diagy[i]
    if valid(ny, nx, g.SizeY(), g.SizeX()) {
      callback(ny, nx)
    }
  }
}

func iterateGroup(g Goban, y, x int,
                  group_callback func(ny, nx int),
                  border_callback func(ny, nx int)) {
  color := g.GetColor(y, x)
  //next := NewListStack()
  next := NewSliceStack(g.SizeX() * g.SizeY())
  marks := g.GetVisitorMarker()
  marks.ClearMarks()
  marks.SetMark(y, x)
  group_callback(y, x)
  next.Push(encode(y, x))
  for !next.Empty() {
    y, x := decode(next.Pop())
    iterateNeighbours(g, y, x, func (ny, nx int) {
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

func iterateAll(g Goban, callback func(y, x int)) {
  for j := 0; j < g.SizeY(); j++ {
    for i := 0; i < g.SizeX(); i++ {
      callback(j, i)
    }
  }
}

func iterateAllColor(g Goban, color Color, callback func(y, x int)) {
  iterateAll(g, func (y, x int) {
    if g.GetColor(y, x) == color {
      callback(y, x)
    }
  })
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
  iterateNeighbours(g, y, x, func (ny, nx int) {
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
  iterateNeighbours(g, y, x, func (ny, nx int) {
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
  iterateAllColor(g, EMPTY, func (y, x int) {
    if !Suicide(g, y, x, color) {
      eye_color, ok := SinglePointEye(g, y, x)
      if !ok || eye_color != color {
        callback(y, x)
      }
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

func addCaptured(state *GameState, color Color, captured int) {
  switch color {
  case WHITE:
    state.captured_black += captured
  case BLACK:
    state.captured_white += captured
  }
}

func Play(state *GameState, y, x int, color Color) {
  state.goban.SetColor(y, x, color)
  iterateNeighbours(state.goban, y, x, func (ny, nx int) {
    if state.goban.GetColor(ny, nx) == Opposite(color) {
      if CountLiberties(state.goban, ny, nx) == 0 {
        captured := RemoveGroup(state.goban, ny, nx)
        addCaptured(state, color, captured)
      }
    }
  })
}

func SinglePointEye(g Goban, y, x int) (Color, bool) {
  histogram := make([]int, 4)
  iterateNeighbours(g, y, x, func (ny, nx int) {
    histogram[g.GetColor(ny, nx)]++
  })
  if histogram[EMPTY] > 0 || (histogram[BLACK] > 0 && histogram[WHITE] > 0) {
    return EMPTY, false
  }
  color := Color(WHITE)
  if histogram[BLACK] > 0 {
    color = BLACK
  }
  diagonals := make([]int, 4)
  iterateDiagonals(g, y, x, func (ny, nx int) {
    diagonals[g.GetColor(ny, nx)]++
  })
  return color, diagonals[color] > diagonals[Opposite(color)]
}

type Position struct {
  y, x int
}

func GetMoveList(g Goban, color Color) []Position {
  moves := make([]Position, 0, g.SizeY() * g.SizeX())
  ValidMoves(g, color, func (y, x int) {
    moves = append(moves, Position{y, x})
  })
  return moves
}

func GetRandomMove(g Goban, color Color) (Position, bool) {
  moves := GetMoveList(g, color)
  if len(moves) == 0 {
    return Position{0, 0}, false
  }
  //fmt.Printf("Color %d, moves %v\n", color, moves)
  return moves[rand.Intn(len(moves))], true
}

func dumpState(state *GameState) {
  conv := map[Color] string {
    EMPTY : ".",
    BLACK : "x",
    WHITE : "o",
  }
  fmt.Printf("--- cap white: %d, cap black %d\n",
             state.captured_white, state.captured_black)
  for j := 0; j < state.goban.SizeY(); j++ {
    for i := 0; i < state.goban.SizeX(); i++ {
      fmt.Printf(conv[state.goban.GetColor(j, i)])
    }
    fmt.Printf("\n")
  }
  fmt.Printf("\n")
}

func PlayRandomGame(state *GameState, color Color) {
  limit := state.goban.SizeY() * state.goban.SizeX() * 3
  for i := 0; i < limit; i++ {
    move, ok := GetRandomMove(state.goban, color)
    if !ok {
      _, ok := GetRandomMove(state.goban, Opposite(color))
      if !ok {
        return
      }
      color = Opposite(color)
      continue
    }
    Play(state, move.y, move.x, color)
    //dumpState(state)
    color = Opposite(color)
  }
}

func EstimatePoints(g Goban) (black, white int) {
  points := make([]int, 4)
  iterateAll(g, func (y, x int) {
    color := g.GetColor(y, x)
    if color == EMPTY {
      eye_color, ok := SinglePointEye(g, y, x)
      if ok {
        points[eye_color] += 1
      }
    } else {
      points[color] += 1
    }
  })
  return points[BLACK], points[WHITE]
}

func Winner(state *GameState) Color {
  black, white := EstimatePoints(state.goban)
  if float32(black) + state.komi > float32(white) {
    return BLACK
  }
  return WHITE
}

type MoveStats struct {
  win, total int
}

func copyState(state *GameState) *GameState {
  new_state := new(GameState)
  new_state.komi = state.komi
  new_state.captured_white = state.captured_white
  new_state.captured_black = state.captured_black
  new_state.goban = state.goban.Copy()
  return new_state
}

type GameResult struct {
  move int
  win bool
}

func launchSinglePlay(state *GameState, moves []Position,
                      color Color, ch chan GameResult) {
  for {
    var result GameResult
    result.move = rand.Intn(len(moves))
    copy_state := copyState(state)
    copy_state.goban.SetColor(moves[result.move].y, moves[result.move].x, color)
    PlayRandomGame(copy_state, color)
    result.win = Winner(copy_state) == color
    ch <- result
  }
}

func GetBestMove(state *GameState, color Color, seconds int) (y, x int) {
  moves := GetMoveList(state.goban, color)
  stats := make([]MoveStats, len(moves))
  processors := runtime.NumCPU()
  runtime.GOMAXPROCS(processors)
  ch := make(chan GameResult, processors)
  for i := 0; i < processors; i++ {
    go launchSinglePlay(state, moves, color, ch)
  }
  timeout := make(chan bool)
  go func() {
    time.Sleep(time.Duration(seconds) * time.Second)
    timeout <- true
  }()
  func() {
    for {
      select {
      case <-timeout:
        return
      case result := <-ch:
        stats[result.move].total += 1
        if result.win {
          stats[result.move].win += 1
        }
      }
    }
  }()
  for i := 0; i < len(moves); i++ {
    fmt.Printf("move %d %d : %d / %d = %f\n",
               moves[i].y, moves[i].x, stats[i].win, stats[i].total,
               float32(stats[i].win) / float32(stats[i].total))
  }
  best := 0
  plays := 0
  for i := 0; i < len(moves); i++ {
    if stats[i].win * stats[best].total > stats[best].win * stats[i].total {
      best = i
    }
    plays += stats[i].total
  }
  fmt.Printf("%f plays/s\n", float32(plays) / float32(seconds))
  fmt.Printf("%d stacks\n", slicestacks)
  return moves[best].y, moves[best].x
}

func NewEmptyGameState(y, x int) *GameState {
  return &GameState{NewArrayGoban(y, x, strings.Repeat(".", y * x)), 6.5, 0, 0}
}

func NewGameState(y, x int, komi float32, goban string) *GameState {
    return &GameState{NewSliceGoban(y, x, goban), komi, 0, 0}
}

