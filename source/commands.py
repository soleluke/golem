import discord
import pprint
from apiclient.discovery import build
from apiclient.errors import HttpError
import datetime
import json
import os
class Commands:
YT_SERVICE_NAME = "youtube"
YT_API_VERSION = "v3"

    def __init__(self):
        self.import_tells()
    def get_unsupported_msg(self,command):
        return command+" is currently unsupported."
    def get_response(self,sender,command,args):
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
        os.remove("../data/tells.json")
        with open("../data/tells.json","w+") as f:
            json.dump(self.tells,f)
    def import_tells(self):
        try:
            with open("../data/tells.json","r") as f:
                try:
                    self.tells = json_load(f)
                except:
                    
                    self.tells = {}
        except:
            self.tells = {}
        pprint.pprint(self.tells)

    def get_youtube(self,search):
        try:
            youtube = build(YT_SERVICE_NAME,YT_API_VERSION,developerKey=self.gkey)
            search_response = youtube.search.list(q = search,part="id,snippet",maxResults=1).execute()
            search_result = search_response.get("items",[])
            vid=search_result["id"]["videoId"]
            return "https://www.youtube.com/watch?v=" + vid
        except HttpError as e:
             print( "An HTTP error %d occurred:\n%s",e.resp.status, e.content)
