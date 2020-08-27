#!/bin/bash
ssh-keygen
eval `ssh-agent`
ssh-add ~/.ssh/id_rsa
cat ~/.ssh/id_rsa.pub