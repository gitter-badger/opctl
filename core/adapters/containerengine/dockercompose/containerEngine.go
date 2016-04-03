package dockercompose

import (
  "github.com/dev-op-spec/engine/core/ports"
  "github.com/dev-op-spec/engine/core/models"
)

func New(
) (containerEngine ports.ContainerEngine, err error) {

  var compositionRoot compositionRoot
  compositionRoot, err = newCompositionRoot()
  if (nil != err) {
    return
  }

  containerEngine = _containerEngine{
    compositionRoot:compositionRoot,
  }

  return

}

type _containerEngine struct {
  compositionRoot compositionRoot
}

func (this _containerEngine) InitOperation(
pathToOperationDir string,
) (err error) {
  return this.compositionRoot.
  InitOperationUseCase().
  Execute(pathToOperationDir)
}

func (this _containerEngine) RunOperation(
pathToOperationDir string,
) (operationRun models.OperationRunDetailedView, err error) {
  return this.compositionRoot.
  RunOperationUseCase().
  Execute(pathToOperationDir)
}
