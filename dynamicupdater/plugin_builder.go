package dynamicupdater

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// CompilePlugin compiles the given source file into a plugin and places it in the plugin directory
func CompilePlugin(sourceFile string) error {
	outputFile := filepath.Join(pluginDir, filepath.Base(sourceFile)+".so")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outputFile, sourceFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to compile plugin: %w", err)
	}
	fmt.Printf("Successfully compiled and placed plugin: %s\n", outputFile)
	return nil
}
