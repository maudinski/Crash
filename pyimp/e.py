# e stands for error
def throwError(*msgs):
    for m in msgs:
        print(m, end='')
    print()
    exit()

def strictlyString(suspect):
    if isFloat(suspect):
        return False
    try:
        int(suspect)
        return False
    except ValueError:
        return True

def isFloat(suspect):
    try:
        float(suspect)
        suspect.index('.')
        return True
    except ValueError:
        return False
