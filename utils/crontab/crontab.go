package crontab

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/atomic"
	"llm-PhotoMagic/go-common/utils/help"
	"sync"
	"time"
)

type Cron struct {
	cron    *cron.Cron
	jobs    sync.Map
	error   error
	running atomic.Bool
	restart atomic.Bool
	mux     sync.Mutex
}

func NewCron() *Cron {
	return &Cron{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (c *Cron) Start() {
	if c.IsError() {
		fmt.Printf("cron start fail error: %s\n", c.error)
		return
	}
	if c.IsRestart() {
		c.Restart()
		return
	}
	if c.IsRunning() {
		return
	}
	c.cron.Start()
	c.running.Store(true)
}

func (c *Cron) Stop() {
	if !c.IsRunning() {
		return
	}
	c.cron.Stop()
	c.running.Store(false)
}

func (c *Cron) Restart() {
	fmt.Println("cron restart start")
	if c.IsError() {
		fmt.Printf("cron start fail error: %s\n", c.error)
		return
	}
	if c.IsRunning() {
		c.cron.Stop()
		time.Sleep(500 * time.Millisecond)
	}
	c.cron.Start()
	c.restart.Store(false)
	fmt.Println("cron restart success")
}

func (c *Cron) IsRunning() bool {
	return c.running.Load()
}
func (c *Cron) IsRestart() bool {
	return c.restart.Load()
}
func (c *Cron) IsError() bool {
	return c.error != nil
}

func (c *Cron) Jobs() []*Job {
	var jobs []*Job
	c.jobs.Range(func(key, value interface{}) bool {
		jobs = append(jobs, value.(*Job))
		return true
	})
	return jobs
}

func (c *Cron) AddFunc(name, spec string, cmd func()) *Cron {
	//覆盖任务
	c.RemoveFunc(name)
	job := NewJob(name, spec, cmd)
	entryID, err := c.cron.AddFunc(spec, job.Running)
	if err != nil {
		fmt.Printf("add cron fail: %s\n", name)
		c.error = err
		return c
	}
	c.jobs.Store(name, job.WithEntryID(entryID))
	return c
}

func (c *Cron) RemoveFunc(name string) *Cron {
	if job, ok := c.jobs.LoadAndDelete(name); ok {
		c.cron.Remove(job.(*Job).entryID)
		c.restart.Store(true)
	}
	return c
}

type State int32

const (
	Ready State = iota + 1
	Running
)

type Job struct {
	entryID     cron.EntryID
	name        string
	spec        string
	cmd         func()
	state       State
	lastRunTime string
	lastDurable string
}

func NewJob(name, spec string, cmd func()) *Job {
	return &Job{
		name:  name,
		spec:  spec,
		cmd:   cmd,
		state: Ready,
	}
}
func (j *Job) WithEntryID(id cron.EntryID) *Job {
	j.entryID = id
	return j
}

func (j *Job) Running() {
	j.state = Running
	j.lastRunTime = time.Now().Format(help.FormatTime)
	start := time.Now().Local().UnixMilli()
	defer func() {
		j.state = Ready
		j.lastDurable = fmt.Sprintf("%0.3f", float64(time.Now().Local().UnixMilli()-start)/1000)
	}()
	j.cmd()
}
