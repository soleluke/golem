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
    sender ="<@"+ message.author.id+">"
    tells = commands.get_tells(sender)
    for tell in tells:
        tmp = await client.send_message(message.channel,tell)
    if message.content.startswith(config.get_prefix()):
        command = get_command(message.content)
        args = message.content.split(' ',1)[1]
        if config.check_command(command):
            tmp = await client.send_message(message.channel,commands.get_response(sender,command,args))
        else:
            tmp = await client.send_message(message.channel,commands.get_unsupported_msg(command))
        tmp = await check_triggers(message)

async def check_triggers(message):
    for response in triggers.get_response(message):
        tmp = await client.send_message(message.channel,response) 
    return
client.run(config.get_token())
