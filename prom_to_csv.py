#!/usr/bin/env python3.6

import csv
import requests
import sys

writeHeader=True
# 사용하지 않음. 참조용
def GetMetrixNames(url):
    response = requests.get('{0}/api/v1/label/__name__/values'.format(url))
    names = response.json()['data']
    #Return metrix names
    return names

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print('Usage: {0} http://localhost:9090'.format(sys.argv[0]))
        sys.exit(1)
    prom_url = sys.argv[1]

    writer = csv.writer(sys.stdout)

    #metrixNames=GetMetrixNames(sys.argv[1])
    metrixNames=["utilization_gpu", "utilization_memory"]
    for metrixName in metrixNames:
        #now its hardcoded for hourly
        response = requests.get(
                       '{0}/api/v1/query'.format(prom_url),
                        params={'query': metrixName+'[1d]'}
                   )
        results = response.json()['data']['result']

        # Build a list of all labelnames used.
        #gets all keys and discard __name__
        labelnames = set()
        for result in results:
           labelnames.update(result['metric'].keys())
        # Canonicalize
        labelnames.discard('__name__')
        labelnames = sorted(labelnames)
        # Write the samples.
        if writeHeader:
            writer.writerow(['name'] + labelnames + ['timestamp', 'value'])
            writeHeader=False
        for result in results:
            labels = []
            for label in labelnames:
                labels.append(result['metric'].get(label, ''))

            for value in result['values']:
                l = [result['metric'].get('__name__', '')] + labels + value
                writer.writerow(l)
