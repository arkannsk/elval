package annotations

import (
	"strconv"
	"strings"
)

// OaResponseAnnotation результат парсинга @oa:response
type OaResponseAnnotation struct {
	StatusCode  int
	Description string
	MediaTypes  []string
}

// OaResponseTarget интерфейс для структуры с response аннотациями
type OaResponseTarget interface {
	GetOaResponses() []OaResponseAnnotation
	SetOaResponses([]OaResponseAnnotation)
}

// OaDiscriminator структура для хранения информации о дискриминаторе
type OaDiscriminator struct {
	PropertyName string
	Mapping      map[string]string
}

// DiscriminatorTarget интерфейс для структуры, которую мы хотим заполнить аннотациями
type DiscriminatorTarget interface {
	GetDiscriminator() *OaDiscriminator
	SetDiscriminator(d *OaDiscriminator)
	GetOaOneOf() []string
	SetOaOneOf([]string)
	GetOaOneOfRefs() []string
	SetOaOneOfRefs([]string)
	GetOaAnyOf() []string
	SetOaAnyOf([]string)
	GetOaAnyOfRefs() []string
	SetOaAnyOfRefs([]string)
}

// ExtractDiscriminatorData извлекает данные дискриминатора и oneOf/anyOf из аннотаций
func ExtractDiscriminatorData(target DiscriminatorTarget, annotations []OaAnnotation) {
	var disc *OaDiscriminator

	for _, ann := range annotations {
		switch ann.Type {
		case "discriminator":
			if disc == nil {
				disc = &OaDiscriminator{Mapping: make(map[string]string)}
			}
			disc.PropertyName = trimQuotes(ann.Value)
		case "discriminator.propertyName":
			if disc == nil {
				disc = &OaDiscriminator{Mapping: make(map[string]string)}
			}
			disc.PropertyName = trimQuotes(ann.Value)

		case "discriminator.mapping":
			if disc == nil {
				disc = &OaDiscriminator{Mapping: make(map[string]string)}
			}
			if parts := strings.SplitN(ann.Value, ":", 2); len(parts) == 2 {
				key := trimQuotes(strings.TrimSpace(parts[0]))
				val := trimQuotes(strings.TrimSpace(parts[1]))
				disc.Mapping[key] = val
			}

		case "oneOf":
			target.SetOaOneOf(parseList(ann.Value))
		case "oneOf-ref":
			target.SetOaOneOfRefs(parseList(ann.Value))
		case "anyOf":
			target.SetOaAnyOf(parseList(ann.Value))
		case "anyOf-ref":
			target.SetOaAnyOfRefs(parseList(ann.Value))
		}
	}

	if disc != nil && disc.PropertyName != "" {
		target.SetDiscriminator(disc)
	}

	// Если target поддерживает response аннотации, обрабатываем их
	if respTarget, ok := target.(OaResponseTarget); ok {
		for _, ann := range annotations {
			if ann.Type == "response" {
				respTarget.SetOaResponses(append(respTarget.GetOaResponses(),
					parseResponseAnnotation(ann.Value)))
			}
		}
	}
}

// parseResponseAnnotation парсит значение аннотации @oa:response
// Формат аннотации: @oa:response "200" "application/json,application/xml"
// После trimQuotes в парсере значение становится:
//   200" "application/json,application/xml
// (trimQuotes удаляет только внешние кавычки)
// Также поддерживается формат без кавычек:
//   @oa:response 200 application/json,application/xml
func parseResponseAnnotation(value string) OaResponseAnnotation {
	// Разбираем значение, учитывая возможные кавычки
	// value может быть: 200" "application/json,application/xml
	// или: 200 application/json,application/xml

	// Удаляем все кавычки для упрощения парсинга
	cleaned := strings.ReplaceAll(value, `"`, "")
	cleaned = strings.ReplaceAll(cleaned, `'`, "")

	parts := strings.SplitN(strings.TrimSpace(cleaned), " ", 2)
	statusCode := 0
	if len(parts) > 0 {
		statusCode, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
	}
	mediaTypes := make([]string, 0)
	if len(parts) > 1 {
		for _, mt := range strings.Split(parts[1], ",") {
			mt = strings.TrimSpace(mt)
			if mt != "" {
				mediaTypes = append(mediaTypes, mt)
			}
		}
	}
	return OaResponseAnnotation{
		StatusCode: statusCode,
		MediaTypes: mediaTypes,
	}
}

func parseList(value string) []string {
	if value == "" {
		return nil
	}
	var result []string
	for _, t := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
