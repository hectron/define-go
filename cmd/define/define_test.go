package main

import (
  "fmt"
  "os"
  "os/exec"
  "testing"
)

type MockLibrary interface {
  Define(string)
}

func TestMainCommandWithoutArgs(t *testing.T) {
  // This will be hit in a subprocess
  if os.Getenv("INVOKED_IN_TEST") == "1" {
    fmt.Println(os.Args)
    runCLI(os.Args)
    return
  }
  // Run the test in a subprocess
  cmd := exec.Command(os.Args[0], "-test.run=TestMainCommandWithoutArgs")
  cmd.Env = append(os.Environ(), "INVOKED_IN_TEST=1")
  err := cmd.Run()

  if e, ok := err.(*exec.ExitError); ok && !e.Success() {
    return
  }
  t.Fatalf("process ran with err %v, want exit status 1", err)
}
