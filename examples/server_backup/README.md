# Discord Server Backup Bot

This example demonstrates how to create a Discord bot that can backup messages from channels and entire servers using slash commands.

## Features

- `/backup-server` - Backs up all text channels in the current server
- `/backup-channel` - Backs up a specific text channel
- `/list-backups` - Lists all available backup files
- `/get-backup` - Retrieves a specific backup file

## Requirements

- Go 1.16 or higher
- A Discord Bot Token
- The bot must have the following permissions:
  - Read Messages/View Channels
  - Read Message History
  - Send Messages
  - Attach Files
  - Use Slash Commands

## Usage

1. Create a new Discord application and bot at https://discord.com/developers/applications
2. Get your bot token
3. Run the bot:

```bash
go run main.go -token "YOUR_BOT_TOKEN"
```

### Optional Flags

- `-guild` - Test guild ID. If not passed, bot registers commands globally
- `-rmcmd` - Remove all commands after shutting down (default: true)
- `-backupdir` - Directory to store backups (default: "backups")

## Backup Format

Backups are stored as JSON files with the following structure:

```json
{
  "guild_id": "123456789",
  "guild_name": "My Server",
  "backup_time": "2024-02-09T12:00:00Z",
  "channels": [
    {
      "id": "987654321",
      "name": "general",
      "type": 0,
      "messages": [
        {
          "id": "111222333",
          "content": "Hello World!",
          "author": {
            "id": "444555666",
            "username": "User"
          },
          "timestamp": "2024-02-09T11:59:00Z"
        }
      ]
    }
  ]
}
```

## Notes

- The bot requires the "Manage Server" permission to use the backup commands
- Backups are stored locally in the specified backup directory
- Large servers with many messages may take some time to backup
- The bot includes rate limiting precautions to avoid Discord API limits 