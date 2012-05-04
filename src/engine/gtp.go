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

type GTP interface {
  BoardSize(size int)
  ClearBoard()
  Play(y, x int, color Color)
  GenMove(color Color) (y, x int, pass bool)
  Komi(komi float32)
}

func (s *GameState) BoardSize(size int) {
  s.goban = NewSliceGoban(size, size)
}

func (s *GameState) ClearBoard() {
  iterateAll(s.goban, func (y, x int) {
    s.goban.SetColor(y, x, EMPTY)
  })
}

func (s *GameState) Play(y, x int, color Color) {
  Play(s, y, x, color)
}

func (s *GameState) Komi(komi float32) {
  s.komi = komi
}

func (s *GameState) GenMove(color Color) (y, x int, pass bool) {
  return GetBestMove(s, color, 2)
}

