package main

import (
	"log"
	"strings"
	"time"

	"net"

	"github.com/hashicorp/mdns"
	"golang.org/x/net/context"
)

const (
	castService = "_googlecast._tcp"
)

// A simple example, showing how to find a Chromecast using mdns, and request its status.

type ChromecastDevice struct {
	FullName string
	Name string
	Host string
	Addr net.IP
	Port int
}

func listChromecastsWithTimeout(ctx context.Context, timeout time.Duration) []ChromecastDevice {
	ctx, cancelFunc := context.WithTimeout(ctx, timeout)
	defer cancelFunc()
	return listChromecasts(ctx)
}

// listChromecasts scans network for registered chromecasts and returns a list
// make sure the context has a deadline set, as scanning will last until context is done.
// TODO: make this threadsafe...
func listChromecasts(ctx context.Context) []ChromecastDevice {

	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	devices := make([]ChromecastDevice, 0, 5)

	go func() {
		for {
			select {
			case entry := <-entriesCh:
				if !strings.Contains(entry.Name, castService) {
					return
				}

				log.Printf("Got new chromecast: %+v\n", entry)

				device := ChromecastDevice{
					FullName: entry.Name,
					Name: findField(entry.InfoFields, "fn="),
					Host: entry.Host,
					Addr: entry.Addr,
					Port: entry.Port,
				}

				devices = append(devices, device)
			//case <- ctx.Done():
			//	return devices
			}

		}
	}()

	go func() {
		mdns.Query(&mdns.QueryParam{
			Service: castService,
			Timeout: time.Second * 5,
			Entries: entriesCh,
		})
	}()

	<-ctx.Done()
	return devices

	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt, os.Kill)
	//
	//// Block until a signal is received.
	//s := <-c
	//fmt.Println("Got signal:", s)
}

func findField(fields []string, startsWith string) string {
	for _, field := range fields {
		if strings.HasPrefix(field, startsWith) {
			return strings.TrimPrefix(field, startsWith)
		}
	}
	return ""
}
