package parser

import (
	"go/ast"
	"log"
	"regexp"
	"strings"
)

// AnnotationParser отвечает за извлечение аннотаций из комментариев
type AnnotationParser struct {
	verbose bool
}

// NewAnnotationParser создаёт новый парсер аннотаций
func NewAnnotationParser(verbose bool) *AnnotationParser {
	return &AnnotationParser{verbose: verbose}
}

// ParseFieldDirectives извлекает @evl:validate директивы из поля
func (ap *AnnotationParser) ParseFieldDirectives(field *ast.Field) []Directive {
	var directives []Directive
	re := regexp.MustCompile(`@evl:validate\s+([a-zA-Z_-]+)(?::([^\s@]+))?`)

	for _, cg := range []*ast.CommentGroup{field.Comment, field.Doc} {
		if cg != nil {
			for _, comment := range cg.List {
				dirs := ap.extractDirectives(comment.Text, re)
				directives = append(directives, dirs...)
			}
		}
	}
	return directives
}

// ParseFieldOaAnnotations извлекает @oa: аннотации из поля
func (ap *AnnotationParser) ParseFieldOaAnnotations(field *ast.Field) []OaAnnotation {
	var oas []OaAnnotation
	re := regexp.MustCompile(`@oa:([a-zA-Z_.-]+)\s+(.+)`)

	for _, cg := range []*ast.CommentGroup{field.Comment, field.Doc} {
		if cg != nil {
			for _, comment := range cg.List {
				oas = append(oas, ap.extractOaAnnotations(comment.Text, re)...)
			}
		}
	}
	return oas
}

// ParseStructOaAnnotations извлекает @oa: аннотации на уровне структуры
func (ap *AnnotationParser) ParseStructOaAnnotations(genDecl *ast.GenDecl, typeSpec *ast.TypeSpec) []OaAnnotation {
	var oas []OaAnnotation
	re := regexp.MustCompile(`@oa:([a-zA-Z_.-]+)\s+(.+)`)

	for _, cg := range []*ast.CommentGroup{genDecl.Doc, typeSpec.Doc} {
		if cg != nil {
			for _, comment := range cg.List {
				oas = append(oas, ap.extractOaAnnotations(comment.Text, re)...)
			}
		}
	}
	return oas
}

// ExtractDiscriminator извлекает discriminator и oneOf/anyOf из аннотаций структуры
func (ap *AnnotationParser) ExtractDiscriminator(s *Struct, annotations []OaAnnotation) {
	var disc *OaDiscriminator

	for _, ann := range annotations {
		if ap.verbose {
			log.Printf("DEBUG: Struct=%s OA[%s]=%q", s.Name, ann.Type, ann.Value)
		}

		switch ann.Type {
		case "discriminator.propertyName":
			if disc == nil {
				disc = &OaDiscriminator{Mapping: make(map[string]string)}
			}
			disc.PropertyName = strings.Trim(ann.Value, `"`)

		case "discriminator.mapping":
			if disc == nil {
				disc = &OaDiscriminator{Mapping: make(map[string]string)}
			}
			if parts := strings.SplitN(ann.Value, ":", 2); len(parts) == 2 {
				key := strings.Trim(strings.TrimSpace(parts[0]), `"`)
				val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
				disc.Mapping[key] = val
			}

		case "oneOf":
			for _, t := range strings.Split(ann.Value, ",") {
				s.OaOneOf = append(s.OaOneOf, strings.TrimSpace(t))
			}
		case "oneOf-ref":
			for _, r := range strings.Split(ann.Value, ",") {
				s.OaOneOfRefs = append(s.OaOneOfRefs, strings.TrimSpace(r))
			}
		case "anyOf":
			for _, t := range strings.Split(ann.Value, ",") {
				s.OaAnyOf = append(s.OaAnyOf, strings.TrimSpace(t))
			}
		case "anyOf-ref":
			for _, r := range strings.Split(ann.Value, ",") {
				s.OaAnyOfRefs = append(s.OaAnyOfRefs, strings.TrimSpace(r))
			}
		}
	}

	if disc != nil && disc.PropertyName != "" {
		s.Discriminator = disc
		if ap.verbose {
			log.Printf("DEBUG: Struct %s has discriminator: %+v", s.Name, disc)
		}
	}
}

// extractDirectives извлекает директивы из текста комментария
func (ap *AnnotationParser) extractDirectives(text string, re *regexp.Regexp) []Directive {
	var directives []Directive
	text = cleanComment(text)

	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			dir := Directive{
				Type: match[1],
				Raw:  match[0],
			}
			if len(match) >= 3 && match[2] != "" {
				param := strings.TrimSpace(strings.Trim(match[2], ` "'`))
				if dir.Type == "pattern" {
					dir.Params = []string{param}
				} else {
					params := strings.Split(param, ",")
					for i, p := range params {
						params[i] = strings.TrimSpace(p)
					}
					dir.Params = params
				}
			}
			directives = append(directives, dir)
		}
	}
	return directives
}

// extractOaAnnotations извлекает OA-аннотации из текста
func (ap *AnnotationParser) extractOaAnnotations(text string, re *regexp.Regexp) []OaAnnotation {
	var oas []OaAnnotation
	text = cleanComment(text)

	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			oas = append(oas, OaAnnotation{
				Type:  match[1],
				Value: strings.TrimSpace(match[2]),
			})
		}
	}
	return oas
}

// cleanComment удаляет маркеры комментариев из строки
func cleanComment(text string) string {
	text = strings.TrimPrefix(text, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	text = strings.TrimSuffix(text, "/")
	return strings.TrimSpace(text)
}
