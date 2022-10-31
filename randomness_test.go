package main

import (
	"testing"
	"fmt"
)

func TestAssignment(t *testing.T) {
	config, _ := readConf("./example-config.yml")
	taylorCount := 0
	lindsayCount := 0
	ryanCount := 0
	for i := 0; i < 100; i++ {
		config.Seed = config.Seed + int64(i)
		head := createAssignments(config)
		if head.Next.Name == "Taylor" {
			taylorCount++
		} else if head.Next.Name == "Lindsay" {
			lindsayCount++
		} else {
			ryanCount++
		}
	}
	fmt.Println(ryanCount)
	fmt.Println(lindsayCount)
	fmt.Println(taylorCount)
}
