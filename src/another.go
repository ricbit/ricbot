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

package main

import . "engine"
import "fmt"

func main() {
  state := NewGameState(6, 6, 0.0, "..xxx." +
                                   ".xooox" +
                                   "xo..ox" +
                                   "xo.oox" +
                                   "xooox." +
                                   ".xxx.x")
  y, x := GetBestMove(state, BLACK, 30)
  fmt.Printf("move %d %d\n", y, x)
}
