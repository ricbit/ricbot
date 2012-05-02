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

//import "engine"
import "io"
import "bufio"
import "strings"
import "bytes"
import "fmt"

type GTP interface {
  Run(reader io.Reader, writer io.Writer)
}

type Session struct {

}

type Handler func ([]string) string

var commands = map[string] Handler {
  "name" : Name,
  "protocol_version" : ProtocolVersion,
  "version" : Version,
}

func (s *Session) Run(reader io.Reader, writer io.Writer) {
  buf := bufio.NewReader(reader)
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
        output := make([]string, 0)
        for command, _ := range commands {
          output = append(output, command)
        }
        fmt.Fprint(writer, "= " + strings.Join(output, "\n") + "\n\n")
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

func ListCommands(args []string) string {
  output := make([]string, 0)
  for command, _ := range commands {
    output = append(output, command)
  }
  return strings.Join(output, "\n")
}
