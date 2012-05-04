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
import "os"

// The game state.
type GameState struct {
  goban Goban
  komi float32
  captured_white, captured_black int
}

// --------------------------
// Iterators over the Goban.

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
  //next := NewSliceStack(g.SizeX() * g.SizeY())
  next := g.GetStack()
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

func CountLiberties(g Goban, y, x int) int {
  liberties := 0
  iterateGroup(g, y, x, func (ny, nx int) {}, func (ny, nx int) {
    if g.GetColor(ny, nx) == EMPTY {
      liberties++
    }
  })
  return liberties
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
  // It's not suicide if you connect to a group with liberties.
  if CountLiberties(g, y, x) > 0 {
    return false
  }
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
  return suicide
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

var dump bool = false

func GetMoveList(g Goban, color Color) []Position {
  moves := make([]Position, 0, g.SizeY() * g.SizeX())
  ValidMoves(g, color, func (y, x int) {
    moves = append(moves, Position{y, x})
  })
  if (dump) {
    fmt.Fprintf(os.Stderr, "# %v\n", moves)
  }
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

func dumpGoban(goban Goban) {
  conv := map[Color] string {
    EMPTY : ".",
    BLACK : "x",
    WHITE : "o",
  }
  for j := 0; j < goban.SizeY(); j++ {
    for i := 0; i < goban.SizeX(); i++ {
      fmt.Fprintf(os.Stderr, conv[goban.GetColor(j, i)])
    }
    fmt.Fprintf(os.Stderr, "\n")
  }
  fmt.Fprintf(os.Stderr, "\n")
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
  if float32(black) > float32(white) + state.komi {
    return BLACK
  }
  return WHITE
}

type MoveStats struct {
  win, total int
  rate float64
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
  stack := NewSliceStack(state.goban.SizeX() * state.goban.SizeY())
  stats := make([]MoveStats, len(moves))
  for i := 0; i < len(stats); i++ {
    stats[i].total = 2
    stats[i].win = 1
    stats[i].rate = 0.5
  }
  for {
    var total float64 = 0.0
    for i := 0; i < len(stats); i++ {
      total += stats[i].rate
    }
    var p = rand.Float64() * total
    var result GameResult
    for i := 0; i < len(stats); i++ {
      result.move = i
      if p < stats[i].rate {
        break
      }
      p -= stats[i].rate
    }
    copy_state := copyState(state)
    copy_state.goban.SetColor(moves[result.move].y, moves[result.move].x, color)
    copy_state.goban.SetStack(stack)
    PlayRandomGame(copy_state, color)
    result.win = Winner(copy_state) == color
    stats[result.move].total += 1
    if result.win {
      stats[result.move].win += 1
    }
    stats[result.move].rate = float64(stats[result.move].win) /
                              float64(stats[result.move].total)
    ch <- result
  }
}

func GetBestMove(state *GameState, color Color, seconds int) (
    y, x int, pass bool) {
  dump = true
  dumpGoban(state.goban)
  moves := GetMoveList(state.goban, color)
  dump = false
  if len(moves) == 0 {
    return 0, 0, true
  }
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
    fmt.Fprintf(os.Stderr, "# move %d %d : %d / %d = %f\n",
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
  fmt.Fprintf(os.Stderr,"# %f plays/s\n", float32(plays) / float32(seconds))
  fmt.Fprintf(os.Stderr,"# %d stacks\n", slicestacks)
  return moves[best].y, moves[best].x, false
}

func NewEmptyGameState(y, x int) *GameState {
  goban := NewArrayGoban(y, x)
  FromString(goban, strings.Repeat(".", y * x))
  return &GameState{goban, 6.5, 0, 0}
}

func NewGameState(y, x int, komi float32, init string) *GameState {
  goban := NewSliceGoban(y, x)
  FromString(goban, init)
  return &GameState{goban, komi, 0, 0}
}

