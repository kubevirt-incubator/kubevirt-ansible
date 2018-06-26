#!/bin/bash

set -e

source hack/common.sh

prefix=kubevirt
tag=v0.7.0-alpha.4
kubeconfig=~/.kube/config

${TESTS_OUT_DIR}/tests.test -kubeconfig=$kubeconfig -tag=$tag -prefix=$prefix -test.timeout 60m
