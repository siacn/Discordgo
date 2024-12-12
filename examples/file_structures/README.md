<img align="right" alt="DiscordGo logo" src="/docs/img/discordgo.svg" width="200">

## DiscordGo 

The main moto of this folder is to represent that how you can use the file structure and create clean bots it will also teach you on how to import functions from other folders , it's basic but it's usefull to use a proper file structure

This example demonstrates how to utilize DiscordGo to create a slash commands bot for confessions 
this sends confessions directly and secertly to the channel designated 

This Bot will respond to "!hello" in any server it's in with "meow meow" in the
sender's DM.

**Join [Discord Gophers](https://discord.gg/0f1SbxBZjYoCtNPP)
Discord chat channel for support.**

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.

From within the file_strcture example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

This example uses bot tokens for authentication only. While user/password is
supported by DiscordGo, it is not recommended for bots.

```
Create a .env File and add ChannelID ApplicationID Token tags and insert the following
Run Command is Simple use ` go run . `
```

The below example shows how to start the bot

```sh
./file_structre go run .
Bot is now running.  Press CTRL-C to exit.
```
