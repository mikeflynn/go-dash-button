# `go-dash-button`
A server to handle Amazon Dash button actions; Written in golang

# What?

This is a simple server that sniffs a preset list of Amazon Dash buttons and triggers a function.

# How Do I Start?

It's not a completely seemless experience. You'll need to clone this repo, and add any Dash button MAC addresses and their action functions to the `main.go` file.

Once your code is in place, compile with `go build` and run the binary `sudo ./go-dash-button` (note that it does require sudo rights).

# How Do I Find the MAC Address of my Dash Button?

Compile the code right away and run the server. Set up your Dash Button so that it joins your wireless network, but don't pick the product. Click the Dash button and you'll see the server log out the MAC address.
