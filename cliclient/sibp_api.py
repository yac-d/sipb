from base64 import b64encode
import requests
import json

class Pastebin:
    def __init__(self, host, port=80, basicauth=None, https=False, pastebinroot="/"):
        self.url = "https://" if https else "http://"
        self.url += "{}:{}{}".format(host, port, pastebinroot)

        self.headers = {}
        if basicauth:
            userAndPass = b64encode(bytes("{}:{}".format(basicauth[0], basicauth[1]), "ascii")).decode("ascii")
            self.headers = {"Authorization": "Basic {}".format(userAndPass)}

    def count(self):
        return int(requests.get(self.url+"retrieve/fileCount", headers=self.headers).text)

    def detailsOfNthNewest(self, n):
        return json.loads(requests.post(self.url+"retrieve", str(n), headers=self.headers).text)

    def upload(self, filepath):
        file = open(filepath, "rb")
        requests.post(self.url+"upload", files={"file": file}, headers=self.headers)
