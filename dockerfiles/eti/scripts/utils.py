import atexit
import os
from contextlib import contextmanager

import config


@contextmanager
def launcher():
    if config.APP_PID_PATH.exists() and config.APP_PID_PATH.is_file():
        exit()
    with open(config.APP_PID_PATH, 'w') as f:
        f.write('1')
    yield


@atexit.register
def goodbye():
    os.remove(path=config.APP_PID_PATH)
    print("Bye.")
