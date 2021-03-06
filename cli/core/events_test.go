package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/pkg/nodeprovider"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/consumenodeapi"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Context("streamEvents", func() {
	Context("Execute", func() {
		It("should call managepackages.GetEventStream", func() {
			/* arrange */
			fakeCliExiter := new(cliexiter.Fake)

			fakeConsumeNodeApi := new(consumenodeapi.Fake)
			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _core{
				consumeNodeApi: fakeConsumeNodeApi,
				cliExiter:      fakeCliExiter,
				nodeProvider:   new(nodeprovider.Fake),
			}

			/* act */
			objectUnderTest.StreamEvents()

			/* assert */
			Expect(fakeConsumeNodeApi.GetEventStreamCallCount()).Should(Equal(1))

		})
		Context("managepackages.GetEventStream errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiter.Fake)
				returnedError := errors.New("dummyError")

				fakeConsumeNodeApi := new(consumenodeapi.Fake)
				fakeConsumeNodeApi.GetEventStreamReturns(nil, returnedError)

				objectUnderTest := _core{
					consumeNodeApi: fakeConsumeNodeApi,
					cliExiter:      fakeCliExiter,
					nodeProvider:   new(nodeprovider.Fake),
				}

				/* act */
				objectUnderTest.StreamEvents()

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
			})
		})
		Context("managepackages.GetEventStream doesn't error", func() {
			Context("channel closes unexpectedly", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeCliExiter := new(cliexiter.Fake)

					fakeConsumeNodeApi := new(consumenodeapi.Fake)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _core{
						consumeNodeApi: fakeConsumeNodeApi,
						cliExiter:      fakeCliExiter,
						nodeProvider:   new(nodeprovider.Fake),
					}

					/* act */
					objectUnderTest.StreamEvents()

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: "Connection to event stream lost", Code: 1}))
				})
			})
		})
	})
})
