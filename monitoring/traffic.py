import argparse
import time
import urllib.request
import urllib.error

parser = argparse.ArgumentParser()
parser.add_argument('--seconds', type=int, default=660)
parser.add_argument('--errors', action='store_true')
args = parser.parse_args()
end = time.monotonic() + args.seconds
while time.monotonic() < end:
    for _ in range(4):
        urllib.request.urlopen('http://localhost:8080/notes').close()
    if args.errors:
        request = urllib.request.Request('http://localhost:8080/notes', data=b'{',
                                         headers={'Content-Type': 'application/json'})
    else:
        request = urllib.request.Request('http://localhost:8080/notes',
                                         data=b'{"title":"traffic","body":"lab8"}',
                                         headers={'Content-Type': 'application/json'})
    try:
        urllib.request.urlopen(request).close()
    except urllib.error.HTTPError as error:
        if error.code != 400 or not args.errors:
            raise
        error.close()
    time.sleep(1)
