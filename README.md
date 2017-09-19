# DVBot

This is the winner in the category "Innovation" at the Open Data Crunch 2016. The DVBot is a Facebook messenger bot for 
the public transport service in Dresden, the DVB.

For running the tests type:

`go test`

For building type:

`go build`

For running the application:

1. `export PORT=8000`
2. `./go-dvbot`

For developing I recommend using `main_test.go` for testing any input from the messenger without deploying it to the messenger.