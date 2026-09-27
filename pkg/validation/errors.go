package validation

var (
	ProductRequestValidate = map[string]string{
		"Name.required":     "Не было передано наименование товара",
		"Name.min":          "Слишком короткое наименование товара",
		"Price.required":    "Не была передана цена товара",
		"Price.gt":          "Цена товара не может быть меньше 0.05",
		"Price.scale":       "Цена товара должна иметь 2 символа после запятой",
		"Quantity.required": "Не было указано количество товара",
		"Descriptions.max":  "Описание товара не может быть длиннее 2000 символов",
		"Images.max":        "Длина ссылки не может быть длиннее 2000 символов",
	}
	ProductDeleteValidate = map[string]string{
		"Name.min":         "Слишком короткое наименование товара",
		"Price.gt":         "Цена товара не может быть меньше 0.05",
		"Price.scale":      "Цена товара должна иметь 2 символа после запятой",
		"Descriptions.max": "Описание товара не может быть длиннее 2000 символов",
		"Images.max":       "Длина ссылки не может быть длиннее 2000 символов",
	}

	UserAuthValidate = map[string]string{
		"Phone.required": "Не был передан номер телефона",
		"Phone.phone":    "Передан некорректный номер телефона",
	}

	AuthConfirmValidate = map[string]string{
		"SessionId.required": "Не был передан идентификатор сессии",
		"SessionId.uuid":     "Недопустимый идентфикатор сессии",
		"Code.required":      "Не был передан код подтверждения",
		"Code.len":           "Неверный код",
	}

	OrderCreateValidate = map[string]string{
		"UserId.required":    "Не был передан идентификатор пользователя",
		"Products.required":  "Не были переданы товары в заказе",
		"Products.dive":      "Не удалось проверить корректность вложенного map",
		"Products.min":       "Не были добавлены товары в заказ",
		"ProductId.required": "Не было указан идентификатор товара",
		"Quantity.required":  "Не было указано количество товара",
	}
)
