package dao_test

import (
	"context"
	"testing"

	"arcticwolf.com/cutler/dao"
	"arcticwolf.com/cutler/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWebserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DAO suite")
}

var _ = Describe("DAO tests", func() {
	var (
		backend *dao.LocalCache
	)
	BeforeEach(func() {
		backend = &dao.LocalCache{}
	})
	Describe("save risk", func() {
		It("does not return an eror", func() {
			_, err := backend.SaveRisk(context.Background(), models.StateAccepted, "MyTitle", "MyDescription")
			Expect(err).ToNot(HaveOccurred())
		})
		It("returns the new risk object with a generated ID", func() {
			result, _ := backend.SaveRisk(context.Background(), models.StateAccepted, "MyTitle", "MyDescription")
			Expect(result.ID).ToNot(BeEmpty())
			Expect(result.State).To(Equal(models.StateAccepted))
			Expect(result.Title).To(Equal("MyTitle"))
			Expect(result.Description).To(Equal("MyDescription"))
		})
	})
	Describe("get all saved risks", func() {
		Context("there are saved risks", func() {
			BeforeEach(func() {
				_, err := backend.SaveRisk(context.Background(), models.StateAccepted, "MyTitle1", "this is risk1")
				Expect(err).ToNot(HaveOccurred())
				_, err = backend.SaveRisk(context.Background(), models.StateClosed, "MyTitle2", "this is risk2")
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns all saved risks", func() {
				result := backend.GetAllRisks(context.Background())
				Expect(len(result)).To(Equal(2))
			})
		})
		Context("there are no saved risks", func() {
			It("returns an empty list", func() {
				result := backend.GetAllRisks(context.Background())
				Expect(result).To(BeEmpty())
			})
		})
	})
	Describe("get risk by ID", func() {
		Context("the ID does not exist", func() {
			It("does not return an error", func() {
				_, err := backend.GetRiskByID(context.Background(), "doesNotExist")
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns an empty pointer", func() {
				result, _ := backend.GetRiskByID(context.Background(), "doesNotExist")
				Expect(result).To(BeNil())
			})
		})
		Context("the ID exists", func() {
			var risk *models.Risk
			BeforeEach(func() {
				var err error
				risk, err = backend.SaveRisk(context.Background(), models.StateAccepted, "MyTitle1", "this is risk1")
				Expect(err).ToNot(HaveOccurred())
			})
			It("does not return an error", func() {
				_, err := backend.GetRiskByID(context.Background(), risk.ID)
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns the risk", func() {
				result, _ := backend.GetRiskByID(context.Background(), risk.ID)
				Expect(result).To(Equal(risk))
			})
		})
	})
})
