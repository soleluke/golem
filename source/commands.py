from apiclient.discovery import build
from apiclient.errors import HttpError
from base_commands import BaseCommands
import contextlib
import datetime
import discord
import json
import os
import random
import pprint
from weather import Weather

YT_SERVICE_NAME = "youtube"
YT_API_VERSION = "v3"
class Commands(BaseCommands):

    def __init__(self,g_token):
        self.import_tells()
        self.import_doots()
        self.import_grabs()
        self.import_places()
        self.backlog = BaseCommands.read_json("../data/backlog.json")
        self.work_tells = BaseCommands.read_json("../data/wtells.json")
        self.gkey = g_token
    def get_response(self,message,command,args):
        sender = "<@"+message.author.id + ">"
        args=args.replace("!","")
        if command == "tell":
            return self.do_tell(sender,args)
        elif command == "ask":
            return self.ask(args)
        elif command == "updoot":
            if args.split(' ',1)[0] != sender:
                return self.doot(args.split(' ',1)[0],1)
            else:
                return "No"
        elif command == "downdoot":
            if args.split(" ",1)[0] != sender:
                return self.doot(args.split(' ',1)[0],-1)
            else:
                return "No"
        elif command == "doots":
            return self.doot(args,0)
        elif command == "grabr":
            return self.get_grab(args)
        elif command == "rgrab":
            return self.get_grab_topic(args)
        elif command == "addplace":
            return self.add_place(args)
        elif command == "suggest":
            return self.suggest_place()
        elif command == "list":
            return self.list_places()
        elif command == "weather":
            return self.get_weather(args)
        elif command == "forecast":
            return self.get_forecast(args)
        elif command == "igrab":
            self.import_grabs()
            return "Refreshed Grabs from File"
        elif command =="bang":
            return self.bang(args)
        elif command =="backlog":
            return self.add_to_backlog(args)
        elif command=="wtell":
            return self.work_tell(sender,args)
        elif command=="home":
            return self.check_work_tells(sender)
        elif command == "help":
            return self.get_help(args)
        else:
            return self.get_unsupported_msg(command)
    def do_tell(self,sender,args):
        dest = args.split(' ',1)[0]
        msg = args.split(' ',1)[1]
        msg = sender+" said '"+msg+"' on " + datetime.datetime.now().strftime("%A, %d %B %Y at %I:%M%p")
        if dest not in self.tells.keys():
            self.tells[dest] = list()
        self.tells[dest].append(msg)
        self.export_tells()
        return "ill tell " + dest + " that next time i see them"
    def get_tells(self,sender):
        if sender in self.tells.keys():
            ret = self.tells[sender]
            self.tells[sender] = list()
            self.export_tells()
            return ret
        else:
            return list()
    def export_tells(self):
        BaseCommands.export_json("../data/tells.json",self.tells)
    def import_tells(self):
        self.tells = BaseCommands.read_json("../data/tells.json")
    def doot(self,person,change):
        if person in self.doots.keys():
            self.doots[person]+=change
            self.export_doots()
        else:
            self.doots[person]=change
            self.export_doots()
        return (person + " has " + str(self.doots[person])+ " doots")
    def export_doots(self):
        BaseCommands.export_json("../data/doots.json",self.doots)
    def import_doots(self):
        self.doots = BaseCommands.read_json("../data/doots.json")
    def add_grab(self,grab):
        grab_source =  "<@"+grab.author.id +">"
        if grab_source not in self.grabs.keys():
            self.grabs[grab_source] = []
        self.grabs[grab_source].append(grab.content)
        self.export_grabs()
    def import_grabs(self):
        self.grabs = BaseCommands.read_json("../data/grabs.json")
    def export_grabs(self):
        BaseCommands.export_json("../data/grabs.json",self.grabs)
    def get_grab_topic(self,thing):
        grabs= [ [{item:x} for x in self.grabs[item] if thing in x ] for item in self.grabs.keys()]
        grabs = [item for sublist in grabs for item in sublist]
        (graba,grab), = random.choice(grabs).items()
        return graba+": "+grab
    def get_grab(self,author):
        if author == "":
            graba= random.choice(list(self.grabs.keys()))
            grab = random.choice(self.grabs[graba])
            return graba+": "+grab
        else:
            grabs = self.grabs[author]
            return author+": "+random.choice(self.grabs[author])
    def ask(self,args):
        try:
            options = args.split(" or ",20)
        except:
            options = ['yes','no']
        if len(options) == 1:
            options = ['yes','no']
        return random.choice(options)
    def import_places(self):
        self.places = BaseCommands.read_json("../data/places.json")
    def export_places(self):
        BaseCommands.export_json("../data/places.json",self.places)
    def add_place(self,place):
        if place not in self.places:
            self.places.append(place)
            self.export_places()
            return "Added " + place
        else:
            return place + " is already on the list"
    def suggest_place(self):
        return random.choice(self.places)
    def list_places(self):
        ret = ""
        for place in self.places:
            ret+=place+"\n"
        return ret
    def get_youtube(self,search):
        try:
            youtube = build(YT_SERVICE_NAME,YT_API_VERSION,developerKey=self.gkey)
            search_response = youtube.search.list(q = search,part="id,snippet",maxResults=1).execute()
            search_result = search_response.get("items",[])
            vid=search_result["id"]["videoId"]
            return "https://www.youtube.com/watch?v=" + vid
        except HttpError as e:
             print( "An HTTP error %d occurred:\n%s",e.resp.status, e.content)
    def get_weather(self,search):
        weather = Weather()
        location = weather.lookup_by_location(search)
        condition = location.condition()
        return location.title()+": " + condition["temp"] + " degrees and " + condition["text"]
    def get_forecast(self,search):
        weather = Weather()
        location = weather.lookup_by_location(search)
        forecast = location.forecast()[0]
        return "Forecast: " + forecast["text"] +", " + forecast["high"] +"H "+ forecast["low"] + "L"
    def bang(self,target):
        bangs = BaseCommands.read_json("../data/bangs.json")
        bang = random.choice(bangs)
        return target + bang
    def add_to_backlog(self,item):
        if self.backlog == None:
            self.backlog = []
        self.backlog.append(item)
        BaseCommands.export_json("../data/backlog.json",self.backlog)
        return "Added to backlog"
    def list_commands(self):
        commands = BaseCommands.read_json("../config/commands.json")
        ret_com = [x +" usage: "+ commands[x] for x in commands.keys()]
        return "\n".join(ret_com)
    def work_tell(self,sender,args):
        target = args.split(' ',1)[0]
        msg = args.split(' ',1)[1]
        if self.work_tells == None:
            self.work_tells = {}
        if target not in self.work_tells.keys():
            self.work_tells[target] = []
        self.work_tells[target].append(sender+" said "+msg)
        BaseCommands.export_json("../data/wtells.json",self.work_tells)
        return "Work tell added"
    def check_work_tells(self,sender):
        ret = "Work tells:"
        if sender not in self.work_tells.keys():
            return "No Work tells for this user"
        for tell in self.work_tells[sender]:
            ret+="\n"+tell
        del self.work_tells[sender]
        BaseCommands.export_json("../data/wtells.json",self.work_tells)
        return ret
    def get_help(self,command):
        if command == "":
            return "usage: help <command> - returns help for <command>. use the 'commands' command to get a list of commands"
        try:
            commands = BaseCommands.read_json("../config/commands.json")
            return "usage: "+ commands[command]
        except:
            return command + " is not currently supported"
