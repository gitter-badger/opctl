package filecopier

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ FileCopier

import (
	"github.com/opspec-io/opctl/util/vfs"
	"github.com/opspec-io/opctl/util/vfs/os"
	"io"
)

type FileCopier interface {
	// copies a fs file from srcPath to dstPath. Creates or overwrites the destination as needed.
	Fs(srcPath string, dstPath string) (err error)
}

func New() FileCopier {
	return fileCopier{
		fs: os.New(),
	}
}

type fileCopier struct {
	fs vfs.Vfs
}

func (this fileCopier) Fs(srcPath string, dstPath string) (err error) {
	srcFile, err := this.fs.Open(srcPath)
	if nil != err {
		return
	}
	defer srcFile.Close()

	srcFileInfo, err := this.fs.Stat(srcPath)
	if nil != err {
		return
	}

	// copy content
	writer, err := this.fs.Create(dstPath)
	if nil != err {
		return
	}
	defer writer.Close()

	// copy mode
	err = this.fs.Chmod(dstPath, srcFileInfo.Mode())
	if nil != err {
		return
	}

	_, err = io.Copy(writer, srcFile)

	err = writer.Sync()

	return
}
