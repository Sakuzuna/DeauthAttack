import subprocess

subprocess.run(['npm', 'init', '-y'])
subprocess.run(['npm', 'install', 'express', 'ts-node', 'typescript', '@types/node', '@types/express'])
subprocess.run(['npm', 'install'])
subprocess.run(['npm', 'run', 'dev'])
