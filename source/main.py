#!/usr/bin/env python3
import asyncio
import config
import commands
import discord
import pprint
import triggers

config = config.Config()
client = discord.Client()
commands = commands.Commands(config.get_gkey())
triggers = triggers.Triggers()
def get_command(message):
    command = message.split(" ",1)[0]
    return command[1:]


@client.event
async def on_ready():
    print('logged in as')
    print(client.user.name)
    print(client.user.id)
    print('-------')
@client.event
async def on_message(message):
    if message.author != client.user:
        tmp = await check_triggers(message)
        if message.channel.is_private == False or message.content.startswith("!help") or message.content.startswith("!list") :

async def check_triggers(message):
    for response in triggers.get_response(message):
        tmp = await client.send_message(message.channel,response) 
    return
client.run(config.get_token())
