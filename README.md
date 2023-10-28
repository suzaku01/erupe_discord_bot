# erupe_discord_bot
This is a discord bot for erupe.

# Functions
- Change channnel name with online players count.
- Stream text to all players.

# Installation
- Add main.go adn sys_session.go to your build correctly. 
- Place `app.py` and `start.py` to somewhere you want.
- Run `pip install Flask`.  
- Run `pip install discord.py`.
- Run `pip install requests`.  

# Edit `start.py`
-Bot token goes to line 7,43,54  
-Discord server ID goes to line 8  
-Discord channel ID goes to line 9  
-Admin user id(s) goes to line 12,13...  

# Usage
Run both `app.py` and `start.py`, then run your erupe.  

Server updates channel name per 1min, but due to API limits, it may take more longer(3-4min sometimes).  

To stream text, use `/adminmessage` command. Example: /adminmessage test message  
