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

package gtp

import "engine"
import "io"
import "bufio"
import "strings"
import "bytes"
import "fmt"
import "strconv"

type Driver interface {
  Run(reader io.Reader, writer io.Writer)
}

type Session struct {
}

// TODO(ricbit): Remove this global var.
var state *engine.GameState

type Handler func ([]string) string

var commands = map[string] Handler {
  "name" : Name,
  "protocol_version" : ProtocolVersion,
  "version" : Version,
  "boardsize" : BoardSize,
  "clear_board" : ClearBoard,
  "play" : Play,
  "genmove" : GenMove,
  "komi" : Komi,
}

func (s *Session) Run(reader io.Reader, writer io.Writer) {
  buf := bufio.NewReader(reader)
  state = new(engine.GameState)
  for {
    line, _, ok := buf.ReadLine()
    if ok != nil {
      return
    }
    args := strings.Split(bytes.NewBuffer(line).String(), " ")
    switch args[0] {
      case "quit":
        return
      case "list_commands":
        fmt.Fprint(writer, "= " + ListCommands() + "\n\n")
      default:
        if handler, ok := commands[args[0]]; ok {
          fmt.Fprint(writer, "= " + handler(args[1:]) + "\n\n")
        } else {
          fmt.Fprint(writer, "?\n\n")
        }
    }
  }
}

func Name(args []string) string {
  return "ricbot"
}

func ProtocolVersion(args []string) string {
  return "2"
}

func Version(args []string) string {
  return "1.0"
}

func ListCommands() string {
  output := make([]string, 0)
  for command, _ := range commands {
    output = append(output, command)
  }
  return strings.Join(output, "\n")
}

func BoardSize(args []string) string {
  size, _ := strconv.Atoi(args[0])
  state.BoardSize(size)
  return ""
}

func ClearBoard(args []string) string {
  state.ClearBoard()
  return ""
}

func stringToColor(s string) engine.Color {
  switch strings.ToLower(s)[0] {
  case 'b':
    return engine.BLACK
  case 'w':
    return engine.WHITE
  }
  return engine.EMPTY
}

func stringToPosition(s string) (y, x int) {
  x = int(strings.ToLower(s)[0]) - int('a')
  temp, _ := strconv.Atoi(s[1:])
  y = temp - 1
  return y, x
}

func Play(args []string) string {
  color := stringToColor(args[0])
  if strings.ToLower(args[1]) == "pass" {
    return ""
  }
  y, x := stringToPosition(args[1])
  state.Play(y, x, color)
  return ""
}

func GenMove(args []string) string {
  color := stringToColor(args[0])
  y, x, pass := state.GenMove(color)
  if pass {
    return "pass"
  }
  state.Play(y, x, color)
  return fmt.Sprintf("%s%d", string(int('a') + x), y + 1)
}

func Komi(args []string) string {
  komi, _ := strconv.ParseFloat(args[0], 32)
  state.Komi(float32(komi))
  return ""
}



