package core

import (
  "github.com/dev-op-spec/engine/core/ports"
)

type compositionRoot interface {
  AddDevOpUseCase() addDevOpUseCase
  AddPipelineUseCase() addPipelineUseCase
  AddStageToPipelineUseCase() addStageToPipelineUseCase
  ListDevOpsUseCase() listDevOpsUseCase
  ListPipelinesUseCase() listPipelinesUseCase
  RunDevOpUseCase() runDevOpUseCase
  RunPipelineUseCase() runPipelineUseCase
  SetDescriptionOfDevOpUseCase() setDescriptionOfDevOpUseCase
  SetDescriptionOfPipelineUseCase() setDescriptionOfPipelineUseCase
}

func newCompositionRoot(
containerEngine ports.ContainerEngine,
filesys ports.Filesys,
) (compositionRoot compositionRoot, err error) {

  yml := _yamlCodec{}

  runDevOpUseCase := newRunDevOpUseCase(containerEngine)

  compositionRoot = &_compositionRoot{
    addDevOpUseCase: newAddDevOpUseCase(filesys, yml, containerEngine),
    addPipelineUseCase: newAddPipelineUseCase(filesys, yml),
    addStageToPipelineUseCase: newAddStageToPipelineUseCase(filesys, yml),
    listDevOpsUseCase: newListDevOpsUseCase(filesys, yml),
    listPipelinesUseCase: newListPipelinesUseCase(filesys, yml),
    runDevOpUseCase: runDevOpUseCase,
    runPipelineUseCase: newRunPipelineUseCase(filesys, yml, runDevOpUseCase),
    setDescriptionOfDevOpUseCase: newSetDescriptionOfDevOpUseCase(filesys, yml),
    setDescriptionOfPipelineUseCase: newSetDescriptionOfPipelineUseCase(filesys, yml),
  }

  return

}

type _compositionRoot struct {
  addDevOpUseCase                 addDevOpUseCase
  addPipelineUseCase              addPipelineUseCase
  addStageToPipelineUseCase       addStageToPipelineUseCase
  listDevOpsUseCase               listDevOpsUseCase
  listPipelinesUseCase            listPipelinesUseCase
  runDevOpUseCase                 runDevOpUseCase
  runPipelineUseCase              runPipelineUseCase
  setDescriptionOfDevOpUseCase    setDescriptionOfDevOpUseCase
  setDescriptionOfPipelineUseCase setDescriptionOfPipelineUseCase
}

func (this _compositionRoot) AddDevOpUseCase() addDevOpUseCase {
  return this.addDevOpUseCase
}

func (this _compositionRoot) AddPipelineUseCase() addPipelineUseCase {
  return this.addPipelineUseCase
}

func (this _compositionRoot) AddStageToPipelineUseCase() addStageToPipelineUseCase {
  return this.addStageToPipelineUseCase
}

func (this _compositionRoot) ListDevOpsUseCase() listDevOpsUseCase {
  return this.listDevOpsUseCase
}

func (this _compositionRoot) ListPipelinesUseCase() listPipelinesUseCase {
  return this.listPipelinesUseCase
}

func (this _compositionRoot) RunDevOpUseCase() runDevOpUseCase {
  return this.runDevOpUseCase
}

func (this _compositionRoot) RunPipelineUseCase() runPipelineUseCase {
  return this.runPipelineUseCase
}

func (this _compositionRoot) SetDescriptionOfDevOpUseCase() setDescriptionOfDevOpUseCase {
  return this.setDescriptionOfDevOpUseCase
}

func (this _compositionRoot) SetDescriptionOfPipelineUseCase() setDescriptionOfPipelineUseCase {
  return this.setDescriptionOfPipelineUseCase
}
