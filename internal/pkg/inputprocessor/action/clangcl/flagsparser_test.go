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

package clangcl

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/reclient/internal/pkg/execroot"
	"github.com/bazelbuild/reclient/internal/pkg/inputprocessor/flags"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	fakeExecRoot = "fake"
)

func TestParseFlags(t *testing.T) {
	er, cleanup := execroot.Setup(t, nil)
	defer cleanup()
	tests := []struct {
		name       string
		command    []string
		workingDir string
		files      map[string][]byte
		want       *flags.CommandFlags
	}{
		{
			name:       "chromium clang-cl with a simple rsp file",
			workingDir: "out/Release",
			files: map[string][]byte{
				filepath.Join("out", "Release", "foo", "args.rsp"):  []byte("simple_input.txt\nmore_inputs.txt\n"),
				filepath.Join("out", "Release", "simple_input.txt"): []byte("put filename in rsp file is a hacky way to upload a random file to RBE"),
				filepath.Join("out", "Release", "more_inputs.txt"):  []byte("you can upload more than one file to RBE with rsp file"),
			},
			command: []string{
				"../../third_party/llvm-build/Release+Asserts/bin/clang-cl.exe",
				"/nologo",
				"/showIncludes:user",
				`-imsvc..\..\third_party\depot_tools\win_toolchain\vs_files\a687d8e2e4114d9015eb550e1b156af21381faac\win_sdk\Include\10.0.19041.0\um`,
				`-DCR_CLANG_REVISION="n358615-fb1aa286-3"`,
				`-I../..`,
				`-fcolor-diagnostics`,
				`-fcrash-diagnostics-dir=../../tools/clang/crashreports`,
				`-fprofile-use=prof.txt`,
				`-Xclang`,
				`-mllvm`,
				`-Xclang`,
				`-instcombine-lower-dbg-declare=0`,
				`/Gy`,
				`/FS`,
				`/bigobj`,
				`/utf-8`,
				`/Zc:twoPhase`,
				`/Zc:sizedDealloc-`,
				`/X`,
				`/D__WRL_ENABLE_FUNCTION_STATICS__`,
				`-fmsc-version=1916`,
				`/guard:cf,nochecks`,
				`/Zc:dllexportInlines-`,
				`-m64`,
				`/Brepro`,
				`-Wno-builtin-macro-redefined`,
				`/W4`,
				`/WX`,
				`/wd4091`,
				`/wd4127`,
				`/Od`,
				`/Ob0`,
				`/GF`,
				`/Z7`,
				`/MDd`,
				`/TP`,
				`/wd4577`,
				`/GR-`,
				"@foo/args.rsp",
				`-I../../buildtools/third_party/libc++/trunk/include`,
				`/c`,
				`/clang:-MD`,
				`/clang:-MF`, `/clang:obj/base/third_party/double_conversion/double_conversion/fixed-dtoa.obj.d`,
				`../../base/third_party/double_conversion/double-conversion/fixed-dtoa.cc`,
				`/Foobj/base/third_party/double_conversion/double_conversion/fixed-dtoa.obj`,
				`/Fdobj/base/third_party/double_conversion/double_conversion_cc.pdb`,
			},
			want: &flags.CommandFlags{
				ExecutablePath: "../../third_party/llvm-build/Release+Asserts/bin/clang-cl.exe",
				TargetFilePaths: []string{
					"simple_input.txt",
					"more_inputs.txt",
					"../../base/third_party/double_conversion/double-conversion/fixed-dtoa.cc",
				},
				IncludeDirPaths: []string{
					"../..",
					"../../buildtools/third_party/libc++/trunk/include",
				},
				WorkingDirectory:      "out/Release",
				EmittedDependencyFile: "obj/base/third_party/double_conversion/double_conversion/fixed-dtoa.obj.d",
				ExecRoot:              er,
				Dependencies: []string{
					"prof.txt",
					"foo/args.rsp",
				},
				Flags: []*flags.Flag{
					{Key: "-nologo"},
					{Key: "-showIncludes:user"},
					{
						Key:    "-imsvc",
						Value:  `..\..\third_party\depot_tools\win_toolchain\vs_files\a687d8e2e4114d9015eb550e1b156af21381faac\win_sdk\Include\10.0.19041.0\um`,
						Joined: true,
					},
					{Key: "-D", Value: `CR_CLANG_REVISION="n358615-fb1aa286-3"`, Joined: true},
					{Key: "-I", Value: "../..", Joined: true},
					{Key: "-fcolor-diagnostics"},
					{Key: "-fcrash-diagnostics-dir=", Value: "../../tools/clang/crashreports", Joined: true},
					{Key: "-fprofile-use=", Value: "prof.txt", Joined: true},
					{Key: "-Xclang", Value: "-mllvm"},
					{Key: "-Xclang", Value: "-instcombine-lower-dbg-declare=0"},
					{Key: "-Gy"},
					{Key: "-FS"},
					{Key: "-bigobj"},
					{Key: "-utf-8"},
					{Key: "-Zc:twoPhase"},
					{Key: "-Zc:sizedDealloc-"},
					{Key: "-X"},
					{Key: "-D", Value: "__WRL_ENABLE_FUNCTION_STATICS__", Joined: true},
					{Key: "-fmsc-version=", Value: "1916", Joined: true},
					{Key: "-guard:", Value: "cf,nochecks", Joined: true},
					{Key: "-Zc:dllexportInlines-"},
					{Key: "-m64"},
					{Key: "-Brepro"},
					{Key: "-W", Value: "no-builtin-macro-redefined", Joined: true},
					{Key: "-W4"},
					{Key: "-WX"},
					{Key: "-wd", Value: "4091", Joined: true},
					{Key: "-wd", Value: "4127", Joined: true},
					{Key: "-Od"},
					{Key: "-Ob0"},
					{Key: "-GF"},
					{Key: "-Z7"},
					{Key: "-MDd"},
					{Key: "-TP"},
					{Key: "-wd", Value: "4577", Joined: true},
					{Key: "-GR-"},
					{Value: "@foo/args.rsp"},
					{Key: "-I", Value: "../../buildtools/third_party/libc++/trunk/include", Joined: true},
					{Key: "-c"},
					{Key: "-clang:", Value: "-MD", Joined: true},
					{Key: "-Fd", Value: "obj/base/third_party/double_conversion/double_conversion_cc.pdb", Joined: true},
				},
				OutputFilePaths: []string{
					"obj/base/third_party/double_conversion/double_conversion/fixed-dtoa.obj.d",
					"obj/base/third_party/double_conversion/double_conversion/fixed-dtoa.obj",
				},
			},
		},
		{
			name:       "chromium clang-cl with dependency file defined by -MF in rsp file",
			workingDir: "out/Release",
			files: map[string][]byte{
				filepath.Join("out", "Release", "foo", "args.rsp"): []byte(
					"/c\n" +
						"simple_input.cc\nother_inputs.txt\n" +
						// -MD tells the compiler to generate a dependency file and includes it in the compiled output with self-described dependencies,
						// -MF tells the compiler the dependency file's name and filepath.
						// Generally, there is only one -MF option in a single clang command.
						"/clang:-MD\n" +
						"/clang:-MF\n/clang:obj/base/third_party/double_conversion/double_conversion/fixed-dtoa.obj.d\n" +
						"../../base/third_party/double_conversion/double-conversion/fixed-dtoa.cc\n" +
						"/Foobj/base/third_party/double_conversion/double_conversion/fixed-dtoa.obj\n" +
						"/Fdobj/base/third_party/double_conversion/double_conversion_cc.pdb"),
				filepath.Join("out", "Release", "simple_input.cc"):  []byte("int main() {return 0;}"),
				filepath.Join("out", "Release", "other_inputs.txt"): []byte("some another random text"),
			},
			command: []string{
				"../../third_party/llvm-build/Release+Asserts/bin/clang-cl.exe",
				"/nologo",
				"/showIncludes:user",
				`-imsvc..\..\third_party\depot_tools\win_toolchain\vs_files\a687d8e2e4114d9015eb550e1b156af21381faac\win_sdk\Include\10.0.19041.0\um`,
				`-DCR_CLANG_REVISION="n358615-fb1aa286-3"`,
				`-I../..`,
				`-I../../buildtools/third_party/libc++/trunk/include`,
				"@foo/args.rsp",
			},
			want: &flags.CommandFlags{
				ExecutablePath: "../../third_party/llvm-build/Release+Asserts/bin/clang-cl.exe",
				TargetFilePaths: []string{
					"simple_input.cc",
					"other_inputs.txt",
					"../../base/third_party/double_conversion/double-conversion/fixed-dtoa.cc",
				},
				IncludeDirPaths: []string{
					"../..",
					"../../buildtools/third_party/libc++/trunk/include",
				},
				WorkingDirectory:      "out/Release",
				EmittedDependencyFile: "obj/base/third_party/double_conversion/double_conversion/fixed-dtoa.obj.d",
				ExecRoot:              er,
				Dependencies: []string{
					"foo/args.rsp",
				},
				Flags: []*flags.Flag{
					{Key: "-nologo"},
					{Key: "-showIncludes:user"},
					{
						Key:    "-imsvc",
						Value:  `..\..\third_party\depot_tools\win_toolchain\vs_files\a687d8e2e4114d9015eb550e1b156af21381faac\win_sdk\Include\10.0.19041.0\um`,
						Joined: true,
					},
					{Key: "-D", Value: `CR_CLANG_REVISION="n358615-fb1aa286-3"`, Joined: true},
					{Key: "-I", Value: "../..", Joined: true},
					{Key: "-I", Value: "../../buildtools/third_party/libc++/trunk/include", Joined: true},
					{Value: "@foo/args.rsp"},
				},
				OutputFilePaths: []string{
					"obj/base/third_party/double_conversion/double_conversion/fixed-dtoa.obj.d",
					"obj/base/third_party/double_conversion/double_conversion/fixed-dtoa.obj",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			execroot.AddFilesWithContent(t, er, test.files)
			got, err := parseFlags(ctx, test.command, test.workingDir, er)
			if err != nil {
				t.Errorf("parseFlags(%v,%v,%v).err = %v, want no error.", test.command, test.workingDir, er, err)
			}

			if diff := cmp.Diff(test.want, got, cmpopts.IgnoreUnexported(flags.Flag{})); diff != "" {
				t.Errorf("Test %v, parseFlags(%v,%v,%v) \nreturned diff, (-want +got): %s", test.name, test.command, test.workingDir, er, diff)
			}
		})
	}
}
