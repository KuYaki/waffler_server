package service

import "errors"

var badOrderParams = errors.New("Incorrect order parameters")

// Function checks for contradicting order parameters
// Bad parameters are of type condition and condition_desc
// Or of type condition and condition (repeats)
func checkOrderParams(orders []string) error {
	for i, order := range orders {
		for j, otherOrder := range orders {
			if order+"_desc" == otherOrder {
				return badOrderParams
			}

			if order == otherOrder && i != j {
				return badOrderParams
			}
		}
	}
	return nil
}

var orderRecords = map[string]string{"score": "score ASC", "score_desc": "score DESC",
	"time": "created_at ASC", "time_desc": "created_at DESC"}

func convertRecordOrder(order []string) ([]string, error) {
	err := checkOrderParams(order)
	if err != nil {
		return nil, err
	}

	newOrder := make([]string, 0, len(order))
	for _, o := range order {
		val, found := orderRecords[o]
		if !found {
			return nil, badOrderParams
		}
		newOrder = append(newOrder, val)
	}
	return newOrder, nil
}

var orderSources = map[string]string{"name": "name ASC", "name_desc": "name DESC", "source": "source_type ASC", "source_desc": "source_type DESC",
	"waffler": "waffler_score ASC", "waffler_desc": "waffler_score DESC", "racizm": "racism_score ASC", "racizm_desc": "racism_score DESC"}

func convertSourceOrder(order []string) ([]string, error) {
	err := checkOrderParams(order)
	if err != nil {
		return nil, err
	}

	newOrder := make([]string, 0, len(order))
	for _, o := range order {
		val, found := orderSources[o]
		if !found {
			return nil, badOrderParams
		}
		newOrder = append(newOrder, val)
	}
	return newOrder, nil
}
