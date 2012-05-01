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

import "container/list"


// A simple stack of ints.
type Stack interface {
  Push(int)
  Pop() int
  Empty() bool
}

// An allocator of Stacks, may be used to cache instances.
type StackAllocator interface {
  GetStack() Stack
  Mark()
  Release()
}

// --------------------------
// Stack implementation using list, slow.

type ListStack struct {
  stack *list.List
}

func (s *ListStack) Push(value int) {
  s.stack.PushBack(value)
}

func (s *ListStack) Pop() int {
  value := s.stack.Front()
  s.stack.Remove(value)
  return value.Value.(int)
}

func (s *ListStack) Empty() bool {
  return s.stack.Len() == 0
}

func NewListStack() *ListStack {
  stack := new(ListStack)
  stack.stack = list.New()
  return stack
}

// --------------------------
// Stack implementation using slices, faster.

type SliceStack struct {
  stack []int
}

func (s *SliceStack) Push(value int) {
  s.stack = append(s.stack, value)
}

func (s *SliceStack) Pop() int {
  value := s.stack[0]
  s.stack = s.stack[1:]
  return value
}

func (s *SliceStack) Empty() bool {
  return len(s.stack) == 0
}

var slicestacks int = 0

func NewSliceStack(capacity int) *SliceStack {
  slicestacks++
  stack := new(SliceStack)
  stack.stack = make([]int, 0, capacity)
  return stack
}

// --------------------------
// A dumb allocator using no caches.

type DumbAllocator struct {
}

func (a *DumbAllocator) GetStack() Stack {
  return NewListStack()
}

func (a *DumbAllocator) Mark() {
}

func (a *DumbAllocator) Release() {
}


