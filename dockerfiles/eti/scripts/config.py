import re
from pathlib import Path

HEADERS = {
    'User-Agent': 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
    'Accept-Language': 'en-US,en;q=0.5',
    'Referer': 'https://gist.github.com/cqr-cryeye/4f0210d3752eb01b8e3e1ec3cc28ec4e/revisions',
    'DNT': '1',
    'Connection': 'keep-alive',
    'Upgrade-Insecure-Requests': '1',
    'Cache-Control': 'max-age=0',
    'TE': 'Trailers',
}

GIST_URL = 'https://gist.github.com/'

BASE_URL = f'{GIST_URL}cqr-cryeye/4f0210d3752eb01b8e3e1ec3cc28ec4e/'

DOCKER_STACKS_PATH = Path('docker-stacks')

HASH_TABLE_PATH = Path('hash_table.json')

APP_PID_PATH = Path('app.pid')

pattern = re.compile(r'cqr-cryeye/4f0210d3752eb01b8e3e1ec3cc28ec4e/archive/(?P<hash>[\w]+).zip')
