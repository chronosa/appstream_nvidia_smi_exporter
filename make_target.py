#!/usr/bin/env python3.6

import json
from pprint import pprint

from config import CONFIG

import boto3
from boto3.session import Session

if __name__ == "__main__":
    s = Session()
    regions = s.get_available_regions('appstream')

    target_list = []
    # TODO-SJ Get Credential Informations
    for region in regions:
        appstream = boto3.client('appstream', region_name=region)

        fleets_meta = appstream.describe_fleets()['Fleets']
        #pprint(fleets_meta)

        for fleet_meta in fleets_meta:
            fleet_name = fleet_meta["Name"]

            assigend_stacks = appstream.list_associated_stacks(FleetName=fleet_name)['Names']
            for stack_name in assigend_stacks:
                sessions_meta = appstream.describe_sessions(
                    StackName=stack_name,
                    FleetName=fleet_name
                )["Sessions"]

                for session_meta in sessions_meta:
                    session_id = session_meta["Id"]
                    user_id = session_meta["UserId"]
                    private_ip = session_meta["NetworkAccessConfiguration"]["EniPrivateIpAddress"]

                    target_list.append({
                        "labels": {
                            "user_id" : user_id,
                            "session_id" : session_id
                        },
                        "targets": [
                            f"{private_ip}:9101"
                        ]
                    })

    with open(CONFIG["TARGET_FILE"], "w") as fp:
        fp.write(json.dumps(target_list, sort_keys=True, indent=4))
