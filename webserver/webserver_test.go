package webserver_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"arcticwolf.com/cutler/models"
	mock_models "arcticwolf.com/cutler/models/mocks"
	"arcticwolf.com/cutler/webserver"
)

var (
	risk1 = models.Risk{
		ID:          "abcd",
		State:       models.StateAccepted,
		Title:       "MyTitle",
		Description: "This is risk1",
	}
	risk2 = models.Risk{
		ID:          "1234",
		State:       models.StateClosed,
		Title:       "MyTitle2",
		Description: "This is risk2",
	}
)

func TestWebserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webserver suite")
}

var _ = Describe("Webserver tests", func() {
	var (
		ctrl    *gomock.Controller
		mockDAO *mock_models.MockDAOInterface
		server  webserver.Server
		req     *http.Request
		w       *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDAO = mock_models.NewMockDAOInterface(ctrl)
		server = webserver.Server{
			Backend: mockDAO,
		}

		w = httptest.NewRecorder()
	})
	AfterEach(func() {
		ctrl.Finish()
	})
	Describe("POST risk requests", func() {
		BeforeEach(func() {
			buf := new(bytes.Buffer)
			_ = json.NewEncoder(buf).Encode(models.CreateRiskRequest{
				State:       risk1.State,
				Title:       risk1.Title,
				Description: risk1.Description,
			})
			req = httptest.NewRequest("POST", "/v1/risk", buf)
		})
		Context("green path", func() {
			BeforeEach(func() {
				mockDAO.EXPECT().SaveRisk(gomock.Any(), risk1.State, risk1.Title, risk1.Description).Return(&risk1, nil)
			})
			It("returns 200 status code", func() {
				server.SaveRisk(w, req)
				Expect(w.Result().StatusCode).To(Equal(200))
			})
			It("returns a risk object", func() {
				server.SaveRisk(w, req)
				var body models.Risk
				_ = json.NewDecoder(w.Body).Decode(&body)
				Expect(body).To(Equal(risk1))
			})
		})
		Context("DAO returns an error", func() {
			BeforeEach(func() {
				mockDAO.EXPECT().SaveRisk(gomock.Any(), risk1.State, risk1.Title, risk1.Description).Return(nil, errors.New("this is an error!"))
			})
			It("returns 500 status code", func() {
				server.SaveRisk(w, req)
				Expect(w.Result().StatusCode).To(Equal(500))
			})
		})
		Context("missing state property", func() {
			BeforeEach(func() {
				buf := new(bytes.Buffer)
				_ = json.NewEncoder(buf).Encode(models.CreateRiskRequest{
					Title:       risk1.Title,
					Description: risk1.Description,
				})
				req = httptest.NewRequest("POST", "/v1/risk", buf)
			})
			It("returns a 400 status code", func() {
				server.SaveRisk(w, req)
				Expect(w.Result().StatusCode).To(Equal(400))
			})
		})
		Context("invalid state property", func() {
			BeforeEach(func() {
				buf := new(bytes.Buffer)
				_ = json.NewEncoder(buf).Encode(models.CreateRiskRequest{
					State:       "notAValidState",
					Title:       risk1.Title,
					Description: risk1.Description,
				})
				req = httptest.NewRequest("POST", "/v1/risk", buf)
			})
			It("returns a 400 status code", func() {
				server.SaveRisk(w, req)
				Expect(w.Result().StatusCode).To(Equal(400))
			})
		})
	})
	Describe("GET risks requests", func() {
		var daoReturn = []models.Risk{risk1, risk2}
		BeforeEach(func() {
			req = httptest.NewRequest("GET", "/v1/risk", nil)
			mockDAO.EXPECT().GetAllRisks(gomock.Any()).Return(daoReturn)
		})
		It("returns a 200 status code", func() {
			server.GetAllRisks(w, req)
			Expect(w.Result().StatusCode).To(Equal(200))
		})
		It("returns a list from the DAO", func() {
			server.GetAllRisks(w, req)
			var body []models.Risk
			_ = json.NewDecoder(w.Body).Decode(&body)
			Expect(body).To(Equal(daoReturn))
		})
	})
	Describe("GET risk by ID requests", func() {
		BeforeEach(func() {
			req = httptest.NewRequest("GET", "/v1/risk/{id}", nil)
		})
		Context("green path", func() {
			BeforeEach(func() {
				req = mux.SetURLVars(req, map[string]string{"id": "abcd"})
				mockDAO.EXPECT().GetRiskByID(gomock.Any(), "abcd").Return(&risk1, nil)
			})
			It("returns a 200 status code", func() {
				server.GetRiskByID(w, req)
				Expect(w.Result().StatusCode).To(Equal(200))
			})
			It("returns a risk from the DAO", func() {
				server.GetRiskByID(w, req)
				var body models.Risk
				_ = json.NewDecoder(w.Body).Decode(&body)
				Expect(body).To(Equal(risk1))
			})
		})
		Context("missing the ID parameter", func() {
			It("returns a 400 status code", func() {
				server.GetRiskByID(w, req)
				Expect(w.Result().StatusCode).To(Equal(400))
			})
		})
		Context("DAO returns an error", func() {
			BeforeEach(func() {
				req = mux.SetURLVars(req, map[string]string{"id": "abcd"})
				mockDAO.EXPECT().GetRiskByID(gomock.Any(), "abcd").Return(nil, errors.New("this is an error!"))
			})
			It("returns a 500 status code", func() {
				server.GetRiskByID(w, req)
				Expect(w.Result().StatusCode).To(Equal(500))
			})
		})
		Context("DAO returns no risk object", func() {
			BeforeEach(func() {
				req = mux.SetURLVars(req, map[string]string{"id": "abcd"})
				mockDAO.EXPECT().GetRiskByID(gomock.Any(), "abcd").Return(nil, nil)
			})
			It("returns a 404 status code", func() {
				server.GetRiskByID(w, req)
				Expect(w.Result().StatusCode).To(Equal(404))
			})
		})
	})
})
