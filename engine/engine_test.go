package engine

import . "launchpad.net/gocheck"
import "testing"

func TestNewArrayGoban(t *testing.T) {
  goban := NewArrayGoban(3, 4, ".o.." +
                               "...." +
                               "..x.")
  if goban.SizeX() != 4 || goban.SizeY() != 3 {
    t.Error("wrong size")
  }
  for j := 0; j < 3; j++ {
    for i := 0; i < 4; i++ {
      switch {
      case i == 1 && j == 0:
        if goban.GetColor(j, i) != WHITE {
          t.Error("GetColor white.")
        }
      case i == 2 && j == 2:
        if goban.GetColor(j, i) != BLACK {
          t.Error("GetColor black.")
        }
      default:
        if goban.GetColor(j, i) != EMPTY {
          t.Error("GetColor empty.")
        }
      }
    }
  }
}

func TestArrayGobanVisitorMarker(t *testing.T) {
  goban := NewArrayGoban(1, 1, ".")
  marker := goban.GetVisitorMarker()
  marker.ClearMarks()
  if marker.IsMarked(0, 0) {
    t.Error("Clear not working")
  }
  marker.SetMark(0, 0)
  if !marker.IsMarked(0, 0) {
    t.Error("SetMark not working")
  }
  marker.SetMark(0, 0)
  if !marker.IsMarked(0, 0) {
    t.Error("SetMark not working")
  }
  if goban.GetColor(0, 0) != EMPTY {
    t.Error("GetColor not working")
  }
}

func TestCountLiberties(t *testing.T) {
  goban := NewArrayGoban(3, 9, "o.o..ooo." +
                               ".o.x.o.o." +
                               "..xx.ooxx")
  if CountLiberties(goban, 0, 0) != 2 {
    t.Error("Liberties in the corner")
  }
  if CountLiberties(goban, 0, 2) != 3 {
    t.Error("Liberties in the side")
  }
  if CountLiberties(goban, 1, 1) != 4 {
    t.Error("Liberties in the middle")
  }
  if CountLiberties(goban, 2, 3) != 5 {
    t.Error("Liberties in the middle group")
  }
  if CountLiberties(goban, 0, 5) != 6 {
    t.Error("Liberties in the eye group")
  }
  if CountLiberties(goban, 2, 8) != 1 {
    t.Error("Liberties in the corner group")
  }
}

func TestSuicide(t *testing.T) {
  goban := NewArrayGoban(3, 9, ".x..x.xxx" +
                               "x.xxx.xoo" +
                               ".xo.x.xo.")
  testcases := []struct {
    y, x int
    color Color
    expected bool
  } {
    {0, 0, WHITE, true},
    {0, 2, WHITE, false},
    {2, 3, WHITE, true},
    {2, 8, WHITE, true},
    {0, 0, BLACK, false},
    {0, 2, BLACK, false},
    {2, 3, BLACK, false},
    {2, 8, BLACK, false},
  }
  for _, tc := range testcases {
    if Suicide(goban, tc.y, tc.x, tc.color) != tc.expected {
      t.Errorf("Error in %d, %d for color %v, expecting %v",
               tc.y, tc.x, tc.color, tc.expected)
    }
  }
}

func comparePoints(t *testing.T, expected, actual []Position) {
  if len(expected) != len(actual) {
    t.Errorf("Different sizes, expected %v actual %v", expected, actual)
    return
  }
  for _, exp_value := range expected {
    found := false
    for _, act_value := range actual {
      if exp_value == act_value {
        found = true
      }
    }
    if !found {
      t.Errorf("Different values, expected %v actual %v", expected, actual)
      return
    }
  }
}

