//go:build linux && !no_bpf
// +build linux,!no_bpf

/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package recorder

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"unsafe"

	"github.com/aquasecurity/libbpfgo"
	"github.com/containers/common/pkg/seccomp"
	libseccomp "github.com/seccomp/libseccomp-golang"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"

	"sigs.k8s.io/security-profiles-operator/internal/pkg/daemon/bpfrecorder"
)

type defaultImpl struct{}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate -header ../../../../hack/boilerplate/boilerplate.generatego.bpf.txt
//counterfeiter:generate . impl
type impl interface {
	LoadBpfRecorder(*bpfrecorder.BpfRecorder) error
	UnloadBpfRecorder(*bpfrecorder.BpfRecorder)
	Command(string, ...string) *exec.Cmd
	CmdStart(*exec.Cmd) error
	CmdPid(*exec.Cmd) uint32
	CmdWait(*exec.Cmd) error
	SyscallsIterator(*bpfrecorder.BpfRecorder) *libbpfgo.BPFMapIterator
	IteratorNext(*libbpfgo.BPFMapIterator) bool
	IteratorKey(*libbpfgo.BPFMapIterator) []byte
	SyscallsGetValue(*bpfrecorder.BpfRecorder, uint32) ([]byte, error)
	GetName(libseccomp.ScmpSyscall) (string, error)
	MarshalIndent(any, string, string) ([]byte, error)
	WriteFile(string, []byte, os.FileMode) error
	Create(string) (*os.File, error)
	CloseFile(*os.File)
	PrintObj(printers.YAMLPrinter, runtime.Object, io.Writer) error
	GoArchToSeccompArch(string) (seccomp.Arch, error)
}

func (*defaultImpl) LoadBpfRecorder(b *bpfrecorder.BpfRecorder) error {
	return b.Load(false)
}

func (*defaultImpl) UnloadBpfRecorder(b *bpfrecorder.BpfRecorder) {
	b.Unload()
}

func (*defaultImpl) Command(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func (*defaultImpl) CmdStart(cmd *exec.Cmd) error {
	return cmd.Start()
}

func (*defaultImpl) CmdPid(cmd *exec.Cmd) uint32 {
	return uint32(cmd.Process.Pid)
}

func (*defaultImpl) CmdWait(cmd *exec.Cmd) error {
	return cmd.Wait()
}

func (*defaultImpl) SyscallsIterator(b *bpfrecorder.BpfRecorder) *libbpfgo.BPFMapIterator {
	return b.Syscalls().Iterator()
}

func (*defaultImpl) IteratorNext(it *libbpfgo.BPFMapIterator) bool {
	return it.Next()
}

func (*defaultImpl) IteratorKey(it *libbpfgo.BPFMapIterator) []byte {
	return it.Key()
}

func (*defaultImpl) SyscallsGetValue(b *bpfrecorder.BpfRecorder, pid uint32) ([]byte, error) {
	return b.Syscalls().GetValue(unsafe.Pointer(&pid))
}

func (*defaultImpl) GetName(s libseccomp.ScmpSyscall) (string, error) {
	return s.GetName()
}

func (*defaultImpl) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (*defaultImpl) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (*defaultImpl) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (*defaultImpl) CloseFile(file *os.File) {
	file.Close()
}

func (*defaultImpl) PrintObj(p printers.YAMLPrinter, obj runtime.Object, w io.Writer) error {
	return p.PrintObj(obj, w)
}

func (*defaultImpl) GoArchToSeccompArch(arch string) (seccomp.Arch, error) {
	return seccomp.GoArchToSeccompArch(arch)
}