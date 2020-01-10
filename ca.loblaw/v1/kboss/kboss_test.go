// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestReplacementTransformer(t *testing.T) {
	th := kusttest_test.MakeEnhancedHarness(t).
		BuildGoPlugin("ca.loblaw", "v1", "kboss")
	defer th.Reset()

	rm := th.LoadAndRunTransformer(`
  apiVersion: ca.loblaw/v1
  kind: kboss
  metadata:
    name: notImportantHere
  file: images.properties
`, `
group: apps
apiVersion: v1
kind: Deployment
metadata:
  name: deploy1
spec:
  template:
    spec:
      containers:
      - image: nginx:1.7.9
        name: nginx-tagged
      - image: nginx:latest
        name: nginx-latest
      - image: foobar:1
        name: replaced-with-digest
      - image: postgres:1.8.0
        name: postgresdb
      initContainers:
      - image: nginx
        name: nginx-notag
      - image: nginx@sha256:111111111111111111
        name: nginx-sha256
      - image: alpine:1.8.0
        name: init-alpine
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: v1
group: apps
kind: Deployment
metadata:
  name: deploy1
spec:
  template:
    spec:
      containers:
      - image: image1
        name: nginx-tagged
      - image: nginx:latest
        name: nginx-latest
      - image: digest:1
        name: replaced-with-digest
      - image: postgres:1.8.0
        name: postgresdb
      initContainers:
      - image: nginx
        name: nginx-notag
      - image: nginx@sha256:111111111111111111
        name: nginx-sha256
      - image: alpine:1.8.0
        name: init-alpine
`)
}
