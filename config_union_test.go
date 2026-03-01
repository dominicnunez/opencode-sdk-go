package opencode

import (
	"encoding/json"
	"errors"
	"testing"
)

// TestConfigAgentBuildPermissionBashUnion tests the ConfigAgentBuildPermissionBashUnion discriminated union
func TestConfigAgentBuildPermissionBashUnion(t *testing.T) {
	t.Run("AsString_ValidString", func(t *testing.T) {
		data := []byte(`"allow"`)
		var u ConfigAgentBuildPermissionBashUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		s, err := u.AsString()
		if err != nil {
			t.Fatal("AsString() returned false for string value")
		}
		if s != ConfigAgentBuildPermissionBashStringAllow {
			t.Errorf("Expected 'allow', got %q", s)
		}
	})

	t.Run("AsMap_ValidMap", func(t *testing.T) {
		data := []byte(`{"cmd1":"allow","cmd2":"deny"}`)
		var u ConfigAgentBuildPermissionBashUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		m, err := u.AsMap()
		if err != nil {
			t.Fatal("AsMap() returned false for map value")
		}
		if len(m) != 2 {
			t.Errorf("Expected map length 2, got %d", len(m))
		}
		if m["cmd1"] != ConfigAgentBuildPermissionBashMapAllow {
			t.Errorf("Expected cmd1='allow', got %q", m["cmd1"])
		}
		if m["cmd2"] != ConfigAgentBuildPermissionBashMapDeny {
			t.Errorf("Expected cmd2='deny', got %q", m["cmd2"])
		}
	})

	t.Run("AsString_WhenMap", func(t *testing.T) {
		data := []byte(`{"cmd":"allow"}`)
		var u ConfigAgentBuildPermissionBashUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		_, err := u.AsString()
		if !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsString() = %v, want ErrWrongVariant", err)
		}
	})

	t.Run("AsMap_WhenString", func(t *testing.T) {
		data := []byte(`"deny"`)
		var u ConfigAgentBuildPermissionBashUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		_, err := u.AsMap()
		if !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsMap() = %v, want ErrWrongVariant", err)
		}
	})
}

// TestConfigProviderOptionsTimeoutUnion tests the ConfigProviderOptionsTimeoutUnion discriminated union
func TestConfigProviderOptionsTimeoutUnion(t *testing.T) {
	t.Run("AsInt_ValidInt", func(t *testing.T) {
		data := []byte(`300000`)
		var u ConfigProviderOptionsTimeoutUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		i, err := u.AsInt()
		if err != nil {
			t.Fatal("AsInt() returned false for int value")
		}
		if i != 300000 {
			t.Errorf("Expected 300000, got %d", i)
		}
	})

	t.Run("AsBool_ValidBool_False", func(t *testing.T) {
		data := []byte(`false`)
		var u ConfigProviderOptionsTimeoutUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		b, err := u.AsBool()
		if err != nil {
			t.Fatal("AsBool() returned false for bool value")
		}
		if b != false {
			t.Errorf("Expected false, got %v", b)
		}
	})

	t.Run("AsBool_ValidBool_True", func(t *testing.T) {
		data := []byte(`true`)
		var u ConfigProviderOptionsTimeoutUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		b, err := u.AsBool()
		if err != nil {
			t.Fatal("AsBool() returned false for bool value")
		}
		if b != true {
			t.Errorf("Expected true, got %v", b)
		}
	})

	t.Run("AsInt_WhenBool", func(t *testing.T) {
		data := []byte(`false`)
		var u ConfigProviderOptionsTimeoutUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		_, err := u.AsInt()
		if !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsInt() = %v, want ErrWrongVariant", err)
		}
	})

	t.Run("AsBool_WhenInt", func(t *testing.T) {
		data := []byte(`60000`)
		var u ConfigProviderOptionsTimeoutUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		_, err := u.AsBool()
		if !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsBool() = %v, want ErrWrongVariant", err)
		}
	})

	t.Run("AsInt_RejectsFloat", func(t *testing.T) {
		data := []byte(`300000.5`)
		var u ConfigProviderOptionsTimeoutUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		_, err := u.AsInt()
		if !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsInt(300000.5) = %v, want ErrWrongVariant", err)
		}
	})

	t.Run("AsInt_RejectsScientificNotation", func(t *testing.T) {
		data := []byte(`3e5`)
		var u ConfigProviderOptionsTimeoutUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		_, err := u.AsInt()
		if !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsInt(3e5) = %v, want ErrWrongVariant", err)
		}
	})

	t.Run("AsInt_NegativeInt", func(t *testing.T) {
		data := []byte(`-1`)
		var u ConfigProviderOptionsTimeoutUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		i, err := u.AsInt()
		if err != nil {
			t.Fatalf("AsInt(-1) unexpected error: %v", err)
		}
		if i != -1 {
			t.Errorf("Expected -1, got %d", i)
		}
	})
}

