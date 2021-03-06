// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package isolation

import (
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

type namespace struct {
	ns         int
	nsToOption map[int]string
}

// NewNamespace creates new instance of namespace isolation; see: man 7 namespaces.
func NewNamespace(mask int) (Decorator, error) {
	if mask&(syscall.CLONE_NEWIPC|syscall.CLONE_NEWNET|syscall.CLONE_NEWNS|syscall.CLONE_NEWPID|syscall.CLONE_NEWUSER|syscall.CLONE_NEWUTS) == 0 {
		return nil, errors.New("Invalid namespace mask")
	}
	optionMap := make(map[int]string)
	optionMap[syscall.CLONE_NEWPID] = "--fork --pid --mount-proc"
	optionMap[syscall.CLONE_NEWIPC] = "--ipc"
	optionMap[syscall.CLONE_NEWNS] = "--mount"
	optionMap[syscall.CLONE_NEWUTS] = "--uts"
	optionMap[syscall.CLONE_NEWNET] = "--net"
	optionMap[syscall.CLONE_NEWUSER] = "--user"
	return &namespace{ns: mask, nsToOption: optionMap}, nil
}

// Decorate implements Decorator.
func (n *namespace) Decorate(command string) (decorated string) {
	var namespaces []string

	for namespace, option := range n.nsToOption {
		if n.ns&namespace == namespace {
			namespaces = append(namespaces, option)
		}
	}

	return "unshare " + strings.Join(namespaces, " ") + " " + command
}
