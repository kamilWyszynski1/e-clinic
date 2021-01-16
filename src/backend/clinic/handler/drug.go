package handler

import (
	"e-clinic/src/backend/clinic"
	"e-clinic/src/backend/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocraft/dbr"
	"github.com/sirupsen/logrus"
)

const queryDrugs = `select d.* from drug d
join composition c on d.id = c.drug
WHERE lower(name) like %s
group by d.id
order by count(*)  desc
OFFSET ? LIMIT ?`

func (h Handler) GetDrugs(prefix string, offset int, limit int) (*clinic.Drugs, int, error) {
	if limit == 0 {
		limit = 10
	}
	var drugs []*models.Drug

	_, err := h.db.SelectBySql(fmt.Sprintf(queryDrugs, fmt.Sprintf("'%s%%'", strings.ToLower(prefix))), offset, limit).Load(&drugs)
	if errors.Is(err, dbr.ErrNotFound) {
		return nil, http.StatusNoContent, nil
	} else if err != nil {
		logrus.WithError(err).Error("failed to query drugs")
		return nil, http.StatusInternalServerError, nil
	}
	if len(drugs) == 0 {
		logrus.Debug("no content")
		return nil, http.StatusNoContent, nil
	}
	return &clinic.Drugs{
		Drugs: drugs,
		Len:   len(drugs),
	}, http.StatusOK, nil
}

const querySubstances = `select s.* from substance s
join composition c on s.id = c.substance
join drug d on c.drug = d.id
where d.id = ?`

func (h Handler) GetDrug(drugID int) (*clinic.DrugWithSubstances, int, error) {
	d, err := models.DrugByID(h.db, drugID)
	if errors.Is(err, dbr.ErrNotFound) {
		return nil, http.StatusNoContent, nil
	} else if err != nil {
		logrus.WithError(err).Error("failed to query drug")
		return nil, http.StatusInternalServerError, nil
	}

	var subs []*models.Substance
	_, err = h.db.SelectBySql(querySubstances, drugID).Load(&subs)
	if errors.Is(err, dbr.ErrNotFound) {
		return nil, http.StatusNoContent, nil
	} else if err != nil {
		logrus.WithError(err).Error("failed to query substances")
		return nil, http.StatusInternalServerError, nil
	}

	return &clinic.DrugWithSubstances{
		Drug:       d,
		Substances: subs,
	}, http.StatusOK, nil
}

const neoReplacementsQuery = `CALL gds.nodeSimilarity.stream('drug')
YIELD node1, node2, similarity
WHERE gds.util.asNode(node1).id = '%d' AND similarity > %f
RETURN 
	gds.util.asNode(node2).id AS id, 
    gds.util.asNode(node2).moc as moc,
    gds.util.asNode(node2).nazwa as nazwa,
    gds.util.asNode(node2).nazwa_powszechna as nazwa_powszechna,
    gds.util.asNode(node2).typ as typ,
    gds.util.asNode(node2).postac as postac,
    similarity
ORDER BY similarity DESCENDING, id, moc, nazwa, nazwa_powszechna, typ, postac
`

func (h Handler) GetReplacement(drugID int, minSimilarity float64) (*clinic.Drugs, int, error) {
	log := h.log.WithField("method", "GetReplacement")
	res, err := h.neoCli.Run(fmt.Sprintf(neoReplacementsQuery, drugID, minSimilarity), nil)
	if err != nil {
		log.WithError(err).Error("failed to run cypher query")
		return nil, http.StatusInternalServerError, nil
	}

	f := func(i interface{}, _ bool) string {
		return i.(string)
	}

	drugs := make([]*models.Drug, 0)
	for res.Next() {
		stringID := f(res.Record().Get("id"))
		id, err := strconv.Atoi(stringID)
		if err != nil {
			log.WithError(err).Error("failed to parse id")
			return nil, http.StatusInternalServerError, nil
		}
		d := &models.Drug{
			ID:                id,
			Name:              f(res.Record().Get("nazwa")),
			TypeOfPreparation: f(res.Record().Get("typ")),
			CommonName:        f(res.Record().Get("nazwa_powszechna")),
			Strength:          f(res.Record().Get("moc")),
			Shape:             f(res.Record().Get("postac")),
		}
		drugs = append(drugs, d)
	}
	return &clinic.Drugs{
		Drugs: drugs,
		Len:   len(drugs),
	}, http.StatusOK, nil
}
