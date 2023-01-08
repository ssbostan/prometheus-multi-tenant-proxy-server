package server

import (
	"net/url"
)

func splitForm(form url.Values, postForm url.Values) (url.Values, url.Values) {
	rawQueryForm := make(url.Values)
	for key, values := range form {
		if _, exists := postForm[key]; exists {
			if len(values) != len(postForm[key]) {
				for _, formValue := range values {
					isPostFormValue := false
					for _, postFormValue := range postForm[key] {
						if formValue == postFormValue {
							isPostFormValue = true
							break
						}
					}
					if isPostFormValue == false {
						rawQueryForm.Add(key, formValue)
					}
				}
			}
		} else {
			rawQueryForm[key] = values
		}
	}
	return rawQueryForm, postForm
}
