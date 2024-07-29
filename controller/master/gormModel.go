package master

import (
	"errors"
	"fmt"
	"test_backend_esb/helper"
	"test_backend_esb/model"

	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FindAndPreloadAll(db *gorm.DB, modelName string, result *[]model.GORMModel, r *http.Request, count *int64) error {

	element, err := helper.CreateGORMModel(modelName)
	if err != nil {
		return err
	}

	// Assert GORMModel interface back to its original struct
	modelType := reflect.TypeOf(element).Elem()

	// Create an empty slice of the model type
	sliceType := reflect.SliceOf(modelType)
	slice := reflect.New(sliceType)

	// Perform query with Preload and clause.Associations
	q := r.URL.Query()

	// dateFields := []string{
	// 	"updated_at",
	// 	"created_at",
	// 	"done_date",
	// 	"pending_date",
	// 	"deadline",
	// 	"published_date",
	// 	"transaction_date",
	// }
	startDate := q.Get("start_date")
	endDate := q.Get("end_date")
	dateField := q.Get("date_field")
	// res := false
	// for _, v := range dateFields {
	// 	if v == dateField {
	// 		res = true
	// 		break
	// 	}
	// }
	if startDate != "" && endDate != "" && dateField != "" {
		sd, err := time.Parse("2006-01-02", startDate)
		// sd = sd.Truncate(24 * time.Hour)
		sds := sd.Format("2006-01-02") + " 00:00:00"
		if err != nil {
			return err
		}
		ed, err := time.Parse("2006-01-02", endDate)
		// ed = ed.Truncate(24 * time.Hour)
		eds := ed.Format("2006-01-02") + " 23:59:59"
		if err != nil {
			return err
		}
		// subquery := initializer.DB.Select("ticket_id, max(created_at) as created_at").Table("lkms").Group("ticket_id")
		// db = db.Joins("join (?) q on q.ticket_id = tickets.id", subquery).Where("q.created_at BETWEEN ? AND ?", sd, ed)
		db = db.Where("`"+dateField+"` BETWEEN ? AND ?", sds, eds)
	}

	filter := q.Get("filter")
	if filter != "" {
		splits := strings.Split(filter, ",")
		var queryBuilder strings.Builder
		var args []interface{}
		queryBuilder.WriteString("(")
		for j, v := range splits {
			if j != 0 {
				queryBuilder.WriteString(" OR ")
			}
			queryBuilder.WriteString("(")
			h := 0
			for i := 0; i < modelType.NumField(); i++ {
				field := modelType.Field(i)

				// Skip fields with unsupported types for string matching
				if field.Type.Kind() != reflect.String {
					continue
				}
				tag := field.Tag.Get("json")
				gormTag := field.Tag.Get("gorm")
				if tag == "-" || gormTag == "-" {
					continue
				}
				// Add an OR clause to the query for each field
				if h != 0 {
					queryBuilder.WriteString(" OR ")
				}
				h++
				queryBuilder.WriteString(fmt.Sprintf("`%s` LIKE ?", tag))

				// Add the filter with wildcards to the args
				args = append(args, fmt.Sprintf("%%%s%%", v))
			}
			queryBuilder.WriteString(")")
		}
		queryBuilder.WriteString(")")
		db = db.Where(queryBuilder.String(), args...)
	}

	order := q.Get("order")
	if order != "" {
		orders := strings.Split(order, ":")
		if len(orders) == 2 {
			var direction string
			a, _ := strconv.Atoi(orders[1])
			if a > 0 {
				direction = "ASC"
			} else {
				direction = "DESC"
			}

			db = db.Order(fmt.Sprintf("%s %s", orders[0], direction))
		}
	}

	if param := q.Get("param"); param != "" {
		var queryBuilder strings.Builder
		var args []interface{}
		params := strings.Split(param, ",")
		fmt.Println(params[0])
		for _, v := range params {
			splits := strings.Split(v, ":")
			fmt.Println(splits)
			if len(splits) != 2 {
				return errors.New("invalid params")
			}

			// Add an OR clause to the query for each field
			if queryBuilder.Len() > 0 {
				queryBuilder.WriteString(" AND ")
			}
			if splits[1] == "null" {
				queryBuilder.WriteString(fmt.Sprintf("%s IS NULL", splits[0]))
			} else if strings.Contains(splits[1], "!") {
				if strings.Contains(splits[1], "null") {
					queryBuilder.WriteString(fmt.Sprintf("%s IS NOT NULL", splits[0]))
				} else {
					queryBuilder.WriteString(fmt.Sprintf("%s NOT LIKE ?", splits[0]))
					_, i := utf8.DecodeRuneInString(splits[1])
					splits[1] = splits[1][i:]
					args = append(args, splits[1])
				}
			} else if strings.Contains(splits[1], "%") {
				queryBuilder.WriteString(fmt.Sprintf("%s LIKE ?", splits[0]))
				args = append(args, splits[1])
			} else if strings.Contains(splits[1], "<>") {
				queryBuilder.WriteString(fmt.Sprintf("%s <> ?", splits[0]))
				args = append(args, splits[1])
			} else {
				queryBuilder.WriteString(fmt.Sprintf("%s = ?", splits[0]))
				args = append(args, splits[1])
			}

		}
		fmt.Println(queryBuilder.String())
		db = db.Where(queryBuilder.String(), args...)
	}

	// Nested Preload Relationship
	preload := q.Get("preload")
	if preload == "" {
		if relationships, ok := helper.RelationshipRegistry[modelType]; ok {
			for _, relationship := range relationships {
				fmt.Println(relationship)
				db = db.Preload(clause.Associations).Preload(relationship)
			}
		} else {
			db = db.Model(element).Preload(clause.Associations)
		}
	} else if preload != "none" {
		rels := strings.Split(preload, ",")
		for _, rel := range rels {
			splits := strings.Split(rel, "-")
			if len(splits) == 2 {
				db = db.Preload(splits[0], splits[1])
			} else {
				db = db.Preload(rel)
			}
		}
	}

	joins := q.Get("join")
	if joins != "" {
		join := strings.Split(joins, ",")
		for _, rel := range join {
			splits := strings.Split(rel, "-")
			if len(splits) == 2 {
				db = db.Joins(splits[0]).Where(splits[1])
			} else {
				db = db.Joins(rel)
			}
		}
	}

	db.Model(element).Count(count)
	pageSize := r.URL.Query().Get("page_size")
	if pageSize == "" {
		pageSize = "20"
	}
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		return err
	}
	if size != -1 {
		db.Scopes(Paginate(r))
	}

	err = db.Find(slice.Interface()).Error
	if err != nil {
		return err
	}

	// Convert the slice of the original struct to a slice of GORMModel
	sliceLen := slice.Elem().Len()
	elements := make([]model.GORMModel, sliceLen)
	for i := 0; i < sliceLen; i++ {
		elements[i] = slice.Elem().Index(i).Interface().(model.GORMModel)
	}

	// Assign the result to the provided result pointer
	*result = elements
	return nil
}

