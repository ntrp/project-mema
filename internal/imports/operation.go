package imports

import (
	"fmt"
	"io"
	"os"
)

type ImportMode string

const (
	ImportModeHardlink ImportMode = "hardlink"
	ImportModeCopy     ImportMode = "copy"
	ImportModeMove     ImportMode = "move"
)

func normalizeImportMode(mode ImportMode) (ImportMode, error) {
	if mode == "" {
		return ImportModeHardlink, nil
	}
	switch mode {
	case ImportModeHardlink, ImportModeCopy, ImportModeMove:
		return mode, nil
	default:
		return "", fmt.Errorf("unsupported import mode: %s", mode)
	}
}

func importFile(source string, target string, mode ImportMode) error {
	mode, err := normalizeImportMode(mode)
	if err != nil {
		return err
	}
	if sameExistingFile(source, target) {
		return nil
	}
	if err := prepareTarget(target); err != nil {
		return err
	}
	switch mode {
	case ImportModeHardlink:
		return hardlinkFile(source, target)
	case ImportModeCopy:
		return copyFile(source, target)
	case ImportModeMove:
		return moveFile(source, target)
	default:
		return fmt.Errorf("unsupported import mode: %s", mode)
	}
}

func prepareTarget(target string) error {
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("target already exists: %s", target)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check target: %w", err)
	}
	probe := target + ".import-probe"
	file, err := os.OpenFile(probe, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("prepare target directory: %w", err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(probe)
		return fmt.Errorf("prepare target directory: %w", err)
	}
	if err := os.Remove(probe); err != nil {
		return fmt.Errorf("prepare target directory: %w", err)
	}
	return nil
}

func hardlinkFile(source string, target string) error {
	if err := os.Link(source, target); err != nil {
		return fmt.Errorf("hardlink %s to %s: %w", source, target, err)
	}
	return nil
}

func copyFile(source string, target string) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer sourceFile.Close()

	targetFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, sourceInfo.Mode().Perm())
	if err != nil {
		return fmt.Errorf("create target: %w", err)
	}
	copied := false
	defer func() {
		_ = targetFile.Close()
		if !copied {
			_ = os.Remove(target)
		}
	}()
	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("copy %s to %s: %w", source, target, err)
	}
	copied = true
	return nil
}

func moveFile(source string, target string) error {
	if err := os.Rename(source, target); err != nil {
		return fmt.Errorf("move %s to %s: %w", source, target, err)
	}
	return nil
}

func sameExistingFile(source string, target string) bool {
	sourceInfo, sourceErr := os.Stat(source)
	targetInfo, targetErr := os.Stat(target)
	return sourceErr == nil && targetErr == nil && os.SameFile(sourceInfo, targetInfo)
}
