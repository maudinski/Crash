# ran by gobuild bash script
from subprocess import call

srcDir = "goimp2"

# ["subDirectory name(inside srcDir)", ["list", "of", "files"]]
dirs = [
    ["", ["crash.go"]],
    ["front", ["data.go", "environementStack.go", "lexer.go", "semanticAnalyzer.go", "expressionParser.go", "parser.go"]],
    ["back", ["generator.go"]],
    ["both", ["ast.go", "nodes.go"]]
]

totalCommand = ["go", "build"]

def addToCommand(dirList):
    global totalCommand, srcDir
    directory = srcDir + "/"
    if dirList[0] != "": # basicall "if there is a sub directory, append it"
        directory += dirList[0] + "/"
    for fn in dirList[1]:
        totalCommand.append(directory + fn)

for d in dirs:
    addToCommand(d)

call(totalCommand)
