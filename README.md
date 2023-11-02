# erupe_discord_bot
This is a discord bot for erupe.

# Functions
- Change channnel name with online players count.
- Stream text to all players.
- Check if serevr is running or not.
- Start erupe from discrod.

# Installation
- Add `main.go` to your build correctly. 
- Place `start.py` to your erupe root.
- Run `pip install Flask`.  
- Run `pip install discord.py`.
- Run `pip install requests`.  

# Edit `start.py`
-Bot token goes to line 8,48


# Usage
Run `start.py`, then run your erupe.  

Server updates channel name per 1min, but due to API limits, it may take more longer(3-4min sometimes).  

To stream text, use `/adminmessage` command. Example: /adminmessage test message. 

Added /restart command to restart erupe when it's crashed or dead. Do not use this while it's alive.  
Added /isalive command to check if the server is running or not.  
Added /seehelp command to see all bot command.  
 
