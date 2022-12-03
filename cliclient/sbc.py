#!/usr/bin/python3

import os
import json
import math
import sys
from datetime import datetime
from prettytable import PrettyTable
from sipb_api import Pastebin

CONFIG_PATH=os.path.expanduser("~/.config/sbcrc")
USAGE="""USAGE: sbc [options] [arguments]
Options:
    -l [n]
        Displays details of n most recent files"""

def create_config_file():
    info = {}
    print("--- Configuration ---")
    print("Enter pastebin server (Ex: https://user:pass@example.com:1234/):")
    info["Server"] = input()
    marshalled = json.dumps(info)
    with open(CONFIG_PATH, "w") as f:
        f.write(marshalled)

def read_config_file():
    with open(CONFIG_PATH) as f:
        return json.loads(f.read())

def pretty_size(byteCnt):
    suffixes = ["B", "KiB", "MiB", "GiB", "TiB"]
    logB1024 = int(math.log(byteCnt, 1024))
    num = byteCnt / (1024 ** logB1024)
    return "{:.2f} {}".format(num, suffixes[logB1024])

def pretty_time(datetimeStr):
    date = datetime.strptime(datetimeStr, "%Y-%m-%dT%H:%M:%SZ")
    return date.strftime("%a, %b %d, %Y at %H:%M")

def display_files(files):
    tab = PrettyTable()
    tab.field_names = ["No", "Name", "Size", "Type", "Time"]
    i = 1
    for file in files:
        tab.add_row([i, file["Name"], pretty_size(file["Size"]), file["Type"], pretty_time(file["Timestamp"])])
    print(tab)

if len(sys.argv) < 2:
    print(USAGE)
    exit(1)

if not os.path.exists(CONFIG_PATH):
    create_config_file()

cfg = read_config_file()
pb = Pastebin(cfg["Server"])

match sys.argv[1]:

    case "-l":
        cnt = pb.count()
        limit = min(int(sys.argv[2]), cnt) if len(sys.argv) > 2 else cnt
        files = [pb.detailsOfNthNewest(i) for i in range(1, limit + 1)]
        display_files(files)

    case other:
        print(USAGE)
