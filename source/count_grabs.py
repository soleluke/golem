#!/usr/bin/env python3
import json
from base_commands import BaseCommands

grabs = BaseCommands.read_json("../data/grabs.json")

count = [ (item,len(grabs[item])) for item in grabs.keys()]

print(count)
