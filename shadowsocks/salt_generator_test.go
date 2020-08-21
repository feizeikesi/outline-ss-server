// Copyright 2020 Jigsaw Operations LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shadowsocks

import (
	"bytes"
	"testing"
)

func TestRandomSaltGenerator(t *testing.T) {
	if err := RandomSaltGenerator.GetSalt(nil); err != nil {
		t.Error(err)
	}
	salt := make([]byte, 16)
	if err := RandomSaltGenerator.GetSalt(salt); err != nil {
		t.Error(err)
	}
	if bytes.Equal(salt, make([]byte, 16)) {
		t.Error("Salt is all zeros")
	}
}

// Test that ServerSaltGenerator recognizes its own salts
func TestServerSaltRecognized(t *testing.T) {
	ssg := NewServerSaltGenerator("test")

	salt := make([]byte, 32)
	if err := ssg.GetSalt(salt); err != nil {
		t.Fatal(err)
	}
	if !ssg.IsServerSalt(salt) {
		t.Error("Server salt was not recognized")
	}
}

// Test that ServerSaltGenerator doesn't recognize random salts
func TestServerSaltUnrecognized(t *testing.T) {
	ssg := NewServerSaltGenerator("test")

	salt := make([]byte, 32)
	if err := RandomSaltGenerator.GetSalt(salt); err != nil {
		t.Fatal(err)
	}
	if ssg.IsServerSalt(salt) {
		t.Error("Client salt was recognized as a server salt")
	}
}

// Test that ServerSaltGenerator produces different output on each call
func TestServerSaltDifferent(t *testing.T) {
	ssg := NewServerSaltGenerator("test")

	salt1 := make([]byte, 32)
	if err := ssg.GetSalt(salt1); err != nil {
		t.Fatal(err)
	}
	salt2 := make([]byte, 32)
	if err := ssg.GetSalt(salt2); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(salt1, salt2) {
		t.Error("salts should be random")
	}
}

// Test that two ServerSaltGenerators derived from the same secret
// produce different outputs and recognize each other's output.
func TestServerSaltSameSecret(t *testing.T) {
	ssg1 := NewServerSaltGenerator("test")
	ssg2 := NewServerSaltGenerator("test")

	salt1 := make([]byte, 32)
	if err := ssg1.GetSalt(salt1); err != nil {
		t.Fatal(err)
	}
	salt2 := make([]byte, 32)
	if err := ssg2.GetSalt(salt2); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(salt1, salt2) {
		t.Error("salts should be random")
	}

	if !ssg1.IsServerSalt(salt2) || !ssg2.IsServerSalt(salt1) {
		t.Error("Cross-recognition failed")
	}
}

// Test that two ServerSaltGenerators derived from different secrets
// do not recognize each other's output.
func TestServerSaltDifferentCiphers(t *testing.T) {
	ssg1 := NewServerSaltGenerator("test1")
	ssg2 := NewServerSaltGenerator("test2")

	salt1 := make([]byte, 32)
	if err := ssg1.GetSalt(salt1); err != nil {
		t.Fatal(err)
	}
	salt2 := make([]byte, 32)
	if err := ssg2.GetSalt(salt2); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(salt1, salt2) {
		t.Error("salts should be random")
	}

	if ssg1.IsServerSalt(salt2) || ssg2.IsServerSalt(salt1) {
		t.Error("Different ciphers should not recognize each other")
	}
}

func TestServerSaltShort(t *testing.T) {
	ssg := NewServerSaltGenerator("test")

	salt20 := make([]byte, 20)
	if err := ssg.GetSalt(salt20); err != nil {
		t.Fatal(err)
	}
	if !ssg.IsServerSalt(salt20) {
		t.Error("Server salt was not recognized")
	}

	salt19 := make([]byte, 19)
	if err := ssg.GetSalt(salt19); err != nil {
		t.Fatal(err)
	}
	if ssg.IsServerSalt(salt19) {
		t.Error("Short salt was marked")
	}

	salt2 := make([]byte, 2)
	if err := ssg.GetSalt(salt2); err != nil {
		t.Fatal(err)
	}
	if ssg.IsServerSalt(salt2) {
		t.Error("Very short salt was marked")
	}
}

func BenchmarkRandomSaltGenerator(b *testing.B) {
	salt := make([]byte, 32)
	for i := 0; i < b.N; i++ {
		if err := RandomSaltGenerator.GetSalt(salt); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkServerSaltGenerator(b *testing.B) {
	ssg := NewServerSaltGenerator("test")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		salt := make([]byte, 32)
		for pb.Next() {
			if err := ssg.GetSalt(salt); err != nil {
				b.Fatal(err)
			}
			if !ssg.IsServerSalt(salt) {
				b.Fatal("Failed to recognize salt")
			}
		}
	})
}
