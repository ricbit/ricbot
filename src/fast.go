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

import "engine"
import "fmt"
import "runtime/pprof"
import "os"

func main() {
  f, _ := os.Create("profile")
  pprof.StartCPUProfile(f)
  state := engine.NewGameState(5, 4, 0.0, "ox.." +
                                          "ox.." +
                                          "ox.x" +
                                          "oxxo" +
                                          "oooo")
  y, x := engine.GetBestMove(state, engine.BLACK, 5)
  fmt.Printf("move %d %d\n", y, x)
  pprof.StopCPUProfile()
}
