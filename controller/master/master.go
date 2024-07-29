package master

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"test_backend_esb/helper"
	"test_backend_esb/initializer"
	"test_backend_esb/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MasterHandle(c *gin.Context) {
	// fmt.Println(c.Request.Method)
	switch c.Request.Method {
	case "GET":
		get(c)
	case "POST":
		post(c)
	case "PATCH":
		patch(c)
	case "DELETE":
		deleteDb(c)
	default:
		return
	}
}

func get(c *gin.Context) {
	var count int64
	r := c.Request
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unknown path or request type",
		})
		return
	}

	if len(path) == 3 {
		var item []model.GORMModel
		err := FindAndPreloadAll(initializer.DB, strings.ToLower(path[2]), &item, c.Request, &count)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":      item,
			"page":      page,
			"page_size": pageSize,
			"count":     &count,
		})
		return
	} else if len(path) == 4 {
		var item []model.GORMModel
		err := FindAndPreload(initializer.DB, strings.ToLower(path[2]), &item, path[3])
		// err = initializer.DB.Where("id = ?", path[3]).First(&item).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "404 Not Found",
			})
			return
		}

		var data model.GORMModel
		if len(item) > 0 {
			data = item[0]
		} else {
			data = nil
		}

		c.JSON(http.StatusOK, gin.H{
			"data":      data,
			"page":      page,
			"page_size": pageSize,
			"count":     count,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unknown path or request type",
		})
		return
	}
}

func batch(c *gin.Context, element string) {
	var items []map[string]interface{}
	err := c.ShouldBind(&items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var modelInstances []model.GORMModel
	for _, data := range items {
		modelInstance, err := helper.CreateGORMModel(element)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		dataID, ok := data["ID"].(float64)
		if ok {
			initializer.DB.First(&modelInstance, int(dataID))
		}

		elem := reflect.ValueOf(modelInstance).Elem()
		elemType := elem.Type()
		for key, value := range data {
			attrib := elem.FieldByName(key)
			if !attrib.IsValid() {
				for i := 0; i < elem.NumField(); i++ {
					field := elem.Field(i)
					if field.IsValid() && field.CanSet() {
						fieldTag := elemType.Field(i).Tag.Get("json")

						if fieldTag == key {
							// Convert the value to the field's type and set it
							if v, ok := ConvertToType(value, field.Type()); ok {
								if field.Kind() == reflect.Ptr {
									if field.IsNil() {
										field.Set(reflect.New(field.Type().Elem()))
									}
									field.Elem().Set(v)
								} else {
									field.Set(v)
								}
							}
							break // Stop looking after a match is found
						}
					}
				}
			} else {
				if v, ok := ConvertToType(value, attrib.Type()); ok {
					attrib.Set(v)
				}
			}
		}

		rs := initializer.DB.Save(modelInstance)
		if rs.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error_db": rs.Error.Error(),
			})
			return
		}
		modelInstances = append(modelInstances, modelInstance)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Ok!",
		"data":    modelInstances,
	})
}

func post(c *gin.Context) {
	r := c.Request
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unknown path or request type",
		})
		return
	} else if len(path) > 3 && path[3] == "batch" {
		batch(c, strings.ToLower(path[2]))
		return
	}
	item, err := helper.CreateGORMModel(strings.ToLower(path[2]))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.ShouldBindJSON(&item)
	fmt.Println(item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	rs := initializer.DB.Create(item)
	if rs.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_db": rs.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Ok!",
		"data":    item,
	})
}

func patch(c *gin.Context) {
	r := c.Request
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unknown path or request type",
		})
		return
	}

	// Get the model name from the URL path
	// modelName := strings.ToLower(path[2])
	item, err := helper.CreateGORMModel(strings.ToLower(path[2]))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var updates map[string]interface{}
	data, exists := c.Get("binder")
	if !exists {
		updates = make(map[string]interface{})
		err = c.ShouldBind(&updates)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		updates = data.(map[string]interface{})
	}
	dataID, ok := updates["ID"].(float64)
	if ok {
		initializer.DB.First(&item, int(dataID))
	}
	// Check if the field you want to set to zero is in the request.
	// If it is, use gorm.Expr to force GORM to update it.
	// for key, value := range updates {
	// Check if the value is zero
	// You might need a more complex check here depending on what types you're expecting

	// }

	elem := reflect.ValueOf(item).Elem()
	elemType := elem.Type()
	for key, value := range updates {
		if value == 0 {
			// Replace the zero value with a GORM expression of zero
			updates[key] = gorm.Expr("?", 0)
		}
		if key == "acl" {
			aclJson, err := json.Marshal(value)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			value = aclJson

		}
		attrib := elem.FieldByName(key)
		if !attrib.IsValid() {
			for i := 0; i < elem.NumField(); i++ {
				field := elem.Field(i)
				if field.IsValid() && field.CanSet() {
					fieldTag := elemType.Field(i).Tag.Get("json")

					if fieldTag == key {
						// Convert the value to the field's type and set it
						if v, ok := ConvertToType(value, field.Type()); ok {
							if field.Kind() == reflect.Ptr {
								if field.IsNil() {
									field.Set(reflect.New(field.Type().Elem()))
								}
								field.Elem().Set(v)
							} else {
								field.Set(v)
							}
						} else {
							if key == "acl" {
								field.Set(v)
							}
						}
						break // Stop looking after a match is found
					}
				}
			}
		} else {
			if v, ok := ConvertToType(value, attrib.Type()); ok {
				attrib.Set(v)
			}
		}
	}
	rs := initializer.DB.Save(item)
	// rs := initializer.DB.Model(item).Where("id = ?", updates["ID"]).Updates(updates)
	if rs.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": rs.Error.Error(),
		})
		return
	}

	// // Update the acl field separately

	// rs = initializer.DB.Model(item).Where("id = ?", updates["ID"]).Update("acl", aclJson)
	// if rs.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": rs.Error.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Ok!",
	})
}

func deleteDb(c *gin.Context) {
	r := c.Request
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unknown path or request type",
		})
		return
	}
	instance := strings.ToLower(path[2])
	item, err := helper.CreateGORMModel(instance)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.ShouldBind(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	rs := initializer.DB.Delete(&item, item)

	if rs.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": rs.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ok!",
		"data":    item,
	})
}
