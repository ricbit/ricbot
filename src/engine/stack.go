package engine

type Stack interface {
  Push(int)
  Pop() int
  Empty() bool
}

