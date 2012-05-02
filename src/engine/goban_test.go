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

import . "launchpad.net/gocheck"

func (s *S) TestArrayGoban(c *C) {
  goban := NewArrayGoban(3, 4)
  checkGoban(c, goban)
}

func (s *S) TestSliceGoban(c *C) {
  goban := NewSliceGoban(3, 4)
  checkGoban(c, goban)
}

func checkGoban(c *C, goban Goban) {
  FromString(goban, ".o.." +
                    "...." +
                    "..x.")
  c.Check(goban.SizeX(), Equals, 4)
  c.Check(goban.SizeY(), Equals, 3)
  for j := 0; j < 3; j++ {
    for i := 0; i < 4; i++ {
      switch {
      case i == 1 && j == 0:
        c.Check(goban.GetColor(j, i), Equals, Color(WHITE))
      case i == 2 && j == 2:
        c.Check(goban.GetColor(j, i), Equals, Color(BLACK))
      default:
        c.Check(goban.GetColor(j, i), Equals, Color(EMPTY))
      }
    }
  }
}

func (s *S) TestArrayGobanVisitorMarker(c *C) {
  goban := NewArrayGoban(1, 1)
  checkGobanVisitorMarker(c, goban)
}

func (s *S) TestSliceGobanVisitorMarker(c *C) {
  goban := NewSliceGoban(1, 1)
  checkGobanVisitorMarker(c, goban)
}

func checkGobanVisitorMarker(c *C, goban Goban) {
  FromString(goban , ".")
  marker := goban.GetVisitorMarker()
  marker.ClearMarks()
  c.Check(marker.IsMarked(0, 0), Equals, false)
  marker.SetMark(0, 0)
  c.Check(marker.IsMarked(0, 0), Equals, true)
  marker.SetMark(0, 0)
  c.Check(marker.IsMarked(0, 0), Equals, true)
  c.Check(goban.GetColor(0, 0), Equals, Color(EMPTY))
}


