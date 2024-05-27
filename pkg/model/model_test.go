package model

import "testing"

func TestGetLang(t *testing.T) {
	t.Run("GetLang Initialize", func(t *testing.T) {
		lang := GetLang("en", "English", nil)
		if lang.Code != "en" {
			t.Errorf("GetLang() = %v; want en", lang)
		}
		if lang.Name != "English" {
			t.Errorf("GetLang() = %v; want English", lang)
		}
		if lang.Meta != nil {
			t.Errorf("GetLang() = %v; want nil", lang)
		}
	})
}
