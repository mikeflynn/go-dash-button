# `go-dash-button`
A server to handle Amazon Dash button actions; Written in golang

## What?

This is a simple server that sniffs a preset list of Amazon Dash buttons and triggers a function.

## How do I start?

It's not a completely seemless experience. You'll need to clone this repo, and add any Dash button MAC addresses and their action functions to the `main.go` file.

Once your code is in place, compile with `go build` and run the binary `sudo ./go-dash-button` (note that it does require sudo rights).

## How do I find the MAC address of my button?

Compile the code right away and run the server. Set up your Dash Button so that it joins your wireless network, but don't pick the product. Click the Dash button and you'll see the server log out the MAC address.

## How can I specify job parameters?

There is an optional config file system built in. It will look for an ini file at `./go-dash-button.ini` but you can use the `-conf` flag to specify any ini file location and use the data in your job like so:

`Config.Section("hue").Key("baseStationIP").String()`
