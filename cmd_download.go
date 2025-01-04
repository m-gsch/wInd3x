package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/freemyipod/wInd3x/pkg/app"
	"github.com/freemyipod/wInd3x/pkg/cache"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var downloadBits = map[string]cache.PayloadKind{
	"wtf":        cache.PayloadKindWTFUpstream,
	"bootloader": cache.PayloadKindBootloaderUpstream,
	"retailos":   cache.PayloadKindRetailOSUpstream,
	"diags":      cache.PayloadKindDiagsUpstream,
}

var downloadCmd = &cobra.Command{
	Use:   "download [kind] [output path]",
	Short: "Download update files from Apple's CDN",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := app.New()
		if err != nil {
			return err
		}
		defer app.Close()

		kind, ok := downloadBits[args[0]]
		if !ok {
			var opts []string
			for k := range downloadBits {
				opts = append(opts, k)
			}
			sort.Strings(opts)
			return fmt.Errorf("invalid kind, must be one of: %s", strings.Join(opts, ", "))
		}

		by, err := cache.Get(app, kind)
		if err != nil {
			return err
		}
		if err := os.WriteFile(args[1], by, 0600); err != nil {
			return err
		}
		glog.Infof("Wrote %s to %s", args[0], args[1])
		return nil
	},
}
