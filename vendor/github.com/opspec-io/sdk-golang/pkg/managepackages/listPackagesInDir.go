package managepackages

//go:generate counterfeiter -o ./fakeListPackagesInDir.go --fake-name fakeListPackagesInDir ./ listPackagesInDir

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

func (this managePackages) ListPackagesInDir(
	dirPath string,
) (
	ops []*model.PackageView,
	err error,
) {

	childFileInfos, err := this.fileSystem.ListChildFileInfosOfDir(dirPath)
	if nil != err {
		return
	}

	for _, childFileInfo := range childFileInfos {
		packageView, err := this.packageViewFactory.Construct(
			path.Join(dirPath, childFileInfo.Name()),
		)
		if nil == err {
			ops = append(ops, &packageView)
		}
	}

	return

}
