A RESTful api on top of a DMX RGB led controller

== Usage ==

$ ./bin/intid -D

$ curl -i -H "Accept: application/json" -X PUT -d '{"frame": [255,0,0, 255,0,0, 255,0,0, 255,0,0, 255,0,0], "duration": 1000}' http://localhost:7231/frame

$ curl -i -H "Accept: application/json" -X PUT -d '{"frame": [0,255,0, 0,255,0, 0,255,0, 0,255,0, 0,255,0], "duration": 1000}' http://localhost:7231/frame

$ curl -i -H "Accept: application/json" -X PUT -d '{"frame": [0,0,255, 0,0,255, 0,0,255, 0,0,255, 0,0,255], "duration": 1000}' http://localhost:7231/frame
