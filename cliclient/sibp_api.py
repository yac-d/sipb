from base64 import b64encode
import requests
import json

class Pastebin:
    def __init__(self, host, port=80, basicauth=None, https=False, pastebinroot="/"):
        self.url = "https://" if https else "http://"
        self.url += "{}:{}{}".format(host, port, pastebinroot)

        self.headers = {}
        if basicauth:
            user_and_pass = b64encode(bytes("{}:{}".format(basicauth[0], basicauth[1]), "ascii")).decode("ascii")
            self.headers = {"Authorization": "Basic {}".format(user_and_pass)}

    def _httpError(self, request_for, response_code):
        raise Exception("{} failed with response code {}".format(request_for, response_code))

    def count(self):
        resp = requests.get(self.url+"retrieve/fileCount", headers=self.headers)
        if not resp.ok:
            self._httpError("Fetching file count", resp.status_code)
        return int(resp.text)

    def detailsOfNthNewest(self, n):
        resp = requests.post(self.url+"retrieve", str(n), headers=self.headers)
        if not resp.ok:
            self._httpError("Fetching details", resp.status_code)
        return json.loads(resp.text)

    def upload(self, filepath):
        file = open(filepath, "rb")
        resp = requests.post(self.url+"upload", files={"file": file}, headers=self.headers)
        if not resp.ok and resp.status_code != 413:
            self._httpError("File upload", resp.status_code)
        return int(resp.text)
