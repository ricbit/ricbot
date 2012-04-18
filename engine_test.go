package engine

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
        if goban.GetPosition(j, i) != WHITE {
          t.Error("GetPosition white.")
        }
      case i == 2 && j == 2:
        if goban.GetPosition(j, i) != BLACK {
          t.Error("GetPosition black.")
        }
      default:
        if goban.GetPosition(j, i) != EMPTY {
          t.Error("GetPosition empty.")
        }
      }
    }
  }
}

func TestArrayGobanVisitorMarker(t *testing.T) {
  goban := NewArrayGoban(1, 1, ".")
  marker := goban.GetVisitorMarker()
  marker.ClearMarks()
  if marker.IsMarked(0,0) {
    t.Error("Clear not working")
  }
  marker.SetMark(0,0)
  if !marker.IsMarked(0,0) {
    t.Error("SetMark not working")
  }
  marker.SetMark(0,0)
  if !marker.IsMarked(0,0) {
    t.Error("SetMark not working")
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
    color Position
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
