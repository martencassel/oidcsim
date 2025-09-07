#!/bin/bash


curl "http://localhost:8080/authorize?response_type=code&client_id=client1&redirect_uri=http://localhost:8080/callback&scope=read&state=xyz"
