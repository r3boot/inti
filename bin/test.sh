#!/bin/bash

curl -i -H "Accept: application/json" -X GET http://localhost:7231/off

curl -i -H "Accept: application/json" -X PUT -d '{"frame": [255,0,0, 255,0,0, 255,0,0, 255,0,0, 255,0,0], "duration": 1000}' http://localhost:7231/frame
curl -i -H "Accept: application/json" -X PUT -d '{"frame": [0,255,0, 0,255,0, 0,255,0, 0,255,0, 0,255,0], "duration": 1000}' http://localhost:7231/frame
curl -i -H "Accept: application/json" -X PUT -d '{"frame": [0,0,255, 0,0,255, 0,0,255, 0,0,255, 0,0,255], "duration": 1000}' http://localhost:7231/frame

sleep 3
curl -i -H "Accept: application/json" -X GET http://localhost:7231/on
sleep 3
curl -i -H "Accept: application/json" -X GET http://localhost:7231/off
