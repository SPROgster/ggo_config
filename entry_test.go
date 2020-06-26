package ggo

import "testing"

func TestParseString(t *testing.T) {
	got := ParseString("")
	if got != nil {
		t.Errorf("Empty string test error\n")
	}

	got = ParseString("#")
	if got != nil {
		t.Errorf("'#' test string error\n")
	}

	got = ParseString("# test 1.2.3.4 # comment")
	if got == nil || got.IsActive || got.Name() != "test" || got.Value != "1.2.3.4" || got.Comment != "comment" {
		t.Errorf("'%s' parse error %v", "# test 1.2.3.4 # comment\n", got)
	}

	got = ParseString("\t  # test   \t     1.2.3.4 #   \t  a comment")
	if got == nil || got.IsActive || got.Name() != "test" || got.Value != "1.2.3.4" || got.Comment != "a comment" {
		t.Errorf("'%s' parse error %v", "# test   \t     1.2.3.4 #   \t  a comment\n", got)
	}

	got = ParseString("## test 1.2.3.4 # comment")
	if got == nil || got.IsActive || got.Name() != "test" || got.Value != "1.2.3.4" || got.Comment != "comment" {
		t.Errorf("'%s' parse error `%s`\n", "## test 1.2.3.4 # comment", got.String())
	}

	got = ParseString("# # test 1.2.3.4 # comment")
	if got == nil || got.IsActive || got.Name() != "test" || got.Value != "1.2.3.4" || got.Comment != "comment" {
		t.Errorf("'%s' parse error %v\n", "# # test 1.2.3.4 # comment", got)
	}

	got = ParseString("#switch off cookie filter")
	if got != nil {
		t.Errorf("'%s' parse error: %v\n", "#switch off cookie filter", got)
	}

	got = ParseString("## TCP")
	if got == nil || got.IsActive || got.Name() != "TCP" || got.Value != "" || got.Comment != "" {
		t.Errorf("'%s' parse error %v\n", "## TCP", got)
	}

}

func TestConfigEntry_String(t *testing.T) {
	e := ParseString("# test 1.2.3.4 # comment")
	got := e.String()
	if got != "# test 1.2.3.4 # comment" {
		t.Errorf("'%s' parse error %v", "# test 1.2.3.4 # comment\n", got)
	}

	e = ParseString("\t  # test   \t     1.2.3.4 #   \t  a comment")
	got = e.String()
	if got != "# test 1.2.3.4 # a comment" {
		t.Errorf("'%s' parse error %v", "# test   \t     1.2.3.4 #   \t  a comment\n", got)
	}
}


