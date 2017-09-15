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
YT_SERVICE_NAME = "youtube"
YT_API_VERSION = "v3"
class Commands(BaseCommands):

    def __init__(self,g_token):
        self.import_tells()
        self.gkey = g_token
        if command == "tell":
            return self.do_tell(sender,args)
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
        try:
        except:
    def get_youtube(self,search):
        try:
            youtube = build(YT_SERVICE_NAME,YT_API_VERSION,developerKey=self.gkey)
            search_response = youtube.search.list(q = search,part="id,snippet",maxResults=1).execute()
            search_result = search_response.get("items",[])
            vid=search_result["id"]["videoId"]
            return "https://www.youtube.com/watch?v=" + vid
        except HttpError as e:
             print( "An HTTP error %d occurred:\n%s",e.resp.status, e.content)
