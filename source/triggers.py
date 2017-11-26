import discord
import json
import pprint
import re
class Triggers:
    def __init__(self):
        self.import_triggers()
    def import_triggers(self):
        try:
            with open("../config/triggers.json","r+") as f:
                try:
                    self.triggers = json.load(f)
                except:
                    self.triggers = {}
        except:
            self.triggers = {}
    def get_response(self,message):
        responses = []
        for trigger in self.triggers.keys():
            rx = re.compile(trigger)
            if rx.search(message.content.lower()):
                responses.append(self.triggers[trigger])
        return responses

