import discord
import pprint
import datetime
import json
import os
class Commands:

    def __init__(self):
        self.import_tells()
    def get_unsupported_msg(self,command):
        return command+" is currently unsupported."
    def get_response(self,sender,command,args):
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

