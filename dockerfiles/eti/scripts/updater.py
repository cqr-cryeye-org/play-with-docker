import hashlib
import io
import json
import subprocess
import timeit
import zipfile
from pathlib import Path

import requests
import yaml

import config
from utils import launcher


def extract_zip(data: bytes):
    """Download a ZIP file and extract its contents in memory yields (filename, file-like object) pairs"""

    with zipfile.ZipFile(io.BytesIO(data)) as zf:
        for z in zf.infolist():
            with zf.open(z) as f:
                yield z.filename, f.read()


def load_hash_table(path: Path):
    with open(path) as f:
        return json.load(fp=f)


def save_hash_table(path: Path, data: dict):
    with open(path, 'w') as f:
        json.dump(data, fp=f, indent=2)


def load_docker_stack(path: Path):
    with open(path, 'r') as stream:
        try:
            return yaml.safe_load(stream)
        except yaml.YAMLError as e:
            print(e)


def get_stack_images(path: Path):
    file_data = load_docker_stack(path=path)
    services = [service for service in file_data.get('services', {}).values()]
    return [service.get('image') for service in services] if services else services


def pull_stack_images(stack_images: list):
    for image in stack_images:
        res = subprocess.run(["docker", "pull", image])
        return res.returncode


def save_stack_images(stack_name: str, stack_images: list):
    return subprocess.run(["docker", "save", "--output", f"/tmp/docker-images/tars/{stack_name}.tar", *stack_images])


def get_zip_path(s: requests.Session, h: str) -> tuple:
    """Get path to zip file with docker stacks. Check if the current hash of the file matches the receive."""

    try:
        resp = s.get(config.BASE_URL)
        if resp.status_code == 200:
            match = config.pattern.search(string=resp.text)
            zip_path, zip_hash = match.group(), *match.groups()
            if zip_hash != h:
                return zip_path, zip_hash
    except Exception as e:
        print(e)
    exit()


def main():
    s = requests.Session()
    s.headers.update(config.HEADERS)

    hash_table = load_hash_table(path=config.HASH_TABLE_PATH) if config.HASH_TABLE_PATH.exists() else {}
    zip_path, zip_hash = get_zip_path(s=s, h=hash_table.get('zip_hash', ''))
    hash_table.update({'zip_hash': zip_hash})

    try:
        resp = s.get(f'{config.GIST_URL}{zip_path}')
        if resp.status_code != 200:
            exit()

        for filename, content in extract_zip(data=resp.content):
            if 'docker-stack' in filename:
                stack_name = filename.split('/')[-1]
                file_hash = hashlib.sha224(content).hexdigest()
                if file_hash != hash_table.get(stack_name):
                    hash_table.update({stack_name: file_hash})
                    with open(config.DOCKER_STACKS_PATH.joinpath(stack_name), 'wb') as f:
                        f.write(content)

    except Exception as e:
        print(e)

    save_hash_table(path=config.HASH_TABLE_PATH, data=hash_table)

    for path in config.DOCKER_STACKS_PATH.iterdir():
        stack_images = get_stack_images(path=path)
        if stack_images:
            print(path)
            print(stack_images)
            code = pull_stack_images(stack_images=stack_images)
            print(code)
            if not code:
                save_stack_images(stack_name=path.name, stack_images=stack_images)
                print('Saved')


if __name__ == '__main__':
    with launcher():
        timeit.timeit(main, number=1)
