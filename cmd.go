package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type FPort string

var fPort FPort

// Parse "10~12,5" to []uint16{10,11,12,5} in arbitrary order
func (f FPort) parse() ([]uint16, error) {
	ports := make([]uint16, 0)
	ps := strings.Split(string(f), ",")
	for _, item := range ps {
		values := strings.Split(item, "~")
		var i int
		var val string
		var pa [2]uint16
		for i, val = range values {
			if i > 1 {
				return nil, fmt.Errorf("%v is invalid format", fPort)
			}
			p, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("%v is invalid format", fPort)
			}
			if p < 0 || p > 65535 {
				return nil, fmt.Errorf("%v contain invalid port", fPort)
			}
			pa[i] = uint16(p)
		}

		if i == 0 {
			ports = append(ports, pa[i])
		} else {
			for j := pa[0]; j <= pa[1]; j++ {
				ports = append(ports, j)
			}
		}
	}
	return ports, nil
}

var rootCmd = &cobra.Command{
	Use:   "portscanner [IP|CIDR]",
	Short: "tcp port scan tool",
	Long:  "A simple tcp port scan tool.",
	Example: `  portscanner 127.0.0.1 -p 80
  portscanner 127.0.0.1 -p 1~1023,1234
  portscanner 127.0.0.1/24 -p 22`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("require least one argument")
		}
		if len(args) > 1 {
			return fmt.Errorf("too much argument given: %v", args)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ports, err := fPort.parse()
		if err != nil {
			return err
		}
		scanner, err := NewScanner(args[0], ports)
		if err != nil {
			return err
		}
		scanner.Start()
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP((*string)(&fPort), "port", "p", "", "(Required) Between 0 and 65535.\nUse ',' or '~' to scan multiple ports. See above examples.")
	rootCmd.MarkPersistentFlagRequired("port")
}
