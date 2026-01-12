package absen_usecase

import (
	"errors"
	"fmt"
	"time"

	absen_model "github.com/DannyAss/users/internal/models/database_model/absen"
	absen_repository "github.com/DannyAss/users/internal/repositories/absen"
)

type (
	dtl = absen_model.AbsenDtl
	hdr = absen_model.AbsenHdr
	dto = absen_model.DTOAbsenRequest
)

type absenUsecase struct {
	repo absen_repository.IAbsemRepository
}

type IAbsenUsecase interface {
	AbsenToday(idClass int, detail dto) error
	GetAbsen(idClass int, mapelID int) bool
}

var (
	InternalServiceError = errors.New("Internal Service Error")
	UnprocessableEntity  = errors.New("Unprocessable Entity")
)

func NewAbsenUsecase(repo absen_repository.IAbsemRepository) IAbsenUsecase {
	return &absenUsecase{repo: repo}
}

func (a *absenUsecase) AbsenToday(idClass int, detail dto) error {
	var mdtl = make([]dtl, 0, len(detail.Detail))

	if idClass == 0 {
		return fmt.Errorf("%w : %s", UnprocessableEntity, "id tidak boleh kosong")
	}

	var catatan *string

	if len(detail.Detail) == 0 {
		msg := "murid masuk semua"
		catatan = &msg
	}

	courseId := uint64(detail.CourseId)

	uIdClass := uint64(idClass)

	mHdr := hdr{
		KelasID: &uIdClass,
		Catatan: catatan,
		Tanggal: time.Now().Format("2006/01/02"),
		MapelID: &courseId,
	}

	for _, c := range detail.Detail {
		absen := dtl{
			StudentID:  c.StudentID,
			Status:     c.Status,
			Keterangan: c.Keterangan,
		}

		mdtl = append(mdtl, absen)
	}

	if err := a.repo.InsertAbsen(mHdr, mdtl); err != nil {
		return fmt.Errorf("%w : %s", InternalServiceError, err.Error())
	}

	return nil
}

func (a *absenUsecase) GetAbsen(idClass int, mapelID int) bool {

	uMapelId := uint64(mapelID)

	uIdClass := uint64(idClass)

	today := time.Now().Format("2006/01/02")

	model := hdr{
		MapelID: &uMapelId,
		KelasID: &uIdClass,
		Tanggal: today,
	}

	cekAbsen := a.repo.GetAbsenToday(model)

	return cekAbsen

}
