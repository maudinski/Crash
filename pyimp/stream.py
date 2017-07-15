# used as an iterator of sorts for the file contents
class Stream:
    def __init__(self, filestring):
        self.filestring = filestring
        self.len = len(filestring)
        self.pos = -1

    def next(self):
        self.pos += 1
        if self.pos == self.len:
            return None
        return self.filestring[self.pos]

    def peek(self):
        if self.pos + 1 == self.len:
            return None
        return self.filestring[self.pos + 1]

    def go_back(self):
        self.pos -= 1