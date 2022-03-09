# Setup

First copy the configuration template to `botconfig.json`:

``` 
$ cp config-template.json botconfig.json
```

The bot will need a valid Discord access token which you can get by creating an
application for the bot in the [Discord Developer Portal dashboard][devportal].
Add a bot user to the application and find the token on the Bot page. Copy the
token and add it to your `botconfig.json` file.

[devportal]: https://discord.com/developers/applications/

At this point you should be able to run the bot (assuming you have a Go
toolchain installed):

``` 
$ go run .
```
