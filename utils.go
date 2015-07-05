package lynd

import (
	"os"
	"os/exec"
	"path"
)

// renameMkdir ensures the destination directory exists before rename.
func renameMkdir(src, dst string) error {
	dir := path.Dir(dst)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return os.Rename(src, dst)
}

// rename copies a file from src to dst. It will create necessary directories
// and will fallback to cp and rename, if operation spans devices.
func rename(src, dst string) error {
	err := renameMkdir(src, dst)
	if err == nil {
		return nil
	}
	dstTmp := dst + "-tmp"
	cmd := exec.Command("cp", "-rf", src, dstTmp)
	err = cmd.Run()
	if err != nil {
		return err
	}
	err = os.Rename(dstTmp, dst)
	if err != nil {
		return err
	}
	return os.Remove(src)
}
