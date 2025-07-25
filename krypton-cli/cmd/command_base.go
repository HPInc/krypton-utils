package cmd

import (
	"flag"
	"time"

	"cli/common"
	"cli/config"
)

const (
	DefaultBatchSize  = 100
	BatchSizeMax      = 1000
	DefaultRunCount   = 1
	DefaultRetryCount = 10
)

type CmdFlags struct {
	fs          *flag.FlagSet
	count       *uint
	batchSize   *uint
	parallel    *bool
	verbose     *bool
	retryCount  *uint
	stdin       *bool
	doc         *string
	apiBasePath *string
}

type CmdBase struct {
	Name        string
	Count       uint
	Parallel    bool
	BatchSize   uint
	Verbose     bool
	RetryCount  uint
	Stdin       bool
	ApiBasePath string
	RunFunc     fnRun
	CmdFlags
}

var (
	wg common.WaitGroupCount
)

// initialize global flags
func (c *CmdBase) BaseInitFlags() *flag.FlagSet {
	c.fs = flag.NewFlagSet("krypton-cli", flag.ExitOnError)
	c.count = c.fs.Uint("count", DefaultRunCount, "how many times to run command")
	c.parallel = c.fs.Bool("parallel", false, "should the runs be in parallel")
	c.batchSize = c.fs.Uint("batch_size", DefaultBatchSize, "number of routines to run. parallel mode only.")
	c.verbose = c.fs.Bool("verbose", false, "verbose logs (default false)")
	c.retryCount = c.fs.Uint("retry_count", DefaultRetryCount, "how many times to retry failures")
	c.stdin = c.fs.Bool("stdin", false, "use stdin for command input")
	c.doc = c.fs.String("doc", "none", "generate documentation (none, shell)")
	c.apiBasePath = c.fs.String("api_base_path", "api/v1", "http api base path")
	return c.fs
}

// parse and get global flag values
func (c *CmdBase) BaseParse(args []string) {
	if err := c.fs.Parse(args); err != nil {
		log.Fatal(err)
	}
	if *c.verbose {
		config.SetVerboseLogging()
	}
	config.SetDocType(*c.doc)
	c.Count = *c.count
	if *c.count == 0 {
		log.Printf("count value %d is invalid. resetting to %d\n",
			*c.count, DefaultRunCount)
		c.Count = DefaultRunCount
	}
	c.Parallel = *c.parallel
	c.BatchSize = *c.batchSize
	if *c.batchSize <= 0 {
		log.Printf("batch_size value %d is invalid. resetting to %d\n",
			c.BatchSize, DefaultBatchSize)
		c.BatchSize = DefaultBatchSize
	} else if c.BatchSize > BatchSizeMax {
		log.Printf("warning: request burst %d is larger than recommended max %d.\n",
			c.BatchSize, BatchSizeMax)
	}
	c.RetryCount = *c.retryCount
	if *c.retryCount == 0 {
		log.Printf("warning: retry count %d is invalid. resetting to %d\n",
			*c.retryCount, DefaultRetryCount)
		c.RetryCount = DefaultRetryCount
	}
	c.Stdin = *c.stdin
	c.ApiBasePath = *c.apiBasePath
	log.Debugf("count: %d, burst size: %d\n", c.Count, c.BatchSize)
}

// wrap parallel runs to handle wait group
func (c *CmdBase) WrapParallel(i, count uint) {
	log.Debugf("Starting parallel run: %d/%d for %s\n", i, count, c.Name)
	defer wg.Done()
	c.RunFunc()
	log.Debugf("Finishing parallel run: %d/%d for %s\n", i, count, c.Name)
}

func (c *CmdBase) RunParallel(count uint) error {
	var i uint
	log.Debugf("%s parallel: count=%d\n", c.Name, count)
	for i = 1; i <= count; i++ {
		wg.Add(1)
		go c.WrapParallel(i, count)
	}
	wg.Wait()
	return nil
}

func (c *CmdBase) RunSerial() error {
	var i uint
	log.Debugf("%s serial: count=%d\n", c.Name, c.Count)
	for i = 1; i <= c.Count; i++ {
		c.RunFunc()
	}
	return nil
}

func (c *CmdBase) Execute() error {
	startTime := time.Now()
	err := c.internalExecute()
	log.Printf("Elapsed: %v, Processed: %d\n",
		time.Since(startTime).String(),
		c.Count)
	if err != nil {
		log.Fatalf("Failed to execute command: %v\n", err)
		return err
	}
	return nil
}

func (c *CmdBase) BaseExecute(cmd Command) error {
	startTime := time.Now()
	err := c.internalExecute()
	log.Printf("Elapsed: %v, Processed: %d\n",
		time.Since(startTime).String(),
		c.Count)
	if err != nil {
		log.Fatalf("Failed to execute command: %v\n", err)
		return err
	}
	return nil
}

func (c *CmdBase) internalExecute() error {
	var err error
	var i uint
	if c.Parallel {
		batch := c.Count / c.BatchSize
		remainder := c.Count % c.BatchSize
		for i = 0; i < batch; i++ {
			err = c.RunParallel(c.BatchSize)
		}
		if remainder > 0 {
			err = c.RunParallel(remainder)
		}
		return err
	} else {
		return c.RunSerial()
	}
}

func (c CmdBase) PrintDefaults() {
	c.fs.PrintDefaults()
}

func (c *CmdBase) GetBaseFlagSet() *flag.FlagSet {
	return c.fs
}

func (c *CmdBase) GetInput() interface{} {
	return nil
}

func (c *CmdBase) ExecuteWithArgs(i interface{}) error {
	return c.Execute()
}
