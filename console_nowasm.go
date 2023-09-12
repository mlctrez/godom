package godom

import (
	"fmt"
)

var _ ConsoleApi = (*consoleValue)(nil)

type consoleValue struct{}

func (c *consoleValue) Log(args ...any)   { fmt.Println(args...) }
func (c *consoleValue) Debug(args ...any) { fmt.Println(args...) }
func (c *consoleValue) Info(args ...any)  { fmt.Println(args...) }
func (c *consoleValue) Error(args ...any) { fmt.Println(args...) }
