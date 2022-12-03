from base64 import b64encode
import requests
import json
import os

class Pastebin:
    def __init__(self, url):
        self.url = url

    def _httpError(self, request_for, response_code):
        raise Exception("{} failed with response code {}".format(request_for, response_code))

    def count(self):
        resp = requests.get(self.url+"retrieve/fileCount")
        if not resp.ok:
            self._httpError("Fetching file count", resp.status_code)
        return int(resp.text)

    def detailsOfNthNewest(self, n):
        resp = requests.post(self.url+"retrieve", str(n))
        if not resp.ok:
            self._httpError("Fetching details", resp.status_code)
        return json.loads(resp.text)

    def upload(self, filepath):
        file = open(filepath, "rb")
        resp = requests.post(self.url+"upload", files={"file": file})
        if not resp.ok and resp.status_code != 413:
            self._httpError("File upload", resp.status_code)
        return int(resp.text) if resp.text else 0

    def downloadNth(self, n, savePath="."):
        details = self.detailsOfNthNewest(n)
        resp = requests.get(self.url + "static/" + details["ID"])
        with open(os.path.join(savePath, details["Name"]), "wb") as f:
            f.write(resp.content)
