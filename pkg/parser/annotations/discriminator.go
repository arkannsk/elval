package annotations

import "strings"

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
