package web

import validator "gopkg.in/go-playground/validator.v9"

func IsRequestValid(m interface{}) (error){
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return err
	}
	return nil
}
