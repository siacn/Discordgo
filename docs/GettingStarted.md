# Getting Started

This page is dedicated to helping you get started on your way to making the
next great Discord bot or client with DiscordGo. Once you've done that please
don't forget to submit it to the 
[Awesome DiscordGo](https://github.com/bwmarrin/discordgo/wiki/Awesome-DiscordGo) list :).


**First, lets cover a few topics so you can make the best choices on how to 
move forward from here.**

#### Bot Application
A bot application is a special program that interacts with the Discord servers
to perform some form of automation or provide some type of service.  Examples 
are things like number trivia games, music streaming, channel moderation, 
sending reminders, playing loud airhorn sounds, comic generators, YouTube 
integration, Twitch integration... You're *almost* only limited by your imagination.

Bot applications require the use of a special Bot account.  These accounts are
tied to your personal user account. Bot accounts cannot login with the normal
user clients and they cannot join servers the same way a user does. They do not 
have access to some user client specific features however they gain access to
many Bot specific features.

To create a new bot account first create yourself a normal user account on 
Discord then visit the [My Applications](https://discord.com/developers/applications/me)
page and click on the **New Application** box.  Follow the prompts from there
to finish creating your account.


**More information about Bot vs Client accounts can be found [here](https://discord.com/developers/docs/topics/oauth2#bot-vs-user-accounts).**

# Requirements

DiscordGo requires Go version 1.4 or higher.  It has been tested to compile and
run successfully on Debian Linux 8, FreeBSD 10, and Windows 7.  It is expected 
that it should work anywhere Go 1.4 or higher works. If you run into problems
please let us know :).

You must already have a working Go environment setup to use DiscordGo.  If you 
are new to Go and have not yet installed and tested it on your computer then 
please visit [this page](https://golang.org/doc/install) first then I highly
recommend you walk though [A Tour of Go](https://tour.golang.org/welcome/1) to
help get your familiar with the Go language.  Also checkout the relevant Go plugin 
for your editor &mdash; they are hugely helpful when developing Go code.

* Vim &mdash; [vim-go](https://github.com/fatih/vim-go)
* Sublime &mdash; [GoSublime](https://github.com/DisposaBoy/GoSublime)
* Atom &mdash; [go-plus](https://atom.io/packages/go-plus)
* Visual Studio &mdash; [vscode-go](https://github.com/Microsoft/vscode-go)


# Install DiscordGo

Like any other Go package the fist step is to `go get` the package.  This will
always pull the latest tagged release from the master branch. Then run 
`go install` to compile and install the libraries on your system.

#### Linux/BSD

Run go get to download the package to your GOPATH/src folder.

```sh
go get github.com/bwmarrin/discordgo
```

Finally, compile and install the package into the GOPATH/pkg folder. This isn't
absolutely required but doing this will allow the Go plugin for your editor to
provide autocomplete for all DiscordGo functions.

```sh
cd $GOPATH/src/github.com/bwmarrin/discordgo
go install
```

#### Windows
Placeholder.


# Creating a new Discord session

After you have loaded your bot token you can then create the bot session using 

```discord, err := discordgo.New("Bot " + "authentication token")```
