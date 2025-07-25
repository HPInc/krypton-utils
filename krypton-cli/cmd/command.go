package cmd

import (
	"cli/logging"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

var log = logging.GetLogger()

func (cmds Commands) Parse(name string, args []string) (Command, error) {
	c, ok := cmds[name]
	if !ok {
		cmds.ListAll()
		return nil, errors.New("No such command")
	}
	found, err := c.Parse(args)
	if err != nil {
		if errors.Is(err, ErrMissingArgs) {
			log.Error("Use -help to see available args")
		} else if !errors.Is(err, ErrParseStdin) {
			c.PrintDefaults()
		}
	}
	return found, err
}

func (cmds Commands) ListAll() {
	fmt.Println("Supported commands are:")
	for k := range cmds {
		fmt.Printf("\x1b[0;32m- %v\x1b[0m\n", k)
	}
}

// parse input as json
// ask command for args struct
// populate and execute
func ExecuteWithStdin(c Command) error {
	input := json.NewDecoder(os.Stdin)
	for {
		i := c.GetInput()
		err := input.Decode(i)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if err = c.ExecuteWithArgs(i); err != nil {
			return err
		}
	}
	return nil
}
