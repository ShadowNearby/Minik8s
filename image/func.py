def run(x,y):
   z = x+y
   return {"z":z} 

def main(userparams):
    x =userparams.get('x')
    y = userparams.get('y')
    return run(x,y)