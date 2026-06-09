package blitzortung

import "testing"

func TestParseStrike_Valid(t *testing.T) {
	payload := []byte(`{"time":1717000000000000000,"lat":51.5,"lon":-0.1,"alt":0.0,"pol":1,"mds":123456,"scs":7}`)
	s, err := ParseStrike(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Time != 1717000000000000000 {
		t.Errorf("Time: got %d, want 1717000000000000000", s.Time)
	}
	if s.Lat != 51.5 {
		t.Errorf("Lat: got %f, want 51.5", s.Lat)
	}
	if s.Lon != -0.1 {
		t.Errorf("Lon: got %f, want -0.1", s.Lon)
	}
	if s.Pol != 1 {
		t.Errorf("Pol: got %d, want 1", s.Pol)
	}
	if s.Mds != 123456 {
		t.Errorf("Mds: got %d, want 123456", s.Mds)
	}
	if s.Scs != 7 {
		t.Errorf("Scs: got %d, want 7", s.Scs)
	}
}

func TestParseStrike_MissingFields(t *testing.T) {
	payload := []byte(`{}`)
	s, err := ParseStrike(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Time != 0 || s.Lat != 0 || s.Lon != 0 {
		t.Errorf("expected zero values for missing fields, got %+v", s)
	}
}

func TestParseStrike_MalformedJSON(t *testing.T) {
	payload := []byte(`not json`)
	_, err := ParseStrike(payload)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}
