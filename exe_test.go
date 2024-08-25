package main

import (
	"testing"
)

func TestExecutionGroupUse(t *testing.T) {
	t.Setenv("DOUBAO_ENDPOINT", "env-endpoint")

	got := (ExecutionGroup{}).Use("default prompt")

	if got.driver != DriverDoubao {
		t.Errorf("ExecutionGroup.Use(empty).driver = %q, want %q", got.driver, DriverDoubao)
	}
	if got.ak != "" {
		t.Errorf("ExecutionGroup.Use(empty).ak = %q, want empty", got.ak)
	}
	if got.sk != "" {
		t.Errorf("ExecutionGroup.Use(empty).sk = %q, want empty", got.sk)
	}
	if got.ep != "env-endpoint" {
		t.Errorf("ExecutionGroup.Use(empty).ep = %q, want %q", got.ep, "env-endpoint")
	}
	if got.pp != "default prompt" {
		t.Errorf("ExecutionGroup.Use(empty).pp = %q, want %q", got.pp, "default prompt")
	}
}

func TestExecutionGroupUseKeepsExplicitValues(t *testing.T) {
	t.Setenv("DOUBAO_ENDPOINT", "env-endpoint")

	in := ExecutionGroup{
		driver: DriverOllama,
		ak:     "flag-ak",
		sk:     "flag-sk",
		ep:     "flag-endpoint",
		pp:     "flag-prompt",
	}
	got := in.Use("default prompt")

	if got.driver != DriverOllama {
		t.Errorf("ExecutionGroup.Use(%+v).driver = %q, want %q", in, got.driver, DriverOllama)
	}
	if got.ak != "flag-ak" {
		t.Errorf("ExecutionGroup.Use(%+v).ak = %q, want %q", in, got.ak, "flag-ak")
	}
	if got.sk != "flag-sk" {
		t.Errorf("ExecutionGroup.Use(%+v).sk = %q, want %q", in, got.sk, "flag-sk")
	}
	if got.ep != "flag-endpoint" {
		t.Errorf("ExecutionGroup.Use(%+v).ep = %q, want %q", in, got.ep, "flag-endpoint")
	}
	if got.pp != "flag-prompt" {
		t.Errorf("ExecutionGroup.Use(%+v).pp = %q, want %q", in, got.pp, "flag-prompt")
	}
}

func TestExecutionGroupAssert(t *testing.T) {
	tests := []struct {
		name    string
		group   ExecutionGroup
		wantErr bool
	}{
		{
			name: "valid ollama",
			group: ExecutionGroup{
				driver: DriverOllama,
				ep:     "llama3.1",
			},
		},
		{
			name: "ollama requires endpoint",
			group: ExecutionGroup{
				driver: DriverOllama,
			},
			wantErr: true,
		},
		{
			name: "valid doubao",
			group: ExecutionGroup{
				driver: DriverDoubao,
				ak:     "alpha",
				sk:     "bravo",
				ep:     "endpoint",
			},
		},
		{
			name: "doubao requires credentials",
			group: ExecutionGroup{
				driver: DriverDoubao,
				ep:     "endpoint",
			},
			wantErr: true,
		},
		{
			name: "invalid driver",
			group: ExecutionGroup{
				driver: "unsupported",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.group.Assert()
			if gotErr := err != nil; gotErr != tt.wantErr {
				t.Errorf("ExecutionGroup.Assert(%+v) error = %v, want error presence = %t", tt.group, err, tt.wantErr)
			}
		})
	}
}
