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
