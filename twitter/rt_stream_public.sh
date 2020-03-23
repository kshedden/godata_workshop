#!/bin/bash

# Replace XYZ below with actual information for your Twitter developer account

consumer_key="XYZ"
consumer_secret="XYZ"
access_token="XYZ"
access_secret="XYZ"

go run rt_stream.go -consumer-key=$consumer_key -consumer-secret=$consumer_secret -access-token=$access_token -access-secret=$access_secret
