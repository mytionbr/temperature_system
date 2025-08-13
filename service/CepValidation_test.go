package service

import "testing"

func TestCepValitation_Success(t *testing.T) {
	got, err := CepValitation("01001-000")
	if err != nil {
		t.Fatalf("expected valid CEP, error %v", err)
	}
	if got != "01001000" {
		t.Fatalf("wanted 01001000, got %s", got)
	}
}

func TestCepValitation_WrongSize(t *testing.T) {
	_, err := CepValitation("1234567")
	if err == nil {
		t.Fatal("Expected error due to invalid length")
	}
}

func TestCepValitation_NonNumeric(t *testing.T) {
	_, err := CepValitation("ABCDEFGH")
	if err == nil {
		t.Fatal("Expected error due to non-numeric characters")
	}
}
