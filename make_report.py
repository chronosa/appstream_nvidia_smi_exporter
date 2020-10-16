#!/usr/bin/env python3.6
from config import CONFIG

import boto3
import pandas as pd

df = pd.read_csv("a")
groupby = df.groupby([df["session_id"], df["user_id"], df["name"]])
#print(groupby['value'].mean())
stats = groupby['value'].agg(['mean', 'max'])
stats.to_csv('abc.csv', sep=',', na_rep='')
