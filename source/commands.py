from apiclient.discovery import build
from apiclient.errors import HttpError
from base_commands import BaseCommands
import contextlib
import datetime
import discord
import json
import os
import random
import re
import pprint
from weather import Weather
from fuzzywuzzy import fuzz

YT_SERVICE_NAME = "youtube"
YT_API_VERSION = "v3"
class Commands(BaseCommands):

    def __init__(self,g_token):
        self.import_tells()
        self.import_doots()
        self.import_grabs()
        self.import_places()
        self.doot_times = {}
        self.backlog = BaseCommands.read_json("../data/backlog.json")
        self.work_tells = BaseCommands.read_json("../data/wtells.json")
        self.reminds = BaseCommands.read_json("../data/reminds.json")
        if self.reminds is None:
            self.reminds = {}
        self.gkey = g_token
    def get_response(self,message,command,args):
        sender = "<@"+message.author.id + ">"
        args=args.replace("!","")
        if command == "tell":
            return self.do_tell(sender,args)
        elif command == "ask":
            return self.ask(args)
        elif command == "coinflip":
            return self.ask("heads or tails")
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
        elif command == "sidedoot":
            return "~~doot doot~~"
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
        elif command =="math":
            return self.math(args)
        elif command =="backlog":
            return self.add_to_backlog(args)
        elif command=="wtell":
            return self.work_tell(sender,args)
        elif command=="home":
            return self.check_work_tells(sender)
        elif command=="dlq":
            return self.add_to_dl_queue(args)
        elif command=="roll":
            return self.roll(args)
        elif command=="remind":
            return self.remind(sender,args)
        elif command == "help":
            return self.get_help(args)
        else:
            return self.get_unsupported_msg(command)
    def get_unsupported_msg(self,command):
        return self.fuzzy_command(command)
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
        if self.tells == None:
            self.tells ={}
    def doot(self,person,change):
#        if change !=0:
#            if person in self.doot_times.keys():
#                delt = datetime.datetime.now() - self.doot_times[person]
#                if delt.total_seconds() < 30:
#                    return "Chill yo"
#            self.doot_times[person] = datetime.datetime.now()
        if person in self.doots.keys():
            self.doots[person]+=change
            self.export_doots()
        else:
            self.doots[person]=change
            self.export_doots()
        return "doot doot"
    def export_doots(self):
        BaseCommands.export_json("../data/doots.json",self.doots)
    def import_doots(self):
        self.doots = BaseCommands.read_json("../data/doots.json")
        if self.doots == None:
            self.doots = {}
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
            grabs = [ [{item:x} for x in self.grabs[item]] for item in self.grabs.keys()]
            grabs = [item for sublist in grabs for item in sublist]
            (graba,grab),=random.choice(grabs).items()
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
    def add_to_dl_queue(self,link):
        queue = BaseCommands.read_json("../data/queue.json")
        queue = [] if queue is None else queue
        queue.append(link)
        BaseCommands.export_json("../data/queue.json",queue)
        return "Added " + link + " to download queue"
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
    def roll(self,args):
        m = re.search(r'\b([0-9]+)d([0-9]+)(D[0-9]+[h,l])?([>,<][0-9]+)?(\+[0-9]+)?(x[0-9]+)?',args)
        if m is None:
            return self.get_help("roll")
        number = int(m.group(1))
        faces = int(m.group(2))
        strDrop = m.group(3)
        strValLimit = m.group(4)
        strMod = m.group(5)
        strRolls = m.group(6)
        if number>100:
            return "I can't roll that many dice"
        if faces>100:
            return "That's a big die"
        if strRolls is not None:
            strRolls = strRolls[1:]
            rolls = int(strRolls)
        else:
            rolls = 1
        results = self.roll_dice_multiple(number,faces,rolls)
        if strValLimit is not None:
            results = [ self.dice_value_limit(x,strValLimit,faces) for x in results ]
        if strDrop is not None:
            results = [ self.drop_dice(x, strDrop) for x in results ]
        if strMod is not None:
            modifier = int(strMod)
        else:
            modifier = 0
        to_ret = "\n".join([ " ".join([str(x) for x in diceSet]) + " ("+str(sum(diceSet)+modifier)+")" for diceSet in results ] )
        return to_ret 
    def roll_dice_multiple(self,number,faces,rolls):
        results=[]
        for i in range(rolls):
            results.append(self.roll_dice(number,faces))
        return results
    def roll_dice(self,number,faces):
        results = []
        for i in range(number):
            results.append(self.roll_die(faces))
        return results
    def roll_die(self,faces):
        options=range(1,faces+1,1)
        return random.choice(options)
    def dice_value_limit(self,diceSet,strValLimit,faces):
        results = diceSet
        m = re.search(r'([<,>])([0-9]+)',strValLimit)
        if m is not None:
            direction = m.group(1)
            limit = int(m.group(2))
            if direction == '<':
                while [x for x in results if x >= limit ] is not []:
                    results = [ self.roll_die(faces) if x >= limit else x for x in results]
            elif direction == '>':
                while [x for x in results if x <= limit ] != []:
                    results = [ self.roll_die(faces) if x <= limit else x for x in results]
        return results
    def drop_dice(self,diceSet,strDrop):
        results = diceSet
        m = re.search(r'D([0-9]+)([h,l])',strDrop)
        if m is not None:
            dice=int(m.group(1))
            side = m.group(2)
            for i in range(dice):
                if side == 'l':
                        results.remove(min(results))            
                elif side == 'h':
                    results.remove(max(results))
        return results
    def remind(self,sender,args):
        if args is not "":
            if sender not in self.reminds.keys():
                self.reminds[sender] = ""
            self.reminds[sender] = self.reminds[sender]+" "+args
            BaseCommands.export_json("../data/reminds.json",self.reminds)
            return "Added remind"
        else:
            cur_reminds = self.reminds[sender]
            self.reminds[sender] = ""
            BaseCommands.export_json("../data/reminds.json",self.reminds)
            return cur_reminds
    def get_help(self,command):
        if command == "":
            return "usage: help <command> - returns help for <command>. use the 'commands' command to get a list of commands"
        try:
            commands = BaseCommands.read_json("../config/commands.json")
            return "usage: "+ commands[command]
        except:
            return command + " is not currently supported\n"+self.fuzzy_command(command)
    def fuzzy_command(self,inp):
        commands = BaseCommands.read_json("../config/commands.json")
        comms = commands.keys()
        best = 0
        ret = ""
        for comm in comms:
            ratio = fuzz.ratio(inp,comm)
            if ratio>best :
                best = ratio
                ret = comm
        return "Closest Command:\n"+commands[ret]
