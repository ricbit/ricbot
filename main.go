package main

import . "./engine"
import "fmt"

func main() {
  state := NewGameState(7, 7, 0.0, "......." +
                                   "o..oox." +
                                   "oxxox.." +
                                   "o.ox..." +
                                   "..oxxx." +
                                   "..o.oo." +
                                   "...o...")
  y, x := GetBestMove(state, BLACK, 150)
  fmt.Printf("move %d %d\n", y, x)
}
