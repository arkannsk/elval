package parser

import (
	"go/ast"

	ann "github.com/arkannsk/elval/pkg/parser/annotations"
)

// AnnotationParser отвечает за извлечение аннотаций из комментариев
type AnnotationParser struct {
	verbose bool
}

// NewAnnotationParser создаёт новый парсер аннотаций
func NewAnnotationParser(verbose bool) *AnnotationParser {
	return &AnnotationParser{verbose: verbose}
}

// CollectCommentTexts собирает все тексты комментариев из группы
func CollectCommentTexts(groups ...*ast.CommentGroup) []string {
	var texts []string
	for _, cg := range groups {
		if cg != nil {
			for _, c := range cg.List {
				texts = append(texts, ann.CleanComment(c.Text))
			}
		}
	}
	return texts
}

// ParseFieldDirectives извлекает @evl:validate директивы из поля
func (ap *AnnotationParser) ParseFieldDirectives(field *ast.Field) []ann.Directive {
	texts := CollectCommentTexts(field.Doc, field.Comment)
	return ann.ParseDirectivesFromTexts(texts)
}

// ParseFieldOaAnnotations извлекает @oa: аннотации из поля
func (ap *AnnotationParser) ParseFieldOaAnnotations(field *ast.Field) []ann.OaAnnotation {
	texts := CollectCommentTexts(field.Doc, field.Comment)
	return ann.ParseOaAnnotationsFromTexts(texts)
}

// ParseStructOaAnnotations извлекает @oa: аннотации из структуры
func (ap *AnnotationParser) ParseStructOaAnnotations(genDecl *ast.GenDecl, typeSpec *ast.TypeSpec) []ann.OaAnnotation {
	texts := CollectCommentTexts(genDecl.Doc, typeSpec.Doc)
	return ann.ParseOaAnnotationsFromTexts(texts)
}

// ExtractDiscriminator извлекает discriminator и oneOf/anyOf из аннотаций структуры
func (ap *AnnotationParser) ExtractDiscriminator(s *Struct, rawAnnotations []ann.OaAnnotation) {
	target := &structAdapter{s: s}
	ann.ExtractDiscriminatorData(target, rawAnnotations)
}

// ProcessFieldAnnotations обрабатывает аннотации поля для генератора
func (ap *AnnotationParser) ProcessFieldAnnotations(rawAnnotations []ann.OaAnnotation) ann.FieldAnnotationResult {
	return ann.ProcessFieldAnnotations(rawAnnotations)
}

// ParseFieldDecorators @evl:decor
func (ap *AnnotationParser) ParseFieldDecorators(field *ast.Field) []ann.Decorator {
	texts := CollectCommentTexts(field.Doc, field.Comment)
	return ann.ParseDecoratorsFromTexts(texts)
}

// structAdapter реализует интерфейс DiscriminatorTarget для parser.Struct
type structAdapter struct {
	s *Struct
}

func (a *structAdapter) GetDiscriminator() *ann.OaDiscriminator  { return a.s.Discriminator }
func (a *structAdapter) SetDiscriminator(d *ann.OaDiscriminator) { a.s.Discriminator = d }
func (a *structAdapter) GetOaOneOf() []string                    { return a.s.OaOneOf }
func (a *structAdapter) SetOaOneOf(v []string)                   { a.s.OaOneOf = v }
func (a *structAdapter) GetOaOneOfRefs() []string                { return a.s.OaOneOfRefs }
func (a *structAdapter) SetOaOneOfRefs(v []string)               { a.s.OaOneOfRefs = v }
func (a *structAdapter) GetOaAnyOf() []string                    { return a.s.OaAnyOf }
func (a *structAdapter) SetOaAnyOf(v []string)                   { a.s.OaAnyOf = v }
func (a *structAdapter) GetOaAnyOfRefs() []string                { return a.s.OaAnyOfRefs }
func (a *structAdapter) SetOaAnyOfRefs(v []string)               { a.s.OaAnyOfRefs = v }
