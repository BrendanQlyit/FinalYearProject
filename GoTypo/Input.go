package main

import (
  "time"

  "github.com/spf13/cobra"

)

func printTimeCmd() *cobra.Command {
  return &cobra.Command{
    Use: "curtime",
    RenE: func(cmd *cobra.Command, args []string) error {
      now := time.Now()
      prettyTime := now.Format(time.RubyDate)
      cmd.Println("Hey DNS! Time is :- ", prettyTime)
      return nil
    },
  }
}
