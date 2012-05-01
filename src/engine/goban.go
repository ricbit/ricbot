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

// Possible contents of each cell on the go board.
const (
  EMPTY = iota
  BLACK
  WHITE
  INVALID
)

type Color int

// Generic Goban (go board) interface.
type Goban interface {
  SizeX() int
  SizeY() int
  GetColor(y, x int) Color
  SetColor(y, x int, color Color)
  GetVisitorMarker() VisitorMarker
  GetStack() Stack
  SetStack(stack Stack)
  Copy() Goban
}

type VisitorMarker interface {
  ClearMarks()
  SetMark(y, x int)
  IsMarked(y, x int) bool
}

// --------------------------
// Goban implementation using a 2D array.

type ArrayGoban struct {
  size_x, size_y int
  board [][]Color
  stack Stack
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
  goban.stack = NewSliceStack(size_x * size_y)
  return goban
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

func (g *ArrayGoban) GetStack() Stack {
  return g.stack
}

func (g *ArrayGoban) SetStack(stack Stack) {
  g.stack = stack
}

// --------------------------
// Goban implementation using a single slice.

type SliceGoban struct {
  size_x, size_y int
  board []Color
  stack Stack
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
  goban.stack = NewSliceStack(size_x * size_y)
  return goban
}

func (g *SliceGoban) Copy() Goban {
  new_goban := new(SliceGoban)
  new_goban.size_x = g.size_x
  new_goban.size_y = g.size_y
  new_goban.board = make([]Color, new_goban.size_y * new_goban.size_x)
  copy(new_goban.board, g.board)
  return new_goban
}

func (g *SliceGoban) SizeX() int {
  return g.size_x
}

func (g *SliceGoban) GetVisitorMarker() VisitorMarker {
  return g
}

func (g *SliceGoban) SizeY() int {
  return g.size_y
}

func (g *SliceGoban) GetColor(y, x int) Color {
  return g.board[y * g.size_x + x] & 0x3
}

func (g *SliceGoban) SetColor(y, x int, color Color) {
  g.board[y * g.size_x + x] = g.board[y * g.size_x + x] & (^0x3) | color
}

func (g *SliceGoban) ClearMarks() {
  for j := 0; j < len(g.board); j++ {
    g.board[j] &= 0x3
  }
}

func (g *SliceGoban) SetMark(y, x int) {
  g.board[y * g.size_x + x] |= 0x4
}

func (g *SliceGoban) IsMarked(y, x int) bool {
  return g.board[y * g.size_x + x] & 0x4 > 0
}

func (g *SliceGoban) GetStack() Stack {
  return g.stack
}

func (g *SliceGoban) SetStack(stack Stack) {
  g.stack = stack
}

