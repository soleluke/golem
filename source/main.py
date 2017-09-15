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
            sender ="<@"+ message.author.id+">"
            tells = commands.get_tells(sender)
            for tell in tells:
                tmp = await client.send_message(message.channel,tell)
            if message.content.startswith(config.get_prefix()):
                command = get_command(message.content)
                if command == "grab":
                    async for log in client.logs_from(message.channel,limit=2):
                        if not log.content.startswith("!grab"):
                            commands.add_grab(log)
                            tmp = await client.send_message(message.channel,"Grab Successful")
                elif command == "list" and config.check_command(command):
                    tmp = await client.send_message(message.author,commands.list_places())
                    tmp2 = await client.send_message(message.channel,"I have PM'd you the list of places")
                else:
                    try:
                        args = message.content.split(' ',1)[1]
                    except:
                        args = ""
                    if config.check_command(command):
                        tmp = await client.send_message(message.channel,commands.get_response(message,command,args))
                    else:
                        tmp = await client.send_message(message.channel,commands.get_unsupported_msg(command))

async def check_triggers(message):
    for response in triggers.get_response(message):
        tmp = await client.send_message(message.channel,response) 
    return
client.run(config.get_token())
