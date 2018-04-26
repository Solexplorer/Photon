package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
)

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.DefaultStreamHandler(os.Stderr)))
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, os.Args[0], "stun.l.google.com:19302")
	}
	flag.Parse()
	addr := flag.Arg(0)
	if len(addr) == 0 {
		//addr = "stun.l.google.com:19302"
		addr = "182.254.155.208:3478"
	}
	c, err := stun.Dial("udp", addr)
	if err != nil {
		log.Crit("dial: %s", err)
	}
	deadline := time.Now().Add(time.Second * 25)
	if err := c.Do(stun.MustBuild(stun.TransactionIDSetter, stun.BindingRequest), deadline, func(res stun.Event) {
		if res.Error != nil {
			log.Crit("res %s", res)
		}
		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err != nil {
			var addr stun.MappedAddress
			err = addr.GetFrom(res.Message)
			if err != nil {
				log.Crit(err.Error())
			}
			log.Info("addr=%s", addr)
		} else {
			log.Info("xoraddr=%s", xorAddr)
		}
	}); err != nil {
		log.Crit("do: %s", err)
	}
	if err := c.Close(); err != nil {
		log.Crit(err.Error())
	}
}