// TestConfigAgentGeneralPermissionBashUnion tests another bash union to ensure pattern consistency
func TestConfigAgentGeneralPermissionBashUnion(t *testing.T) {
	t.Run("AsString_ValidString", func(t *testing.T) {
		data := []byte(`"ask"`)
		var u ConfigAgentGeneralPermissionBashUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		s, err := u.AsString()
		if err != nil {
			t.Fatal("AsString() returned false for string value")
		}
		if s != ConfigAgentGeneralPermissionBashStringAsk {
			t.Errorf("Expected 'ask', got %q", s)
		}
	})

	t.Run("AsMap_ValidMap", func(t *testing.T) {
		data := []byte(`{"test":"ask"}`)
		var u ConfigAgentGeneralPermissionBashUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		m, err := u.AsMap()
		if err != nil {
			t.Fatal("AsMap() returned false for map value")
		}
		if m["test"] != ConfigAgentGeneralPermissionBashMapAsk {
			t.Errorf("Expected test='ask', got %q", m["test"])
		}
	})
}

// TestConfigPermissionBashUnion tests the top-level permission bash union
func TestConfigPermissionBashUnion(t *testing.T) {
	t.Run("AsString_ValidString", func(t *testing.T) {
		data := []byte(`"deny"`)
		var u ConfigPermissionBashUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		s, err := u.AsString()
		if err != nil {
			t.Fatal("AsString() returned false for string value")
		}
		if s != ConfigPermissionBashStringDeny {
			t.Errorf("Expected 'deny', got %q", s)
		}
	})

	t.Run("AsMap_EmptyMap", func(t *testing.T) {
		data := []byte(`{}`)
		var u ConfigPermissionBashUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		m, err := u.AsMap()
		if err != nil {
			t.Fatal("AsMap() returned false for empty map")
		}
		if len(m) != 0 {
			t.Errorf("Expected empty map, got length %d", len(m))
		}
	})
}

// TestConfigAgentPlanPermissionBashUnion tests plan bash union
func TestConfigAgentPlanPermissionBashUnion(t *testing.T) {
	t.Run("InvalidJSON", func(t *testing.T) {
		data := []byte(`{invalid}`)
		var u ConfigAgentPlanPermissionBashUnion
		err := json.Unmarshal(data, &u)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})

	t.Run("AsString_Null", func(t *testing.T) {
		data := []byte(`null`)
		var u ConfigAgentPlanPermissionBashUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		_, err := u.AsString()
		if !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsString() = %v, want ErrWrongVariant for null", err)
		}
	})
}

// TestConfigModeBuildPermissionBashUnion tests mode build bash union
func TestConfigModeBuildPermissionBashUnion(t *testing.T) {
	t.Run("AsMap_ComplexMap", func(t *testing.T) {
		data := []byte(`{"cmd1":"allow","cmd2":"ask","cmd3":"deny"}`)
		var u ConfigModeBuildPermissionBashUnion
		if err := json.Unmarshal(data, &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		m, err := u.AsMap()
		if err != nil {
			t.Fatal("AsMap() returned false for map value")
		}
		if len(m) != 3 {
			t.Errorf("Expected map length 3, got %d", len(m))
		}
	})
}

// TestConfigModePlanPermissionBashUnion tests mode plan bash union
func TestConfigModePlanPermissionBashUnion(t *testing.T) {
	t.Run("RoundTrip_String", func(t *testing.T) {
		original := `"allow"`
		var u ConfigModePlanPermissionBashUnion
		if err := json.Unmarshal([]byte(original), &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		s, err := u.AsString()
		if err != nil {
			t.Fatal("AsString() returned false")
		}
		if s != ConfigModePlanPermissionBashStringAllow {
			t.Errorf("Expected 'allow', got %q", s)
		}
	})

	t.Run("RoundTrip_Map", func(t *testing.T) {
		original := `{"test":"deny"}`
		var u ConfigModePlanPermissionBashUnion
		if err := json.Unmarshal([]byte(original), &u); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		m, err := u.AsMap()
		if err != nil {
			t.Fatal("AsMap() returned false")
		}
		if m["test"] != ConfigModePlanPermissionBashMapDeny {
			t.Errorf("Expected test='deny', got %q", m["test"])
		}
	})
}

func TestZeroValueUnion_AsReturnsErrWrongVariant(t *testing.T) {
	t.Run("ConfigAgentBuildPermissionBashUnion", func(t *testing.T) {
		var u ConfigAgentBuildPermissionBashUnion
		if _, err := u.AsString(); !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsString() = %v, want ErrWrongVariant", err)
		}
		if _, err := u.AsMap(); !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsMap() = %v, want ErrWrongVariant", err)
		}
	})

	t.Run("ConfigProviderOptionsTimeoutUnion", func(t *testing.T) {
		var u ConfigProviderOptionsTimeoutUnion
		if _, err := u.AsInt(); !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsInt() = %v, want ErrWrongVariant", err)
		}
		if _, err := u.AsBool(); !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsBool() = %v, want ErrWrongVariant", err)
		}
	})

	t.Run("PermissionPattern", func(t *testing.T) {
		var p PermissionPattern
		if _, err := p.AsString(); !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsString() = %v, want ErrWrongVariant", err)
		}
		if _, err := p.AsArray(); !errors.Is(err, ErrWrongVariant) {
			t.Errorf("AsArray() = %v, want ErrWrongVariant", err)
		}
	})
}