func FindAndPreload(db *gorm.DB, modelName string, result *[]model.GORMModel, id string) error {
	element, err := helper.CreateGORMModel(modelName)
	if err != nil {
		return err
	}

	// Assert GORMModel interface back to its original struct
	modelType := reflect.TypeOf(element).Elem()

	// Create an empty slice of the model type
	sliceType := reflect.SliceOf(modelType)
	slice := reflect.New(sliceType)

	// Nested Preload Relationship
	if relationships, ok := helper.RelationshipRegistry[modelType]; ok {
		for _, relationship := range relationships {
			fmt.Println(relationship)
			db = db.Preload(clause.Associations).Preload(relationship)
		}
	} else {
		db = db.Model(element).Preload(clause.Associations)
	}

	// Perform query with Preload and clause.Associations
	db.Where("id = ?", id).Find(slice.Interface())

	// Convert the slice of the original struct to a slice of GORMModel
	sliceLen := slice.Elem().Len()
	elements := make([]model.GORMModel, sliceLen)
	for i := 0; i < sliceLen; i++ {
		elements[i] = slice.Elem().Index(i).Interface().(model.GORMModel)
	}

	// Assign the result to the provided result pointer
	*result = elements
	return nil
}

func CreateGORMInstance(modelName string) (interface{}, error) {
	// Map the model names to their respective types
	modelType := helper.TypeRegistry[modelName]

	// Create a new instance of the model type
	modelInstance := reflect.New(modelType).Interface()

	return modelInstance, nil
}

// Function to convert a value to a specified type
func ConvertToType(value interface{}, targetType reflect.Type) (reflect.Value, bool) {
	// Implement conversion logic based on your specific needs
	// This is a simplified example; you may need to handle different types and cases
	switch targetType.Kind() {
	case reflect.Int64:
		if v, ok := value.(float64); ok {
			return reflect.ValueOf(int64(v)), true
		}
	case reflect.Float32:
		if v, ok := value.(float64); ok {
			return reflect.ValueOf(float32(v)), true
		}
	case reflect.Pointer:
		// fmt.Println("pointer")
		valType := targetType.Elem()
		// fmt.Println(valType)
		// Convert the value to the target type
		return ConvertToType(value, valType)
	case reflect.Bool:
		if v, ok := value.(bool); ok {
			return reflect.ValueOf(bool(v)), true
		}
	case reflect.Int:
		if v, ok := value.(float64); ok {
			return reflect.ValueOf(int(v)), true
		}
	case reflect.Uint:
		if v, ok := value.(float64); ok {
			return reflect.ValueOf(uint(v)), true
		}
	case reflect.Struct:
		if targetType == reflect.TypeOf(time.Time{}) {
			switch v := value.(type) {
			case string:
				t, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return reflect.Value{}, false
				}
				return reflect.ValueOf(t), true
				// Add other cases if necessary, e.g., converting from other formats or types
			}
		}
	case reflect.String:
		if v, ok := value.(string); ok {
			return reflect.ValueOf(string(v)), true
		}
	}

	return reflect.ValueOf(value), false
}
