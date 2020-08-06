#! /usr/bin/env bash

set -euo pipefail

files=$(find . \( -name "*.go" -or -name "*.tf" \) -type f -not -path './batect/caches/*' -not -path './vendor/*')

desired_header="\
// Copyright 2019-$(date +%Y) Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// and the Commons Clause License Condition v1.0 (the \"Condition\");
// you may not use this file except in compliance with both the License and Condition.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// You may obtain a copy of the Condition at
//
//     https://commonsclause.com/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License and the Condition is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See both the License and the Condition for the specific language governing permissions and
// limitations under the License and the Condition."

header_lines=$(echo "$desired_header" | wc -l)

declare -a non_compliant_files=()

for file in $files; do
  echo -n "Checking $file..."
  current_header=$(head -n "$header_lines" "$file")

  if [[ "$current_header" != "$desired_header" ]]; then
    non_compliant_files+=("$file")
    echo "$(tput setaf 1)not ok!$(tput sgr0)"
  else
    echo "$(tput setaf 2)ok$(tput sgr0)"
  fi
done

echo

if [[ "${#non_compliant_files[@]}" -ne "0" ]]; then
  echo "The following files are missing the required license header:"

  for file in "${non_compliant_files[@]}"; do
    echo " - $file"
  done

  echo
  echo "The required license header is:"
  echo "$desired_header"

  exit 1
else
  echo "All files are compliant."
fi