func TestValidMoves(t *testing.T) {
  goban := NewArrayGoban(3, 4, ".xox" +
                               "xo.x" +
                               ".xx.")
  expected_white := []Position {
    {0, 0},
  }
  expected_black := []Position {
    {0, 0}, {1, 2}, {2, 0}, {2, 3},
  }
  actual_white := GetMoveList(goban, WHITE)
  actual_black := GetMoveList(goban, BLACK)
  comparePoints(t, expected_white, actual_white)
  comparePoints(t, expected_black, actual_black)
}

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type S struct{}
var _ = Suite(&S{})

func ToString(g Goban) string {
  output := ""
  conv := map[Color] string {
    EMPTY : ".",
    BLACK : "x",
    WHITE : "o",
  }
  for j := 0; j < g.SizeY(); j++ {
    for i := 0; i < g.SizeX(); i++ {
      output += conv[g.GetColor(j,i)]
    }
  }
  return output
}

func (s *S) TestRemoveGroup(c *C) {
  goban1 := NewArrayGoban(3, 4, "xxox" +
                                "xo.x" +
                                ".xx.")
  c.Check(RemoveGroup(goban1, 0, 0), Equals, 3)
  c.Check(ToString(goban1), Equals, "..ox" +
                                    ".o.x" +
                                    ".xx.")

  goban2 := NewArrayGoban(3, 4, "xxox" +
                                "xo.x" +
                                ".xx.")
  c.Check(RemoveGroup(goban2, 0, 2), Equals, 1)
  c.Check(ToString(goban2), Equals, "xx.x" +
                                    "xo.x" +
                                    ".xx.")
}

func (s *S) TestPlay(c *C) {
  state1 := NewGameState(3, 4, "xxox" +
                               "xo.x" +
                               ".xx.")
  Play(state1, 1, 2, BLACK)
  c.Check(ToString(state1.goban), Equals, "xx.x" +
                                          "x.xx" +
                                          ".xx.")
  c.Check(state1.captured_white, Equals, 2)
  c.Check(state1.captured_black, Equals, 0)

  state2 := NewGameState(3, 4, "xxox" +
                               "xo.x" +
                               ".xx.")
  Play(state2, 2, 0, WHITE)
  c.Check(ToString(state2.goban), Equals, "..ox" +
                                          ".o.x" +
                                          "oxx.")
  c.Check(state2.captured_white, Equals, 0)
  c.Check(state2.captured_black, Equals, 3)
}

func (s *S) TestSinglePointEye(c *C) {
  goban := NewArrayGoban(3, 4, "..ox" +
                               "oxxx" +
                               "x.x.")
  color, ok := SinglePointEye(goban, 2, 3)
  c.Check(ok, Equals, true)
  c.Check(color, Equals, Color(BLACK))

  color, ok = SinglePointEye(goban, 0, 0)
  c.Check(ok, Equals, false)

  color, ok = SinglePointEye(goban, 2, 1)
  c.Check(ok, Equals, false)
}

func (s *S) TestEstimatePoints(c *C) {
  goban := NewArrayGoban(3, 4, ".o.x" +
                               "ooxx" +
                               "xxx.")
  black, white := EstimatePoints(goban)
  c.Check(black, Equals, 7)
  c.Check(white, Equals, 4)
}

func (s *S) TestGetRandomMove(c *C) {
  goban := NewArrayGoban(1, 4, "....")
  histogram := make([]int, 4)
  lots := 10000
  for i := 0; i < lots; i++ {
    move, ok := GetRandomMove(goban, WHITE)
    c.Check(ok, Equals, true)
    histogram[move.x]++
  }
  mean := float32(lots) / float32(4)
  variance := float32(0.0)
  for _, val := range histogram {
    tmp := float32(val) - mean
    variance += tmp * tmp
  }
  variance /= float32(lots)
  if variance > 0.1 {
    c.Errorf("Variance is too high: %f\n", variance)
  }
}

func (s *S) TestGetRandomMoveEmpty(c *C) {
  goban := NewArrayGoban(1, 4, "xoxo")
  _, ok := GetRandomMove(goban, WHITE)
  c.Check(ok, Equals, false)
}
