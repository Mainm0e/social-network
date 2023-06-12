#!/bin/bash
# This Bash script works only on Mac OS
# ChatGTP will help you find what is command for your operating system
# This script will open 2 new tabs in Terminal and run the Go and React programs

# Open a new Terminal window
# this one for testing
# osascript -e 'tell application "Terminal" to do script "ls"'

# Open a new tab and run the Go program
osascript -e 'tell application "Terminal" to do script "cd jsBranch/social-network; bash start_backend.sh"'
# Please change the path to the project folder         #^^^^^^^^^^^^^^^^^^^^^^^^^^#

# Open a new tab and execute npm start
osascript -e 'tell application "Terminal" to do script "cd jsBranch/social-network; bash start_frontend.sh"'
# Please change the path to the project folder         #^^^^^^^^^^^^^^^^^^^^^^^^^^#
