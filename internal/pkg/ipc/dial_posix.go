// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !windows

package ipc

import (
	"context"
	"os/exec"
	"strings"

	"google.golang.org/grpc"
)

const lsofOutputPrefix = "n/"

// DialContext connects to the serverAddress for grpc.  Returns immediately and will make the connection lazily when needed.
// This is the recommended approach to dialing from the grpc documentation (https://github.com/grpc/grpc-go/blob/master/Documentation/anti-patterns.md)
func DialContext(ctx context.Context, serverAddr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, serverAddr, grpc.WithInsecure(), grpc.WithMaxMsgSize(GrpcMaxMsgSize))
}

// DialContextWithBlock connects to the serverAddress for grpc, but waits until the connection is made (via WithBlock) until returning.
// This is NOT the recommended approach to dialing, but is needed for bootstrap which relies on WithBlock as a check that reproxy has started successfully.
// Also needed for reproxy connecting to the depsscannerservice.
// TODO(b/290804932): Remove the dependence on WithBlock as a startup check.
func DialContextWithBlock(ctx context.Context, serverAddr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, serverAddr, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithMaxMsgSize(GrpcMaxMsgSize))
}

// GetAllReproxySockets returns all unix sockets where an reproxy service is listening.
func GetAllReproxySockets(ctx context.Context) ([]string, error) {
	lsofOutput, err := execLsof("-U", "-c", "reproxy", "-a", "-Fn")

	if err != nil {
		// Lsof completes with exit code 1 when there is no output
		if err.(*exec.ExitError).ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}
	return parseSockets(lsofOutput), nil
}

func parseSockets(lsofOutput string) []string {
	var sockets []string
	seenSockets := map[string]bool{}
	for _, line := range strings.Split(lsofOutput, "\n") {
		if socket := parseSocketLine(line); socket != "" && !seenSockets[socket] {
			sockets = append(sockets, "unix://"+socket)
			seenSockets[socket] = true
		}
	}
	return sockets
}

func parseSocketLine(line string) string {
	for _, s := range strings.Split(line, " ") {
		if strings.HasPrefix(s, lsofOutputPrefix) {
			return strings.TrimLeft(s, "n")
		}
	}
	return ""
}

func execLsof(args ...string) (string, error) {
	command := "lsof"
	args = append([]string{"-w"}, args...)
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
