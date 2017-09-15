import contextlib
import json
import os
class BaseCommands:
    def __init__(self):
        return
    def get_unsupported_msg(self,command):
        return command+" is currently unsupported."
    def get_response(self,message,command,args):
        return get_unsupported_msg(self,command)
    def read_json(filename):
        try:
            with open(filename,"r+") as f:
                return json.load(f)
        except:
            return None
    def export_json(filename,thing):
        with contextlib.suppress(FileNotFoundError):
            os.remove(filename)
        with open(filename,"a+") as f:
                json.dump(thing,f)
