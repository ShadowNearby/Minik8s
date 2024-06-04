def run(z):
    return {"str": str(z)}

def main(userparams):
    z = userparams.get('z')
    return run(z)
