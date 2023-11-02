import discord
import requests
import aiohttp
import json
from discord.ext import commands, tasks
import subprocess

TOKEN = ""
GUILD_ID =   # サーバーIDを指定
CHANNEL_ID =   # チャンネルIDを指定
# 許可するユーザーのIDリスト
ALLOWED_USERS = [
    '',
    '', 
    ''  
]
intents = discord.Intents.all()
bot = commands.Bot(command_prefix='/',intents=intents)
    
@bot.event
async def on_ready():
    print(f'{bot.user.name} has connected to Discord!')
    update_channel_name.start()  # タスクを開始

@tasks.loop(minutes=1)  # 1分ごとに実行。必要に応じて調整してください。
async def update_channel_name():
    response = requests.get("http://localhost:9999/getplayers")
    if response.status_code == 200:
        players_data = response.json()
        total_players = players_data['totalPlayers']
        channel = bot.get_channel(CHANNEL_ID)
        await channel.edit(name=f"Online Players: {total_players}")
        


@bot.command()
async def adminmessage(ctx, *, message_content):
    if str(ctx.author.id) not in ALLOWED_USERS:
        await ctx.send("You are not authorized to use this command.")
        return
    
    await send_to_game_server(message_content)  # awaitを使用して非同期関数を呼び出す
    #await ctx.send(f'メッセージ "{message_content}" がゲームサーバーに送信されました。')
    
async def send_to_game_server(message):
    url = 'http://127.0.0.1:9999/send'
    headers = {
        'Authorization': 'Bearer your_TOKEN',
        'Content-Type': 'application/json'  # JSON形式を指定
    }
    data = json.dumps({'message': message})  # データをJSON文字列に変換

    async with aiohttp.ClientSession() as session:
        async with session.post(url, headers=headers, data=data) as response:
            if response.status != 200:
                print(f"Error: {response.status}, {await response.text()}")

@bot.command()
async def isalive(ctx):
    if str(ctx.author.id) not in ALLOWED_USERS:
        await ctx.send("You are not authorized to use this command.")
        return
        
    server_url = "http://localhost:9999/isalive"
    async with aiohttp.ClientSession() as session:
        try:
            async with session.get(server_url) as response:
                if response.status == 200:
                    await ctx.send("Server is running!")
                else:
                    await ctx.send("Server is not running.")
        except Exception as e:
            await ctx.send("Server is not running")

@bot.command()
async def restart(ctx):
    if str(ctx.author.id) not in ALLOWED_USERS:
        await ctx.send("You are not authorized to use this command.")
        return
        
    try:
        # 新しいコマンドプロンプトウィンドウで `go run .` を実行
        subprocess.Popen(["cmd", "/c", "start", "cmd", "/k", "go run ."], shell=True)
        #await ctx.send('コマンドプロンプトで `go run .` を実行しました。')
    except Exception as e:
        await ctx.send(f'An error occurred while executing the command: {e}')
        
@bot.command()
async def seehelp(ctx):
    if str(ctx.author.id) not in ALLOWED_USERS:
        await ctx.send("You are not authorized to use this command.")
        return

    #await ctx.send("Hello!\nHow are you?")
    await ctx.send("`/adminmessage your_text`:Stream text to all connected players.\n`/isalive`:Check if server is running or not.\n`/restart`:Restart erupe. Only available when erupe is not running.")

if __name__ == "__main__":
    bot.run(TOKEN)
