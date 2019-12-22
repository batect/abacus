#! /usr/bin/env bash

set -euo pipefail

# Unfortunately, domains.google doesn't have an API right now.
# So instead of just configuring the records required to delegate the zone to GCP's Cloud DNS,
# we have to manually configure it. This scripts checks that we configured it correctly.

ZONE_NAME=app-zone

function main() {
  echo "Retrieving zone configuration..."
  zoneConfig=$(retrieveZoneConfig)
  dnsName=$(echo "$zoneConfig" | jq -r '.dnsName')

  checkNameServers "$zoneConfig" "$dnsName"
  checkDSRecord "$dnsName"
}

function checkNameServers() {
  echo "Checking name server configuration..."
  zoneConfig=$1
  dnsName=$2
  desiredNameServers=$(echo "$zoneConfig" | jq -r '.nameServers | .[]' | sort)
  currentNameServers=$(dig "$dnsName" NS +short | sort)

  echo "Desired name servers for $dnsName are:"
  echo "$desiredNameServers"
  echo
  echo "Current name servers for $dnsName are:"
  echo "$currentNameServers"
  echo

  if ! diff <(echo "$desiredNameServers") <(echo "$currentNameServers"); then
    echo
    echo "Name servers are not configured correctly!"
    exit 1
  else
    echo "Name servers are configured correctly."
    echo
  fi
}

function checkDSRecord() {
  echo "Checking DS record..."
  key=$(gcloud dns dns-keys list --zone "$ZONE_NAME" --project "$GOOGLE_PROJECT" --format json | jq '.[] | select(.type == "keySigning")')
  keyTag=$(echo "$key" | jq -r '.keyTag')
  algorithm=$(echo "$key" | jq -r '.algorithm')
  digestInfo=$(echo "$key" | jq -r '.digests[0]')
  digest=$(echo "$digestInfo" | jq -r '.digest')
  digestType=$(echo "$digestInfo" | jq -r '.type')

  desiredDSRecord=$(constructDSRecord "$keyTag" "$algorithm" "$digestType" "$digest")
  currentDSRecord=$(dig "$dnsName" DS +short +nosplit)

  echo "Desired DS record for $dnsName is:"
  echo "$desiredDSRecord"
  echo
  echo "Current DS record for $dnsName is:"
  echo "$currentDSRecord"
  echo

  if ! diff <(echo "$desiredDSRecord") <(echo "$currentDSRecord"); then
    echo
    echo "DS record is not configured correctly!"
    exit 1
  else
    echo "DS record is configured correctly."
  fi
}

function retrieveZoneConfig() {
  gcloud dns managed-zones describe "$ZONE_NAME" --project "$GOOGLE_PROJECT" --format json
}

function constructDSRecord() {
  if [[ "$2" != "rsasha256" ]]; then
    echo "Unkown algorithm: $2" > /dev/stderr
    echo "Update this script based on table on https://cloud.google.com/community/tutorials/dnssec-cloud-dns-domains#get_dnskey_information"
    exit 1
  fi

  if [[ "$3" != "sha256" ]]; then
    echo "Unkown digest type: $3" > /dev/stderr
    echo "Update this script based on table on https://cloud.google.com/community/tutorials/dnssec-cloud-dns-domains#get_dnskey_information"
    exit 1
  fi

  echo "$1 8 2 $4"
}

main
