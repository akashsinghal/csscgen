#!/usr/bin/env bash

# NOTE: This script assumes you are pushing to a registry that is already logged. Please run `docker login <registry host name>` prior.

set -e

usage() { printf "Usage: $0 [-r <registry host name> ] [-n <number of referrers to create per subject>] [-s <number of subjects to create]\n" 1>&2; exit 1; }

while getopts "r:n:s:" o; do
  case "${o}" in
    r) registry=${OPTARG} ;;
    n) num_referrers=${OPTARG} ;;
    s) num_subjects=${OPTARG} ;;
    *) usage ;;
  esac
done
repo="${num_subjects}-containers-${num_referrers}-referrers"
# check all parameters are specified
if [ -z "${registry}" ] || [ -z "${num_referrers}" ] || [ -z "${num_subjects}" ]; then
  usage
fi

# for each subject
for ((i=1;i<=${num_subjects};i++)); do
  # build a unique scratch dockerfile and build image
  docker build github.com/wabbit-networks/net-monitor --build-arg="TEXT=repository ${repo} image ${i}" --build-arg="SLEEP=120m" -t ${registry}/${repo}:${i}
  # push new image to specified registry, repo, and tag
  docker push ${registry}/${repo}:${i}
  docker image rm ${registry}/${repo}:${i}
  # add specified number of referrers to the image
  sleep 2s
  for ((j=1;j<=${num_referrers};j++)); do
      notation sign --signature-format cose --key wabbit-networks-io-pipeline ${registry}/${repo}:${i}
  done
done
set +e
