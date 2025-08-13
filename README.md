# NetWatch

A network interface watcher for windows

## What is this?

Does your app need to know when a network interface connects? When it disconnects? Maybe you need to do things when an IP address changes?
Does your app run on windows?
Well I think this library might be for you.

This library also includes some useful network utilities, such as listing all interfaces that satisfy a specific set of flags,
or getting the IPv4 address of a given interface

## How do I use this?

First you'll want to grab this library

```shell
go get github.com/mbaklor/netwatch
```

next you probably want an app that needs this functionality, here's a simple demo

```go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/mbaklor/netwatch"
)

func main() {
	// Create a network monitor
	nm := netwatch.NewNetMonitor()

	// we need to register our monitor to hook into OS notifications
	// passing in true will request a notification to fire when we register, to indicate everything is working
	err := nm.Register(true)
	if err != nil {
		// we need to clean up in case one of the calls in Register succeeded and the other failed
		err1 := nm.Unregister()
		panic(errors.Join(err, err1))
	}

	// now we listen for notifications for the next 10 seconds, before finally unregistering and exiting the program
	for {
		select {
		case n := <-nm.MonitorNotificationChan:
			fmt.Println(n)
		case <-time.After(time.Second * 10):
			fmt.Println("Done")
			err := nm.Unregister()
			if err != nil {
				fmt.Println("error unregistering", err)
			}
			return
		}
	}
}
```

## Roadmap

I don't know if I'll ever end up taking the time to add cross platform compatibility, but it's on the list of things I'd like to do, so:

- [ ] MacOS
- [ ] Linux
