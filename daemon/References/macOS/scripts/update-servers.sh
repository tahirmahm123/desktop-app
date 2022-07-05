#!/bin/bash
cd "$( dirname "${BASH_SOURCE[0]}" )"

# if servers.json not exists or it was updated more that 1 hour - update it
# CAREFUL! (file can be updated from Git. Therefore, would be not possible to update it from website during 60 mins )
#if [[ ! -r "etc/servers.json" || $(find "etc/servers.json" -mmin +60) ]]; then

#fi
