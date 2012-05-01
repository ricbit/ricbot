package main

import . "./engine"
import "fmt"

func main() {
  state := NewGameState(6, 6, "..xxx." +
                              ".xooox" +
                              "xo..ox" +
                              "xo.oox" +
                              "xooox." +
                              ".xxx.x")
  y, x := GetBestMove(state, BLACK, 1000)
  fmt.Printf("move %d %d\n", y, x)
}
