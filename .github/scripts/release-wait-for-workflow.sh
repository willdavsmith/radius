# ------------------------------------------------------------
# Copyright 2023 The Radius Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#    
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ------------------------------------------------------------

set -ex

# REPOSITORY is the repository to search in
# (e.g. radius)
REPOSITORY=$1

if [[ -z "$REPOSITORY" ]]; then
    echo "Error: REPOSITORY is not set."
    exit 1
fi

# TAG_NAME is the Radius release tag name
# (e.g. v0.29.0, v0.29.0-rc1)
TAG_NAME=$2

if [[ -z "$TAG_NAME" ]]; then
    echo "Error: TAG_NAME is not set."
    exit 1
fi

# WORKFLOW_NAME is the GitHub workflow name
# e.g. "Build and Test"
WORKFLOW_NAME=$3

if [[ -z "$WORKFLOW_NAME" ]]; then
    echo "Error: WORKFLOW_NAME is not set."
    exit 1
fi

MAX_RETRIES=30
RETRY_INTERVAL=60

for ((i=0; i<MAX_RETRIES; i++)); do
    RUN_ID=$(gh run list --limit 1 --workflow "$WORKFLOW_NAME" -b $TAG_NAME --repo radius-project/$REPOSITORY --json databaseId --jq '.[0].databaseId')
    echo "RUN_ID: ${RUN_ID}"

    RUN_STATUS=$(gh run view $RUN_ID --json status --jq '.status')
    echo "RUN_STATUS: ${RUN_STATUS}"

    RUN_CONCLUSION=$(gh run view $RUN_ID --json conclusion --jq '.conclusion')
    echo "RUN_CONCLUSION: ${RUN_CONCLUSION}"

    if [[ "${RUN_STATUS}" == "completed" && "${RUN_CONCLUSION}" == "success" ]]; then
        echo "Run status is: ${RUN_STATUS} and conclusion is: ${RUN_CONCLUSION}."
        break
    fi

    if [[ $i -eq $((MAX_RETRIES - 1)) ]]; then
        echo "Error: Maximum retries reached. Run status: ${RUN_STATUS}, Run conclusion: ${RUN_CONCLUSION}."
        exit 1
    fi

    echo "Retrying in ${RETRY_INTERVAL} seconds..."
    sleep $RETRY_INTERVAL
done
