package usecase

import (
	"errors"
	"net/netip"
	"testing"
	"time"

	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
)

type stubWidgetParser struct {
	authData *service.TelegramAuthData
	err      error
}

func (s *stubWidgetParser) Parse(params map[string]any) (*service.TelegramAuthData, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.authData, nil
}

type stubHashVerifier struct {
	err error
}

func (s *stubHashVerifier) Verify(query string, hash string, botToken string) error {
	return s.err
}

func TestLoginByWidget_verifyChallenge(t *testing.T) {
	uc := &LoginByWidget{}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "empty challenge", input: "", wantErr: true},
		{name: "valid challenge", input: "challenge-123", wantErr: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := uc.verifyChallenge(testCase.input)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, ErrInvalidInput) {
					t.Fatalf("expected ErrInvalidInput, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestLoginByWidget_verifyIP(t *testing.T) {
	uc := &LoginByWidget{}

	tests := []struct {
		name    string
		input   netip.Addr
		wantErr bool
	}{
		{name: "unspecified ip", input: netip.MustParseAddr("0.0.0.0"), wantErr: true},
		{name: "valid ipv4", input: netip.MustParseAddr("127.0.0.1"), wantErr: false},
		{name: "valid ipv6", input: netip.MustParseAddr("::1"), wantErr: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := uc.verifyIP(testCase.input)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, ErrInvalidInput) {
					t.Fatalf("expected ErrInvalidInput, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestLoginByWidget_parseAndVerifyAuthData(t *testing.T) {
	now := time.Now()
	validAuthData := &service.TelegramAuthData{
		Raw:      "auth_date=123&id=1",
		Hash:     "hash",
		AuthDate: now,
		User: &service.TelegramUserData{
			Id:        1,
			FirstName: "John",
		},
	}

	tests := []struct {
		name      string
		params    map[string]any
		parser    *stubWidgetParser
		verifier  *stubHashVerifier
		freshness time.Duration
		wantErr   bool
	}{
		{
			name:      "empty params",
			params:    map[string]any{},
			parser:    &stubWidgetParser{authData: validAuthData},
			verifier:  &stubHashVerifier{},
			freshness: 5 * time.Minute,
			wantErr:   true,
		},
		{
			name:      "parser error",
			params:    map[string]any{"id": 1},
			parser:    &stubWidgetParser{err: errors.New("parse error")},
			verifier:  &stubHashVerifier{},
			freshness: 5 * time.Minute,
			wantErr:   true,
		},
		{
			name:      "expired auth data",
			params:    map[string]any{"id": 1},
			parser:    &stubWidgetParser{authData: &service.TelegramAuthData{Raw: "r", Hash: "h", AuthDate: now.Add(-10 * time.Minute), User: validAuthData.User}},
			verifier:  &stubHashVerifier{},
			freshness: 5 * time.Minute,
			wantErr:   true,
		},
		{
			name:      "invalid hash",
			params:    map[string]any{"id": 1},
			parser:    &stubWidgetParser{authData: validAuthData},
			verifier:  &stubHashVerifier{err: errors.New("bad hash")},
			freshness: 5 * time.Minute,
			wantErr:   true,
		},
		{
			name:      "success",
			params:    map[string]any{"id": 1},
			parser:    &stubWidgetParser{authData: validAuthData},
			verifier:  &stubHashVerifier{},
			freshness: 5 * time.Minute,
			wantErr:   false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			uc := &LoginByWidget{
				widgetDataParser:  testCase.parser,
				authHashVerifier:  testCase.verifier,
				authDataFreshness: testCase.freshness,
			}

			_, err := uc.parseAndVerifyAuthData(testCase.params, "bot-token")
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, ErrInvalidInput) {
					t.Fatalf("expected ErrInvalidInput, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}
