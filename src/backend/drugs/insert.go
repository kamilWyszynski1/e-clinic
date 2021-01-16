package drugs

import (
	"e-clinic/src/backend/models"
	"strconv"

	"github.com/gocraft/dbr"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func InsertToDb(db *dbr.Session, log logrus.FieldLogger, drug *ProduktLeczniczy) {
	id, err := strconv.Atoi(drug.ID)
	if err != nil {
		log.WithError(err).Error("faield to parse id")
		return
	}
	d := models.Drug{
		ID:                id,
		Name:              drug.NazwaProduktu,
		TypeOfPreparation: drug.RodzajPreparatu,
		CommonName:        drug.NPS,
		Strength:          drug.Moc,
		Shape:             drug.Postac,
	}
	if err := d.Insert(db); err != nil {
		log.Error(err)
	}
	for _, sub := range drug.SubstancjeCzynne.SubstancjaCzynna {
		s := models.Substance{
			ID:   uuid.NewV4(),
			Name: sub,
		}
		if err := s.Insert(db); err != nil {
			log.Error(err)
		}

		c := models.Composition{
			ID:        uuid.NewV4(),
			Drug:      d.ID,
			Substance: s.ID,
		}
		if err := c.Insert(db); err != nil {
			log.Error(err)
		}
	}

}
