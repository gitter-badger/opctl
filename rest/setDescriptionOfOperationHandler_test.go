package rest

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "net/http/httptest"
  "net/http"
  "strings"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core/models"
  "bytes"
  "encoding/json"
)

var _ = Describe("setDescriptionOfOperationHandler", func() {

  Context("ServeHTTP() method", func() {
    It("should return StatusCode of 400 if projectUrl is malformed in Request", func() {

      /* arrange */
      objectUnderTest := setDescriptionOfOperationHandler{}
      recorder := httptest.NewRecorder()
      m := mux.NewRouter()
      m.Handle(setDescriptionOfOperationRelUrlTemplate, objectUnderTest)

      providedProjectUrl := "%%invalidProjectUrl%%"
      providedOperationName := "validOperationName"
      providedSetDescriptionOfOperationReqJson, err := json.Marshal(models.SetDescriptionOfOperationReq{})
      if (nil != err) {
        Fail(err.Error())
      }

      httpReq, err := http.NewRequest(http.MethodGet, "", bytes.NewReader(providedSetDescriptionOfOperationReqJson))
      if (nil != err) {
        Fail(err.Error())
      }

      // brute force a request with malformed projectUrl
      httpReq.URL.Path = strings.Replace(setDescriptionOfOperationRelUrlTemplate, "{projectUrl}", providedProjectUrl, 1)
      httpReq.URL.Path = strings.Replace(httpReq.URL.Path, "{operationName}", providedOperationName, 1)

      /* act */
      m.ServeHTTP(recorder, httpReq)

      /* assert */
      Expect(recorder.Code).To(Equal(http.StatusBadRequest))

    })
  })
})