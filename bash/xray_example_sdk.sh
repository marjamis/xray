#!/bin/bash

function generateTraceId {
  XRAY_VERSION=1
  HEX_TIMESTAMP=$(printf '%x\n' $(date +%s))
  RANDOM_HEX_DATA=$(dd if=/dev/random bs=12 count=1 2>/dev/null | od -An -tx1 | tr -d ' \t\n')
  echo "$XRAY_VERSION-$HEX_TIMESTAMP-$RANDOM_HEX_DATA"
}

function generateSegmentId {
  dd if=/dev/random bs=8 count=1 2>/dev/null | od -An -tx1 | tr -d ' \t\n'
}

function generateSegment {
  START_TIME=$(date +%s.000)
  END_TIME=$(date +%s.100)
  echo \
    '{
      "trace_id": "'$TRACE_ID'",
      "id": "'$1'",
      "start_time": '$START_TIME',
      "end_time": '$END_TIME',
      "name": "'$2'",
      "metadata": {
        "info": {
          "version": "0.0.1"
        }
      },
      "http": {
        "request": {
          "url": "http://10.0.0.2/",
          "method": "GET",
          "user_agent": "curl",
          "client_ip": "::ffff:10.0.0.3"
        },
        "response": {
          "status": 200
        }
      }
    }' | jq -c .
}

function generateSubSegment {
  SEGMENT_ID=$(generateSegmentId)
  START_TIME=$(date +%s.000)
  END_TIME=$(date +%s.100)
  echo \
    '{
      "trace_id": "'$TRACE_ID'",
      "id": "'$SEGMENT_ID'",
      "start_time": '$START_TIME',
      "end_time": '$END_TIME',
      "name": "'$2'",
      "parent_id": "'$1'",
      "subsegments": [
        '$(generatePartialSubsegment 200)',
        '$(generatePartialSubsegment 200)',
        '$(generatePartialSubsegment 503)'
      ],
      "metadata": {
        "debug": {
          "intro": "yes"
        }
      }
    }' | jq -c .
}

function generatePartialSubsegment {
  START_TIME=$(date +%s.000)
  END_TIME=$(date +%s.100)
  echo \
    '{
      "id": "'$(generateSegmentId)'",
      "name": "internalMethod",
      "start_time": '$START_TIME',
      "end_time": '$END_TIME',
      "http": {
        "request": {
          "url": "http://10.0.0.3:3000/api/",
          "method": "GET",
          "user_agent": "custom",
          "client_ip": "::ffff:10.0.0.4"
        },
        "response": {
          "status": '$1'
        }
      }
    }' | jq -c .
}

function daemon {
  echo
  HEADER='{"format":"json","version":1}'
  TRACE_DATA="$HEADER\n$@"
  UDP_IP="127.0.0.1"
  UDP_PORT=2000
  echo $TRACE_DATA
  # Bash alias to send network traffic. localhost by default doesn't work everywhere so I used 127.0.0.1
  # BUG fix the partial data being sent issue
  echo -n $TRACE_DATA > /dev/udp/127.0.0.1/2000
}

function api {
  echo
  echo $@
  aws xray put-trace-segments --trace-segment-documents $@
}

function call {
  METHOD=$1
  TRACE_ID=$(generateTraceId)

  PARENT_SEGMENT_ID=$(generateSegmentId)
  $METHOD $(generateSegment $PARENT_SEGMENT_ID example.com) $(generateSubSegment $PARENT_SEGMENT_ID api.example.com) $(generateSubSegment $PARENT_SEGMENT_ID new-api.example.com)

  PARENT_SEGMENT_ID=$(generateSegmentId)
  $METHOD $(generateSegment $PARENT_SEGMENT_ID test.example.com) $(generateSubSegment $PARENT_SEGMENT_ID api.test.example.com)
}

if [ "$#" -lt 1 ]; then daemon; fi

while getopts "ad" opt; do
  case $opt in
    a) call api ;;
    d) call daemon ;;
    \?) echo "Invalid option: -$OPTARG" >&2 ;;
  esac
done
