def run(x, y):
    z = x + y
    x = x - y
    y = y - x
    print(z)
    return x, y, z

def main(userparams):
    x = userparams.get('x')
    y = userparams.get('y')
    return run(x,y)