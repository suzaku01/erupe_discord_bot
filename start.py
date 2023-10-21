import discord
import requests
import aiohttp
import json
from discord.ext import commands, tasks

TOKEN = "TOKEN"
GUILD_ID = SERVER_ID  # サーバーIDを指定
CHANNEL_ID = Channel_ID  # チャンネルIDを指定
# 許可するユーザーのIDリスト
ALLOWED_USERS = [
    'USER_ID1',  # 例のID。実際のユーザーIDに置き換えてください
    'USER_ID2'  # 例のID。実際のユーザーIDに置き換えてください
]
intents = discord.Intents.all()
bot = commands.Bot(command_prefix='/',intents=intents)
    
@bot.event
async def on_ready():
    print(f'{bot.user.name} has connected to Discord!')
    update_channel_name.start()  # タスクを開始

@tasks.loop(minutes=1)  # 1分ごとに実行。必要に応じて調整してください。
async def update_channel_name():
    response = requests.get("http://localhost:9998/get_players")
    if response.status_code == 200:
        players_data = response.json()
        channel = bot.get_channel(CHANNEL_ID)
        await channel.edit(name=f"Online Players: {players_data['players']}")

@bot.command()
async def adminmessage(ctx, *, message_content):
    if str(ctx.author.id) not in ALLOWED_USERS:
        await ctx.send("あなたはこのコマンドを使用する権限がありません。")
        return
    
    await send_to_game_server(message_content)  # awaitを使用して非同期関数を呼び出す
    #await ctx.send(f'メッセージ "{message_content}" がゲームサーバーに送信されました。')
    
async def send_to_game_server(message):
    url = 'http://localhost:9999/send'
    headers = {
        'Authorization': 'Bearer TOKEN',
        'Content-Type': 'application/json'  # JSON形式を指定
    }
    data = json.dumps({'message': message})  # データをJSON文字列に変換

    async with aiohttp.ClientSession() as session:
        async with session.post(url, headers=headers, data=data) as response:
            if response.status != 200:
                print(f"Error: {response.status}, {await response.text()}")

if __name__ == "__main__":
    bot.run(TOKEN)