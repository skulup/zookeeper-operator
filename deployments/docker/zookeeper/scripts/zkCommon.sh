#!/usr/bin/env bash

#
# Copyright 2020 Skulup Ltd, Open Collaborators
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

RETRIES=20

POD_NAME=$(hostname -s)
CONFIG_DIR=$DATA_DIR/conf
MYID_FILE=$DATA_DIR/myid
STATIC_CONFIG_FILE=$CONFIG_DIR/zoo.cfg
DYNAMIC_CONFIG_FILE=$CONFIG_DIR/zoo.cfg.dynamic

CLIENT_HOST="127.0.0.1"

function zkServerConfig() {
  role=$1
  HOST="$POD_NAME.$SERVICE_NAME"
  echo "$HOST:$QUORUM_PORT:$LEADER_PORT:$role;0.0.0.0:$CLIENT_PORT"
}

function zkClientUrl() {
  set +e
  nslookup "$SERVICE_NAME" &>/dev/null
  if [[ $? -eq 0 ]]; then
    set -e
    echo "$SERVICE_NAME:$CLIENT_PORT"
  else
    retries=$RETRIES
    while [ $retries -gt 0 ]; do
      sleep 2
      echo "zkClientUrl() retry countdown: $retries" >&2
      nslookup "$SERVICE_NAME" &>/dev/null
      if [[ $? -eq 0 ]]; then
        echo "$SERVICE_NAME:$CLIENT_PORT"
        return
      fi
      retries=$((retries - 1))
    done
    set -e
    echo "zkClientUrl() failed: unable to lookup client host($SERVICE_NAME)"
    exit 1
  fi
}

function checkEnsemblePresence() {
  set +e
  ## Check if there is already an existing ensemble
  LOOKUP_RESULT=$(nslookup "$SERVICE_NAME")
  if [[ $? -eq 0 ]]; then
    return 0
  elif echo "$LOOKUP_RESULT" | grep -q "server can't find $SERVICE_NAME"; then
    # If this node is not the first i.e `$1 -ne 1` in the ensemble server sequence,
    # it means we the first may already be running. Since we failed, it's likely due
    # to DNS update delay for the ensemble service name. Below, we sleep for a bit
    # and retry the DNS resolution if possible
    RECURSIVE_RETRIES=${2:-$RETRIES}
    if [[ $1 -ne 1 && $RECURSIVE_RETRIES -gt 0 ]]; then
      echo "The ensemble service $LOOKUP_RESULT is not yet available. retrying in 2 seconds. retry-countdown: $RECURSIVE_RETRIES" >&2
      sleep 2
      nextRetry=$((RECURSIVE_RETRIES - 1))
      checkEnsemblePresence "$1" $nextRetry
      return $?
    fi
    echo "could not detect any existing ensemble:: $LOOKUP_RESULT ::" >&2
    return 1
  else ## DNS lookup failed; do a sleep-then retry loop for a finite time
    retries=$RETRIES
    while [ $retries -gt 0 ]; do
      sleep 2
      retries=$((retries - 1))
      echo "checkEnsemblePresence() retry-countdown: $retries" >&2
      nslookup "$SERVICE_NAME" &>/dev/null
      if [[ $? -eq 0 ]]; then
        return 0
      fi
    done
    return 1
  fi
}
