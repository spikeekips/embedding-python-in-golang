import json

print '[python] > %s is imported' % __file__

def run () :
    print '[python] >\t run `run` function'
    return json.dumps(
            range(10, ),
        )


