package crontab

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCron(t *testing.T) {
	c := NewCron()
	c.AddFunc("test1", "*/3 * * * * *", func() {
		fmt.Println("in")
	}).Start()
	time.Sleep(10 * time.Second)
}

func TestCron_AddFunc(t *testing.T) {
	c := NewCron()
	c.AddFunc("test1", "*/2 * * * * *", func() {
		fmt.Println("in")
	}).Start()
	fmt.Println("here")
	time.Sleep(10 * time.Second)
	c.AddFunc("test1", "*/4 * * * * *", func() {
		fmt.Println("in2")
	}).Restart()
	time.Sleep(10 * time.Second)
}
