#!/bin/bash

# Start the first process
/free5gc/bin/upf -c /free5gc/config/upfcfg.yaml &

# Start the second process
/go-tt/go-tt &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?