package bundle_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMacBundleAssetsExist(t *testing.T) {
	required := []string{
		filepath.FromSlash("macos/GoPomodoro.app/Contents/Info.plist"),
		filepath.FromSlash("macos/GoPomodoro.app/Contents/Resources/AppIcon.icns"),
	}

	for _, path := range required {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Required asset not found: %s", path)
		}
	}

}
