package core

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/pkg/nodeprovider"
	"github.com/opspec-io/opctl/util/clicolorer"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/clioutput"
	"github.com/opspec-io/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/consumenodeapi"
	"github.com/opspec-io/sdk-golang/pkg/managepackages"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
	"path/filepath"
	"time"
)

var _ = Context("runOp", func() {
	Context("Execute", func() {
		Context("vos.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeVos := new(vos.Fake)
				expectedError := errors.New("dummyError")
				fakeVos.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					managePackages: new(managepackages.Fake),
					cliExiter:      fakeCliExiter,
					nodeProvider:   new(nodeprovider.Fake),
					vos:            fakeVos,
				}

				/* act */
				objectUnderTest.RunOp([]string{}, "dummyName")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("vos.Getwd doesn't error", func() {
			It("should call managepackages.GetPackage w/ expected args", func() {
				/* arrange */
				fakeManagePackages := new(managepackages.Fake)

				fakeConsumeNodeApi := new(consumenodeapi.Fake)
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

				fakeCliExiter := new(cliexiter.Fake)

				providedPkgRef := "dummyPkgName"
				wdReturnedFromVos := "dummyWorkDir"

				fakeVos := new(vos.Fake)
				fakeVos.GetwdReturns(wdReturnedFromVos, nil)
				expectedPath := filepath.Join(wdReturnedFromVos, ".opspec", providedPkgRef)

				objectUnderTest := _core{
					managePackages:    fakeManagePackages,
					consumeNodeApi:    fakeConsumeNodeApi,
					cliExiter:         fakeCliExiter,
					cliParamSatisfier: new(cliparamsatisfier.Fake),
					nodeProvider:      new(nodeprovider.Fake),
					vos:               fakeVos,
				}

				/* act */
				objectUnderTest.RunOp([]string{}, providedPkgRef)

				/* assert */
				Expect(fakeManagePackages.GetPackageArgsForCall(0)).Should(Equal(expectedPath))
			})
			Context("managepackages.GetPackage errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeCliExiter := new(cliexiter.Fake)
					returnedError := errors.New("dummyError")

					fakeManagePackages := new(managepackages.Fake)
					fakeManagePackages.GetPackageReturns(model.PackageView{}, returnedError)

					objectUnderTest := _core{
						managePackages:    fakeManagePackages,
						cliExiter:         fakeCliExiter,
						cliParamSatisfier: new(cliparamsatisfier.Fake),
						nodeProvider:      new(nodeprovider.Fake),
						vos:               new(vos.Fake),
					}

					/* act */
					objectUnderTest.RunOp([]string{}, "dummyName")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
				})
			})
			Context("managepackages.GetPackage doesn't error", func() {
				It("should call paramSatisfier.Satisfy w/ expected args", func() {
					/* arrange */
					param1Name := "DUMMY_PARAM1_NAME"
					arg1Value := &model.Data{String: "dummyParam1Value"}

					providedArgs := []string{fmt.Sprintf("%v=%v", param1Name, arg1Value.String)}

					expectedParams := map[string]*model.Param{
						param1Name: {
							String: &model.StringParam{},
						},
					}

					fakeManagePackages := new(managepackages.Fake)
					fakeManagePackages.GetPackageReturns(
						model.PackageView{
							Inputs: expectedParams,
						},
						nil,
					)

					// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
					fakeConsumeNodeApi := new(consumenodeapi.Fake)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

					fakeCliParamSatisfier := new(cliparamsatisfier.Fake)

					objectUnderTest := _core{
						managePackages:    fakeManagePackages,
						consumeNodeApi:    fakeConsumeNodeApi,
						cliExiter:         new(cliexiter.Fake),
						cliParamSatisfier: fakeCliParamSatisfier,
						nodeProvider:      new(nodeprovider.Fake),
						vos:               new(vos.Fake),
					}

					/* act */
					objectUnderTest.RunOp(providedArgs, "dummyPkgName")

					/* assert */
					actualArgs, actualParams := fakeCliParamSatisfier.SatisfyArgsForCall(0)
					Expect(actualArgs).To(Equal(providedArgs))
					Expect(actualParams).To(Equal(expectedParams))
				})
				It("should call consumeNodeApi.StartOp w/ expected args", func() {
					/* arrange */
					pwd := "dummyWorkDir"
					fakeVos := new(vos.Fake)
					fakeVos.GetwdReturns(pwd, nil)

					providedPkgRef := "dummyOp"
					expectedArgs := model.StartOpReq{
						Args: map[string]*model.Data{
							"dummyArg1Name": {String: "dummyArg1Value"},
						},
						PkgRef: path.Join(pwd, ".opspec", providedPkgRef),
					}

					// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
					fakeConsumeNodeApi := new(consumenodeapi.Fake)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

					fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
					fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args)

					objectUnderTest := _core{
						managePackages:    new(managepackages.Fake),
						consumeNodeApi:    fakeConsumeNodeApi,
						cliExiter:         new(cliexiter.Fake),
						cliParamSatisfier: fakeCliParamSatisfier,
						nodeProvider:      new(nodeprovider.Fake),
						vos:               fakeVos,
					}

					/* act */
					objectUnderTest.RunOp([]string{}, providedPkgRef)

					/* assert */
					actualArgs := fakeConsumeNodeApi.StartOpArgsForCall(0)
					Expect(actualArgs).To(Equal(expectedArgs))
				})
				Context("consumeNodeApi.StartOp errors", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						fakeCliExiter := new(cliexiter.Fake)
						returnedError := errors.New("dummyError")

						fakeManagePackages := new(managepackages.Fake)
						fakeManagePackages.GetPackageReturns(model.PackageView{}, nil)

						fakeConsumeNodeApi := new(consumenodeapi.Fake)
						fakeConsumeNodeApi.StartOpReturns("dummyOpId", returnedError)

						objectUnderTest := _core{
							managePackages:    fakeManagePackages,
							consumeNodeApi:    fakeConsumeNodeApi,
							cliExiter:         fakeCliExiter,
							cliParamSatisfier: new(cliparamsatisfier.Fake),
							nodeProvider:      new(nodeprovider.Fake),
							vos:               new(vos.Fake),
						}

						/* act */
						objectUnderTest.RunOp([]string{}, "dummyPkgName")

						/* assert */
						Expect(fakeCliExiter.ExitArgsForCall(0)).
							Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
					})
				})
				Context("consumeNodeApi.StartOp doesn't error", func() {
					It("should call consumeNodeApi.GetEventStream w/ expected args", func() {
						/* arrange */
						fakeManagePackages := new(managepackages.Fake)
						fakeManagePackages.GetPackageReturns(model.PackageView{}, nil)
						rootOpIdReturnedFromStartOp := "dummyRootOpId"
						expectedEventFilter := &model.GetEventStreamReq{
							Filter: &model.EventFilter{
								RootOpIds: []string{rootOpIdReturnedFromStartOp},
							},
						}

						fakeConsumeNodeApi := new(consumenodeapi.Fake)
						fakeConsumeNodeApi.StartOpReturns(rootOpIdReturnedFromStartOp, nil)
						eventChannel := make(chan model.Event)
						close(eventChannel)
						fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

						objectUnderTest := _core{
							managePackages:    fakeManagePackages,
							consumeNodeApi:    fakeConsumeNodeApi,
							cliExiter:         new(cliexiter.Fake),
							cliParamSatisfier: new(cliparamsatisfier.Fake),
							nodeProvider:      new(nodeprovider.Fake),
							vos:               new(vos.Fake),
						}

						/* act */
						objectUnderTest.RunOp([]string{}, "dummyPkgName")

						/* assert */
						Expect(fakeConsumeNodeApi.GetEventStreamArgsForCall(0)).
							Should(Equal(expectedEventFilter))
					})
					Context("consumeNodeApi.GetEventStream errors", func() {
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeCliExiter := new(cliexiter.Fake)
							returnedError := errors.New("dummyError")

							fakeManagePackages := new(managepackages.Fake)
							fakeManagePackages.GetPackageReturns(model.PackageView{}, nil)

							fakeConsumeNodeApi := new(consumenodeapi.Fake)
							fakeConsumeNodeApi.GetEventStreamReturns(nil, returnedError)

							objectUnderTest := _core{
								managePackages:    fakeManagePackages,
								consumeNodeApi:    fakeConsumeNodeApi,
								cliExiter:         fakeCliExiter,
								cliParamSatisfier: new(cliparamsatisfier.Fake),
								nodeProvider:      new(nodeprovider.Fake),
								vos:               new(vos.Fake),
							}

							/* act */
							objectUnderTest.RunOp([]string{}, "dummyPkgName")

							/* assert */
							Expect(fakeCliExiter.ExitArgsForCall(0)).
								Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
						})
					})
					Context("consumeNodeApi.GetEventStream doesn't error", func() {
						Context("event channel closes", func() {
							It("should call exiter w/ expected args", func() {
								/* arrange */
								fakeCliExiter := new(cliexiter.Fake)

								fakeManagePackages := new(managepackages.Fake)
								fakeManagePackages.GetPackageReturns(model.PackageView{}, nil)

								fakeConsumeNodeApi := new(consumenodeapi.Fake)
								eventChannel := make(chan model.Event)
								close(eventChannel)
								fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

								objectUnderTest := _core{
									managePackages:    fakeManagePackages,
									consumeNodeApi:    fakeConsumeNodeApi,
									cliExiter:         fakeCliExiter,
									cliParamSatisfier: new(cliparamsatisfier.Fake),
									nodeProvider:      new(nodeprovider.Fake),
									vos:               new(vos.Fake),
								}

								/* act */
								objectUnderTest.RunOp([]string{}, "dummyPkgName")

								/* assert */
								Expect(fakeCliExiter.ExitArgsForCall(0)).
									Should(Equal(cliexiter.ExitReq{Message: "Event channel closed unexpectedly", Code: 1}))
							})
						})
						Context("event channel doesn't close", func() {
							Context("event received", func() {
								rootOpId := "dummyRootOpId"
								Context("OpEndedEvent", func() {
									Context("Outcome==SUCCEEDED", func() {
										It("should call exiter w/ expected args", func() {
											/* arrange */
											opEndedEvent := model.Event{
												Timestamp: time.Now(),
												OpEnded: &model.OpEndedEvent{
													OpId:     rootOpId,
													PkgRef:   "dummyPkgRef",
													Outcome:  model.OpOutcomeSucceeded,
													RootOpId: rootOpId,
												},
											}

											fakeCliExiter := new(cliexiter.Fake)

											fakeManagePackages := new(managepackages.Fake)
											fakeManagePackages.GetPackageReturns(model.PackageView{}, nil)

											fakeConsumeNodeApi := new(consumenodeapi.Fake)
											eventChannel := make(chan model.Event, 10)
											eventChannel <- opEndedEvent
											defer close(eventChannel)
											fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)
											fakeConsumeNodeApi.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

											objectUnderTest := _core{
												managePackages:    fakeManagePackages,
												cliColorer:        clicolorer.New(),
												consumeNodeApi:    fakeConsumeNodeApi,
												cliExiter:         fakeCliExiter,
												cliOutput:         new(clioutput.Fake),
												cliParamSatisfier: new(cliparamsatisfier.Fake),
												nodeProvider:      new(nodeprovider.Fake),
												vos:               new(vos.Fake),
											}

											/* act/assert */
											objectUnderTest.RunOp([]string{}, "dummyPkgName")
											Expect(fakeCliExiter.ExitArgsForCall(0)).
												Should(Equal(cliexiter.ExitReq{Code: 0}))
										})
									})
									Context("Outcome==KILLED", func() {
										It("should call exiter w/ expected args", func() {
											/* arrange */
											opEndedEvent := model.Event{
												Timestamp: time.Now(),
												OpEnded: &model.OpEndedEvent{
													OpId:     rootOpId,
													PkgRef:   "dummyPkgRef",
													Outcome:  model.OpOutcomeKilled,
													RootOpId: rootOpId,
												},
											}

											fakeCliExiter := new(cliexiter.Fake)

											fakeManagePackages := new(managepackages.Fake)
											fakeManagePackages.GetPackageReturns(model.PackageView{}, nil)

											fakeConsumeNodeApi := new(consumenodeapi.Fake)
											eventChannel := make(chan model.Event, 10)
											eventChannel <- opEndedEvent
											defer close(eventChannel)
											fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)
											fakeConsumeNodeApi.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

											objectUnderTest := _core{
												managePackages:    fakeManagePackages,
												cliColorer:        clicolorer.New(),
												consumeNodeApi:    fakeConsumeNodeApi,
												cliExiter:         fakeCliExiter,
												cliOutput:         new(clioutput.Fake),
												cliParamSatisfier: new(cliparamsatisfier.Fake),
												nodeProvider:      new(nodeprovider.Fake),
												vos:               new(vos.Fake),
											}

											/* act/assert */
											objectUnderTest.RunOp([]string{}, "dummyPkgName")
											Expect(fakeCliExiter.ExitArgsForCall(0)).
												Should(Equal(cliexiter.ExitReq{Code: 137}))
										})

									})
									Context("Outcome==FAILED", func() {
										It("should call exiter w/ expected args", func() {
											/* arrange */
											opEndedEvent := model.Event{
												Timestamp: time.Now(),
												OpEnded: &model.OpEndedEvent{
													OpId:     rootOpId,
													PkgRef:   "dummyPkgRef",
													Outcome:  model.OpOutcomeFailed,
													RootOpId: rootOpId,
												},
											}

											fakeCliExiter := new(cliexiter.Fake)

											fakeManagePackages := new(managepackages.Fake)
											fakeManagePackages.GetPackageReturns(model.PackageView{}, nil)

											fakeConsumeNodeApi := new(consumenodeapi.Fake)
											eventChannel := make(chan model.Event, 10)
											eventChannel <- opEndedEvent
											defer close(eventChannel)
											fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)
											fakeConsumeNodeApi.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

											objectUnderTest := _core{
												managePackages:    fakeManagePackages,
												cliColorer:        clicolorer.New(),
												consumeNodeApi:    fakeConsumeNodeApi,
												cliExiter:         fakeCliExiter,
												cliOutput:         new(clioutput.Fake),
												cliParamSatisfier: new(cliparamsatisfier.Fake),
												nodeProvider:      new(nodeprovider.Fake),
												vos:               new(vos.Fake),
											}

											/* act/assert */
											objectUnderTest.RunOp([]string{}, "dummyPkgName")
											Expect(fakeCliExiter.ExitArgsForCall(0)).
												Should(Equal(cliexiter.ExitReq{Code: 1}))
										})
									})
									Context("Outcome==?", func() {
										It("should call exiter w/ expected args", func() {
											/* arrange */
											opEndedEvent := model.Event{
												Timestamp: time.Now(),
												OpEnded: &model.OpEndedEvent{
													OpId:     rootOpId,
													PkgRef:   "dummyPkgRef",
													Outcome:  "some unexpected outcome",
													RootOpId: rootOpId,
												},
											}

											fakeCliExiter := new(cliexiter.Fake)

											fakeManagePackages := new(managepackages.Fake)
											fakeManagePackages.GetPackageReturns(model.PackageView{}, nil)

											fakeConsumeNodeApi := new(consumenodeapi.Fake)
											eventChannel := make(chan model.Event, 10)
											eventChannel <- opEndedEvent
											defer close(eventChannel)
											fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)
											fakeConsumeNodeApi.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

											objectUnderTest := _core{
												managePackages:    fakeManagePackages,
												cliColorer:        clicolorer.New(),
												consumeNodeApi:    fakeConsumeNodeApi,
												cliExiter:         fakeCliExiter,
												cliOutput:         new(clioutput.Fake),
												cliParamSatisfier: new(cliparamsatisfier.Fake),
												nodeProvider:      new(nodeprovider.Fake),
												vos:               new(vos.Fake),
											}

											/* act/assert */
											objectUnderTest.RunOp([]string{}, "dummyPkgName")
											Expect(fakeCliExiter.ExitArgsForCall(0)).
												Should(Equal(cliexiter.ExitReq{Code: 1}))
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})
})
