package main

import . "./engine"
import "fmt"

func main() {
  //state := NewEmptyGameState(5, 5)
  state := NewGameState(7, 7, "......." +
                              "o..oox." +
                              "oxxox.." +
                              "o.ox..." +
                              "..oxxx." +
                              "..o.oo." +
                              "...o...")
  //PlayRandomGame(state, BLACK)
  y, x := GetBestMove(state, BLACK, 10000)
  fmt.Printf("move %d %d\n", y, x)
}
